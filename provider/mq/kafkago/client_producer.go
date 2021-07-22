package kafkago

import (
	"fmt"
	"github.com/hdget/sdk/types"
	"github.com/segmentio/kafka-go"
)

type ProducerClient struct {
	*BaseClient
	Config *ProducerConfig
	Option *PublishOption
	Writer *kafka.Writer
}

func (k *Kafka) newProducerClient(name string, options map[types.MqOptionType]types.MqOptioner) (*ProducerClient, error) {
	// 获取匹配的路由配置
	config := k.getProducerConfig(name)
	if config == nil {
		return nil, fmt.Errorf("no matched producer config for: %s", name)
	}
	return &ProducerClient{
		BaseClient: k.newBaseClient(name, options),
		Config:     config,
		Option:     getPublishOption(options),
	}, nil
}

// connect balance策略
func (pc *ProducerClient) connect(brokers []string) error {
	var implBalancer kafka.Balancer
	switch pc.Config.Balance {
	case "roundrobin": // equally distributes messages across all available partitions.
		implBalancer = &kafka.RoundRobin{}
	case "leastbytes": // routes messages to the partition that has received the least amount of data.
		implBalancer = &kafka.LeastBytes{}
	case "hash": //
		if pc.Option.HashFunc != nil {
			implBalancer = &kafka.Hash{Hasher: pc.Option.HashFunc}
		}
	case "crc32": // uses the CRC32 hash function to determine which partition to route messages to
		implBalancer = &kafka.CRC32Balancer{}
	case "murmur2": // uses Murmur2 hash function to determine which partition to route messages to
		implBalancer = &kafka.Murmur2Balancer{}
	default:
		implBalancer = &kafka.RoundRobin{}
	}

	// 有可能
	if implBalancer == nil {
		return ErrInvalidBalancer
	}

	// make a writer that produces to topic using the least-bytes distribution
	pc.Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: implBalancer,
	}
	return nil
}
