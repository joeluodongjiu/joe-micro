package orm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"joe-micro/lib/log"
)

var db *gorm.DB

type dbInfo struct {
	Address      string `json:"address"`
	Port         int    `json:"port"`
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
	DbName       string `json:"db_name"`
}

/*host: 192.168.0.162
user: root
pwd:  gogocuri
dbname: wanqu2
port: 3306*/

func init() {
	log.Info("mysql  链接中。。。")
	var v dbInfo
	v.UserName = "root"
	v.Port = 3306
	v.UserPassword = "gogocuri"
	v.DbName = "wanqu2"
	v.Address = "192.168.0.162"

	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		v.UserName, v.UserPassword, v.Address, v.Port, v.DbName))
	if err != nil {
		log.Fatalf("failed to connect database：", err)
	}
	// 关闭复数表名，如果设置为true，`User`表的表名就会是`user`，而不是`users`
	db.SingularTable(true)
	// 启用Logger，显示详细日志
	db.LogMode(true)
	//自定义日志
	db.SetLogger(log.NewGormLogger())
	//连接池
	db.DB().SetMaxIdleConns(50)
	db.DB().SetMaxOpenConns(200)
	log.Info("mysql 链接成功")
}

//get db
func GetDB() *gorm.DB {
	return db
}
