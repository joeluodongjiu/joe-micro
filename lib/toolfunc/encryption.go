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
	m5.Write([]byte(password))
	m5.Write([]byte(salt))
	return salt+hex.EncodeToString(m5.Sum(nil))
}

//随机生成7位的盐值
func GenerateSalt()  (salt string) {
	var  basisStr= []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	//邀请码生成
	var str string
	for i := 0; i < 7; i++ {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(35) //35位字符
		str += string(basisStr[num])
	}
	return str
}
