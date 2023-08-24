Enterprise ready, robust and easy extensible sdk which help to quickly develop backedn services.

### SDK

SDK中支持的环境`env`定义
- local: 本地环境
- dev: 开发环境
- test: 测试环境
- pre: 预发布环境
- sim: 仿真环境
- prod: 生产环境

### SDK用户指南

#### SDK快速上手

下面是一个简单的使用hdsdk的数据库能力的示例，代码非常简单，但支持多环境多配置源初始化加载各种配置项，同时数据库支持主备或指定数据库源

```go
type rootConf struct {
    sdk.Config `mapstructure:",squash"`
}

var config rootConf
err := hdsdk.NewConfig("app", "test").Load(&config)
if err != nil {
    utils.LogFatal("load config", "err", err)
}

err = hdsdk.Initialize(config)
if err != nil {
    utils.LogFatal("sdk initialize", "err", err)
}

var total int64
err = hdsdk.Mysql.My().Get(&total, `SELECT COUNT(1) FROM table`)
if err != nil  {
	hdsdk.Logger.Error("db get total", "err", err)
}
```

#### SDK使用

##### 第一步：定义配置结构体
1. 通常我们必须通过自定义一个继承自`sdk.Config`的结构体来包括sdk配置信息。这里请注意，`sdk.Config`的tag必须要加上`mapstructure:",squash"`tag, 因为sdk配置能力底层是通过`viper`来实现，而`viper`是通过`mapstructure`这个库进行配置信息的`encode/decode`的，加上这个tag告诉`viper`在读取配置信息的时候将该结构体的字段作为提升到父结构中
```go
type rootConf struct {
    sdk.Config `mapstructure:",squash"`
    // 其他的配置项
    ...
}
```

2. 其他配置项可以定义在结构体的其他地方并支持嵌套，注意一定要用`mapstructure`tag来定义和配置文件中一致的配置项名字, e,g: `mapstructure:"debug"`
```go
type rootConf struct {
	sdk.Config  `mapstructure:",squash"`
	App confApp `mapstructure:"app"`
}

type confApp struct {
	Debug bool    `mapstructure:"debug"`
	Wxmp confWxmp `mapstructure:"wxmp"`
}

type confWxmp struct {
	AppId string `mapstructure:"app_id"`
}

```

##### 第二步： 生成配置文件

sdk段落中的全部为sdk自身能力的配置项，例如Etcd能力，Redis能力，MySQL能力，其他的段落可以用来为应用程序进行个性化的定制， 例如：

```toml
[sdk]
    [sdk.log]
        # 当前支持日志级别: "trace", "debug", "info", "warn", "error", "fatal", "panic"
        level = "debug"
        # 日志文件名称
        filename = "base.log"
        # 日志结转配置
        [sdk.log.rotate]
            # 日志最大保存时间7天(单位hour)
            max_age = 720
            # 日志切割时间间隔24小时（单位hour)
            rotation_time=24
    [sdk.etcd]
        url = "http://192.168.1.1:2379"
    [sdk.mysql]
        [sdk.mysql.default]
            database = "hdmall_base"
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

##### 第三步：初始化配置并加载配置
1. 通常我们使用`Load()`函数来加载配置，`Load()`函数会先尝试加载本地配置，默认行为会在某些环境下会再尝试加载远程配置
2. 相同配置项优先级从高到低为：
- 环境配置（高）
- 文件配置（中）
- 远程配置（低）

同一个配置变量如果在多个配置源中出现，最终生效的为高优先级配置源中的值, 例如文件配置中`app.wxmp.app_id=1`，远程配置中`app.wxmp.app_id=0`, 那么最终生效的值为0
3. 一般来说，我们会指定环境`env`，如果`env`为空，则默认加载仅包含日志配置的最小配置
```go
var conf rootConf
err := NewConfig("test", "local").Load(&conf)
if err != nil {
    utils.LogFatal("sdk initialize", "err", err)
}
```

在代码中，我们通过`NewConfig(app, env)`实例化再通过`Load()`或`LoadLocal()`或`LoadRemote()`函数加载应用程序的所有配置信息，然后unmarshal成我们自定义的配置结构实例
- app:  加载配置的时候必须指定应用的名字
- env:  加载什么环境的配置, 可以为空，如果为空，仅初始化日志配置
- args: 通过不同选项函数来自定义加载配置的行为

```
var conf rootConf
err := NewConfig("test", "local").Load(&conf)
if err != nil {
    utils.LogFatal("sdk initialize", "err", err)
}
```

SDK配置初始化时通过不同的选项函数改变默认行为，例如下面就指定了加载指定的配置文件

```go
options := make([]hdsdk.ConfigOption, 0)
if configFile != "" {
    options = append(options, hdsdk.WithConfigFile(configFile))
}

var conf rootConf
err := NewConfig("test", "local", options...).Load(&conf)
if err != nil {
    utils.LogFatal("sdk initialize", "err", err)
}
```

##### 第四步：初始化SDK
在我们加载配置信息以后，我们需要通过已经初始化好的配置结构体来初始化SDK的各项能力
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

#### 配置参考
1. 配置类型
- 默认配置文件类型为`toml`类型，可以`NewConfig()`的时候通过`WithConfigType()`选项函数来改变
- 相关函数
- WithConfigType()： 指定配置文件的类型，支持类型同`viper`的支持类型，默认为`toml`
 
2. 远程配置
- 除无环境和local环境外，运行Load()函数会默认加载远程配置，默认行为可以`NewConfig()`的时候通过`WithDisableRemoteEnvs()`选项函数来改变 
- 加载远程配置时会优先以本地配置`sdk.etcd.url`定义的值来加载远程配置，如果无`sdk.etcd.url`本地配置默认使用URL: http://127.0.0.1:2379
- 远程配置加载成功，默认开启远程配置监控，如果远程配置项有变化，5秒钟后生效。该默认行为可以`NewConfig()`的时候通过`WithWatch()`选项函数来改变
- 可以单独使用LoadRemote()函数来仅加载远程配置而忽略本地配置，注意这个时候本地配置的`sdk.etcd.url`不生效，默认从`http://127.0.0.1:2379`加载，默认行为可以`NewConfig()`的时候通过`WithRemote()`选项函数来改变 
- 相关函数
  * WithRemote()： 定义远程配置提供者的信息
  * WithDisableRemoteEnvs()： 定义在哪些环境不开启远程配置功能
  * WithWatch()： 用来开关远程配置监控以及定义远程配置更改生效时间

3. 本地配置
- 可以单独使用LoadLocal()函数仅获取本地配置， 例如从本地配置文件和环境变量中加载配置， 注意这个时候不会处理远程配置
- 默认从当前目录以及当前目录的上级目录下搜索`setting/app/<app>/<app>.<env>.toml`， 默认搜索条件可以`NewConfig()`的时候通过`WithConigDir()`和`WithConfigFilename()`选项函数来改变
- 如果需要指定配置文件，可以`NewConfig()`的时候通过`WithConfigFile()`选项函数来指定, 注意这个时候不会进行配置文件智能搜索
- 相关函数
  * WithConfigDir()： 增加智能搜索的目录，可以有多个，搜索的时候会自动加入指定目录和指定目录的上级目录
  * WithConfigFilename()： 在指定的目录下搜索文件名，不需要带文件名后缀

4. 测试用例配置

有时候我们在单元测试的时候缺少缺少配置文件，我们可以尝试直接定义配置内容并读取加载
```go
const configTestMysql = `
[sdk]
    [sdk.log]
        filename = "demo.log"
        [sdk.log.rotate]
            # 最大保存时间7天(单位hour)
            max_age = 168
            # 日志切割时间间隔24小时（单位hour)
            rotation_time=24
    [sdk.mysql]
        [sdk.mysql.default]
            database = "mysql"
            user = "root"
            password = ""
            host = "127.0.0.1"
            port = 3306

var conf testConf
err := NewConfig("test", "local").ReadString(configTestMysql).Load(&conf)
if err != nil {
    utils.LogFatal("sdk initialize", "err", err)
}
```

### SDK使用
- 日志
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

- 数据库
    * MySQL: 请参考[MySQL能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/db/mysql)
  
- 缓存
    * Redis: 请参考[Redis能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/cache/redis)

- 消息队列
    * RabbitMq: 请参考[RabbitMQ能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/mq/rabbitmq)
    * Kafka: 请参考[Kafka能力介绍](https://github.com/hdget/hdsdk/tree/main/provider/mq/kafka)

