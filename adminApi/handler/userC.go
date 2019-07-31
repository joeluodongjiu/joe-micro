package handler

import (
	"github.com/ahmetb/go-linq/v3"
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

type UserController struct{} //用户操作控制器

type login struct {
	Username string `json:"username" binding:"required" ` //用户名
	Password string `json:"password" binding:"required" ` // 密码
}

// @Title 登录接口
// @Summary 用户登录
// @Tags user   用户操作
// @Accept  json
// @Produce  json
// @Param   body    body    handler.login   false   "body"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok,data:{token:}}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Router /user/login [post]
func (UserController) Login(c *gin.Context) {
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
// @Tags user   用户操作
// @Accept  json
// @Produce  json
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user/logout [get]
func (UserController) Logout(c *gin.Context) {
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
// @Tags user   用户操作
// @Accept  json
// @Produce  json
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user/info [get]
func (UserController) Info(c *gin.Context) {
	//用户uid
	uid := c.GetString(USER_UID_KEY)
	//根据用户uid 获取菜单权限
	// 根据用户ID获取用户权限菜单
	var menuData []model.Menu
	var err error
	// 如何是超级管理员则可以看所有菜单
	log.Warn(uid)
	if uid == SUPER_ADMIN_ID {
		err = orm.Find(&model.Menu{}, &menuData, "parent_id asc", "sequence asc")
		if err != nil {
			log.Error(err)
			resErrSrv(c)
			return
		}
	} else {
		menu := model.Menu{}
		err = menu.GetMenusByUid(uid, &menuData)
		if err != nil {
			log.Error(err)
			resErrSrv(c)
			return
		}
	}
	//整理菜单为父子级展示
	var menus []MenuModel
	if len(menuData) > 0 {
		var topmenuid = menuData[0].ParentID
		if topmenuid == "" {
			topmenuid = menuData[0].ID
		}
		menus = setMenu(menuData, topmenuid)
	}
	user := model.AdminUser{}
	notFound, err := orm.FirstByID(&user, uid)
	if err != nil {
		if notFound {
			resBusinessP(c, "没有此条记录")
			return
		}
		log.Error(err)
		resErrSrv(c)
		return
	}
	resData := map[string]interface{}{
		"user":  user,
		"menus": menus,
	}
	resSuccess(c, resData)
}

type MenuModel struct {
	Path      string      `json:"path"`      // 路由
	Component string      `json:"component"` // 可以对应前端的控制器的 name
	Name      string      `json:"name"`      // 菜单名称
	Hidden    bool        `json:"hidden"`    // 是否隐藏
	Children  []MenuModel `json:"children"`  // 子级菜单
}

type editPwdReq struct {
	OldPassword string `json:"old_password" binding:"required" `  //旧密码
	NewPassword string `json:"new_password" validate:"required" ` //新密码
}

// 获取用户信息及可访问的权限菜单
// @Title 用户修改密码
// @Summary 用户修改密码
// @Tags user   用户操作
// @Accept  json
// @Produce  json
// @Param   body    body    handler.editPwdReq    true     "修改密码"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user/edit_pwd [post]
func (UserController) EditPwd(c *gin.Context) {
	// 用户ID
	uid := c.GetString(USER_UID_KEY)
	reqData := editPwdReq{}
	err := c.ShouldBind(&reqData)
	if err != nil {
		log.Warn(err)
		resBadRequest(c, err)
		return
	}
	user := model.AdminUser{}
	notFound, err := orm.FirstByID(&user, uid)
	if err != nil {
		if notFound {
			resBusinessP(c, "没有此条记录")
			return
		}
		log.Error(err)
		resErrSrv(c)
		return
	}
	if toolfunc.EncUserPwd(reqData.OldPassword, user.Salt) != user.Password {
		resBusinessP(c, "原密码输入不正确")
		return
	}
	new_password := toolfunc.EncUserPwd(reqData.NewPassword, user.Salt)
	NewAdminUser := model.AdminUser{Password: new_password}
	err = orm.Updates(user, NewAdminUser)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}

// 递归菜单
func setMenu(menus []model.Menu, parentID string) (out []MenuModel) {
	var menuArr []model.Menu
	linq.From(menus).Where(func(c interface{}) bool {
		return c.(model.Menu).ParentID == parentID
	}).OrderBy(func(c interface{}) interface{} {
		return c.(model.Menu).Sequence
	}).ToSlice(&menuArr)
	if len(menuArr) == 0 {
		return
	}
	for _, v := range menuArr {
		menu := MenuModel{
			Path:      v.URL,
			Component: v.Code,
			Name:      v.Name,
			Children:  []MenuModel{}}
		if v.MenuType == 3 {
			menu.Hidden = true
		}
		//查询是否有子级
		menuChildren := setMenu(menus, v.ID)
		if len(menuChildren) > 0 {
			menu.Children = menuChildren
		}
		/*		if v.MenuType == 2 {
					// 添加子级首页
					menuIndex := MenuModel{
						Path:      "index",
						Component: v.Code,
						Name:      v.Name,
						Children:  []MenuModel{}}
					menu.Children = append(menu.Children, menuIndex)
					menu.Name = menu.Name + "index"
				}*/
		out = append(out, menu)
	}
	return
}
