### RabbitMQ配置

#### Capability
- rabbitmq.Capability

#### 配置
```
[sdk.rabbitmq]
    host = "127.0.0.1"          # 主机，必填
    username = ""               # RabbitMQ连接用户名, 不填默认为guest
    password = ""               # RabbitMQ连接密码，不填默认为guest
    port = 5672                 # RabbitMQ端口，不填默认为5672
    vhost = "/"                 # RabbitMQ VHOST，不填默认为"/"
    requeue_in_failure=true     # 消息消费失败是否重新加入队列， 不填默认为true
    channel_pool_size=10        # channel池大小，不填默认为10
    prefetch_count=2            # 预取数，不填默认为2
```

### 使用指南

#### 发送端

```go
    pub, err := hdsdk.Mq().NewPublisher()
    if err != nil {
        return err
    }
    
	// 发送即时消息, 不创建exchange
    err = pub.Publish(context.Background(), "test", [][]byte("test1", "test2"))
    if err != nil {
        return err
    }

    // 发送即时消息, 创建exchange:order, topic为:close
    err = pub.Publish(context.Background(), "order:close", [][]byte("test1", "test2"))
    if err != nil {
        return err
    }

    // 发送10秒的延时消息, 不创建exchange
    err = pub.Publish(context.Background(), "test@delay", [][]byte("test1", "test2"), 10)
    if err != nil {
        return err
    }

    // 发送10秒的延时消息, 创建exchange:order, topic为:cancel
    err = pub.Publish(context.Background(), "order:cancel@delay", [][]byte("test1", "test2"), 10)
    if err != nil {
        return err
    }
```

#### 接收端

```
    sub, err := hdsdk.Mq().NewSubscriber()
    if err != nil {
        return err
    }
    
	// 订阅即时消息
	var messages <-chan *mq.Message
	go func() {
        messages, err = sub.Subscribe(context.Background(), "test")
        if err != nil {
            return err
        }
	}
    
    for {
      case msg := <-messages:
            // 处理消息
    }
```

```    
    // 订阅延时消息
    var messages <-chan *mq.Message
	go func() {
        messages, err = sub.Subscribe(context.Background(), "test@delay")
        if err != nil {
            return err
        }
	}
    
    for {
      case msg := <-messages:
            // 处理消息
    }
```
  