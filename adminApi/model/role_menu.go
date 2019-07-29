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
	RoleID int `gorm:"column:role_id;not null;"` // 角色ID
	MenuID int `gorm:"column:menu_id;not null;"` // 菜单ID
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
	bc.CreateAt = time.Now()
	bc.UpdateAt = time.Now()
	return nil
}

// 更新前
func (bu *RoleMenu) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdateAt = time.Now()
	return nil
}