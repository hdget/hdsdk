### RabbitMQ配置

- default: 缺省的RabbitMQ客户端
- items: 其他的RabbitMQ客户端
 
 ```
[sdk.rabbitmq]
    [sdk.rabbitmq.default]
        host="127.0.0.1"                       <--- RabbitMQ连接地址
        username="testuser"                    <--- RabbitMQ连接的用户名
        password="testpassword"                <--- RabbitMQ连接的密码
        port=5672                              <--- RabbitMQ连接的端口         
        vhost="/"                              <--- RabbitMQ连接的vhost
        [[sdk.rabbitmq.default.consumers]]     <--- RabbitMQ消费者配置，包括exchange_name、exchange_type...
            name = "consumer1"                 <--- 消费者必须用name来区分
            exchange_name="testexchange"
            exchange_type="direct"             <--- 当前支支持direct,fanout,topic三种交换方式
            queue_name = "testqueue"           <--- 队列名，如果queue_name为空，则会自动生成一个随机名字的队列
            routing_keys = [""]                <--- 在topic模式下，可以指定多个routing_key
        [[sdk.rabbitmq.default.producers]]     <--- RabbitMQ生产者配置，包括exchange_name、exchange_type
            name = "producer1"                 <--- 消费者必须用name来区分
            exchange_name="testexchange"
            exchange_type="direct"             <--- 当前支支持direct,fanout,topic三种交换方式
    [[sdk.rabbitmq.items]]
            name = "extra_rabbitmq"
            host="127.0.0.1"
            username="testuser"
            password="testpassword"
            port=5672
            vhost="/"
            [[sdk.rabbitmq.items.consumers]]
                ...
```
> 1. 除了default的RabbitMQ，如果还需要使用其他RabbitMQ服务，可以在`[[sdk.rabbitmq.items]]`中配置，用`name`来区分
> 3. consumers配置中必须指定`name`来区分不同交换配置, 在创建consumer或者producer的时候使用`name`来获取相关配置信息
> 4. consumers配置实际使用中建议必须设置queue_name，否则server会自动生成一个随机队列，在下次重启consumer或者producer的时候，无法知道这个随机队列的名字，从而无法继续处理
> 5. 在producer发送消息的时候，如果在clients.routing_keys中指定了路由键，需要在发送时同样指定key,否则发送会失败。如果clients.routing_keys中包含`""`路由键，则发送的时候可以忽略key

### RabbitMQ使用指南
  
#### 获取RabbitMQ客户端

- 获取缺省RabbitMQ客户端: `sdk.Rabbitmq.My()`
- 获取指定名字的RabbitMQ客户端: `sdk.Rabbitmq.By(name)`
    
#### Rabbitmq创建producer并发送消息

- 首先指定client的name通过Producer()获取producer实例
- 调用producer实例的Publish()方法发送消息

```
p, err := sdk.Rabbitmq.My().Producer("producer1")
if err != nil {
	sdk.Logger.Error("create producer", "err", err)
    return err
}

err = p.Publish([]byte("test"))
if err != nil {
	sdk.Logger.Fatal("publish", "err", err)
}
```

> 这里sdk中保证了发送消息是可靠的，如果日志打印出`publishing is not acked`错误请及时检查代码

注意： 如果在topic模式，或者其他指定了非空的routing_key的情况下，Publish函数的第二个值必须为能够匹配上的routing_key
```
...
err = p.Publish([]byte("test"), "routing_key1")
if err != nil {
	sdk.Logger.Fatal("publish", "err", err)
}
```

#### Rabbitmq创建consumer并接收处理消息

1. 首先指定client的name通过Consumer()获取consumer实例
2. 定义消息内容的处理函数, 消息函数为`func(data []byte) error`格式
3. 调用consumer实例的Consume()方法接收并处理消息

```
func msgProcess(data []byte) types.MqMsgAction {
	fmt.Println(utils.BytesToString(data))
	return types.Ack
}

func xxx() {
    c, err := sdk.Rabbitmq.My().Consumer("cosumer1", msgProcess)
	if err != nil {
		sdk.Logger.Fatal("create consumer", "err", err)
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
mq := Rabbitmq.My()
options := mq.GetDefaultOptions()
queueOption := options[types.MqOptionQueue].(*rabbitmq.QueueOption)
queueOption.Durable = false
options[types.MqOptionQueue] = queueOption
p, err := mq.Producer("client1", msgProcess, options)
... 
```

> 这里sdk保证了各种connection或channel错误，包括网络故障，RabbitMQ服务重启都可以重连恢复