package intf

import (
	kv "github.com/sagikazarmark/crypt/config"
)

type KvProvider interface {
	Get(key string) ([]byte, error)
	List(key string) (kv.KVPairs, error)
	Set(key string, value []byte) error
	Watch(key string, stop chan bool) <-chan *kv.Response
}
