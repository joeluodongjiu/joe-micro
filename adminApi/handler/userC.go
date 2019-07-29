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
	"strconv"
	"time"
)

type User struct{}

type login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u User) Login(c *gin.Context) {
	var req login
	if err := c.Bind(&req); err != nil {
          resBadRequest(c,err)
          return
	}
	//获取用户记录
	user := model.AdminUser{}
	where := model.AdminUser{}
	where.UserName = req.Username
	notFound, err := orm.First(&where,&user)
	if err != nil {
		if  notFound {
			resBusinessP(c,"没有此条记录")
			return
		}
		log.Warn(err)
		resErrSrv(c)
		return
	}

	//校验密码
	if toolfunc.EncUserPwd(req.Password, user.Salt) != user.Password {
		resBusinessP(c,"密码错误")
		return
	}

	if user.Status != 1 {
		resBusinessP(c,"该用户被封禁")
	}
	//颁发token
	token, err := jwt.CreateToken(user.ID)
	if err != nil {
		resErrSrv(c)
		return
	}
	//casbin 处理
	err = casbin.CsbinAddRoleForUser(user.ID)
	if err != nil {
		resErrSrv(c)
		return
	}
	//token存进cache
	err=cache.Put(strconv.Itoa(user.ID), token, config.C.Jwt.TimeOut * time.Hour)
	if err != nil {
		resErrSrv(c)
		return
	}

	//返回参数
	resData := make(map[string]string)
	resData["token"] = token

	resSuccess(c,resData)
}


func(u User)Logout(c *gin.Context){
	uid := c.GetInt(USER_UID_KEY)
	//从缓存中删除uid
	err:=cache.Delete(strconv.Itoa(uid))
	if err != nil {
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
	return
}