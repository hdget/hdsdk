package dapr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
	"strings"
)

type ServiceModule struct {
	app       string
	version   int                              // 版本号
	namespace string                           // 命名空间
	name      string                           // 服务模块名
	client    string                           // 客户端
	handlers  map[string]*serviceModuleHandler // 定义的方法
}

const ContentTypeJson = "application/json"

// as we use gogoprotobuf which doesn't has protojson.Message interface
//var jsonpb = protojson.MarshalOptions{
//	EmitUnpopulated: true,
//}
//var jsonpbMarshaler = jsonpb.Marshaler{EmitDefaults: true}

// Invoke 调用dapr服务
func Invoke(appId string, version int, namespace, method string, data interface{}, args ...string) ([]byte, error) {
	fullMethodName := getFullMethodName(version, namespace, "", method)
	return realInvoke(appId, fullMethodName, data, args...)
}

// InvokeWithClient 调用dapr服务
func InvokeWithClient(appId string, version int, namespace, client, method string, data interface{}, args ...string) ([]byte, error) {
	fullMethodName := getFullMethodName(version, namespace, client, method)
	return realInvoke(appId, fullMethodName, data, args...)
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

// RegisterHandlers register handlers to global registry
func RegisterHandlers(app string, module interface{}, methods map[string]common.ServiceInvocationHandler, registry map[string]*ServiceModule) error {
	if registry == nil {
		return errors.New("registry is nil")
	}

	// 给定结构体实例获取结构体名字
	moduleName := utils.GetStructName(module)
	if moduleName == "" {
		return fmt.Errorf("invalid module, module:%v", module)
	}

	svcModule, err := newServiceModule(app, moduleName)
	if err != nil {
		return errors.Wrap(err, "new service module")
	}

	// 注册handlers
	err = svcModule.registerHandlers(methods)
	if err != nil {
		return err
	}

	registry[moduleName] = svcModule
	return nil
}

// GetInvocationHandlers 从注册中心获取Dapr InvocationHandler
func GetInvocationHandlers(registry map[string]*ServiceModule) map[string]common.ServiceInvocationHandler {
	invocationHandlers := make(map[string]common.ServiceInvocationHandler)
	for _, module := range registry {
		for _, handler := range module.handlers {
			invocationHandlers[handler.method] = handler.invocationHandler
		}
	}
	return invocationHandlers
}

// getFullMethodName 构造version:namespace:client:realMethod的方法名
// 为了进行namespace和版本号的区分，组装method=version:namespace:clientName:realMethod
// 其中client可能为空，这个说明该接口可以给任何client使用
func getFullMethodName(version int, namespace, client, method string) string {
	tokens := []string{fmt.Sprintf("v%d", version), namespace, method}
	if client != "" {
		tokens = []string{fmt.Sprintf("v%d", version), namespace, client, method}
	}

	return strings.Join(tokens, ":")
}

// newServiceModule 从函数的receiver即moduleName中按v<version>_<namespace>的格式解析出API版本号和命名空间
func newServiceModule(app, moduleName string) (*ServiceModule, error) {
	var partVersion, partNamespace, partClient string
	tokens := strings.Split(moduleName, "_")
	countTokens := len(tokens)
	if countTokens <= 1 {
		return nil, errors.New("invalid module, it should be: v<number>_<namespace> or v<number>_<namespace>_<client>")
	}

	// 解析version, namespace和client
	if countTokens == 2 {
		partVersion = tokens[0]
		partNamespace = tokens[1]
	} else if countTokens > 2 {
		partVersion = tokens[0]
		partNamespace = strings.Join(tokens[1:countTokens-1], "_")
		partClient = tokens[countTokens-1]
	}

	// 校验version和namespace, client可以为空
	strVersion := partVersion[1:]
	if !strings.HasPrefix(partVersion, "v") || strVersion == "" {
		return nil, errors.New("invalid version, it should be: v<number>")
	}

	if partNamespace == "" {
		return nil, fmt.Errorf("invalid namespace, moduleName: %s", moduleName)
	}

	return &ServiceModule{
		app:       app,
		name:      moduleName,
		version:   cast.ToInt(strVersion),
		namespace: partNamespace,
		client:    partClient,
		handlers:  make(map[string]*serviceModuleHandler),
	}, nil
}
