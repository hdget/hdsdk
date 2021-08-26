Enterprise ready, robust and easy extensible sdk which help to quickly develop backedn services.

### SDK概览
- log: 日志能力
    * zerolog日志能力
- db: 数据库能力
    * mysql数据库能力
- cache: 缓存能力
    * redis缓存能力
- mq: 消息队列能力
    * RabbitMq消息队列能力
    * Kafka消息队列能力
- dts: 数据同步能力
    * Aliyun DTS数据同步能力

### SDK配置
在配置文件中，sdk配置有如下格式，当前支持:
- sdk.log:    日志能力配置
- sdk.mysql:  mysql数据库能力配置
- sdk.redis:  redis缓存能力配置
- sdk.rabbitmq:  rabbitmq消息队列配置
- sdk.kafka: kafka消息队列配置

> 日志能力`sdk.log`的配置是最基本需求的配置，不管使用sdk的时候是否具有其他能力, 日志配置信息必须要包含

```
[sdk]
    [sdk.<capability1>]
       ...
    [sdk.<capability2>]
       ...
```
- 日志能力配置，当前使用zerolog进行日志输出
    ```
    [sdk.log]
        # 当前支持日志级别: "trace", "debug", "info", "warn", "error", "fatal", "panic"
        level = "debug"
        # 日志文件名称
        filename = "demo.log"
        # 日志结转配置
        [sdk.log.rotate]
            # 日志最大保存时间7天(单位hour)
            max_age = 168
            # 日志切割时间间隔24小时（单位hour)
            rotation_time=24
    ```

- 数据库能力
    * MySQL: 请参考[MySQL能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/db/mysql)
 
- 缓存能力配置，当前只支持redis的相关配置信息
    * Redis: 请参考[Redis能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/cache/redis)
    
- 消息队列能力配置，当前只支持redis的相关配置信息
    * RabbitMq: 请参考[RabbitMQ能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/mq/rabbitmq)
    * Kafka: 请参考[RabbitMQ能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/mq/kafka)

### SDK用户指南

#### 一、 支持的环境定义
- local: 本地环境
- dev:   开发环境
- test:  测试环境
- pre:   预发布环境
- sim:   仿真环境
- prod:  生产环境

#### 二、 SDK配置加载

sdk的配置需要通过LoadConfig方法进行加载，加载结果为[viper](https://github.com/spf13/viper)实例
```
LoadConfig(app, cliEnv, cliFile string, args ...ConfigOption) *viper.Viper
```

在代码中，我们通过LoadConfig方法加载应用程序的所有配置信息，然后unmarshal成我们自定义的配置结构实例。
LoadConfig的时候可以通过自定义`ConfigOption`来进行默认配置的修改
- app:     加载配置的时候指定的应用的名字(必须指定)
- cliEnv:  加载什么环境的配置(必须指定)
- cliFile: 指定的配置文件路径(必须指定)
- args:    加载配置时候的选项

默认加载选项：
- env:     如果不指定cliEnv或者指定无效的cliEnv, 默认的env会设置为"prod", 我们也可以通过环境变量`HDGET_RUNTIME`来设置env
- rootdir: 配置信息的根目录，默认为"setting", 在LoadConfig的时候通过自定义`ConfigOption`可以进行自定义

sdk的底层能力配置支持从不同的源进行定义，当前支持:
- etcd:    配置信息读取路径: `<rootdir>/<app>/<env>`
- 文件系统: 配置信息读取路径: `<rootdir>/app/<app>.<env>.toml`

#### 三、 嵌入SDK配置项到应用的配置结构中

通常我们通过自定义一个继承自sdk.Config的结构体来包括sdk配置信息。
```
type MyConfig struct {
    sdk.Config `mapstructure:",squash"`
    // 其他的配置项
    ...
}
```
这里注意，`sdk.Config`的tag必须要加上`mapstructure:",squash"`，
因为viper是用过mapstructure这个包来读取配置信息的，加上这个tag告诉viper在读取配置信息的时候将该结构体的字段提到父结构中

#### 四、初始化SDK
在第一步我们加载配置信息以后，我们需要unmarshal配置信息到我们的自定义数据结构，然后来初始化sdk
```
v := sdk.LoadConfig("demo", "local", "")

var conf MyConfig
err := v.Unmarshal(&conf)
if err != nil {
    log.Fatalf("unmarshal config, error=%v", err)
}

err = sdk.Initialize(&conf)
if err != nil {
    log.Fatalf("msg=\"sdk initialize\" error=\"%v\"", err)
}
```

#### 五、使用SDK

- 日志能力

    sdk.Log输出的时候第一个message参数必须要填，后续按照`key/value`的格式指定额外需要输出的信息，如果有错误信息，`key`必须指定为`err`
    e,g:
    ```
    sdk.Log.Debug("message content", "err", errors.New("testerr"), "key1", 1, "key2", "value2")
    ```
    支持的日志输出级别有:
    - sdk.Log.Trace
    - sdk.Log.Info
    - sdk.Log.Debug
    - sdk.Log.Warn
    - sdk.Log.Error
    - sdk.Log.Fatal
    - sdk.Log.Panic

- 数据库能力
    * MySQL: 请参考[MySQL能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/db/mysql)
  
- 缓存能力
    * Redis: 请参考[Redis能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/cache/redis)

- 消息队列能力
    * RabbitMq: 请参考[RabbitMQ能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/mq/rabbitmq)
    * Kafka: 请参考[Kafka能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/mq/kafka)

