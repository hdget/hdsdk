package captcha

type CaptchaGenerator interface {
	Generate() (string, string, error)
}

type CaptchaStore interface {
	Set(captchaId string, value string, expires int) error
	Get(captchaId string, clear bool) (string, error)
	Verify(captchaId, captcha string, clear bool) (bool, error)
}

type baseGenerator struct {
	option *captchaOption
	store  CaptchaStore
}

func newGenerator() *baseGenerator {
	return &baseGenerator{
		option: newOption(),
		store:  Store(),
	}
}
