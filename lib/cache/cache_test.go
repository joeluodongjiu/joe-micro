package cache

import (
	"fmt"
	"testing"
	"time"
)

func  TestCache(t *testing.T){
	fmt.Println(Cache.Put("astaxie", 1, 10*time.Second))
	fmt.Println(GetInt(Cache.Get("astaxie")))
	fmt.Println(Cache.IsExist("astaxie"))
	fmt.Println(Cache.Delete("astaxie"))
}
