package toolfunc

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

//用户密码加密
func EncUserPwd(password, salt string) (hash string) {
	m5 := md5.New()
	//加密串
	m5.Write([]byte("user")) //以后可以改成适合的加密串,服务端保留,不暴露出去
	//密码
	m5.Write([]byte(password))
	//盐
	m5.Write([]byte(salt))
	return salt + hex.EncodeToString(m5.Sum(nil))
}

//随机生成7位的盐值
func GenerateSalt() (salt string) {
	var basisStr = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	//邀请码生成
	var str string
	for i := 0; i < 7; i++ {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(35) //35位字符
		str += string(basisStr[num])
	}
	return str
}
