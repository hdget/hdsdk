package captcha

import (
	"github.com/hdget/hdsdk/v2"
	"github.com/mojocn/base64Captcha"
)

type imageCaptchaGenerator struct {
	*baseGenerator
	driver base64Captcha.Driver
}

func NewImageCaptcha(options ...Option) CaptchaGenerator {
	m := &imageCaptchaGenerator{
		baseGenerator: newGenerator(""),
	}

	for _, option := range options {
		option(m.option)
	}

	m.driver = base64Captcha.NewDriverDigit(
		m.option.height,
		m.option.width,
		m.option.length,
		0.2,
		1)

	return m
}

func (m imageCaptchaGenerator) Generate() (string, string, error) {
	s := &imageCaptchaStore{
		expires: m.option.expires,
		store:   Store(m.name),
	}

	c := base64Captcha.NewCaptcha(m.driver, s)
	captchaId, imgContent, _, err := c.Generate()
	return captchaId, imgContent, err
}

type imageCaptchaStore struct {
	expires int
	store   CaptchaStore
}

func (r imageCaptchaStore) Set(captchaId string, value string) error {
	return r.store.Set(captchaId, value, r.expires)
}

func (r imageCaptchaStore) Get(captchaId string, clear bool) string {
	val, err := r.store.Get(captchaId, clear)
	if err != nil {
		hdsdk.Logger().Error("base64 get captcha", "captchaId", captchaId, "err", err)
		return ""
	}
	return val
}

func (r imageCaptchaStore) Verify(captchaId, answer string, clear bool) bool {
	if answer == "" {
		return false
	}

	captcha := r.Get(captchaId, clear)
	if captcha == "" {
		return false
	}

	return answer == captcha
}
