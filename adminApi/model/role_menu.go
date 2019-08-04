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


// 设置角色菜单权限
func (RoleMenu) SetRole(roleid string, menuids []string) error {
	tx := orm.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(&RoleMenu{RoleID: roleid}).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(menuids) > 0 {
		for _, mid := range menuids {
			rm := new(RoleMenu)
			rm.RoleID = roleid
			rm.MenuID = mid
			if err := tx.Create(rm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}