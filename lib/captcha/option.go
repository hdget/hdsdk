package captcha

type captchaOption struct {
	length  int
	expires int
	height  int
	width   int
}

type Option func(*captchaOption)

const (
	captchaLength    = 4
	captchaExpires   = 180 // seconds
	captchaImgHeight = 36
	captchaImgWidth  = 80
)

func newOption() *captchaOption {
	return &captchaOption{
		length:  captchaLength,
		expires: captchaExpires,
		height:  captchaImgHeight,
		width:   captchaImgWidth,
	}
}

func WithLength(length int) Option {
	return func(opt *captchaOption) {
		opt.length = length
	}
}

func WithExpires(expires int) Option {
	return func(opt *captchaOption) {
		opt.expires = expires
	}
}

func WithHeight(height int) Option {
	return func(opt *captchaOption) {
		opt.height = height
	}
}

func WithWidth(width int) Option {
	return func(opt *captchaOption) {
		opt.width = width
	}
}
