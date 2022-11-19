package controller

import (
	"encoding/json"
	"fabric-learn/fabric-did/contract-sdk/Model"
	"fabric-learn/fabric-did/contract-sdk/Utils"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// http://localhost:8080/get?id=
func GetDocument(contract *gateway.Contract, id string) (time.Duration, Model.DID, error) {
	start := time.Now()
	result, err := contract.EvaluateTransaction(Model.QueryDocument, id)
	if err != nil {
		return 0, Model.DID{}, fmt.Errorf("获取DID失败，err -> %v", err.Error())
	}
	var d Model.DID
	json.Unmarshal(result, &d)
	elapsed := time.Since(start)
	return elapsed, d, nil
}

func CreateSelfDocument(doc Model.DID, contract *gateway.Contract) (time.Duration, error) {
	start := time.Now()
	var publicPem string
	for _, auth := range doc.Authenticaions {
		if auth.Kid == doc.Proof.Creator {
			publicPem = auth.PublicPem
			break
		}
	}

	if publicPem == "" {
		return 0, fmt.Errorf("%s 未匹配任何公钥", doc.Proof.Creator)
	}

	if ok, err := Utils.VerifyElliptic(doc, doc.Proof.Signature[0], doc.Proof.Signature[1], publicPem); err != nil {
		return 0, err
	} else if !ok {
		return 0, fmt.Errorf("验证签名失败")
	} else {
		_, err = contract.SubmitTransaction(
			Model.CreateDocument,
			doc.ID,
			doc.Proof.Method,
			doc.Proof.Creator,
			string(doc.Proof.Signature[0]),
			string(doc.Proof.Signature[1]),
			doc.Authenticaions[0].PublicPem,
		)
		if err != nil {
			return 0, fmt.Errorf("发起交易失败-> %v", err.Error())
		}
		elapsed := time.Since(start)

		return elapsed, nil
	}
}

func CreatePubDocument(doc Model.DID, contract *gateway.Contract) (time.Duration, error) {
	start := time.Now()

	idArr := strings.Split(doc.Proof.Creator, "#")
	creator := idArr[0]

	_, creatorDoc, err := GetDocument(contract, creator)
	fmt.Println("获取creator doc 成功")
	if err != nil {
		return 0, err
	}

	var publicPem string
	for _, auth := range creatorDoc.Authenticaions {
		if auth.Kid == doc.Proof.Creator {
			publicPem = auth.PublicPem
			break
		}
	}

	if publicPem == "" {
		return 0, fmt.Errorf("未匹配~~")
	}

	if ok, err := Utils.VerifyElliptic(doc, doc.Proof.Signature[0], doc.Proof.Signature[1], publicPem); err != nil {
		return 0, err
	} else if !ok {
		return 0, fmt.Errorf("验证签名失败")
	} else {
		_, err = contract.SubmitTransaction(
			Model.CreateDocument,
			doc.ID,
			doc.Proof.Method,
			doc.Proof.Creator,
			string(doc.Proof.Signature[0]),
			string(doc.Proof.Signature[1]),
			doc.Authenticaions[0].PublicPem,
		)
		if err != nil {
			return 0, fmt.Errorf("发起交易失败-> %v", err.Error())
		}
		elapsed := time.Since(start)

		return elapsed, nil
	}
}
