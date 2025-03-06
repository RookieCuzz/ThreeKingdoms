package logic

import (
	"ThreeKingdoms/src/constant"
	"ThreeKingdoms/src/gamedatabase"
	"ThreeKingdoms/src/services/common"
	"ThreeKingdoms/src/services/models"
	"ThreeKingdoms/src/services/web/model"
	"ThreeKingdoms/src/utils"
	"log"
	"time"
)

var DefaultAccountLogicBean = &AccountLogicStruct{}

type AccountLogicStruct struct {
}

func (s AccountLogicStruct) Register(request *model.RegisterRequestStruct) error {

	//一般web 服务 错误格式自定义
	username := request.Username

	user := &models.User{}
	//查看用户是否已经注册
	ok, err := gamedatabase.Engine.Table(user).Where("username = ?", username).Get(user)
	if err != nil {
		log.Println("注册查询失败", err)
		return common.New(constant.DBError, "数据库异常")
	}

	if ok {

		return common.New(constant.UserExist, "用户已存在")
	} else {
		//注册逻辑 插入数据库
		user.Mtime = time.Now()
		user.Ctime = time.Now()
		user.Username = username
		//加密盐
		user.Passcode = utils.RandSeq(6)

		user.Passwd = utils.Password(request.Password, user.Passcode)

		user.Hardware = request.Hardware

		_, err := gamedatabase.Engine.Table(user).Insert(user)
		if err != nil {
			log.Println("注册失败", err)
			return common.New(constant.DBError, "数据库异常")
		}
		log.Println("注册成功")
		return nil
	}

}
