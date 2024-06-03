package wxoa

// Signature signature接口返回结果
type Signature struct {
	AppID     string `json:"appId"`
	Ticket    string `json:"ticket"`
	Noncestr  string `json:"noncestr"`
	Url       string `json:"url"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// Ticket 类型
type Ticket struct {
	Value     string `json:"ticket,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}
