package casbin

import (
	"github.com/casbin/casbin"
	"joe-micro/adminApi/model"
	"joe-micro/lib/log"
	"joe-micro/lib/orm"
	"strconv"
)

const (
	PrefixUserID = "u_"
	PrefixRoleID = "r_"
)

var enforcer *casbin.Enforcer

func init() {
	/************************************/
	/********** casbin  权限管理  ********/
	/************************************/
	casbinModel := `[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub) == true && keyMatch2(r.obj, p.obj) == true && regexMatch(r.act, p.act) == true || r.sub == "r_1"`
	var err error
	enforcer, err = casbin.NewEnforcerSafe(
		casbin.NewModel(casbinModel),
	)
	if err != nil {
		log.Error(err)
		return
	}
	var roles []model.Role
	//查询所有角色
	err = orm.Find(&model.Role{}, &roles)
	if err != nil {
		return
	}
	for _, role := range roles {
		setRolePermission(enforcer, role.ID)
	}
	return
}

// 删除角色
func CsbinDeleteRole(roleids []int) {
	if enforcer == nil {
		return
	}
	for _, rid := range roleids {
		enforcer.DeletePermissionsForUser(PrefixRoleID + strconv.Itoa(rid))
		enforcer.DeleteRole(PrefixRoleID + strconv.Itoa(rid))
	}
}

// 设置角色权限
func CsbinSetRolePermission(roleid int) {
	if enforcer == nil {
		return
	}
	enforcer.DeletePermissionsForUser(PrefixRoleID + strconv.Itoa(roleid))
	setRolePermission(enforcer, roleid)
}

// 为每个角色赋值权限
func setRolePermission(enforcer *casbin.Enforcer, roleid int) {
	var rolemenus []model.RoleMenu
	err := orm.Find(&model.RoleMenu{RoleID: roleid}, &rolemenus)
	if err != nil {
		log.Info(err)
		return
	}
	for _, rolemenu := range rolemenus {
		menu := model.Menu{}
		where := model.Menu{}
		where.ID = rolemenu.MenuID
		_, err = orm.First(&where, &menu)
		if err != nil {
			log.Info(err)
			return
		}
		if menu.MenuType == 3 {
			enforcer.AddPermissionForUser(PrefixRoleID+strconv.Itoa(roleid), "/api/admin"+menu.URL,menu.OperateType)
		}
	}
}

// 检查用户是否有权限
func CsbinCheckPermission(userID, url, methodtype string) (bool, error) {
	return enforcer.EnforceSafe(PrefixUserID+userID, url, methodtype)
}

// 给用户添加角色,可以在登录是赋值到内存,也可以启动项目时加进内存
func CsbinAddRoleForUser(userid int) (err error) {
	if enforcer == nil {
		return
	}
	uid := PrefixUserID + strconv.Itoa(userid)
	enforcer.DeleteRolesForUser(uid)
	var adminsroles []model.AdminUserRoles
	err = orm.Find(&model.AdminUserRoles{UserID: userid}, &adminsroles)
	if err != nil {
		return
	}
	for _, adminsrole := range adminsroles {
		enforcer.AddRoleForUser(uid, PrefixRoleID+strconv.Itoa(adminsrole.RoleID))
	}
	return
}
