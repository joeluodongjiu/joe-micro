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
	Memo     string `gorm:"column:memo;size:64;" json:"memo" form:"memo"`                 // 备注
	Name     string `gorm:"column:name;size:32;not null;" json:"name" form:"name"`        // 名称
	Sequence int    `gorm:"column:sequence;not null;" json:"sequence" form:"sequence"`    // 排序值
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
	bc.CreateAt = time.Now()
	bc.UpdateAt = time.Now()
	return nil
}

// 更新前
func (bu *Role) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdateAt = time.Now()
	return nil
}
