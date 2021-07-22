### Kafka配置

- default: 缺省的Kafka客户端
- items:   额外的Kafka客户端
 
```
[sdk.log]
        level = "debug"
         filename = "demo.log"
         [sdk.log.rotate]
             max_age = 168
             rotation_time=24
            
[sdk.kafka]
    [sdk.kafka.default]
        brokers = ["192.168.0.123:9092"]            <--- kafka连接的brokers
        [[sdk.kafka.default.consumers]]             <--- kafka消费端配置
            name = "consumer1"                      <--- 消费端必须用name来区分
            topic ="testtopic1"
            group_id=""                             <--- 如果要在同一个消费组，可以设置group_id
            user=""
            password=""
        [[sdk.kafka.default.producers]]             <--- kafka生产端配置
            name = "producer1"                      <--- 生产端必须用name来区分
            topics  = ["testtopic1", "testtopic2"]  <--- 支持同时往多个topic发送
            balance = ""                            <--- 支持balance的策略:roundrobin,leastbytes,hash,crc32,murmur2     

    [[sdk.kafka.items]]
            name = "other_kafka"
            brokers = ["192.168.0.123:9092"]        <--- kafka连接的brokers
            [[sdk.kafka.items.consumers]]
                name = "consumer2"
                topic ="testtopic2"
                group_id=""
```
> 1. 除了default的kafka，还需要使用其他Kafka服务，可以配置在`[[sdk.kafka.items]]`中，必须使用`name`来区分不同kafka服务
> 3. consumer/producer配置中必须指定`name`来区分不同配置, 在创建consumer或者producer的时候使用`name`来获取交换配置信息
> 4. 注意在退出程序的时候需要调用consumer/producer的Close()方法来保证消息不丢失

### Kafka使用指南
  
#### 获取Kafka客户端

- 获取缺省Kafka客户端: `sdk.Kafka.My()`
- 获取指定名字的Kafka客户端: `sdk.Kafka.By(name)`
    
#### Kafka创建producer并发送消息

- 首先指定client的name通过Producer()获取producer实例
- 调用producer实例的Publish()方法发送消息

```
p, err := sdk.Kafka.My().CreateProducer("producer1")
if err != nil {
	sdk.Logger.Error("kafka create producer", "err", err)
    return err
}

err = p.Publish([]byte("test"))
if err != nil {
	sdk.Logger.Fatal("publish", "err", err)
}
```

#### Kafka创建consumer并接收处理消息

1. 首先指定consumer的name通过CreateConsumer()创建consumer实例
2. 定义消息内容的处理函数, 消息函数为`func(data []byte) types.MqMsgAction`格式
3. 调用consumer实例的Consume()方法接收并处理消息

```
func msgProcess(data []byte) error {
	fmt.Println(utils.BytesToString(data))
	return nil
}

func handle() {
    c, err := sdk.Kafka.My().CreateConsumer("consumer1", msgProcess)
	if err != nil {
		sdk.Logger.Fatal("kafka create consumer", "err", err)
	}

    c.Consume()
}
```

#### 自定义发送或者接收的选项
1. 首先调用GetDefaultOptions()函数获取到系统所有默认选项, 默认选项有QueueOption, ExchangeOption, PublishOption, ConsumeOption四种
2. 修改指定选项中的默认值
2. 在创建producer或者consumer的时候将修改后的Options作为参数传入

```
...
mq := Kafka.My()
options := mq.GetDefaultOptions()
consumeOption := options[types.MqOptionQueue].(*Kafka.ConsumeOption)
consumeOption.MinBytes = 2
options[types.MqOptionConsume] = consumeOption
p, err := mq.CreateProducer("producer1", msgProcess, options)
... 
```
