package model

import (
	"github.com/jinzhu/gorm"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
	"time"
)

//角色表
type Role struct {
	orm.CommonModel
	Memo     string `gorm:"column:memo;size:64;" json:"memo" form:"memo"`              // 备注
	Name     string `gorm:"column:name;size:32;not null;" json:"name" form:"name"`     // 名称
	Sequence uint64 `gorm:"column:sequence;not null;" json:"sequence" form:"sequence"` // 排序值
}

// 表名
func (Role) TableName() string {
	return "admin_role"
}

// 创建前
func (bc *Role) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("id", toolfunc.GenerateUUID())
	if err != nil {
		return err
	}
	bc.CreatedAt = orm.JsonTime(time.Now())
	bc.UpdatedAt = orm.JsonTime(time.Now())
	return nil
}

// 更新前
func (bu *Role) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdatedAt = orm.JsonTime(time.Now())
	return nil
}


// 删除角色及关联数据
func (Role) Delete(roleids []string) error {
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
	if err := tx.Where("id in (?)", roleids).Delete(&Role{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("role_id in (?)", roleids).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
