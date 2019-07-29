package model

import (
	"github.com/jinzhu/gorm"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
	"time"
)

//菜单表
type Menu struct {
	orm.CommonModel
	Status      uint8  `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"`             // 状态(1:启用 2:不启用)
	Memo        string `gorm:"column:memo;size:64;" json:"memo" form:"memo"`                                    // 备注
	ParentID    uint64 `gorm:"column:parent_id;not null;" json:"parent_id" form:"parent_id"`                    // 父级ID
	URL         string `gorm:"column:url;size:72;" json:"url" form:"url"`                                       // 菜单URL
	Name        string `gorm:"column:name;size:32;not null;" json:"name" form:"name"`                           // 菜单名称
	Sequence    int    `gorm:"column:sequence;not null;" json:"sequence" form:"sequence"`                       // 排序值
	MenuType    uint8  `gorm:"column:menu_type;type:tinyint(1);not null;" json:"menu_type" form:"menu_type"`    // 菜单类型 1模块2菜单3操作
	Code        string `gorm:"column:code;size:32;not null;unique_index:uk_menu_code;" json:"code" form:"code"` // 菜单代码
	OperateType string `gorm:"column:operate_type;size:32;not null;" json:"operate_type" form:"operate_type"`   // 操作类型 read/write
}

// 表名
func (Menu) TableName() string {
	return "admin_menu"
}

// 创建前
func (bc *Menu) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("id", toolfunc.GenerateUUID())
	if err != nil {
		return err
	}
	bc.CreateAt = time.Now()
	bc.UpdateAt = time.Now()
	return nil
}

// 更新前
func (bu *Menu) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdateAt = time.Now()
	return nil
}
