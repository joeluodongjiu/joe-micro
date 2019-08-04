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
	Name        string `gorm:"column:name;size:32;not null;" json:"name" form:"name"`                           // 菜单名称
	Status      uint8  `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"`             // 状态(1:启用 2:不启用)
	ParentID    string `gorm:"column:parent_id;not null;" json:"parent_id" form:"parent_id"`                    // 父级ID
	URL         string `gorm:"column:url;size:72;" json:"url" form:"url"`                                       // 菜单URL
	Sequence    uint64 `gorm:"column:sequence;not null;" json:"sequence" form:"sequence"`                       // 排序值
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
	bc.CreatedAt = orm.JsonTime(time.Now())
	bc.UpdatedAt = orm.JsonTime(time.Now())
	return nil
}

// 更新前
func (bu *Menu) BeforeUpdate(scope *gorm.Scope) error {
	bu.UpdatedAt = orm.JsonTime(time.Now())
	return nil
}

// 获取该用户权限下所有菜单
func (Menu) GetMenusByUid(uid string, menus *[]Menu) (err error) {
	sql := `select * from admin_menu
	      where id in (
					select menu_id from admin_role_menu where 
				  role_id in (select role_id from admin_user_roles where user_id=?)
				)`
	err = orm.GetDB().Raw(sql, uid).Find(menus).Error
	return
}


//todo 删除菜单及关联数据
func (Menu) Delete(menuids []string) error {
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
	for _, menuid := range menuids {
		if err := deleteMenuRecurve(tx, menuid); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where("menu_id in (?)", menuids).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id in (?)", menuids).Delete(&Menu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func deleteMenuRecurve(db *gorm.DB, parentID string) error {
	where := &Menu{}
	where.ParentID = parentID
	var menus []Menu
	dbslect := db.Where(&where)
	if err := dbslect.Find(&menus).Error; err != nil {
		return err
	}
	for _, menu := range menus {
		if err := db.Where("menu_id = ?", menu.ID).Delete(&RoleMenu{}).Error; err != nil {
			return err
		}
		if err := deleteMenuRecurve(db, menu.ID); err != nil {
			return err
		}
	}
	if err := dbslect.Delete(&Menu{}).Error; err != nil {
		return err
	}
	return nil
}


// 获取菜单有权限的操作列表
func (Menu) GetMenuButton(uid  string, menuCode string, btns *[]string) (err error) {
	sql := `select operate_type from admin_menu
	      where id in (
					select menu_id from admin_role_menu where 
					menu_id in (select id from admin_menu where parent_id in (select id from admin_menu where code=?))
					and role_id in (select role_id from admin_user_roles where user_id=?)
				)`
	err = orm.GetDB().Raw(sql, menuCode, uid).Pluck("operate_type", btns).Error
	return
}