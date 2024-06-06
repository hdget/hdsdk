package dapr

type APIer interface {
	Invoke(appId string, moduleVersion int, module, method string, data any, args ...string) ([]byte, error)
	Lock(lockStore, lockOwner, resource string, expiryInSeconds int) error
	Unlock(lockStore, lockOwner, resource string) error
	Publish(pubSubName, topic string, data interface{}, args ...bool) error
	SaveState(storeName, key string, value interface{}) error
	GetState(storeName, key string) ([]byte, error)
	DeleteState(storeName, key string) error
}

type apiImpl struct {
}

func Api() APIer {
	return &apiImpl{}
}
