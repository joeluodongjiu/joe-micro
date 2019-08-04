package handler

import (
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/model"
	"joe-micro/lib/log"
	"joe-micro/lib/orm"
)

type MenuController struct{} //菜单管理控制器

// 分页数据
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

// 更新
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

//新增
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
	//go InitMenu(menu)
	resSuccess(c, gin.H{"id": menu.ID})
}

// 删除数据
func (MenuController) Delete(c *gin.Context) {
	var ids idsReq
	err := c.ShouldBind(&ids)
	if err != nil  {
		resBadRequest(c, err.Error())
		return
	}
	menu := model.Menu{}
	//应该不让删除
	err = menu.Delete(ids.Ids)  //删除菜单数据
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}


/*// 新增菜单后自动添加菜单下的常规操作
func InitMenu(model sys.Menu) {
	if model.MenuType != 2 {
		return
	}
	add := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/create", Name: "新增", Sequence: 1, MenuType: 3, Code: model.Code + "Add", OperateType: "add"}
	models.Create(&add)
	del := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/delete", Name: "删除", Sequence: 2, MenuType: 3, Code: model.Code + "Del", OperateType: "del"}
	models.Create(&del)
	view := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/detail", Name: "查看", Sequence: 3, MenuType: 3, Code: model.Code + "View", OperateType: "view"}
	models.Create(&view)
	update := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/update", Name: "编辑", Sequence: 4, MenuType: 3, Code: model.Code + "Update", OperateType: "update"}
	models.Create(&update)
	list := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/list", Name: "分页api", Sequence: 5, MenuType: 3, Code: model.Code + "List", OperateType: "list"}
	models.Create(&list)
}*/

// 获取一个用户的菜单有权限的操作列表
func (MenuController) MenuButtonList(c *gin.Context) {
	// 用户ID
	uid,exist := c.GetQuery("uid")
	if !exist {
		resBadRequest(c,"缺少uid参数")
		return
	}
	menuCode,exist := c.GetQuery("menuCode")
	if !exist {
		resBadRequest(c,"缺少菜单码")
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
