package aliyun

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"hash"
	"io"
	"net/url"
	"path"
	"time"
)

type OssPolicyConfig struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type Signature struct {
	AccessKeyId string
	Host        string
	Expire      int64
	Signature   string
	Directory   string
	Policy      string
}

type AliOss struct {
	AccessKey    string
	AccessSecret string
	Domain       string
	Endpoint     string
}

const DefaultExpireTime = 600

func NewAliOss(domain, endpoint, accessKey, accessSecret string) *AliOss {
	return &AliOss{
		Domain:       domain,
		Endpoint:     endpoint,
		AccessKey:    accessKey,
		AccessSecret: accessSecret,
	}
}

// Upload file
func (a *AliOss) Upload(bucket, dir, fileName string, data []byte) (string, error) {
	// 获取存储空间
	client, err := oss.New(a.Endpoint, a.AccessKey, a.AccessSecret)
	if err != nil {
		return "", err
	}

	buk, err := client.Bucket(bucket)
	if err != nil {
		return "", err
	}

	// 上传Byte数组
	absPath := path.Join(dir, fileName)
	err = buk.PutObject(absPath, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	// 将domain和文件路径结合起来
	u, err := url.Parse(a.Domain)
	if err != nil {
		return "", errors.Wrap(err, "parse oss domain")
	}
	u.Path = path.Join(u.Path, absPath)
	return u.String(), nil
}

// GenSignature 生成oss直传token
func (a *AliOss) GenSignature(dir, filename string) (*Signature, error) {
	expiresIn := time.Now().Unix() + DefaultExpireTime
	policyData, err := getPolicyData(dir, expiresIn)
	if err != nil {
		return nil, err
	}

	// create post policy json
	stdPolicyData := base64.StdEncoding.EncodeToString(policyData)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(a.AccessSecret))
	_, err = io.WriteString(h, stdPolicyData)
	if err != nil {
		return nil, err
	}

	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return &Signature{
		AccessKeyId: a.AccessKey,
		Host:        a.Domain,
		Expire:      expiresIn,
		Signature:   signedStr,
		Directory:   dir,
		Policy:      stdPolicyData,
	}, nil
}

func getPolicyData(dir string, expiresIn int64) ([]byte, error) {
	strExpireTime := time.Unix(expiresIn, 0).UTC().Format("2006-01-02T15:04:05Z")

	// 指定此次上传的文件名必须以user-dir开头
	condition := []string{"starts-with", "$key", dir}
	config := OssPolicyConfig{
		Expiration: strExpireTime,
		Conditions: [][]string{
			condition,
		},
	}

	// calculate signature
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return data, nil
}
