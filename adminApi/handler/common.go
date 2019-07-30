package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SUCCESS_CODE   = 0     //成功的状态码
	BAD_REQUEST    = 1     //参数错误
	BUSINESS_ERR   = 2     //业务错误
	PERMISSION_ERR = 4     //权限错误
	FAIL_CODE      = -1    //失败的状态码
	USER_UID_KEY   = "UID" //页面UUID键名
	SUPER_ADMIN_ID = "1"     //超级管理员
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
func resBadRequest(c *gin.Context, err error) {
	ret := ResponseModel{Code: BAD_REQUEST, Message: "参数绑定错误: " + err.Error()}
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
	Total uint64      `json:"total"`
	Res   interface{} `json:"items"`
}

type ResponsePage struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    ResponsePageData `json:"data"`
}

// 响应成功-分页数据
func resSuccessPage(c *gin.Context, total uint64, list interface{}) {
	ret := ResponsePage{Code: SUCCESS_CODE, Message: "ok", Data: ResponsePageData{Total: total, Res: list}}
	resJSON(c, http.StatusOK, &ret)
}
