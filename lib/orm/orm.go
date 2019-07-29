package orm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"joe-micro/lib/config"
	"joe-micro/lib/log"
)

var db *gorm.DB

func init() {
	log.Info("mysql  链接中。。。")

	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.C.Mysql.UserName, config.C.Mysql.Password, config.C.Mysql.Address, config.C.Mysql.Port, config.C.Mysql.DbName))
	if err != nil {
		log.Fatalf("failed to connect database：%v", err)
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


