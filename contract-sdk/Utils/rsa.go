package Utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fabric-learn/fabric-did/contract-sdk/Model"
	"fmt"
	"os"
)

func GenerateRSAKeys(priFilename, pubFilename string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	priFile, err := os.Create(priFilename)
	if err != nil {
		return err
	}
	err = pem.Encode(priFile, &block)
	if err != nil {
		return nil
	}
	defer priFile.Close()

	pubKey := privateKey.PublicKey
	dePubkey, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		return err
	}
	block = pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: dePubkey,
	}
	pubFile, err := os.Create(pubFilename)
	if err != nil {
		return err
	}
	err = pem.Encode(pubFile, &block)
	if err != nil {
		return nil
	}
	defer pubFile.Close()
	return nil
}

func SigntureRSA(did Model.DID, priPem string) ([]byte, error) {
	block, _ := pem.Decode([]byte(priPem))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the priate key")
	}

	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key")
	}
	if did.Proof.Signature != nil {
		return nil, fmt.Errorf("doc proof is not empty")
	}

	didAsBytes, err := json.Marshal(did)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed")
	}

	hash := sha256.New()
	hash.Write(didAsBytes)
	hashText := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, hashText)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// VerifyRSA verify data by rsa
func VerifyRSA(publicPem string, doc Model.DID, signature []byte) error {
	block, _ := pem.Decode([]byte(publicPem))
	if block == nil {
		return fmt.Errorf("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse DER encoded public key")
	}
	publicKey := pub.(*rsa.PublicKey)
	// 代认证的数据
	var data Model.DID
	data.ID = doc.ID
	data.Authenticaions = doc.Authenticaions

	verifyBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to Marshal doc ")
	}

	hash := sha256.New()
	hash.Write(verifyBytes)
	hashText := hash.Sum(nil)

	// verify
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashText, signature)
	if err != nil {
		return fmt.Errorf("%s proof signture incorrect, error -> %v\n", doc.ID, err)
	}
	return nil
}
