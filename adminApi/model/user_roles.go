package model

import (
	"github.com/jinzhu/gorm"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
	"time"
)

type AdminUserRoles struct {
	orm.CommonModel
	UserID string `gorm:"column:user_id;not null;"` // 管理员ID
	RoleID string `gorm:"column:role_id;not null;"` // 角色ID
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
	bc.CreateAt = orm.JsonTime(time.Now())
	bc.UpdateAt = orm.JsonTime(time.Now())
	return nil
}

// 更新前
func (bu *AdminUserRoles) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdateAt = orm.JsonTime(time.Now())
	return nil
}
