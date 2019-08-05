package handler

import (
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/model"
	"joe-micro/lib/log"
	"joe-micro/lib/orm"
)

type MenuController struct{} //菜单管理控制器

// 获取菜单列表
// @Summary 获取菜单列表
// @Tags   menu   菜单管理
// @Accept  json
// @Produce  json
// @Param   page         query    int       false     "页码,默认为1"
// @Param   num          query    int       false     "返回条数,默认为10"
// @Param   sort         query    string    false     "排序字段,默认为createdAt"
// @Param   key          query    string    false     "搜索关键字"
// @Param   orderType    query    string    false     "排序规则,默认为DESC"
// @Param   beginAt      query    string    false     "开始时间"
// @Param   endAt        query    string    false     "结束时间"
// @Success 200 {array}   model.Menu 	           "菜单列表"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /menu/list [get]
func (MenuController) List(c *gin.Context) {
	reqData := ListReq{}
	err := reqData.getListQuery(c)
	if err != nil {
		log.Warn(err)
		resBadRequest(c, err.Error())
		return
	}
	parent_id := c.Query("parent_id")
	var whereOrder []orm.PageWhere
	if reqData.Key != "" {
		v := "%" + reqData.Key + "%"
		var arr []interface{}
		arr = append(arr, v)
		arr = append(arr, v)
		whereOrder = append(whereOrder, orm.PageWhere{Where: "name like ? or code like ?", Value: arr})
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
	if parent_id != "" {
		var arr []interface{}
		arr = append(arr, parent_id)
		whereOrder = append(whereOrder, orm.PageWhere{Where: "parent_id = ?", Value: arr})
	}
	list := make([]model.Menu, 0)
	var indexPage orm.IndexPage
	indexPage.Page = reqData.Page
	indexPage.Num = reqData.Num
	err = orm.GetPage(&model.Menu{}, &model.Menu{}, &list, &indexPage, reqData.Sort, whereOrder...)
	if err != nil {
		log.Warn(err)
		resErrSrv(c)
		return
	}
	resSuccessPage(c, indexPage, list)
}

// 详情
func (MenuController) Detail(c *gin.Context) {
	menuID, exist := c.GetQuery("id")
	if !exist {
		resBadRequest(c, "id不能为空")
		return
	}
	var menu model.Menu
	where := model.Menu{}
	where.ID = menuID
	_, err := orm.First(&where, &menu)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccess(c, menu)
}

type menuUpdateReq struct {
	ID          string `gorm:"column:id;primary_key"  json:"id" binding:"required"` //id
	Name        string `gorm:"column:name;" json:"name" `                           // 菜单名称
	Status      uint8  `gorm:"column:status;" json:"status" `                       // 状态(1:启用 2:不启用)
	ParentID    string `gorm:"column:parent_id;" json:"parent_id" `                 // 父级ID
	URL         string `gorm:"column:url;" json:"url" `                             // 菜单URL
	Sequence    uint64 `gorm:"column:sequence;" json:"sequence" `                   // 排序值
	MenuType    uint8  `gorm:"column:menu_type;" json:"menu_type" `                 // 菜单类型 1模块2菜单3操作
	Code        string `gorm:"column:code;size:32;" json:"code" `                   // 菜单代码
	OperateType string `gorm:"column:operate_type;" json:"operate_type" `           // 操作类型 read/write
}

// 更新菜单
// @Summary 更新菜单
// @Tags   menu   菜单管理
// @Accept  json
// @Produce  json
// @Param   body         body    handler.menuUpdateReq       false     "菜单信息"
// @Success 200 {object}   handler.ResponseModel 	"返回ok"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /menu/update [post]
func (MenuController) Update(c *gin.Context) {
	reqData := menuUpdateReq{}
	err := c.Bind(&reqData)
	if err != nil {
		resBadRequest(c, err.Error())
		return
	}
	where := model.Menu{}
	where.ID = reqData.ID
	err = orm.Updates(&where, &reqData)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}

type menuCreateReq struct {
	Name        string `gorm:"column:name;" json:"name" binding:"required"`                   // 菜单名称
	Status      uint8  `gorm:"column:status;" json:"status" binding:"required|int=1,2"`       // 状态(1:启用 2:不启用)
	ParentID    string `gorm:"column:parent_id;" json:"parent_id" `                           // 父级ID
	URL         string `gorm:"column:url;" json:"url" `                                       // 菜单URL
	Sequence    uint64 `gorm:"column:sequence;" json:"sequence" `                             // 排序值
	MenuType    uint8  `gorm:"column:menu_type;" json:"menu_type" binding:"required|int=1,3"` // 菜单类型 1模块2菜单3操作
	Code        string `gorm:"column:code;size:32;" json:"code" `                             // 菜单代码
	OperateType string `gorm:"column:operate_type;" json:"operate_type" binding:"required"`   // 操作类型 read/write
}

// 创建菜单
// @Summary 创建菜单
// @Tags   menu   菜单管理
// @Accept  json
// @Produce  json
// @Param   body         body    handler.menuCreateReq       false     "菜单信息"
// @Success 200 {object}   handler.ResponseModel 	           "返回菜单id"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /menu/create [post]
func (MenuController) Create(c *gin.Context) {
	reqData := menuCreateReq{}
	err := c.ShouldBind(&reqData)
	if err != nil {
		resBadRequest(c, err.Error())
		return
	}
	var menu model.Menu
	menu.Name = reqData.Name
	menu.Status = reqData.Status
	menu.ParentID = reqData.ParentID
	menu.URL = reqData.URL
	menu.Sequence = reqData.Sequence
	menu.MenuType = reqData.MenuType
	menu.Code = reqData.Code
	menu.OperateType = reqData.OperateType
	err = orm.Create(&menu)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	go InitMenu(menu)
	resSuccess(c, gin.H{"id": menu.ID})
}

// 删除菜单
// @Summary 删除菜单
// @Tags   menu   菜单管理
// @Accept  json
// @Produce  json
// @Param   body         body    handler.idsReq       false     "菜单id列表"
// @Success 200 {object}   handler.ResponseModel 	           "返回成功"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /menu/delete [post]
func (MenuController) Delete(c *gin.Context) {
	var ids idsReq
	err := c.ShouldBind(&ids)
	if err != nil {
		resBadRequest(c, err.Error())
		return
	}
	menu := model.Menu{}
	//应该不让删除
	err = menu.Delete(ids.Ids) //删除菜单数据
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}

// 新增菜单后自动添加菜单下的常规操作
func InitMenu(menu model.Menu) {
	if menu.MenuType != 2 {
		return
	}
	read := model.Menu{Status: 1, ParentID: menu.ID, URL: menu.URL + "/*", Name: "读", Sequence: 1, MenuType: 3, Code: menu.Code + "read", OperateType: "read"}
	err := orm.Create(&read)
	if err != nil {
		log.Error(err)
		return
	}
	write := model.Menu{Status: 1, ParentID: menu.ID, URL: menu.URL + "/*", Name: "写", Sequence: 2, MenuType: 3, Code: menu.Code + "write", OperateType: "write"}
	err = orm.Create(&write)
	if err != nil {
		log.Error(err)
		return
	}
}

// 获取一个用户的菜单有权限的操作列表
// @Summary 获取一个用户的菜单有权限的操作列表
// @Tags   menu   菜单管理
// @Accept  json
// @Produce  json
// @Param   uid      query    string       true     "用户uid"
// @Param   menuCode      query    string       true     "菜单code"
// @Success 200 {object}   handler.ResponseModel 	           "返回权限列表string数组"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /menu/menubuttonlist [get]
func (MenuController) MenuButtonList(c *gin.Context) {
	// 用户ID
	uid, exist := c.GetQuery("uid")
	if !exist {
		resBadRequest(c, "缺少uid参数")
		return
	}
	menuCode, exist := c.GetQuery("menuCode")
	if !exist {
		resBadRequest(c, "缺少菜单码")
		return
	}
	btnList := []string{}
	if uid == SUPER_ADMIN_ID {
		//管理员
		btnList = append(btnList, "read")
		btnList = append(btnList, "write")
	} else {
		menu := model.Menu{}
		err := menu.GetMenuButton(uid, menuCode, &btnList)
		if err != nil {
			log.Error(err)
			resErrSrv(c)
			return
		}
	}
	resSuccess(c, &btnList)
}
