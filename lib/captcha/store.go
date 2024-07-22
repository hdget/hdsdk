package captcha

import (
	"fmt"
	"github.com/hdget/hdsdk/v2"
	"github.com/pkg/errors"
)

type redisCaptchaStore struct {
}

const (
	captchaRedisKeyTemplate = "captcha:%s"
)

func Store() CaptchaStore {
	return &redisCaptchaStore{}
}

func (r redisCaptchaStore) Set(captchaId string, value string, expires int) error {
	err := hdsdk.Redis().My().SetEx(r.getKey(captchaId), value, expires)
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

	return answer == captcha, nil
}

func (r redisCaptchaStore) getKey(id string) string {
	return fmt.Sprintf(captchaRedisKeyTemplate, id)
}
