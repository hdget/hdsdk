package dapr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"reflect"
	"strings"
)

const ContentTypeJson = "application/json"

// as we use gogoprotobuf which doesn't has protojson.Message interface
//var jsonpb = protojson.MarshalOptions{
//	EmitUnpopulated: true,
//}
//var jsonpbMarshaler = jsonpb.Marshaler{EmitDefaults: true}

// Invoke 调用dapr服务
func Invoke(version int, appId, namespace, method string, data interface{}, args ...string) ([]byte, error) {
	fullMethodName := getFullMethodName(version, namespace, "", method)
	return realInvoke(appId, fullMethodName, data, args...)
}

// InvokeWithClient 调用dapr服务
func InvokeWithClient(version int, appId, namespace, client, method string, data interface{}, args ...string) ([]byte, error) {
	fullMethodName := getFullMethodName(version, namespace, client, method)
	return realInvoke(appId, fullMethodName, data, args...)
}

// InvokeWithDaprClient 需要传入daprClient去调用
func InvokeWithDaprClient(daprClient client.Client, appId, methodName string, data interface{}, args ...string) ([]byte, error) {
	if daprClient == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	var value []byte
	switch t := data.(type) {
	case string:
		value = utils.StringToBytes(t)
	case []byte:
		value = t
	default:
		v, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "marshal invoke data")
		}
		value = v
	}

	// 添加额外的meta信息
	ctx := context.Background()
	if len(args) > 0 {
		md := metadata.Pairs(args...)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	content := &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	}

	ret, err := daprClient.InvokeMethodWithContent(ctx, appId, methodName, "post", content)
	if err != nil {
		return nil, errors.Wrapf(err, "dapr invoke method with content, app:%s, method: %s, content: %s", appId, methodName, utils.BytesToString(value))
	}

	return ret, nil
}

func realInvoke(appId, fullMethodName string, data interface{}, args ...string) ([]byte, error) {
	var value []byte
	switch t := data.(type) {
	case string:
		value = utils.StringToBytes(t)
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
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer daprClient.Close()
	// 添加额外的meta信息
	ctx := context.Background()
	if len(args) > 0 {
		md := metadata.Pairs(args...)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	content := &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	}

	resp, err := daprClient.InvokeMethodWithContent(ctx, appId, fullMethodName, "post", content)
	if err != nil {
		return nil, errors.Wrapf(err, "dapr invoke method with content, app:%s, method: %s, content: %s", appId, fullMethodName, utils.BytesToString(value))
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

// RegisterHandlers register namespace's method to global registry
func RegisterHandlers(holder interface{}, methods map[string]common.ServiceInvocationHandler, registry map[string]map[string]common.ServiceInvocationHandler) error {
	if registry == nil {
		return errors.New("registry is nil")
	}
	namespace := getNamespaceName(holder)
	if namespace != "" {
		registry[namespace] = methods
	}
	return nil
}

// ParseHandlers parse handlers from registry
func ParseHandlers(registry map[string]map[string]common.ServiceInvocationHandler) map[string]common.ServiceInvocationHandler {
	handlers := make(map[string]common.ServiceInvocationHandler)
	for namespace, methods := range registry {
		for methodName, fn := range methods {
			tokens := strings.Split(namespace, "_")
			tokens = append(tokens, methodName)
			handlers[strings.Join(tokens, ":")] = fn
		}
	}
	return handlers
}

func getNamespaceName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

// getFullMethodName 构造version:namespace:client:realMethod的方法名
// 为了进行namespace和版本号的区分，组装method=version:moduleName:clientName:realMethod
// 其中client可能为空，这个说明该接口可以给任何client使用
func getFullMethodName(version int, namespace, client, method string) string {
	tokens := []string{fmt.Sprintf("v%d", version), namespace, method}
	if client != "" {
		tokens = []string{fmt.Sprintf("v%d", version), namespace, client, method}
	}

	return strings.Join(tokens, ":")
}
