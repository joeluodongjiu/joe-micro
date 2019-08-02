package model

import (
	"github.com/jinzhu/gorm"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
	"time"
)

//角色菜单关联表
type RoleMenu struct {
	orm.CommonModel
	RoleID string `gorm:"column:role_id;not null;"  json:"role_id"` // 角色ID
	MenuID string `gorm:"column:menu_id;not null;"  json:"menu_id"` // 菜单ID
}

// 表名
func (RoleMenu) TableName() string {
	return "admin_role_menu"
}

// 创建前
func (bc *RoleMenu) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("id", toolfunc.GenerateUUID())
	if err != nil {
		return err
	}
	bc.CreatedAt = orm.JsonTime(time.Now())
	bc.UpdatedAt = orm.JsonTime(time.Now())
	return nil
}

// 更新前
func (bu *RoleMenu) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdatedAt =orm.JsonTime(time.Now())
	return nil
}
