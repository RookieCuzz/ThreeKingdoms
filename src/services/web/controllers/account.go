package controllers

import (
	"ThreeKingdoms/src/constant"
	"ThreeKingdoms/src/services/common"
	"ThreeKingdoms/src/services/web/logic"
	"ThreeKingdoms/src/services/web/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var DefaultAccountController = &AccountControllerStruct{}

type AccountControllerStruct struct {
}

func (accountControllerBean *AccountControllerStruct) Register(context *gin.Context) {

	/*
		1. 获取请求的参数
		2. 根据用户名 查询数据库是否有存在该用户,若没有则继续注册
		3. 告诉前端,注册成功即可

	*/

	request := &model.RegisterRequestStruct{}
	err := context.ShouldBind(&request)
	if err != nil {
		log.Println("参数格式不合法")
		context.JSON(http.StatusOK, common.Error(constant.InvalidParam, "参数不合法"))
		return
	}
	err = logic.DefaultAccountLogicBean.Register(request)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusOK, common.Error(err.(*common.MyError).Code(), "注册业务出错"))
		return
	}
	context.JSON(http.StatusOK, common.Success(constant.OK, "注册成功"))

}
