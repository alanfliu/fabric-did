package main

import (
	"encoding/json"
	"fabric-learn/fabric-did/contract-sdk/Model"
	"fabric-learn/fabric-did/contract-sdk/Utils"
	"fabric-learn/fabric-did/contract-sdk/controller"
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	contract, err := Utils.SDKInit()
	fmt.Println("SDK初始化成功，成功接入Fabric网络")
	if err != nil {
		fmt.Printf("初始化SDK失败->%v", err.Error())
		return
	}

	time := CreateSelfDocument("004", contract)
	if time == 0 {
		return
	}
	fmt.Println("创建私有DID耗时:", time)
	time = CreatePubDocument("lf2", "004", 0, contract)
	if time == 0 {
		return
	}
	fmt.Println("创建公有DID耗时:", time)
	time, err = GetAndSaveDocument("did:cid:lf2", contract)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("获取DID耗时:", time)
}

// CreatePubDocument 创建公有DID
func CreatePubDocument(id, issuer string, kid int, contract *gateway.Contract) time.Duration {
	priPemFilename := fmt.Sprintf("./keystore/%s-private.pem", id)
	pubPemFilename := fmt.Sprintf("./keystore/%s-public.pem", id)
	// 创建公共DID
	err := Utils.GenerateECCKeys(priPemFilename, pubPemFilename)
	if err != nil {
		fmt.Println("创建密钥失败， err ->", err.Error())
		return 0
	}
	publiePem := Utils.ReadFile(pubPemFilename)
	doc := Model.DID{
		ID: fmt.Sprintf("did:cid:%s", id),
		Authenticaions: []Model.Authenticaion{
			{
				Kid:        fmt.Sprintf("did:cid:%s#key-1", id),
				Method:     "ECC",
				Controller: fmt.Sprintf("did:cid:%s", id),
				PublicPem:  string(publiePem),
			},
		},
	}
	// "./keystore/1000-private.pem" 读取issuer的私钥
	priPem := Utils.ReadFile(fmt.Sprintf("./keystore/%s-private.pem", issuer))
	// issuer 私钥签名
	rText, sTest, err := Utils.SigByECC(doc, string(priPem))
	if err != nil {
		fmt.Println("签名失败 err ->", err.Error())
		return 0
	}
	sig := make([][]byte, 0)
	sig = append(sig, rText)
	sig = append(sig, sTest)

	// proof
	doc.Proof = Model.Proof{
		Creator:   fmt.Sprintf("did:self:%s#key-%d", issuer, kid),
		Method:    "ECC",
		Signature: sig,
	}

	elapsed, err := controller.CreatePubDocument(doc, contract)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("创建公共DID总耗时：", elapsed)
	return elapsed
}

// CreateSelfDocument 创建私有did
func CreateSelfDocument(id string, contract *gateway.Contract) time.Duration {
	pripemFileName := fmt.Sprintf("./keystore/%s-private.pem", id)
	publicFileName := fmt.Sprintf("./keystore/%s-public.pem", id)
	ID := fmt.Sprintf("did:self:%s", id)
	Kid := fmt.Sprintf("did:self:%s#key-1", id)

	// 创建一组密钥
	err := Utils.GenerateECCKeys(pripemFileName, publicFileName)
	if err != nil {
		fmt.Println("创建密钥失败， err ->", err.Error())
		return 0
	}

	// 创建一个私有DID
	publiePem := Utils.ReadFile(publicFileName)
	did := Model.DID{
		ID: ID,
		Authenticaions: []Model.Authenticaion{
			{
				Kid:        Kid,
				Controller: ID,
				Method:     "ECC",
				PublicPem:  string(publiePem),
			},
		},
	}

	priPem := Utils.ReadFile(pripemFileName)
	rText, sTest, err := Utils.SigByECC(did, string(priPem))
	if err != nil {
		fmt.Println("签名失败 err ->", err.Error())
		return 0
	}
	sig := make([][]byte, 0)
	sig = append(sig, rText)
	sig = append(sig, sTest)

	did.Proof = Model.Proof{
		Creator:   Kid,
		Method:    "ECC",
		Signature: sig,
	}

	// 創建DID
	elapsedByCreateSelf, err := controller.CreateSelfDocument(did, contract)
	if err != nil {
		fmt.Println("err ->", err.Error())
		return 0
	}
	fmt.Println("创建成功")
	return elapsedByCreateSelf
}

// 获取DID, 并写入文件中
func GetAndSaveDocument(id string, contract *gateway.Contract) (time.Duration, error) {
	elapsed, doc, err := controller.GetDocument(contract, id)
	if err != nil {
		return 0, fmt.Errorf("获取DID文档失败,err -> %v", err.Error())
	}
	docAsBytes, err := json.Marshal(doc)
	if err != nil {
		return 0, err
	}

	// 生存JSON文件
	err = os.WriteFile(fmt.Sprintf("./document/%s.json", id), docAsBytes, os.ModeAppend)
	if err != nil {
		return 0, err
	}

	return elapsed, nil
}
