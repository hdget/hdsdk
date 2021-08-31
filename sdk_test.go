package hdsdk

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdsdk/provider/mq/rabbitmq"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils"
	"github.com/hdget/hdsdk/utils/alidts"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

type TestConf struct {
	Config `mapstructure:",squash"`
}

const TEST_CONFIG_LOG = `
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

const TEST_CONFIG_MYSQL = `
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

const TEST_CONFIG_REDIS = `
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

const TEST_CONFIG_RABBITMQ = `
[sdk]
	[sdk.log]
        filename = "demo.log"
		[sdk.log.rotate]
			# 最大保存时间7天(单位hour)
        	max_age = 168
        	# 日志切割时间间隔24小时（单位hour)
        	rotation_time=24
	[sdk.aliyun]
		[sdk.aliyun.default]
			host="192.168.0.114"
			username="guest"
			password="guest"
			port=5672
			vhost="/"
			[[sdk.aliyun.default.consumers]]
				name="consumer1"
				exchange_name="testexchange"
				exchange_type="direct"
				queue_name = "testqueue"
				routing_keys = [""]
			[[sdk.aliyun.default.producers]]
				name="producer1"
				exchange_name="testexchange"
				exchange_type="direct"`

const TEST_CONFIG_KAFKA = `
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

const TEST_CONFIG_ALIYUN_DTS = `
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

	// Logger.Fatal("msg content", "err", errors.New("new error"), "key1 ", "value1 123")
	Logger.Panic("msg content", "err", errors.New("new error"), "key1 ", "value1 123")
}

// nolint:errcheck
func TestLogger(t *testing.T) {
	v := LoadConfig("test", "local", "")

	// try merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_LOG)))

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

	// Logger.Fatal("msg content", "err", errors.New("new error"), "key1 ", "value1 123")
	Logger.Panic("msg content", "err", errors.New("new error"), "key1 ", "value1 123")
}

// nolint:errcheck
func TestMysql(t *testing.T) {
	v := LoadConfig("demo", "local", "")

	// merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_MYSQL)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.Fatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.Fatal("sdk initialize", "err", err)
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

// nolint:errcheck
func TestRedis(t *testing.T) {
	v := LoadConfig("demo", "local", "")

	// merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_REDIS)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.Fatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.Fatal("sdk initialize", "err", err)
	}

	Redis.My().Set("key1", 1)
	Redis.My().Set("key2", "strvalue")
	Redis.My().Set("key3", 123)

	value1, _ := Redis.My().GetInt("key1")
	assert.Equal(t, value1, 1)

	value2, _ := Redis.My().GetString("key2")
	assert.Equal(t, value2, "strvalue")

	value3, _ := Redis.My().GetInt64("key3")
	assert.Equal(t, value3, int64(123))

	k1exist, _ := Redis.My().Exists("key1")
	assert.Equal(t, k1exist, true)

	Redis.My().Del("key1")
	k1exist, _ = Redis.My().Exists("key1")
	assert.Equal(t, k1exist, false)

	Redis.My().Expire("key2", 3)
	k2exist, _ := Redis.My().Exists("key2")
	assert.Equal(t, k2exist, true)
	time.Sleep(time.Second * 4)
	k2exist, _ = Redis.My().Exists("key2")
	assert.Equal(t, k2exist, false)

	Redis.My().Incr("key3")
	v3, _ := Redis.My().GetInt("key3")
	assert.Equal(t, v3, 124)

	err = Redis.My().Ping()
	assert.Equal(t, err, nil)

	Redis.My().SetEx("key4", 456, 3)
	k4exist, _ := Redis.My().Exists("key4")
	assert.Equal(t, k4exist, true)
	time.Sleep(time.Second * 4)
	k4exist, _ = Redis.My().Exists("key4")
	assert.Equal(t, k4exist, false)

	Redis.My().HSet("key5", "field1", 111)
	k5f1, _ := Redis.My().HGetInt("key5", "field1")
	assert.Equal(t, k5f1, 111)

	Redis.My().HSet("key5", "field2", "field2value")
	k5f2, _ := Redis.My().HGetString("key5", "field2")
	assert.Equal(t, k5f2, "field2value")

	k5all, _ := Redis.My().HGetAll("key5")
	assert.Equal(t, k5all, map[string]string{
		"field1": "111",
		"field2": "field2value",
	})

	Redis.My().HMSet("key6", map[string]interface{}{
		"field1": "v1",
		"field2": "v2",
		"field3": "v3",
	})
	k6values, _ := Redis.My().HMGet("key6", []string{"field1", "field2"})
	assert.Equal(t, k6values[0], utils.StringToBytes("v1"))
	assert.Equal(t, k6values[1], utils.StringToBytes("v2"))

	Redis.My().HDels("key6", []interface{}{"field1", "field2"})
	k61v, _ := Redis.My().HGet("key6", "field1")
	assert.Equal(t, len(k61v), 0)
	k63v, _ := Redis.My().HGet("key6", "field3")
	assert.Equal(t, k63v, utils.StringToBytes("v3"))

	Redis.By("extra1").Set("key7", 333.01)
	k7v, _ := Redis.By("extra1").GetFloat64("key7")
	assert.Equal(t, k7v, 333.01)
}

// nolint:errcheck
func TestRabbitmqSend(t *testing.T) {
	v := LoadConfig("demo", "local", "")

	// merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_RABBITMQ)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.Fatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.Fatal("sdk initialize", "err", err)
	}

	p, err := Rabbitmq.My().CreateProducer("producer1")
	if err != nil {
		Logger.Fatal("create producer", "err", err)
	}

	for i := 0; i < 1000; i++ {
		s := fmt.Sprintf("%d", i)
		err = p.Publish([]byte(s))
		if err != nil {
			Logger.Error("publish", "last", p.GetLastConfirmedId(), "err", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func msgProcess(data []byte) types.MqMsgAction {
	fmt.Println(utils.BytesToString(data))
	return types.Ack
}

// nolint:errcheck
func TestRabbitmqRecv(t *testing.T) {
	v := LoadConfig("demo", "local", "")

	// merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_RABBITMQ)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.Fatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.Fatal("sdk initialize", "err", err)
	}

	mq := Rabbitmq.My()
	options := mq.GetDefaultOptions()
	qosOption := options[types.MqOptionQos].(*rabbitmq.QosOption)
	qosOption.PrefetchCount = 2
	options[types.MqOptionQos] = qosOption
	c, err := mq.CreateConsumer("consumer1", msgProcess, options)
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
	v := LoadConfig("demo", "local", "")

	// merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_KAFKA)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.Fatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.Fatal("sdk initialize", "err", err)
	}

	p, err := Kafka.My().CreateProducer("producer1")
	if err != nil {
		utils.Fatal("kafka create producer", "err", err)
	}
	defer p.Close()

	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%d", i)
		err = p.Publish([]byte(s))
		if err != nil {
			utils.Fatal("kafka producer publish", "err", err)
		}
	}
}

// nolint:errcheck
func TestKafkaRecv(t *testing.T) {
	v := LoadConfig("demo", "local", "")

	// merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_KAFKA)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.Fatal("unmarshal democonf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.Fatal("sdk initialize", "err", err)
	}

	c, err := Kafka.My().CreateConsumer("consumer1", msgProcess)
	if err != nil {
		utils.Fatal("kafka create consumer", "err", err)
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
	dts, err := alidts.New()
	if err != nil {
		utils.Fatal("new alidts", "err", err)
	}

	data, err := ioutil.ReadFile("alidts.dump")
	if err != nil {
		utils.Fatal("open alidts data", "err", err)
	}

	_, err = dts.Parse(data)
	if err != nil {
		utils.Fatal("alidts getrecord", "err", err)
	}
}

func TestUtilsAlidts(t *testing.T) {
	dts, err := alidts.New()
	if err != nil {
		utils.Fatal("new alidts", "err", err)
	}

	data, err := ioutil.ReadFile("alidts.dump")
	if err != nil {
		utils.Fatal("open alidts data", "err", err)
	}

	r, err := dts.Parse(data)
	if err != nil {
		utils.Fatal("alidts getrecord", "err", err)
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

func parseDtsData(data []byte) *alidts.DtsRecord {
	dts, err := alidts.New()
	if err != nil {
		utils.Print("error", "err new alidts")
		return nil
	}

	r, err := dts.Parse(data)
	if err != nil {
		utils.Print("error", "err parse alidts data")
		return nil
	}
	return r
}

// nolint:errcheck
func TestDts(t *testing.T) {
	v := LoadConfig("demo", "local", "")

	// merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_ALIYUN_DTS)))

	// 将配置信息转换成对应的数据结构
	var conf TestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.Fatal("unmarshal demo conf", "err", err)
	}

	err = Initialize(&conf)
	if err != nil {
		utils.Fatal("sdk initialize", "err", err)
	}

	c, err := Kafka.By("xxx").CreateConsumer("syncdata", dtsHandler)
	if err != nil {
		utils.Fatal("create consumer", "err", err)
	}

	c.Consume()
}
