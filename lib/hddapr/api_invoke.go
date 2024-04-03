package hddapr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"strings"
)

type apiServiceInvocation struct {
}

const ContentTypeJson = "application/json"

func NewApiServiceInvocation() *apiServiceInvocation {
	return &apiServiceInvocation{}
}

// as we use gogoprotobuf which doesn't has protojson.Message interface
//var jsonpb = protojson.MarshalOptions{
//	EmitUnpopulated: true,
//}
//var jsonpbMarshaler = jsonpb.Marshaler{EmitDefaults: true}

// Invoke 调用dapr服务
func (a apiServiceInvocation) Invoke(appId string, moduleVersion int, module, method string, data any, args ...string) ([]byte, error) {
	var value []byte
	switch t := data.(type) {
	case string:
		value = hdutils.StringToBytes(t)
	case []byte:
		value = t
	default:
		v, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "marshal invoke data")
		}
		value = v
	}

	daprClient, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return nil, errors.New("dapr client is null, handlerName resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer daprClient.Close()
	// 添加额外的meta信息
	ctx := context.Background()
	if len(args) > 0 {
		md := metadata.Pairs(args...)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	fullMethodName := a.GetServiceInvocationName(moduleVersion, module, method)
	resp, err := daprClient.InvokeMethodWithContent(ctx, appId, fullMethodName, "post", &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetMetaValues get grpc meta values
func GetMetaValues(ctx context.Context, key string) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md.Get(key)
}

// GetMetaValue get the first grpc meta value
func GetMetaValue(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// GetServiceInvocationName 构造version:module:realMethod的方法名
func (apiServiceInvocation) GetServiceInvocationName(version int, module, method string) string {
	return strings.Join([]string{fmt.Sprintf("v%d", version), module, method}, ":")
}
