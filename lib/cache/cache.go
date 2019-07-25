package cache

import (
	"fmt"
	cache2 "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"joe-micro/lib/config"
	"joe-micro/lib/log"
)

var Cache cache2.Cache

func init() {
	log.Info("redis  链接中。。。")
	var err error
	Cache, err = cache2.NewCache("redis", fmt.Sprintf(`{"key":"%v","conn":"%v","dbNum":"%v","password":"%v"}`,
		config.C.Redis.Key, config.C.Redis.Host+":"+config.C.Redis.Port, config.C.Redis.Db, config.C.Redis.Auth))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("redis链接成功")
}

// GetString convert interface to string.
func GetString(v interface{}) string {
	return cache2.GetString(v)
}

// GetInt convert interface to int.
func GetInt(v interface{}) int {
	return cache2.GetInt(v)
}

// GetInt64 convert interface to int64.
func GetInt64(v interface{}) int64 {
	return cache2.GetInt64(v)
}

// GetFloat64 convert interface to float64.
func GetFloat64(v interface{}) float64 {
	return cache2.GetFloat64(v)
}

// GetBool convert interface to bool.
func GetBool(v interface{}) bool {
	return cache2.GetBool(v)
}
