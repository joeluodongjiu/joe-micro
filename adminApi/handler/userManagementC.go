package handler

import (
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/model"
	"joe-micro/lib/log"
	"joe-micro/lib/orm"
)

type UserManagementController struct{} //用户管理控制器

// 获取admin用户列表
// @Summary 获取admin用户列表
// @Tags   user_mana   用户管理模块
// @Accept  json
// @Produce  json
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
