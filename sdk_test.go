package hdsdk

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdsdk/lib/aliyun"
	"github.com/hdget/hdsdk/provider/mq/rabbitmq"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

type TestConf struct {
	Config `mapstructure:",squash"`
}

const configTestLog = `
[sdk]
    [sdk.log]
        # "debug", "info", "warn", "error"
        level = "debug"
        filename = "demo.log"

        [sdk.log.rotate]
            # 最大保存时间7天(单位hour)
            max_age = 168
            # 日志切割时间间隔24小时（单位hour)
            rotation_time=24`

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
			password = "password"
			host = "127.0.0.1"
			port = 3306
		[sdk.mysql.master]
			database = "mysql"
			user = "root"
			password = "password"
			host = "127.0.0.1"
			port = 3306
		[[sdk.mysql.slaves]]
			database = "mysql"
			user = "root"
			password = "password"
			host = "127.0.0.1"
			port = 3306
		[[sdk.mysql.items]]
			name = "xxx"
			database = "mysql"
			user = "root"
			password = "password"
			host = "127.0.0.1"
			port = 3306
`

const configTestRedis = `
[sdk]
	[sdk.log]
        filename = "demo.log"
		[sdk.log.rotate]
			# 最大保存时间7天(单位hour)
        	max_age = 168
        	# 日志切割时间间隔24小时（单位hour)
        	rotation_time=24
	[sdk.redis]
		[sdk.redis.default]
			host = "127.0.0.1"
        	port = 6379
        	password = ""
        	db = 0
		[[sdk.redis.items]]
			name = "extra1"
			host = "127.0.0.1"
        	port = 6379
        	password = ""
        	db = 0
`

const configTestRabbitmq = `
[sdk]
	[sdk.log]
        filename = "demo.log"
		[sdk.log.rotate]
			# 最大保存时间7天(单位hour)
        	max_age = 168
        	# 日志切割时间间隔24小时（单位hour)
        	rotation_time=24
	[sdk.rabbitmq]
		[sdk.rabbitmq.default]
			host="127.0.0.1"
			username="guest"
			password="guest"
			port=5672
			vhost="/"`

//[[sdk.rabbitmq.default.consumers]]
//	name="consumer1"
//	exchange_name="testexchange"
//	exchange_type="direct"
//	queue_name = "testqueue"
//	routing_keys = [""]
//[[sdk.rabbitmq.default.producers]]
//	name="producer1"
//	exchange_name="testexchange"
//	exchange_type="direct"`

const configTestRabbitmqDelay = `
[sdk]
	[sdk.log]
        filename = "demo.log"
		[sdk.log.rotate]
			# 最大保存时间7天(单位hour)
        	max_age = 168
        	# 日志切割时间间隔24小时（单位hour)
        	rotation_time=24
	[sdk.rabbitmq]
		[sdk.rabbitmq.default]
			host="127.0.0.1"
			username="guest"
			password="guest"
			port=5672
			vhost="/"
			[[sdk.rabbitmq.default.consumers]]
				name="consumer1"
				exchange_name="exchange_delay"
				exchange_type="delay:topic"
				queue_name = "queue1"
				routing_keys = ["close"]
			[[sdk.rabbitmq.default.consumers]]
				name="consumer2"
				exchange_name="exchange_delay"
				exchange_type="delay:topic"
				queue_name = "queue1"
				routing_keys = ["delivery"]
			[[sdk.rabbitmq.default.producers]]
				name="producer1"
				exchange_name="exchange_delay"
				exchange_type="delay:topic"`

const configTestKafka = `
[sdk]
	[sdk.log]
        filename = "demo.log"
		[sdk.log.rotate]
			# 最大保存时间7天(单位hour)
        	max_age = 168
        	# 日志切割时间间隔24小时（单位hour)
        	rotation_time=24
	[sdk.kafka]
		[sdk.kafka.default]
			brokers=["192.168.0.114:9094"]
			[[sdk.kafka.default.producers]]
				name = "producer1"	
				topics=["testtopic1", "testtopic2"]
			[[sdk.kafka.default.consumers]]
				name = "consumer1"	
				group_id="testgroup"
				topic="testtopic1"
`

const configTestAliyunDts = `
[sdk]
	[sdk.log]
        filename = "demo.log"
		[sdk.log.rotate]
			# 最大保存时间7天(单位hour)
        	max_age = 168
        	# 日志切割时间间隔24小时（单位hour)
        	rotation_time=24
	[sdk.kafka]
		[[sdk.kafka.items]]
			name = "xxx"
			brokers=["127.0.0.1:18003"]
			[[sdk.kafka.items.consumers]]
				name = "syncdata"	
				user = "testuser"
                password = "testpassword"
                group_id = "testgroup"
                topic = "testtopic"
`

const configTestNeo4j = `
[sdk]
	[sdk.log]
        filename = "demo.log"
		[sdk.log.rotate]
			# 最大保存时间7天(单位hour)
        	max_age = 168
        	# 日志切割时间间隔24小时（单位hour)
        	rotation_time=24
	[sdk.neo4j]
		virtual_uri = "neo4j://test.newaigou.com:7687"
		username = "neo4j"
		password = "123456"
		#[[sdk.neo4j.servers]]
		#	host = "xxx"
		#	port = 1234
`

// nolint:errcheck
func TestEmptyLogger(t *testing.T) {
	err := Initialize(nil)
	if err != nil {
		log.Fatalf("msg=\"sdk initialize\" error=\"%v\"", err)
	}

	e1 := errors.New("e1")
	e2 := errors.Wrap(e1, "e2")

	Logger.Info("msg content")
	Logger.Info("msg content", "err")
	Logger.Info("msg content", "err", nil)
	Logger.Info("msg content", "err", "error message")
	Logger.Info("msg content", "err", errors.New("new error"))
	Logger.Info("msg content", "err", e2)
	Logger.Info("msg content", "err", errors.New("new error"), "key1 ")
	Logger.Info("msg content", "err", errors.New("new error"), "key1 ", 123)
	Logger.Info("msg content", "err", errors.New("new error"), "key1 ", "value1 123")

	Logger.Warn("msg content")
	Logger.Warn("msg content", "err")
	Logger.Warn("msg content", "err", nil)
	Logger.Warn("msg content", "err", errors.New("new error"))
	Logger.Warn("msg content", "err", e2)
	Logger.Warn("msg content", "err", errors.New("new error"), "key1 ")
	Logger.Warn("msg content", "err", errors.New("new error"), "key1 ", 123)
	Logger.Warn("msg content", "err", errors.New("new error"), "key1 ", "value1 123")

	Logger.Debug("msg content")
	Logger.Debug("msg content", "err")
	Logger.Debug("msg content", "err", nil)
	Logger.Debug("msg content", "err", errors.New("new error"))
	Logger.Debug("msg content", "err", e2)
	Logger.Debug("msg content", "err", errors.New("new error"), "key1 ")
	Logger.Debug("msg content", "err", errors.New("new error"), "key1 ", 123)
	Logger.Debug("msg content", "err", errors.New("new error"), "key1 ", "value1 123")

	Logger.Error("msg content")
	Logger.Error("msg content", "err")
	Logger.Error("msg content", "err", nil)
	Logger.Error("msg content", "err", errors.New("new error"))
	Logger.Error("msg content", "err", e2)
	Logger.Error("msg content", "err", errors.New("new error"), "key1 ")
	Logger.Error("msg content", "err", errors.New("new error"), "key1 ", 123)
	Logger.Error("msg content", "err", errors.New("new error"), "key1 ", "value1 123")

	// Logger.LogFatal("msg content", "err", errors.New("new error"), "key1 ", "value1 123")
	Logger.Panic("msg content", "err", errors.New("new error"), "key1 ", "value1 123")
}

// nolint:errcheck
func TestLogger(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// try merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestLog)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unmarshal democonf, error=%v", err)
	}

	err = Initialize(&conf)
	if err != nil {
		log.Fatalf("msg=\"sdk initialize\" error=\"%v\"", err)
	}

	e1 := errors.New("e1")
	e2 := errors.Wrap(e1, "e2")

	Logger.Info("msg content")
	Logger.Info("msg content", "err")
	Logger.Info("msg content", "err", nil)
	Logger.Info("msg content", "err", errors.New("new error"))
	Logger.Info("msg content", "err", e2)
	Logger.Info("msg content", "err", errors.New("new error"), "key1 ")
	Logger.Info("msg content", "err", errors.New("new error"), "key1 ", 123)
	Logger.Info("msg content", "err", errors.New("new error"), "key1 ", "value1 123")

	Logger.Warn("msg content")
	Logger.Warn("msg content", "err")
	Logger.Warn("msg content", "err", nil)
	Logger.Warn("msg content", "err", errors.New("new error"))
	Logger.Warn("msg content", "err", e2)
	Logger.Warn("msg content", "err", errors.New("new error"), "key1 ")
	Logger.Warn("msg content", "err", errors.New("new error"), "key1 ", 123)
	Logger.Warn("msg content", "err", errors.New("new error"), "key1 ", "value1 123")

	Logger.Debug("msg content")
	Logger.Debug("msg content", "err")
	Logger.Debug("msg content", "err", nil)
	Logger.Debug("msg content", "err", errors.New("new error"))
	Logger.Debug("msg content", "err", e2)
	Logger.Debug("msg content", "err", errors.New("new error"), "key1 ")
	Logger.Debug("msg content", "err", errors.New("new error"), "key1 ", 123)
	Logger.Debug("msg content", "err", errors.New("new error"), "key1 ", "value1 123")

	Logger.Error("msg content")
	Logger.Error("msg content", "err")
	Logger.Error("msg content", "err", nil)
	Logger.Error("msg content", "err", errors.New("new error"))
	Logger.Error("msg content", "err", e2)
	Logger.Error("msg content", "err", errors.New("new error"), "key1 ")
	Logger.Error("msg content", "err", errors.New("new error"), "key1 ", 123)
	Logger.Error("msg content", "err", errors.New("new error"), "key1 ", "value1 123")

	// Logger.LogFatal("msg content", "err", errors.New("new error"), "key1 ", "value1 123")
	Logger.Panic("msg content", "err", errors.New("new error"), "key1 ", "value1 123")
}

// nolint:errcheck
func TestMysql(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestMysql)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	var total1 int
	err = Mysql.My().Get(&total1, "SELECT count(1) FROM db")
	Logger.Debug("get total from default db", "total", total1, "err", err)

	var total2 int
	err = Mysql.Master().Get(&total2, "SELECT count(1) FROM db")
	Logger.Debug("get total from master db", "total", total2, "err", err)

	var total3 int
	err = Mysql.Slave(0).Get(&total3, "SELECT count(1) FROM db")
	Logger.Debug("get total from slave db", "total", total3, "err", err)

	var total4 int
	err = Mysql.By("xxx").Get(&total4, "SELECT count(1) FROM db")
	Logger.Debug("get total from extra db", "total", total4, "err", err)
}

//const lusHasStock = `
//for i,v in pairs(KEYS) do
//	local ret = redis.call("Get", v)
//	if( ret - ARGV[i] <= 0 ) then
//		error("not enough stock")
//	end;
//end;
//return 1
//`

const luaDeduckStock = `
for i,v in pairs(KEYS) do 
	local current = redis.call('GET', v)
	local delta = current - ARGV[i]
	if( delta >= 0 ) then
		local leftStock = redis.call("DECRBY", v, ARGV[i])
	else
		error("not enough stock")
	end;
end;
return 1
`

//const luaRevertStock = `
//for i,v in pairs(KEYS) do
//	redis.call("INCRBY", v, ARGV[i])
//end;
//return 1
//`

// nolint:errcheck
func TestRedis(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestRedis)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	result, err := Redis.My().Eval(luaDeduckStock, []interface{}{"stock:123:25:2343", "stock:123:25:2342"}, []interface{}{100, 200})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("eval result:%v", result)

	_ = Redis.My().Set("key1", 1)
	_ = Redis.My().Set("key2", "strvalue")
	_ = Redis.My().Set("key3", 123)

	value1, _ := Redis.My().GetInt("key1")
	assert.Equal(t, value1, 1)

	value2, _ := Redis.My().GetString("key2")
	assert.Equal(t, value2, "strvalue")

	value3, _ := Redis.My().GetInt64("key3")
	assert.Equal(t, value3, int64(123))

	k1exist, _ := Redis.My().Exists("key1")
	assert.Equal(t, k1exist, true)

	_ = Redis.My().Del("key1")
	k1exist, _ = Redis.My().Exists("key1")
	assert.Equal(t, k1exist, false)

	_ = Redis.My().Expire("key2", 3)
	k2exist, _ := Redis.My().Exists("key2")
	assert.Equal(t, k2exist, true)
	time.Sleep(time.Second * 4)
	k2exist, _ = Redis.My().Exists("key2")
	assert.Equal(t, k2exist, false)

	_ = Redis.My().Incr("key3")
	v3, _ := Redis.My().GetInt("key3")
	assert.Equal(t, v3, 124)

	err = Redis.My().Ping()
	assert.Equal(t, err, nil)

	_ = Redis.My().SetEx("key4", 456, 3)
	k4exist, _ := Redis.My().Exists("key4")
	assert.Equal(t, k4exist, true)
	time.Sleep(time.Second * 4)
	k4exist, _ = Redis.My().Exists("key4")
	assert.Equal(t, k4exist, false)

	_, _ = Redis.My().HSet("key5", "field1", 111)
	k5f1, _ := Redis.My().HGetInt("key5", "field1")
	assert.Equal(t, k5f1, 111)

	_, _ = Redis.My().HSet("key5", "field2", "field2value")
	k5f2, _ := Redis.My().HGetString("key5", "field2")
	assert.Equal(t, k5f2, "field2value")

	k5all, _ := Redis.My().HGetAll("key5")
	assert.Equal(t, k5all, map[string]string{
		"field1": "111",
		"field2": "field2value",
	})

	_ = Redis.My().HMSet("key6", map[string]interface{}{
		"field1": "v1",
		"field2": "v2",
		"field3": "v3",
	})
	k6values, _ := Redis.My().HMGet("key6", []string{"field1", "field2"})
	assert.Equal(t, k6values[0], utils.StringToBytes("v1"))
	assert.Equal(t, k6values[1], utils.StringToBytes("v2"))

	_, _ = Redis.My().HDels("key6", []interface{}{"field1", "field2"})
	k61v, _ := Redis.My().HGet("key6", "field1")
	assert.Equal(t, len(k61v), 0)
	k63v, _ := Redis.My().HGet("key6", "field3")
	assert.Equal(t, k63v, utils.StringToBytes("v3"))

	_ = Redis.By("extra1").Set("key7", 333.01)
	k7v, _ := Redis.By("extra1").GetFloat64("key7")
	assert.Equal(t, k7v, 333.01)

	_ = Redis.My().Del("key8")
	_ = Redis.My().LPush("key8", 1, 2)
	k8v, _ := Redis.My().LRangeInt64("key8", 0, 5)
	assert.Equal(t, k8v[0], 1)
	assert.Equal(t, k8v[1], 2)
}

// nolint:errcheck
func TestRabbitmqSend(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestRabbitmq)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal democonf", "err", err)
	}

	fmt.Println(v.AllKeys())

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	params := map[string]interface{}{
		"exchangeName": "default",
		"exchangeType": "topic",
	}
	p, err := Rabbitmq.My().CreateProducer(params)
	if err != nil {
		Logger.Fatal("create producer", "err", err)
	}

	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%d", i)
		err = p.Publish([]byte(s), "test")
		if err != nil {
			Logger.Error("publish", "last", p.GetLastConfirmedId(), "err", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func TestRabbitmqSendDelay(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestRabbitmqDelay)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	params := map[string]interface{}{
		"exchangeName": "delay",
		"exchangeType": "delay:topic",
	}
	p, err := Rabbitmq.My().CreateProducer(params)
	if err != nil {
		Logger.Fatal("create producer", "err", err)
	}

	err = p.PublishDelay([]byte("1"), int64(60000), "close")
	if err != nil {
		Logger.Error("publish", "last", p.GetLastConfirmedId(), "err", err)
	}

	err = p.PublishDelay([]byte("2"), int64(70000), "close")
	if err != nil {
		Logger.Error("publish", "last", p.GetLastConfirmedId(), "err", err)
	}
}

func msgProcess(data []byte) types.MqMsgAction {
	fmt.Println(time.Now(), utils.BytesToString(data))
	return types.Ack
}

// nolint:errcheck
func TestRabbitmqRecv(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestRabbitmq)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	qosOption := Rabbitmq.My().GetDefaultOptions()[types.MqOptionQos].(*rabbitmq.QosOption)
	qosOption.PrefetchCount = 2
	params := map[string]interface{}{
		"exchangeName": "default",
		"exchangeType": "topic",
		"routingKeys":  []string{"test"},
	}
	c, err := Rabbitmq.My().CreateConsumer(msgProcess, params, qosOption)
	if err != nil {
		Logger.Fatal("create consumer", "err", err)
	}

	//go func() {
	//	time.Sleep(3 * time.Second)
	//	c.close()
	//}()

	c.Consume()
}

// nolint:errcheck
func TestRabbitmqRecvDelay(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestRabbitmq)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	qosOption := Rabbitmq.My().GetDefaultOptions()[types.MqOptionQos].(*rabbitmq.QosOption)
	qosOption.PrefetchCount = 2
	params := map[string]interface{}{
		"queueName":    "close",
		"exchangeName": "delay",
		"exchangeType": "delay:topic",
		"routingKeys":  []string{"close"},
	}
	c, err := Rabbitmq.My().CreateConsumer(msgProcess, params, qosOption)
	if err != nil {
		Logger.Fatal("create consumer", "err", err)
	}

	//go func() {
	//	time.Sleep(3 * time.Second)
	//	c.close()
	//}()

	c.Consume()
}

// nolint:errcheck
func TestKafkaSend(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestKafka)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	p, err := Kafka.My().CreateProducer(nil)
	if err != nil {
		utils.LogFatal("kafka create producer", "err", err)
	}
	defer p.Close()

	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%d", i)
		err = p.Publish([]byte(s))
		if err != nil {
			utils.LogFatal("kafka producer publish", "err", err)
		}
	}
}

// nolint:errcheck
func TestKafkaRecv(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestKafka)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	c, err := Kafka.My().CreateConsumer(msgProcess, nil)
	if err != nil {
		utils.LogFatal("kafka create consumer", "err", err)
	}
	defer c.Close()

	c.Consume()
}

func BenchmarkHamba(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseByHamba()
	}
}

func parseByHamba() {
	dts, err := aliyun.New()
	if err != nil {
		utils.LogFatal("new alidts", "err", err)
	}

	data, err := os.ReadFile("alidts.dump")
	if err != nil {
		utils.LogFatal("open alidts data", "err", err)
	}

	_, err = dts.Parse(data)
	if err != nil {
		utils.LogFatal("alidts getrecord", "err", err)
	}
}

func TestUtilsAlidts(t *testing.T) {
	dts, err := aliyun.New()
	if err != nil {
		utils.LogFatal("new alidts", "err", err)
	}

	data, err := os.ReadFile("alidts.dump")
	if err != nil {
		utils.LogFatal("open alidts data", "err", err)
	}

	r, err := dts.Parse(data)
	if err != nil {
		utils.LogFatal("alidts getrecord", "err", err)
	}

	fmt.Println(r)
}

func dtsHandler(data []byte) types.MqMsgAction {
	r := parseDtsData(data)
	fmt.Printf("%v %s %s.%s [%s]",
		time.Unix(r.SourceTimeStamp, 0),
		r.SourceTxId,
		r.Database,
		r.Table,
		r.Operation,
	)
	return types.Ack
}

func parseDtsData(data []byte) *aliyun.DtsRecord {
	dts, err := aliyun.New()
	if err != nil {
		utils.LogError("err new alidts")
		return nil
	}

	r, err := dts.Parse(data)
	if err != nil {
		utils.LogError("err parse alidts data")
		return nil
	}
	return r
}

// nolint:errcheck
func TestDts(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestAliyunDts)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal demo conf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	c, err := Kafka.By("xxx").CreateConsumer(dtsHandler, nil)
	if err != nil {
		utils.LogFatal("create consumer", "err", err)
	}

	c.Consume()
}

// nolint:errcheck
func TestNeo4j(t *testing.T) {
	v := NewConfig("test", "local").Load()

	// merge config from string
	_ = v.MergeConfig(bytes.NewReader(utils.StringToBytes(configTestNeo4j)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.LogFatal("unmarshal demo conf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.LogFatal("sdk initialize", "err", err)
	}

	works := []neo4j.TransactionWork{
		graphDeleteAllPersons(),
		graphAddPerson("A"),
		graphAddPerson("B"),
		graphAddPerson("C"),
		graphAddReferralRelation("A", "B"),
		graphAddReferralRelation("B", "C"),
		graphDeletePerson("B"),
	}

	_, err = Neo4j.Exec(works)
	if err != nil {
		utils.LogFatal("neo4j exec", "err", err)
	}

	result, err := Neo4j.Select("MATCH (a:Person) RETURN a")
	if err != nil {
		utils.LogFatal("neo4j select", "err", err)
	}
	fmt.Println(result)

	type Person struct {
		Name string `json:"name"`
	}
	ddd, err := Neo4j.Get("MATCH (a:Person {name: $Name}) RETURN a", &Person{Name: "A"})
	if err != nil {
		utils.LogFatal("neo4j get", "err", err)
	}
	fmt.Println(ddd)
}

// 找到匹配的a节点，然后detach命令会删除a节点相关的所有关系
func graphDeleteAllPersons() neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		var result, err = tx.Run(
			"MATCH (a:Person) DETACH DELETE a", nil)

		if err != nil {
			return nil, err
		}

		return result.Consume()
	}
}

func graphDeletePerson(name string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		// 1. 删除
		result, err := tx.Run(
			"MATCH (a:Person)-[:REFERRAL]-(b:Person {name: $name})-[:REFERRAL]-(c:Person) RETURN a,c", map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		var from, to dbtype.Node
		if result.Next() {
			from = result.Record().Values[0].(dbtype.Node)
			to = result.Record().Values[1].(dbtype.Node)
		}

		// 2. 删除节点
		result, err = tx.Run(
			"MATCH (a:Person {name: $name}) DETACH DELETE a", map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		// 3. 增加边
		if from.Id > 0 && to.Id > 0 {
			result, err = tx.Run(
				"MATCH (a:Person {name: $from}),(b:Person {name: $to}) MERGE (a)-[:REFERER]->(b)", map[string]interface{}{"from": from.Props["name"], "to": to.Props["name"]})
			if err != nil {
				return nil, err
			}
		}
		return result.Consume()
	}
}

func graphAddPerson(person1 string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		var result, err = tx.Run(
			"CREATE (a:Person {name: $name1})", map[string]interface{}{"name1": person1})

		if err != nil {
			return nil, err
		}

		return result.Consume()
	}
}

func graphAddReferralRelation(person1 string, person2 string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		var result, err = tx.Run(
			"MATCH (a:Person {name: $name1}),"+
				"(b:Person {name: $name2}) "+
				"MERGE (a)-[:REFERRAL]->(b)", map[string]interface{}{"name1": person1, "name2": person2})

		if err != nil {
			return nil, err
		}

		return result.Consume()
	}
}
