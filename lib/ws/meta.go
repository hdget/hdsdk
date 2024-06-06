package ws

import "github.com/gin-gonic/gin"

const keyMeta = "meta"

// GetMetaKvs 获取meta的[key, value] slice
func GetMetaKvs(c *gin.Context) []string {
	return c.GetStringSlice(keyMeta)
}

// AddMetaKvs 添加信息到meta中去
func AddMetaKvs(c *gin.Context, key, value string) {
	c.Set(keyMeta, append(GetMetaKvs(c), key, value))
}
