package worktool

import (
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
)

const (
	socketType = 2 //通讯类型
	apiUrl     = "https://worktool.asrtts.cn/wework/sendRawMessage"
)

type MsgKind int

const (
	MsgKindText MsgKind = 203
	MsgKindFile MsgKind = 218
)

type TextMessage struct {
	Kind      MsgKind  `json:"type"`            // 消息类型
	Targets   []string `json:"titleList"`       // 发送目标，可以为联系人或者群
	Content   string   `json:"receivedContent"` // 文本内容 (\n换行)
	AtPeoples []string `json:"atList"`          // @人列表
}

type File struct {
	Name    string `json:"objectName"` // 对象名，这里代表文件名
	Url     string `json:"fileUrl"`    // 文件路径
	Kind    string `json:"fileType"`   // 文件类型 image等
	Comment string `json:"extraText"`  // 附件流言 可选填
}

type FileMessage struct {
	Kind    MsgKind  `json:"type"`
	Targets []string `json:"titleList"` // 发送人列表
	File
}

type MessageBody struct {
	SocketType int   `json:"socketType"`
	List       []any `json:"list"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type ApiWorkTool interface {
	SendFile(targets []string, file File) error
	SendText(targets []string, content string, atPeoples ...string) error
}

type workToolImpl struct {
	robotId string // 机器人id
}

func New(robotId string) ApiWorkTool {
	return &workToolImpl{robotId: robotId}
}

func (w *workToolImpl) SendText(targets []string, content string, atPeoples ...string) error {
	msg := &TextMessage{
		Kind:      MsgKindText,
		Targets:   targets,
		Content:   content,
		AtPeoples: atPeoples,
	}

	return w.send(msg)
}

func (w *workToolImpl) SendFile(targets []string, file File) error {
	msg := FileMessage{
		Kind:    MsgKindFile,
		Targets: targets,
		File:    file,
	}

	return w.send(msg)
}

func (w *workToolImpl) send(messages ...any) error {
	body := MessageBody{
		SocketType: socketType,
		List:       messages,
	}

	var apiResponse Response
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("robotId", w.robotId).
		SetBody(body).
		SetResult(&apiResponse).
		Post(apiUrl)
	if err != nil {
		return err
	}

	// http status code != 200
	if resp.StatusCode() != 200 {
		return errors.New(hdutils.BytesToString(resp.Body()))
	}

	// api response code != 200 means error
	if apiResponse.Code != 200 {
		return errors.New(apiResponse.Message)
	}

	return nil
}
