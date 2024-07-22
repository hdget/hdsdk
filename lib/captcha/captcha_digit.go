package captcha

import (
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pkg/errors"
)

type digitCaptchaManager struct {
	*baseGenerator
}

func NewDigitCaptcha(options ...Option) CaptchaGenerator {
	m := &digitCaptchaManager{
		baseGenerator: newGenerator(),
	}

	for _, option := range options {
		option(m.option)
	}

	return m
}

func (m digitCaptchaManager) Generate() (string, string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", "", errors.Wrap(err, "generate captcha id")
	}

	captchaValue, err := gonanoid.Generate("0123456789", m.option.length)
	if err != nil {
		return "", "", errors.Wrap(err, "generate captcha")
	}

	err = m.store.Set(uuid.String(), captchaValue, m.option.expires)
	if err != nil {
		return "", "", errors.Wrap(err, "store set captcha")
	}

	return uuid.String(), "", nil
}
