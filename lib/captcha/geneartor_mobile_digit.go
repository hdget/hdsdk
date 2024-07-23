package captcha

import (
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pkg/errors"
)

type mobileDigitCaptchaGenerator struct {
	*baseGenerator
}

func NewMobileDigitCaptcha(mobile string, options ...Option) CaptchaGenerator {
	m := &mobileDigitCaptchaGenerator{
		baseGenerator: newGenerator(mobile),
	}

	for _, option := range options {
		option(m.option)
	}

	return m
}

func (m mobileDigitCaptchaGenerator) Generate() (string, string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", "", errors.Wrap(err, "generate captcha id")
	}

	captchaValue, err := gonanoid.Generate("0123456789", m.option.length)
	if err != nil {
		return "", "", errors.Wrap(err, "generate captcha")
	}

	// 保存验证码的时候加入generator前缀
	err = Store(m.name).Set(uuid.String(), captchaValue, m.option.expires)
	if err != nil {
		return "", "", errors.Wrap(err, "store set captcha")
	}

	return uuid.String(), captchaValue, nil
}
