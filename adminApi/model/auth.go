package model

import (
	"github.com/jinzhu/gorm"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
)

type Admin_user struct {
	orm.CommonModel
	Username string `gorm:"Column:username" json:"username"`
	Password string `gorm:"Column:password" json:"password"`
	Status   int    `gorm:"Column:status"  json:"status"`
	RealName string `gorm:"Column:real_name" json:"real_name"`
	Email    string `gorm:"Column:email"  json:"email"`
	Phone    string `gorm:"Column:phone" json:"phone"`
	Salt     string `gorm:"Column:salt" json:"salt"`
}


// 设置admin_user的表名为`admin_user`
func (Admin_user) TableName() string {
	return "admin_user"
}


func (au *Admin_user) BeforeCreate(scope *gorm.Scope) error {
	err :=scope.SetColumn("id", toolfunc.GenerateUUID())
	if  err!=nil {
		return err
	}
	err = scope.SetColumn("salt",toolfunc.GenerateSalt())
	if  err!=nil {
		return err
	}
	return nil
}

func GetByUsername(username string) (userM Admin_user,has bool , err error) {
	err = orm.GetDB().Where("username = ?", username).First(&userM).Error
	has = orm.IsFound(err)
	return userM,has , err
}

func CreateOne(user Admin_user) error {
    return  orm.GetDB().Create(&user).Error
}
