package cache

import (
	"fmt"
	"testing"
	"time"
)

func  TestCache(t *testing.T){
	fmt.Println(myCache.Put("astaxie", 1, 10*time.Second))
	fmt.Println(GetInt(myCache.Get("astaxie")))
	fmt.Println(myCache.IsExist("astaxie"))
	fmt.Println(myCache.Delete("astaxie"))
}
