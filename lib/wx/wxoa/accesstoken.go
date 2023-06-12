package wxoa

import "github.com/hdget/hdsdk/lib/wx/common"

func (w *implWxoa) GetAccessToken(appId, appSecret string) (string, error) {
	// 尝试从缓存中获取access token
	cachedAccessToken, _ := _cache.GetAccessToken(appId)
	if cachedAccessToken != "" {
		return cachedAccessToken, nil
	}

	// 如果从缓存中获取不到，尝试请求access token
	wxAccessToken, err := common.RequestAccessToken(appId, appSecret)
	if err != nil {
		return "", err
	}

	err = _cache.SetAccessToken(appId, wxAccessToken.AccessToken, wxAccessToken.ExpiresIn-1000)
	if err != nil {
		return "", err
	}

	return wxAccessToken.AccessToken, nil
}
