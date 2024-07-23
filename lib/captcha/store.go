package captcha

import (
	"fmt"
	"github.com/hdget/hdsdk/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type redisCaptchaStore struct {
	generator string // 验证码生成者
}

const (
	maxFailures                = 3
	captchaRedisKeyTemplate    = "captcha:%s"
	luaCaptchaIncreaseFailures = `
local key = KEYS[1]
local maxFailures = tonumber(ARGV[1])
local failures = tonumber(redis.call('HGET', key, "failures"))
if failures < maxFailures-1 then
	redis.call('HMSET', key, "failures", failures+1)
else
	redis.call('DEL', key)
end;
`
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

func (r redisCaptchaStore) Set(captchaId string, captchaValue string, expires int) error {
	err := hdsdk.Redis().My().HMSet(r.getStoreKey(captchaId), map[string]interface{}{
		"value":    r.getStoreValue(captchaValue),
		"failures": 0,
	})
	if err != nil {
		return errors.Wrap(err, "store set captcha")
	}

	err = hdsdk.Redis().My().Expire(r.getStoreKey(captchaId), expires)
	if err != nil {
		return errors.Wrap(err, "store expire captcha")
	}

	return nil
}

func (r redisCaptchaStore) Get(captchaId string, clear bool) (string, error) {
	val, err := hdsdk.Redis().My().HGet(r.getStoreKey(captchaId), "value")
	if err != nil {
		return "", errors.Wrap(err, "store get captcha")
	}

	if clear {
		err = hdsdk.Redis().My().Del(r.getStoreKey(captchaId))
		if err != nil {
			return "", errors.Wrap(err, "store clear captcha")
		}
	}
	return cast.ToString(val), nil
}

func (r redisCaptchaStore) Verify(captchaId, answer string, clear bool) (bool, error) {
	storeValue, err := r.Get(captchaId, clear)
	if err != nil {
		return false, err
	}

	verified := false
	if storeValue != "" && r.getStoreValue(answer) == storeValue {
		verified = true
	}

	if !verified {
		if !clear {
			err = r.increaseFailures(captchaId)
			if err != nil {
				return false, err
			}
		}
	}
	return verified, nil
}

func (r redisCaptchaStore) getStoreKey(id string) string {
	return fmt.Sprintf(captchaRedisKeyTemplate, id)
}

func (r redisCaptchaStore) getStoreValue(captchaValue string) string {
	return r.generator + captchaValue
}

func (r redisCaptchaStore) increaseFailures(captchaId string) error {
	_, err := hdsdk.Redis().My().Eval(luaCaptchaIncreaseFailures, []any{r.getStoreKey(captchaId)}, []any{maxFailures})
	if err != nil {
		return errors.Wrap(err, "store increase captcha")
	}
	return nil
}
