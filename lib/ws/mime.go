package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var (
	HttpContentTypes = map[MimeType]string{
		MimeTypeBinary: "application/octet-stream",
		MimeTypePDF:    "application/pdf",
		MimeTypeMSWord: "application/msword",
		MimeTypeJPEG:   "image/jpeg",
	}
)

type MimeType int

const (
	_ MimeType = iota
	MimeTypeBinary
	MimeTypeJPEG
	MimeTypePDF
	MimeTypeMSWord
)

func getContentType(fileType MimeType) string {
	contentType, exist := HttpContentTypes[fileType]
	if !exist {
		contentType = HttpContentTypes[MimeTypeBinary]
	}
	return contentType
}

// Download a mime file
func Download(c *gin.Context, mimeType MimeType, filename string, content []byte) error {
	attachment := fmt.Sprintf("attachment; filename=%s", filename)
	contentType := getContentType(mimeType)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", attachment)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))

	//回写到web 流媒体 形成下载
	_, err := c.Writer.Write(content)
	if err != nil {
		return err
	}
	return nil
}
