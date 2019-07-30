package handler

import (
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/model"
	"joe-micro/adminApi/model/casbin"
	"joe-micro/lib/cache"
	"joe-micro/lib/config"
	"joe-micro/lib/jwt"
	"joe-micro/lib/log"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
	"time"
)

type User struct{}

type login struct {
	Username string `json:"username" binding:"required" swaggo:"true,用户名"`
	Password string `json:"password" binding:"required" swaggo:"true,密码"`
}

// @Title 登录接口
// @Summary 用户登录
// @Accept  json
// @Produce  json
// @Param   body    body    handler.login      true       "json"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok,data:{token:}}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Router /api/admin/user/login [post]
func (u User) Login(c *gin.Context) {
	var req login
	if err := c.Bind(&req); err != nil {
		resBadRequest(c, err)
		return
	}
	//获取用户记录
	user := model.AdminUser{}
	where := model.AdminUser{}
	where.UserName = req.Username
	notFound, err := orm.First(&where, &user)
	if err != nil {
		if notFound {
			resBusinessP(c, "没有此条记录")
			return
		}
		log.Error(err)
		resErrSrv(c)
		return
	}

	//校验密码
	if toolfunc.EncUserPwd(req.Password, user.Salt) != user.Password {
		resBusinessP(c, "密码错误")
		return
	}

	if user.Status != 1 {
		resBusinessP(c, "该用户被封禁")
	}
	//颁发token
	token, err := jwt.CreateToken(user.ID)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	//casbin 处理
	err = casbin.CsbinAddRoleForUser(user.ID)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	//token存进cache
	err = cache.Put(user.ID, token, config.C.Jwt.TimeOut*time.Hour)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}

	//返回参数
	resData := make(map[string]string)
	resData["token"] = token

	resSuccess(c, resData)
}

// @Title 登出接口
// @Summary 用户登出
// @Accept  json
// @Produce  json
// @Param   token   header    string        true    "token"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Router /api/admin/user/logout [get]
func (u User) Logout(c *gin.Context) {
	uid := c.GetString(USER_UID_KEY)
	//从缓存中删除uid
	err := cache.Delete(uid)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
	return
}

// 获取用户信息及可访问的权限菜单
// @Title 获取用户信息及可访问的权限菜单
// @Summary 获取用户信息及可访问的权限菜单
// @Accept  json
// @Produce  json
// @Param   token   header    string        true    "token"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Router /api/admin/user/info [get]
func (u User) Info(c *gin.Context) {
    //用户uid
    uid:= c.GetString(USER_UID_KEY)
    //根据用户uid 获取菜单权限
	// 根据用户ID获取用户权限菜单
	var menuData []model.Menu
	var err error
    // 如何是超级管理员则可以看所有菜单
    if uid == SUPER_ADMIN_ID{
		err=orm.Find(&model.Menu{}, &menuData, "parent_id asc", "sequence asc")
		if err!=nil {
			log.Error(err)
			resErrSrv(c)
			return
		}
	}
	resSuccess(c,menuData)
}

// 修改用户密码
func (u User) EditPwd(c *gin.Context) {

}
