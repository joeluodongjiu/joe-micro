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
	bc.CreatedAt = orm.JsonTime(time.Now())
	bc.UpdatedAt = orm.JsonTime(time.Now())
	return nil
}

// 更新前
func (bu *AdminUserRoles) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdatedAt = orm.JsonTime(time.Now())
	return nil
}

// 分配用户角色
func (AdminUserRoles) SetRole(uid string, roleids []string) error {
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
	if err := tx.Where(&AdminUserRoles{UserID: uid}).Unscoped().Delete(&AdminUserRoles{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(roleids) > 0 {
		for _, rid := range roleids {
			rm := new(AdminUserRoles)
			rm.RoleID = rid
			rm.UserID = uid
			if err := tx.Create(rm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}