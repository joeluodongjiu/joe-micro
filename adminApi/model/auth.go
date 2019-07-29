package model

import (
	"github.com/jinzhu/gorm"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
	"time"
)

type AdminUser struct {
	orm.CommonModel
	UserName string `gorm:"column:username;size:32;not null;" json:"user_name" form:"user_name"`     // 用户名
	RealName string `gorm:"column:real_name;size:32;" json:"real_name" form:"real_name"`             // 真实姓名
	Password string `gorm:"column:password;type:char(32);not null;" json:"password" form:"password"` // 密码
	Email    string `gorm:"column:email;size:64;" json:"email" form:"email"`                         // 邮箱
	Phone    string `gorm:"column:phone;type:char(20);" json:"phone" form:"phone"`                   // 手机号
	Status   uint8  `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"`     // 状态(1:正常 2:未激活 3:暂停使用)
	Salt     string `gorm:"Column:salt" json:"salt"`
}

// 设置admin_user的表名为`admin_user`
func (AdminUser) TableName() string {
	return "admin_user"
}

func (bc *AdminUser) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("id", toolfunc.GenerateUUID())
	if err != nil {
		return err
	}
	err = scope.SetColumn("salt", toolfunc.GenerateSalt())
	if err != nil {
		return err
	}
	bc.CreateAt = time.Now()
	bc.UpdateAt = time.Now()
	return nil
}

// 更新前
func (bu *AdminUser) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdateAt = time.Now()
	return nil
}
