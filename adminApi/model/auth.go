package model

import (
	"github.com/jinzhu/gorm"
	"joe-micro/lib/orm"
)

type Admin_user struct {
	orm.CommonModel
	Username string `gorm:"Column:username"`
	Password string `gorm:"Column:password"`
	Status   int    `gorm:"Column:status"`
	RealName string `gorm:"Column:real_name"`
	Email    string `gorm:"Column:email"`
	Phone    string `gorm:"Column:phone"`
	Salt     string `gorm:"Column:salt"`
}



func GetByUsername(username string) (userM Admin_user,has bool , err error) {
	err = orm.GetDB().Where("username = ?", username).First(&userM).Error
	has = gorm.IsRecordNotFoundError(err)
	return userM,has , err
}
