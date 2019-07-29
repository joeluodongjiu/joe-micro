package casbin

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"joe-micro/lib/log"
	"strconv"
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
	m = g(r.sub, p.sub) == true && keyMatch2(r.obj, p.obj) == true && regexMatch(r.act, p.act) == true || r.sub == "1"`
	var err error
	enforcer, err = casbin.NewEnforcerSafe(
		casbin.NewModel(casbinModel),
	)
	if err != nil {
		return
	}
	var roles []sys.Role
	err = models.Find(&sys.Role{}, &roles)
	if err != nil {
		return
	}
	for _, role := range roles {
		setRolePermission(enforcer, role.ID)
	}
	return
}



// 角色-URL导入
func InitCsbinEnforcer() (err error) {

}

// 删除角色
func CsbinDeleteRole(roleids []uint64) {
	if enforcer == nil {
		return
	}
	for _, rid := range roleids {
		enforcer.DeletePermissionsForUser(PrefixRoleID + convert.ToString(rid))
		enforcer.DeleteRole(PrefixRoleID + convert.ToString(rid))
	}
}

// 设置角色权限
func CsbinSetRolePermission(roleid uint64) {
	if enforcer == nil {
		return
	}
	enforcer.DeletePermissionsForUser(PrefixRoleID + convert.ToString(roleid))
	setRolePermission(Enforcer, roleid)
}

// 设置角色权限
func setRolePermission(enforcer *casbin.Enforcer, roleid uint64) {
	var rolemenus []sys.RoleMenu
	err := models.Find(&sys.RoleMenu{RoleID: roleid}, &rolemenus)
	if err != nil {
		return
	}
	for _, rolemenu := range rolemenus {
		menu := sys.Menu{}
		where := sys.Menu{}
		where.ID = rolemenu.MenuID
		_, err = models.First(&where, &menu)
		if err != nil {
			return
		}
		if menu.MenuType == 3 {
			enforcer.AddPermissionForUser(PrefixRoleID+convert.ToString(roleid), "/api"+menu.URL, "GET|POST")
		}
	}
}

// 检查用户是否有权限
func CsbinCheckPermission(userID, url, methodtype string) (bool, error) {
	return enforcer.EnforceSafe(PrefixUserID+userID, url, methodtype)
}

// 用户角色处理
func CsbinAddRoleForUser(userid uint64)(err error){
	if enforcer == nil {
		return
	}
	uid:=PrefixUserID+convert.ToString(userid)
	enforcer.DeleteRolesForUser(uid)
	var adminsroles []sys.AdminsRole
	err = models.Find(&sys.AdminsRole{AdminsID: userid}, &adminsroles)
	if err != nil {
		return
	}
	for _, adminsrole := range adminsroles {
		enforcer.AddRoleForUser(uid, PrefixRoleID+convert.ToString(adminsrole.RoleID))
	}
	return
}
