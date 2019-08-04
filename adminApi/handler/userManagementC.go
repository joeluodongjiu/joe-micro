package handler

import (
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/model"
	"joe-micro/adminApi/model/casbin"
	"joe-micro/lib/log"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
)

type UserManagementController struct{} //用户管理控制器

// 获取admin用户列表
// @Summary 获取admin用户列表
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   page         query    int       false     "页码,默认为1"
// @Param   num          query    int       false     "返回条数,默认为10"
// @Param   sort         query    string    false     "排序字段,默认为createdAt"
// @Param   key          query    string    false     "搜索关键字"
// @Param   orderType    query    string    false     "排序规则,默认为DESC"
// @Param   beginAt      query    string    false     "开始时间"
// @Param   endAt        query    string    false     "结束时间"
// @Success 200 {array}   model.AdminUser 	"用户列表"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/list [get]
func (UserManagementController) List(c *gin.Context) {
	reqData := ListReq{}
	err := reqData.getListQuery(c)
	if err != nil {
		log.Warn(err)
		resBadRequest(c, err.Error())
		return
	}
	var whereOrder []orm.PageWhere
	if reqData.Key != "" {
		v := "%" + reqData.Key + "%"
		var arr []interface{}
		arr = append(arr, v)
		arr = append(arr, v)
		whereOrder = append(whereOrder, orm.PageWhere{Where: " username like ? or real_name like ? ", Value: arr})
	}
	if reqData.BeginAt != "" {
		var arr []interface{}
		arr = append(arr, reqData.BeginAt)
		whereOrder = append(whereOrder, orm.PageWhere{Where: " createAt > ? ", Value: arr})
	}
	if reqData.EndAt != "" {
		var arr []interface{}
		arr = append(arr, reqData.EndAt)
		whereOrder = append(whereOrder, orm.PageWhere{Where: " createAt < ? ", Value: arr})
	}
	list := make([]model.AdminUser, 0)
	var indexPage orm.IndexPage
	indexPage.Page = reqData.Page
	indexPage.Num = reqData.Num
	err = orm.GetPage(&model.AdminUser{}, &model.AdminUser{}, &list, &indexPage, reqData.Sort, whereOrder...)
	if err != nil {
		log.Warn(err)
		resErrSrv(c)
		return
	}
	resSuccessPage(c, indexPage, list)
}

// 获取admin用户详情
// @Summary 获取admin用户详情
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   uid    query    string    true     "用户uid"
// @Success 200 {object}   model.AdminUser 	   "用户详情"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/detail [get]
func (UserManagementController) Detail(c *gin.Context) {
	uid, exist := c.GetQuery("uid")
	if !exist {
		resBadRequest(c, "uid不能为空")
		return
	}
	user := model.AdminUser{}
	where := model.AdminUser{}
	where.ID = uid
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
	resSuccess(c, user)
}



// 删除admin用户
// @Summary 删除admin用户
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   body    body    handler.idsReq    true     "id 列表"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/delete [post]
func (UserManagementController) Delete(c *gin.Context) {
	var uids idsReq
	err := c.ShouldBind(&uids)
	if err != nil  {
		resBadRequest(c, err.Error())
		return
	}
	admin := model.AdminUser{}
	err = admin.Delete(uids.Ids) //删除用户 和 其他的一些关联数据
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}

type updateAdminUserReq struct {
	ID       string `gorm:"column:id" json:"id"  binding:"required"`                       //用户uid
	UserName string `gorm:"column:username;size:32;not null;" json:"username" `            // 用户名
	Password string `gorm:"column:password;not null;" json:"password" binding:"password" ` // 密码
	RealName string `gorm:"column:real_name;size:32;" json:"real_name" `                   // 真实姓名
	Email    string `gorm:"column:email;size:64;" json:"email" `                           // 邮箱
	Phone    string `gorm:"column:phone;type:char(20);" json:"phone" `                     // 手机号
	Status   uint8  `gorm:"column:status" json:"status"  binding:"int=1,2"`                // 状态(1:启用  2.禁用)
}

// 更新admin用户
// @Summary 更新admin用户
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   uid     query   string         true         "用户uid"
// @Param   body    body    handler.updateAdminUserReq    true     "用户资料"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/update [post]
func (UserManagementController) Update(c *gin.Context) {
	reqData := updateAdminUserReq{}
	err := c.ShouldBind(&reqData)
	if err != nil {
		resBadRequest(c, err.Error())
		return
	}
	where := model.AdminUser{}
	where.ID = reqData.ID
	oldAdminUser := model.AdminUser{}
	notFound, err := orm.First(&where, &oldAdminUser)
	if err != nil {
		if notFound {
			resBusinessP(c, "没有此条记录")
			return
		}
		log.Error(err)
		resErrSrv(c)
		return
	}
	if reqData.Password != "" {
		//更新密码
		reqData.Password = toolfunc.EncUserPwd(reqData.Password, oldAdminUser.Salt)
	}
	err = orm.Updates(where, reqData)
	if err != nil {
		log.Warn(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}

type createAdminUserReq struct {
	UserName string `gorm:"column:username" json:"username" binding:"required"`          // 用户名
	Password string `gorm:"column:password" json:"password"  binding:"required" `        // 密码
	RealName string `gorm:"column:real_name" json:"real_name" `                          // 真实姓名
	Email    string `gorm:"column:email" json:"email" `                                  // 邮箱
	Phone    string `gorm:"column:phone" json:"phone" `                                  // 手机号
	Status   uint8  `gorm:"column:status" json:"status"  binding:"required,min=1,max=2"` // 状态(1:启用  2.禁用)
}

// 创建admin用户
// @Summary 创建admin用户
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   body    body    handler.createAdminUserReq    true     "用户资料"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/create [post]
func (UserManagementController) Create(c *gin.Context) {
	reqData := createAdminUserReq{}
	err := c.ShouldBind(&reqData)
	if err != nil {
		resBadRequest(c, err.Error())
		return
	}
	adminUser := model.AdminUser{}
	adminUser.Status = reqData.Status
	adminUser.UserName = reqData.UserName
	adminUser.RealName = reqData.RealName
	adminUser.Phone = reqData.Phone
	adminUser.Email = reqData.Email
	adminUser.Salt = toolfunc.GenerateSalt()
	adminUser.Password = toolfunc.EncUserPwd(reqData.Password, adminUser.Salt)
	err = orm.Create(&adminUser)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resData := map[string]interface{}{"id": adminUser.ID}
	resSuccess(c, resData)
}

// 获取用户下的角色ID列表
// @Summary 获取用户下的角色ID列表
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   uid     query   string         true         "用户uid"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/users_roleid_list [get]
func (UserManagementController) UsersRoleIDList(c *gin.Context) {
	uid, exist := c.GetQuery("uid")
	if !exist {
		resBadRequest(c, "uid不能为空")
		return
	}
	var roleList []string
	where := model.AdminUserRoles{UserID: uid}
	err := orm.PluckList(&model.AdminUserRoles{}, &where, "role_id", &roleList)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccess(c, roleList)
}

// 为用户添加角色权限
// @Summary 为用户添加角色权限
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   uid     query   string         true         "用户uid"
// @Param   body    body    handler.idsReq    true         "角色id"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/set_role [post]
func (UserManagementController) SetRole(c *gin.Context) {
	uid, exist := c.GetQuery("uid")
	if !exist {
		resBadRequest(c, "uid不能为空")
		return
	}
	var roleids idsReq
	err := c.Bind(&roleids)
	if err != nil {
		resBadRequest(c, err.Error())
		return
	}
	userRole := model.AdminUserRoles{}
	err = userRole.SetRole(uid, roleids.Ids) //添加记录
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	//给用户添加角色 添加权限
	if  err:= casbin.CasbinAddRoleForUser(uid);err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}
