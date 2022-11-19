package Utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fabric-learn/fabric-did/contract-sdk/Model"
	"fmt"
	"math/big"
	"os"
)

// 获取哈希
func getHash(data []byte) []byte {
	hash256 := sha256.New()
	hash256.Write(data)
	return hash256.Sum(nil)
}

func GenerateECCKeys(priFilename, pubFilename string) error {
	// 1. 使用 eccdsa 生存密钥对
	privatekey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return err
	}
	// 2. 将私钥本地化,使用x509进行序列化
	privateAsBytes, err := x509.MarshalECPrivateKey(privatekey)
	if err != nil {
		return err
	}
	//3. 将它转换成pem的格式编码
	block := pem.Block{
		Type:  "ECC PRIVATE KEY",
		Bytes: privateAsBytes,
	}
	// 在本地创建pem文件
	privateFile, err := os.Create(priFilename)
	if err != nil {
		return err
	}
	defer privateFile.Close()
	// 进行pem编码
	err = pem.Encode(privateFile, &block)
	if err != nil {
		return err
	}

	// 公钥同理
	publicKey := privatekey.PublicKey
	publicDerBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	publicPemBlock := pem.Block{
		Type:  "Ecc PublicKey",
		Bytes: publicDerBytes,
	}
	publicFile, err := os.Create(pubFilename)
	if err != nil {
		return err
	}
	defer publicFile.Close()
	err = pem.Encode(publicFile, &publicPemBlock)
	if err != nil {
		return err
	}
	return nil
}

func SigByECC(did Model.DID, priPem string) (rText, sText []byte, err error) {
	block, _ := pem.Decode([]byte(priPem))
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private.pem, err -> %v", err.Error())
	}
	if did.Proof.Signature != nil {
		return nil, nil, fmt.Errorf("did proof is not empty")
	}

	didAsBytes, err := json.Marshal(did)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal did, err -> %v", err.Error())
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, getHash(didAsBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sig , err -> %v", err.Error())
	}
	rText, _ = r.MarshalText()
	sText, _ = s.MarshalText()
	return
}

func VerifyElliptic(doc Model.DID, rText, sText []byte, publicPem string) (bool, error) {
	block, _ := pem.Decode([]byte(publicPem))
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public pem. err -> %v", err.Error())
	}
	publicKey := pubInterface.(*ecdsa.PublicKey)

	// 验证数据
	var data Model.DID
	data.ID = doc.ID
	data.Authenticaions = doc.Authenticaions
	verifyBytes, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Errorf("failed to marshal veridy data. err -> %v", err.Error())
	}

	hashText := getHash(verifyBytes)

	var r, s big.Int
	r.UnmarshalText(rText)
	s.UnmarshalText(sText)

	ok := ecdsa.Verify(publicKey, hashText, &r, &s)
	return ok, nil
}
