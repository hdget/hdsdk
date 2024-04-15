### Redis配置

- default: 缺省的Redis客户端
- items: 额外的Redis客户端列表
 
 ```
[sdk.redis]
    [sdk.redis.default]
        host = "127.0.0.1"   <--- redis的服务器
        port = 6379          <--- redis的服务端口
        password = ""        <--- redis的连接密码
        db = 0               <--- redis的连接db
    [[sdk.redis.items]]
        name = "extra1"      <--- 需要通过指定name来区分使用的redis连接
        host = "127.0.0.1"
        port = 6379
        password = ""
        db = 0
    [[sdk.redis.items]]
        name = "extra2"
        host = "127.0.0.1"
        port = 6379
        password = ""
        db = 0
```
> 在配置其他Redis连接的时候需要定义在`[[sdk.redis.items]]`中，同时必须指定`name`

### Redis使用指南
  
#### 获取初始化的Redis客户端
- 获取缺省Redis客户端: `sdk.Redis.My()`
- 获取指定名字的Redis客户端: `sdk.Redis.By(string)`
    
#### 支持的Redis接口

##### 常规接口
- 删除某个key
 
    `func Del(keys []string) error`
    
- 删除多个key

    `func Dels(keys []string) error`
    
- 检查某个key是否存在

    `func Exists(key string) (bool, error)`
    
- 让某个key过期

    `func Expire(key string, expire int) error`
    
- 在某个key上加1

    `func Incr(key string) error`

- 批量发送命令到Redis并一次性执行

    `func Pipeline(commands []*CacheCommand) (reply interface{}, err error)`
     
- 检查redis是否存活

    `func Ping() error`

##### string类型
	
- 设置指定key的值为value
    
    `func Set(key string, value interface{}) error`
    
- 设置指定key的值为value,并同时设置时间为expire(单位为秒)

    `func	SetEx(key string, value interface{}, expire int) error`
    
- 获取指定key的值

    `func Get(key string) ([]byte, error)`
    
- 获取指定key的int值, 如果该key的值不为int或者不能转换成int会返回错误

    `func GetInt(key string) (int, error)`
    
- 获取指定key的int64值,如果该key的值不为int或者不能转换成int会返回错误
    
    `func GetInt64(key string) (int64, error)`
    
- 获取指定key的float64值,如果该key的值不为int64或者不能转换成int64会返回错误

    `func GetFloat64(key string) (float64, error)`
    
- 获取指定key的string值,如果该key的值不为string或者不能转换成string会返回错误

    `func GetFloat64(key string) (float64, error)`

##### HashMap类型

- 获取在哈希表中指定key的所有字段和值
    
    `func HGetAll(key string) (map[string]string, error)`

- 获取存储在哈希表中指定字段的值,返回为[]byte
    
    `func HGet(key string, field string) ([]byte, error)`

- 获取存储在哈希表中指定字段的int值,如果该key的值不为int或者不能转换成int会返回错误

    `func HGetInt(key string, field string) (int, error)`

- 获取存储在哈希表中指定字段的int64值,如果该key的值不为int64或者不能转换成int64会返回错误

    `func HGetInt64(key string, field string) (int64, error)`

- 获取存储在哈希表中指定字段的float64值,如果该key的值不为int64或者不能转换成float64会返回错误

    `func HGetFloat64(key string, field string) (float64, error)`
    
- 获取存储在哈希表中指定字段的string值,如果该key的值不为string或者不能转换成string会返回错误

    `func HGetString(key string, field string) (string, error)`
    
- 获取所有给定字段的值, 返回值为[]byte数组

    `func HMGet(key string, fields []string) ([][]byte, error)`
    
- 将哈希表key中的字段field的值设为value

    `func HSet(key string, field interface{}, value interface{}) (int, error)`

- 同时将多个field-value(域-值)对设置到哈希表key中, field-value以map[string]interface{}的形式组装

    `func HMSet(key string, args map[string]interface{}) error`

- 删除一个哈希表字段

    `func HDel(key string, field interface{}) (int, error)`
    
- 删除多个哈希表字段

    `HDels(key string, fields []interface{}) (int, error)`

##### 集合类型

- 判断 member 元素是否是集合 key 的成员

    `func SIsMember(key string, member interface{}) (bool, error)`
    
- 向集合添加一个或多个成员, members可以为一个slice

    `func SAdd(key string, members interface{}) error`
    
- 移除集合中一个或多个成员

    `func SRem(key string, members interface{}) error`
    
- 返回给定所有集合的交集, 其中每个key对应的存储数据必须为set数据类型

    `func SInter(keys []string) ([]string, error)`
    
- 返回所有给定集合的并集, 其中每个key对应的存储数据必须为set数据类型

    `func SUnion(keys []string) ([]string, error)`

- 返回第一个集合与其他集合之间的差异

    `func SDiff(keys []string) ([]string, error)`

- 返回集合中的所有成员
    
    `func SMembers(key string) ([]string, error)`

##### zset有序集合

- 向有序集合添加一个或多个成员，或者更新已存在成员的分数

    `func ZAdd(key string, score int64, member interface{}) error`
    
- 获取有序集合的成员数

    `func ZCard(key string) (int, error)`
    
- 通过索引区间返回有序集合指定区间内的成员

    `func ZRange(key string, min, max int64) (map[string]string, error)`
    
- 通过分数返回有序集合指定区间内的成员
    
    `ZRangeByScore(key string, min, max interface{}) ([]string, error)`
    
- 移除有序集合中给定的分数区间的所有成员

    `func ZRemRangeByScore(key string, min, max interface{}) error`
    
- 返回有序集中，成员的分数值

    `func ZScore(key string, member interface{}) (int64, error)`

- 计算给定的一个或多个有序集的交集并将结果集存储在新的有序集合destination key中

    `func ZInterstore(destKey string, keys ...interface{}) (int64, error)`

##### list

- 移除列表的最后一个元素，返回值为移除的元素

    `func RPop(key string) ([]byte, error)`


