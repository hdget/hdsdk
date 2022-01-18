package typwx

// WxoaSignature signature接口返回结果
type WxoaSignature struct {
	AppID     string `json:"appId"`
	Ticket    string `json:"ticket"`
	Noncestr  string `json:"noncestr"`
	Url       string `json:"url"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}
