package dapr

import (
	"github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
)

// Lock 锁
func (a apiImpl) Lock(lockStore, lockOwner, resource string, expiryInSeconds int) error {
	daprClient, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	resp, err := daprClient.TryLockAlpha1(a.ctx, a.normalize(lockStore), &client.LockRequest{
		LockOwner:       lockOwner,
		ResourceID:      resource,
		ExpiryInSeconds: int32(expiryInSeconds),
	})
	if err != nil {
		return errors.Wrap(err, "try lock")
	}

	if !resp.Success {
		return errors.New("lock failed")
	}

	return nil
}

// Unlock 取消锁
func (a apiImpl) Unlock(lockStore, lockOwner, resource string) error {
	daprClient, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	resp, err := daprClient.UnlockAlpha1(a.ctx, a.normalize(lockStore), &client.UnlockRequest{
		LockOwner:  lockOwner,
		ResourceID: resource,
	})
	if err != nil {
		return errors.Wrap(err, "try lock")
	}

	if resp.StatusCode != 0 {
		return errors.New(resp.Status)
	}

	return nil
}
