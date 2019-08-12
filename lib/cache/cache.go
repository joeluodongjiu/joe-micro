package cache

import (
	"fmt"
	cache2 "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"joe-micro/lib/config"
	"joe-micro/lib/log"
	"time"
)

var myCache cache2.Cache

func init() {
	log.Info("redis  链接中。。。")
	var err error
	myCache, err = cache2.NewCache("redis", fmt.Sprintf(`{"key":"%v","conn":"%v","dbNum":"%v","password":"%v","maxIdle":"30"}`,
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

// get cached value by key.
func Get(key string) interface{} {
	return myCache.Get(key)
}

// GetMulti is a batch version of Get.
func GetMulti(keys []string) []interface{} {
	return myCache.GetMulti(keys)
}

// set cached value with key and expire time.
func Put(key string, val interface{}, timeout time.Duration) error {
	return myCache.Put(key, val, timeout)
}

// delete cached value by key.
func Delete(key string) error {
	return myCache.Delete(key)
}

// check if cached value exists or not.
func IsExist(key string) bool {
	return myCache.IsExist(key)
}

/*// increase cached int value by key, as a counter.
Incr(key string) error
// decrease cached int value by key, as a counter.
Decr(key string) error

// clear all cache.
ClearAll() error
// start gc routine based on config string settings.
StartAndGC(config string) error*/
