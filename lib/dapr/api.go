package dapr

import (
	"context"
	"github.com/dapr/go-sdk/client"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

type APIer interface {
	Invoke(appId string, moduleVersion int, module, method string, data any) ([]byte, error)
	Lock(lockStore, lockOwner, resource string, expiryInSeconds int) error
	Unlock(lockStore, lockOwner, resource string) error
	Publish(pubSubName, topic string, data interface{}, args ...bool) error
	SaveState(storeName, key string, value interface{}) error
	GetState(storeName, key string) ([]byte, error)
	DeleteState(storeName, key string) error
	GetConfigurationItems(configStore string, keys []string) (map[string]*client.ConfigurationItem, error)
	SubscribeConfigurationItems(ctx context.Context, configStore string, keys []string, handler client.ConfigurationHandleFunction) (string, error)
	GetBulkState(storeName string, keys any) (map[string][]byte, error)
}

type apiImpl struct {
	ctx context.Context
}

func Api(kvs ...string) APIer {
	ctx := context.Background()
	if len(kvs) > 0 {
		md := metadata.Pairs(kvs...)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return &apiImpl{
		ctx: ctx,
	}
}

func TenantApi(tid int64) APIer {
	return Api(MetaKeyTid, cast.ToString(tid))
}
