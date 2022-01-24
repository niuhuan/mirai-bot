package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/niuhuan/mirai-bot/config"
	"github.com/niuhuan/mirai-bot/utils"
	logger "github.com/sirupsen/logrus"
	"time"
)

const (
	setIfNotExist     = "NX" // 不存在则执行
	setWithExpireTime = "PX" // 过期时间(秒)  PX 毫秒
)

var (
	RdPool *redis.Pool
	Nil    = redis.ErrNil
)

func InitRedis() {
	uri := fmt.Sprintf(
		"%s:%d",
		config.Config.Database.Redis.Hostname,
		config.Config.Database.Redis.Port,
	)
	RdPool = &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", uri)
		},
		MaxIdle:     10,
		MaxActive:   20,
		IdleTimeout: 1000,
	}
}

func Test() {
	conn := RdPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	utils.PanicNotNil(err)
}

func DelKey(key string) {
	conn := RdPool.Get()
	defer conn.Close()
	conn.Do("DEL", key)
}

func SetString(key string, value string, duration time.Duration) bool {
	conn := RdPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value, setWithExpireTime, duration.Milliseconds())
	if err != nil {
		logger.Info(err)
	}
	return err == nil
}

func GetString(key string) string {
	conn := RdPool.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		logger.Info(err)
		return ""
	}
	return value
}

func GetStringError(key string) (string, error) {
	conn := RdPool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func SetByteArray(key string, value []byte, duration time.Duration) error {
	conn := RdPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value, setWithExpireTime, duration.Milliseconds())
	return err
}

func GetByteArray(key string) []byte {
	conn := RdPool.Get()
	defer conn.Close()
	value, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		logger.Info(err)
		return nil
	}
	return value
}

func SetInt(key string, value int, duration time.Duration) bool {
	conn := RdPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value, setWithExpireTime, duration.Milliseconds())
	if err != nil {
		logger.Info(err)
	}
	return err == nil
}

func GetInt(key string) (int, error) {
	conn := RdPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("GET", key))
}

func SetBool(key string, value bool, duration time.Duration) error {
	conn := RdPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value, setWithExpireTime, duration.Milliseconds())
	return err
}

func GetBoolErr(key string) (bool, error) {
	conn := RdPool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("GET", key))
}
