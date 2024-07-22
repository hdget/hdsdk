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
	name   string
	option *captchaOption
}

func newGenerator(generator string) *baseGenerator {
	return &baseGenerator{
		option: newOption(),
		name:   generator,
	}
}
