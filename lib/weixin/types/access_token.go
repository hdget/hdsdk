package types

type WxAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WxAccessTokenResult struct {
	WxAccessToken
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
