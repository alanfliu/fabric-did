// [{"id":"did:self:002#key-1","method":"ECC","controller":"did:self:002","publicpem":"-----BEGIN Ecc PublicKey-----\nMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQA+2bpHwT6qZ+v5BUmy4/KQlZOf3qU\nMf1SlVMatqOb2efLrKX3pR3I0+LAchjsNnMeWnFLoDfoGKEK0/6l6FpxUBUAC7c+\nna5wCnKPvaOaB0KqXufQv1O96NmSttMWS931cE0G1cr83GsZGwJ6KBFA1zuRwFZQ\n/B9Dzatfd2OcZQyeFtw=\n-----END Ecc PublicKey-----\n"}],"proof":{"creator":"did:self:002#key-1","method":"ECC","signature":["MjU0OTE4OTAyNjgyMjgwODg3ODM2ODY5NzEzMjY5NzQyODk2NDA2ODUwNDI4NzYxNjE0NjU3NjM4MjI0MTY0ODQ4MjE5NDg3MzQ5NjgxMDc2NTI4NjI0ODI1Njg1NDU1MDExNTUzOTA0ODI1MTM2ODEyNjE0ODQ5NTU2NTY0ODk1ODQzNDU5MDQ3NDE2MDE5NDEyNzI2NzQ5MzMyMA==","NDM3OTc1MjQwMTQzMTk5MzYyMTM5OTg1NjczMTkyMzMzMzk0MzA3NzIwMzEyNjgxNTAwMjYwNjIzNzE1MjAyNTU3NjExNjE0OTE3ODA4NDExMjY0ODc0NzI4NTEyMDQxMDg3Mjk5MTE0ODY2NjA2MDE3MzA0NjU0Mzc1MDE5NTg5MzgxNTQ0MzgwMDA3NDU3MzkwODYwOTAwMTUw"]}}

package main

import (
	"encoding/json"
	"fmt"
)

type Authenticaion struct {
	Kid        string `json:"id"`         // 加密密钥ID
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

func main() {
	var str = `{"id":"did:self:002","authenticaions":[{"id":"did:self:002#key-1","method":"ECC","controller":"did:self:002","publicpem":"-----BEGIN Ecc PublicKey-----\nMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQA+2bpHwT6qZ+v5BUmy4/KQlZOf3qU\nMf1SlVMatqOb2efLrKX3pR3I0+LAchjsNnMeWnFLoDfoGKEK0/6l6FpxUBUAC7c+\nna5wCnKPvaOaB0KqXufQv1O96NmSttMWS931cE0G1cr83GsZGwJ6KBFA1zuRwFZQ\n/B9Dzatfd2OcZQyeFtw=\n-----END Ecc PublicKey-----\n"}],"proof":{"creator":"did:self:002#key-1","method":"ECC","signature":["MjU0OTE4OTAyNjgyMjgwODg3ODM2ODY5NzEzMjY5NzQyODk2NDA2ODUwNDI4NzYxNjE0NjU3NjM4MjI0MTY0ODQ4MjE5NDg3MzQ5NjgxMDc2NTI4NjI0ODI1Njg1NDU1MDExNTUzOTA0ODI1MTM2ODEyNjE0ODQ5NTU2NTY0ODk1ODQzNDU5MDQ3NDE2MDE5NDEyNzI2NzQ5MzMyMA==","NDM3OTc1MjQwMTQzMTk5MzYyMTM5OTg1NjczMTkyMzMzMzk0MzA3NzIwMzEyNjgxNTAwMjYwNjIzNzE1MjAyNTU3NjExNjE0OTE3ODA4NDExMjY0ODc0NzI4NTEyMDQxMDg3Mjk5MTE0ODY2NjA2MDE3MzA0NjU0Mzc1MDE5NTg5MzgxNTQ0MzgwMDA3NDU3MzkwODYwOTAwMTUw"]}}`

	var doc DID
	if err := json.Unmarshal([]byte(str), &doc); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("doc is ->", doc)
}
