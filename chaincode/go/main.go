package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type Authenticaion struct {
	Kid        string `json:"kid"`        // 加密密钥ID
	Method     string `json:"method"`     // 加密算法
	Controller string `json:"controller"` // 持有者DID标识符
	PublicPem  string `json:"publicpem"`  // 公钥证书
}

type Proof struct {
	Creator   string   `json:"creator"`   // issuer DID标识符
	Method    string   `json:"method"`    // 加密算法
	Signature [][]byte `json:"signature"` // 签名值
}

type DID struct {
	ID             string          `json:"id"`             // did标识符
	Authenticaions []Authenticaion `json:"authenticaions"` // 公钥集合
	Proof          Proof           `json:"proof"`          // 证明
}

type SmartContract struct {
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "createDocument" {
		err := s.createDocument(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte("create document success"))
	} else if fn == "queryDocument" {
		did, err := s.QueryDocument(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(did)
	} else {
		return shim.Error(fmt.Sprintf("不支持%s方法", fn))
	}
}

// id string, method string, creator string, rText string, sText string, publicpems ...string
func (s *SmartContract) createDocument(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) <= 5 {
		return fmt.Errorf("Incorrect arguments. Expecting id, method, creator, rText, sText, pubpems")
	}

	id := args[0]
	method := args[1]
	creator := args[2]
	rText := args[3]
	sText := args[4]
	publicPems := args[5:]

	// 判断链上是否存在 did
	exit, err := stub.GetState(id)
	if err != nil {
		return fmt.Errorf("failed to get state from world")
	}
	if exit != nil {
		return fmt.Errorf("%s already exit", id)
	}

	// 根据参数组装DID文档
	var auths []Authenticaion
	// 拼接 authentions
	for index, pem := range publicPems {
		kid := fmt.Sprintf("%s#key-%d", id, index)
		auth := Authenticaion{
			Kid:        kid,
			Method:     method,
			Controller: id,
			PublicPem:  pem,
		}
		auths = append(auths, auth)
	}
	// 拼接 singature
	sign := make([][]byte, 0)
	rTextAsBytes := []byte(rText)
	sTextAsBytes := []byte(sText)
	sign = append(sign, rTextAsBytes)
	sign = append(sign, sTextAsBytes)

	proof := Proof{
		Creator:   creator,
		Method:    method,
		Signature: sign,
	}

	doc := DID{
		ID:             id,
		Authenticaions: auths,
		Proof:          proof,
	}

	docAsBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal doc.err=%v", err.Error())
	}

	err = stub.PutState(id, docAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put state to world.err=%v", err.Error())
	}
	return nil
}

func (s *SmartContract) QueryDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a id")
	}

	id := args[0]
	docAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from world. err=%v", err.Error())
	}

	if docAsBytes == nil {
		return nil, fmt.Errorf("%s does not exit", id)
	}

	return docAsBytes, nil
}

func main() {
	if err := shim.Start(new(SmartContract)); err != nil {
		fmt.Printf("Error starting SimpleContract chaincode: %s", err)
	}
}
