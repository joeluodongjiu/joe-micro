package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"joe-micro/lib/orm"
	"net/http"
	"strconv"
)

const (
	SUCCESS_CODE   = 0     //成功的状态码
	BAD_REQUEST    = 1     //参数错误
	BUSINESS_ERR   = 2     //业务错误
	PERMISSION_ERR = 4     //权限错误
	FAIL_CODE      = -1    //失败的状态码
	USER_UID_KEY   = "UID" //页面UUID键名
	SUPER_ADMIN_ID = "1"   //超级管理员
)

type ResponseModel struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

// 响应JSON数据
func resJSON(c *gin.Context, status int, v interface{}) {
	c.JSON(status, v)
}

// 响应成功
func resSuccess(c *gin.Context, v interface{}) {
	ret := ResponseModel{Code: SUCCESS_CODE, Message: "ok", Data: v}
	resJSON(c, http.StatusOK, &ret)
}

// 响应成功
func resSuccessMsg(c *gin.Context) {
	ret := ResponseModel{Code: SUCCESS_CODE, Message: "ok"}
	resJSON(c, http.StatusOK, &ret)
}

//参数错误
func resBadRequest(c *gin.Context, msg string) {
	ret := ResponseModel{Code: BAD_REQUEST, Message: "参数绑定错误: " + msg}
	resJSON(c, http.StatusOK, &ret)
}

//业务错误
func resBusinessP(c *gin.Context, msg string) {
	ret := ResponseModel{Code: BUSINESS_ERR, Message: msg}
	resJSON(c, http.StatusOK, &ret)
}

// 响应错误-服务端故障
func resErrSrv(c *gin.Context) {
	ret := ResponseModel{Code: FAIL_CODE, Message: "服务端故障"}
	resJSON(c, http.StatusOK, &ret)
}

// 响应失败
func resFailCode(c *gin.Context, msg string, code int) {
	ret := ResponseModel{Code: code, Message: msg}
	resJSON(c, http.StatusOK, &ret)
}

// 响应错误-用户端故障
func resErrCli(c *gin.Context, err error) {
	ret := ResponseModel{Code: FAIL_CODE, Message: err.Error()}
	resJSON(c, http.StatusOK, &ret)
}

type ResponsePageData struct {
	IndexPage orm.IndexPage   `json:"index"`
	Res   interface{} `json:"res"`
}



type ResponsePage struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    ResponsePageData `json:"data"`
}

// 响应成功-分页数据
func resSuccessPage(c *gin.Context, indexPage orm.IndexPage, list interface{}) {
	ret := ResponsePage{Code: SUCCESS_CODE, Message: "ok", Data: ResponsePageData{IndexPage: indexPage, Res: list}}
	resJSON(c, http.StatusOK, &ret)
}

type ListReq struct {
	Page      uint64 `json:"page" form:"page" `          //页数
	Num       uint64 `json:"num"  form:"num" `           //数量
	Key       string `json:"key" form:"key" `            //搜索关键字
	Sort      string `json:"sort" form:"sort"`           // 排序字段
	OrderType string `json:"orderType" form:"orderType"` //排序规则
}

func (l *ListReq) getListQuery(c *gin.Context) (err error) {
	query := c.Query("page")
	if query == "" {
		query = "1"
	}
	l.Page, err = strconv.ParseUint(query, 0, 64)
	if err != nil {
		return err
	}
	query = c.Query("num")
	if query == "" {
		query = "10"
	}
	l.Num, err = strconv.ParseUint(query, 0, 64)
	if err != nil {
		return err
	}
	l.Key = c.Query("key")
	query = c.Query("orderType")
	if query == "" {
		query = "createdAt"
	}
	l.Sort = query
	query = c.Query("orderType")
	if query == "" {
		query = "DESC"
	}
	l.OrderType = query
	if l.OrderType != "ASC" && l.OrderType != "DESC" {
		return errors.New("orderType 参数错误")
	}
	l.Sort = l.Sort + "  " + l.OrderType
	return
}
