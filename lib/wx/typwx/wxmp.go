package typwx

type WxmpUserInfo struct {
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	Language  string `json:"language"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		AppId     string `json:"appid"`
	} `json:"watermark"`
}

type WxmpMobileInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		AppId     string      `json:"appid"`
		Timestamp interface{} `json:"timestamp"`
	} `json:"watermark"`
}

type WxmpMobileResult struct {
	Errcode   int            `json:"errcode"`
	Errmsg    string         `json:"errmsg"`
	PhoneInfo WxmpMobileInfo `json:"phone_info"`
}

// WxmpSession wechat miniprogram login session
type WxmpSession struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

type WxAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WxAccessTokenResult struct {
	WxAccessToken
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// 微信小程序码
type CommonWxaCodeParam struct {
	// 要打开的小程序版本。正式版为 release，体验版为 trial，开发版为 develop
	EnvVersion string `json:"env_version"`
	// 二维码的宽度，单位 px。最小 280px，最大 1280px
	Width int `json:"width"`
	// auto_color 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调
	AutoColor bool `json:"auto_color"`
	// auto_color 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"} 十进制表示
	LineColor struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
	} `json:"line_color"`
	// 是否需要透明底色，为 true 时，生成透明底色的小程序码
	IsHyaline bool `json:"is_hyaline"`
}

type LimitedWxaCodeParam struct {
	// 扫码进入的小程序页面路径，最大长度 128 字节，不能为空；
	// 对于小游戏，可以只传入 query 部分，来实现传参效果，如：传入 "?foo=bar"，
	// 即可在 wx.getLaunchOptionsSync 接口中的 query 参数获取到 {foo:"bar"}。
	Path string `json:"path"`
	*CommonWxaCodeParam
}

type UnLimitedWxaCodeParam struct {
	// 最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，其它字符请自行编码为合法字符（因不支持%，中文无法使用 urlencode 处理，请使用其他编码方式）
	Scene string `json:"scene"`
	// 页面 page，例如 pages/index/index，根路径前不要填加 /，不能携带参数（参数请放在scene字段里），如果不填写这个字段，默认跳主页面
	Page      string `json:"page"`
	CheckPath bool   `json:"check_path"`
	*CommonWxaCodeParam
}
