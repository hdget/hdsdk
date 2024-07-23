package captcha

import (
	"fmt"
	"github.com/hdget/hdsdk/v2"
	"github.com/pkg/errors"
)

type redisCaptchaStore struct {
	generator string // 验证码生成者
}

const (
	captchaRedisKeyTemplate = "captcha:%s"
)

func Store(args ...string) CaptchaStore {
	var generator string
	if len(args) > 0 {
		generator = args[0]
	}

	return &redisCaptchaStore{
		generator: generator,
	}
}

func (r redisCaptchaStore) Set(captchaId string, value string, expires int) error {
	err := hdsdk.Redis().My().SetEx(r.getKey(captchaId), r.generator+value, expires)
	if err != nil {
		return errors.Wrap(err, "store set captcha")
	}
	return nil
}

func (r redisCaptchaStore) Get(captchaId string, clear bool) (string, error) {
	val, err := hdsdk.Redis().My().GetString(r.getKey(captchaId))
	if err != nil {
		return "", errors.Wrap(err, "store get captcha")
	}

	if clear {
		err = hdsdk.Redis().My().Del(r.getKey(captchaId))
		if err != nil {
			return "", errors.Wrap(err, "store clear captcha")
		}
	}
	return val, nil
}

func (r redisCaptchaStore) Verify(captchaId, answer string, clear bool) (bool, error) {
	if answer == "" {
		return false, errors.New("empty answer")
	}

	captcha, err := r.Get(captchaId, clear)
	if err != nil {
		return false, err
	}

	if captcha == "" {
		return false, errors.New("empty captcha")
	}

	return r.generator+answer == captcha, nil
}

func (r redisCaptchaStore) getKey(id string) string {
	return fmt.Sprintf(captchaRedisKeyTemplate, id)
}
