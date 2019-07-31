package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/model"
	"joe-micro/lib/log"
	"joe-micro/lib/orm"
	"joe-micro/lib/toolfunc"
)

type UserManagementController struct{} //用户管理控制器

// 获取admin用户列表
// @Summary 获取admin用户列表
// @Tags   user_mana   用户管理模块
// @Tag.description  用户管理模块
// @Accept  json
// @Produce  json
// @Param   page         query    int       false     "页码,默认为1"
// @Param   num          query    int       false     "返回条数,默认为10"
// @Param   sort         query    string    false     "排序字段,默认为createAt"
// @Param   key          query    string    false     "搜索关键字"
// @Param   orderType    query    string    false     "排序规则,默认为DESC"
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
		resBadRequest(c, err)
		return
	}
	var whereOrder []orm.PageWhere
	if reqData.Key != "" {
		v := "%" + reqData.Key + "%"
		var arr []interface{}
		arr = append(arr, v)
		arr = append(arr, v)
		whereOrder = append(whereOrder, orm.PageWhere{Where: "user_name like ? or real_name like ?", Value: arr})
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
		resBadRequest(c, errors.New("缺少参数"))
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

type ids struct {
	Ids []string `json:"ids"  validate:"required" ` //id 列表
}

// 删除admin用户
// @Summary 删除admin用户
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   body    body    handler.ids    true     "id 列表"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/delete [post]
func (UserManagementController) Delete(c *gin.Context) {
	var uids ids
	err := c.ShouldBind(&uids)
	if err != nil || len(uids.Ids) == 0 {
		resBadRequest(c, err)
		return
	}
	admin := model.AdminUser{}
	err = admin.Delete(uids.Ids) //删除用户 和 其他的一些关联数据
	if err != nil {
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}

type adminUserReq struct {
	UserName string `gorm:"column:username;size:32;not null;" json:"username" form:"user_name"` // 用户名
	Password string `gorm:"column:password;not null;" json:"password" form:"password"`          // 密码
	RealName string `gorm:"column:real_name;size:32;" json:"real_name" form:"real_name"`        // 真实姓名
	Email    string `gorm:"column:email;size:64;" json:"email" form:"email"`                    // 邮箱
	Phone    string `gorm:"column:phone;type:char(20);" json:"phone" form:"phone"`              // 手机号
	Status   uint64 `gorm:"column:status" json:"status" form:"status" binding:"max=2"`          // 状态(1:启用  2.禁用)
}

// 更新admin用户
// @Summary 更新admin用户
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   uid     query   string         true         "用户uid"
// @Param   body    body    handler.adminUserReq    true     "用户资料"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/update [post]
func (UserManagementController) Update(c *gin.Context) {
	uid, exist := c.GetQuery("uid")
	if !exist {
		resBadRequest(c, errors.New("uid不能为空"))
		return
	}
	adminUser := adminUserReq{}
	err := c.ShouldBind(&adminUser)
	if err != nil {
		resBadRequest(c, err)
		return
	}
	where := model.AdminUser{}
	where.ID = uid
	oldAdminUser := model.AdminUser{}
	notFound, err := orm.First(&where, &oldAdminUser)
	if err != nil {
		if notFound {
			resBusinessP(c, "没有此条记录")
			return
		}
		log.Warn(err)
		resErrSrv(c)
		return
	}
	if adminUser.Password != "" {
		//更新密码
		adminUser.Password = toolfunc.EncUserPwd(adminUser.Password, oldAdminUser.Salt)
	}
	err = orm.Updates(where, adminUser)
	if err != nil {
		log.Warn(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}

// 创建admin用户
// @Summary 创建admin用户
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   body    body    handler.adminUserReq    true     "用户资料"
// @Success 200 {object}  handler.ResponseModel 	"{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/create [post]
func (UserManagementController) Create(c *gin.Context) {
	reqData := adminUserReq{}
	err := c.ShouldBind(&reqData)
	if err != nil {
		resBadRequest(c, err)
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
		log.Warn(err)
		resErrSrv(c)
		return
	}
	resData := map[string]interface{}{"id": adminUser.ID}
	resSuccess(c, resData)
}
