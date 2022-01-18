package typwx

// WxErrResponse 生成微信的AccessToken
type WxErrResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
