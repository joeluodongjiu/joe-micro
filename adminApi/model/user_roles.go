package model

import (
	"github.com/jinzhu/gorm"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
	"time"
)

type AdminUserRoles struct {
	orm.CommonModel
	UserID int `gorm:"column:user_id;not null;"` // 管理员ID
	RoleID int `gorm:"column:role_id;not null;"` // 角色ID
}

func (AdminUserRoles) TableName() string {
	return "admin_user_roles"
}

// 创建前
func (bc *AdminUserRoles) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("id", toolfunc.GenerateUUID())
	if err != nil {
		return err
	}
	bc.CreateAt = time.Now()
	bc.UpdateAt = time.Now()
	return nil
}

// 更新前
func (bu *AdminUserRoles) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdateAt = time.Now()
	return nil
}
