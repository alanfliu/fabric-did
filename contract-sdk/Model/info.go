package Model

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
