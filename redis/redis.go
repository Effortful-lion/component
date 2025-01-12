package redis

import (
	"os"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis RedisConf `yaml:"redis"`
}

// redis的连接

type RedisConf struct{
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Password string `yaml:"password"`
	DB int `yaml:"db"`
	PoolSize int `yaml:"poolsize"`
}

var Rdb *redis.Client

func InitRedis(){
	var conf Config
	dataBytes, err := os.ReadFile("./conf/redis.yaml")
	if err != nil {
		panic(err)
	}
	// 解析yaml,将数据反序列化到conf结构体
	err = yaml.Unmarshal(dataBytes, &conf)
	if err != nil {
		panic(err)
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr: conf.Redis.Host + ":" + conf.Redis.Port,
		Password: conf.Redis.Password,
		DB: conf.Redis.DB,
		PoolSize: conf.Redis.PoolSize,
	})
}

// 获得rdb对象
func GetRedis() *redis.Client {
	return Rdb
}

// 设置key
func SetKey(key string, value string) error {
	return Rdb.Set(key, value, 0).Err()
}

// 获得key
func GetKey(key string) (string, error) {
	return Rdb.Get(key).Result()
}
