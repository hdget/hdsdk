Enterprise ready, robust and easy extensible sdk which help to quickly develop backend services.

## SDK用户指南

### SDK约定

SDK中支持的环境`env`定义
- local: 本地环境
- dev: 开发环境
- test: 测试环境
- pre: 预发布环境
- sim: 仿真环境
- prod: 生产环境


### SDK能力介绍
#### 1. 配置加载
- sdk段落中的全部为sdk自身能力的配置项，例如Redis能力，MySQL能力，其他的段落可以用来为应用程序进行个性化的定制， 例如：

```toml
[sdk]
    [sdk.log]
        # 当前支持日志级别: "trace", "debug", "info", "warn", "error", "fatal", "panic"
        level = "debug"
        # 日志文件名称
        filename = "app.log"
        # 日志结转配置
        [sdk.log.rotate]
            # 日志最大保存时间7天(单位hour)
            max_age = 720
            # 日志切割时间间隔24小时（单位hour)
            rotation_time=24

    [sdk.mysql]
        [sdk.mysql.default]
            database = "database"
            user = "root"
            host = "127.0.0.1"
            password = ""
            port = 3306
   
    [sdk.redis.default]
            host = "127.0.0.1"
            port = 6379
            password = ""
            db = 0

[app]
    debug=false
    # 微信小程序
    [app.wxmp]
        app_id="xxxxxxxxxxxxxx"

```

#### 2. 日志输出
- sdk日志输出为结构化日志输出, 第一个参数必须要填，后续按照`key/value`的格式指定额外需要输出的信息，如果有错误信息，`key`必须指定为`err`或`error`，以下几个特殊的key值为系统占用，请勿使用:
  * msg
  * message
  * time
  * err
  * error

- 示例
    ```
    sdk.Logger().Debug("message content", "err", errors.New("testerr"), "key1", 1, "key2", "value2")
    ```
  
- 日志输出级别
  - sdk.Logger().Trace
  - sdk.Logger().Info
  - sdk.Logger().Debug
  - sdk.Logger().Warn
  - sdk.Logger().Error
  - sdk.Logger().Fatal
  - sdk.Logger().Panic

- 数据库
  * MySQL: 请参考[MySQL能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/db/mysql)

- 缓存
  * Redis: 请参考[Redis能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/cache/redis)

- 消息队列
  * RabbitMq: 请参考[RabbitMQ能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/mq/rabbitmq)
  * Kafka: 请参考[Kafka能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/mq/kafka)

### SDK快速上手

下面是一个简单的使用hdsdk的数据库能力的示例，代码非常简单，但支持多环境多配置源初始化加载各种配置项，同时数据库支持主备或指定数据库源

```go
import "github.com/hdget/hdsdk/v2/provider/db/sqlx-mysql"

err := hdsdk.New("app", "test").Initialize(sqlx_mysql.Capability)
if err != nil {
    log.Fatal(err)
}
```

### SDK使用

#### 第一步： 在配置文件中定义配置
```toml
[sdk]
    [sdk.log]
        # 当前支持日志级别: "trace", "debug", "info", "warn", "error", "fatal", "panic"
        level = "debug"
        # 日志文件名称
        filename = "app.log"
        # 日志结转配置
        [sdk.log.rotate]
            # 日志最大保存时间7天(单位hour)
            max_age = 720
            # 日志切割时间间隔24小时（单位hour)
            rotation_time=24

    [sdk.mysql]
        [sdk.mysql.default]
            database = "database"
            user = "root"
            host = "127.0.0.1"
            password = ""
            port = 3306
   
    [sdk.redis.default]
            host = "127.0.0.1"
            port = 6379
            password = ""
            db = 0

[app]
    debug=false
    # 微信小程序
    [app.wxmp]
        app_id="xxxxxxxxxxxxxx"

```

#### 第二步：初始化SDK实例并可以选择加载其他APP配置
 
```go
var config AppConfig
err := hdsdk.New(app, env).LoadConfig(&config).Initialize()
if err != nil {
    log.Fatal(err)
}
```

如果代码逻辑中不需要APP配置，则可以忽略加载APP配置

```go
var config AppConfig
err := hdsdk.New(app, env).Initialize()
if err != nil {
    log.Fatal(err)
}
```

如果SDK中需要Redis能力，可以在初始化的时候指定`redigo.capability`

```go
import "github.com/hdget/hdsdk/v2/provider/redis/redigo"

var config AppConfig
err := hdsdk.New(app, env).Initialize(redigo.Capability)
if err != nil {
    log.Fatal(err)
}
```

如果SDK中需要原生SQL的操作能力，可以在初始化的时候指定`sqlx_mysql.Capability`

```go
import (
    "github.com/hdget/hdsdk/v2/provider/redis/redigo"
	"github.com/hdget/hdsdk/v2/provider/db/sqlx-mysql"
)

var config AppConfig
err := hdsdk.New(app, env).Initialize(
	    redigo.Capability,
        sqlx_mysql.Capability,
	)
if err != nil {
    log.Fatal(err)
}
```

在代码中，我们通过`New(app, env)`实例化SDK再通过`LoadConfig`函数加载应用程序的所有配置信息，然后unmarshal成我们自定义的配置结构实例
- app:  加载配置的时候必须指定应用的名字
- env:  加载什么环境的配置, 可以为空，如果为空，则默认加载PROD环境的配置
