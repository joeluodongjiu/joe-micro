package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/model"
	"joe-micro/lib/log"
	"joe-micro/lib/orm"
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
	if err!=nil {
		if notFound {
			resBusinessP(c, "没有此条记录")
			return
		}
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccess(c,user)
}



type   ids struct {
	Ids   []string  `json:"ids"  validate:"required" ` //id 列表
}
// 删除admin用户
// @Summary 删除admin用户
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
// @Param   body    body    handler.ids    true     "id 列表"
// @Success 200 {object}   model.AdminUser 	        "用户详情"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/delete [post]
func  (UserManagementController) Delete(c *gin.Context){
	var uids  ids
	err := c.Bind(&uids)
	if err != nil || len(uids.Ids) == 0 {
		resBadRequest(c, err)
		return
	}
	admin:=model.AdminUser{}
	err = admin.Delete(uids.Ids)  //删除用户 和 其他的一些关联数据
	if err != nil {
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}



// 更新admin用户
// @Summary 更新admin用户
// @Tags   user_mana   用户管理模块
// @tag.description  用户管理模块
// @Accept  json
// @Produce  json
// @Param   body    body    model.AdminUser    true     "用户资料"
// @Success 200 {object}   model.AdminUser 	        "用户详情"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/update [post]
func (UserManagementController) Update(c *gin.Context) {
	adminUser := model.AdminUser{}
	err := c.Bind(&adminUser)
	if err != nil {
		resBadRequest(c,err)
		return
	}
	where := model.AdminUser{}
	where.ID = adminUser.ID
	err = orm.Updates(where,adminUser)
	if err != nil {
		log.Warn(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}