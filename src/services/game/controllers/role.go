package controllers

import (
	"ThreeKingdoms/src/constant"
	"ThreeKingdoms/src/gamedatabase"
	"ThreeKingdoms/src/net"
	"ThreeKingdoms/src/services/common"
	"ThreeKingdoms/src/services/game/gameconfig"
	"ThreeKingdoms/src/services/game/logic"
	"ThreeKingdoms/src/services/game/model"
	"ThreeKingdoms/src/services/game/model/data"
	"ThreeKingdoms/src/utils"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

var DefaultRoleController = &RoleControllerStruct{}

type RoleControllerStruct struct{}

func (roleController *RoleControllerStruct) SetupRouter(router *net.RouterStruct) {
	group := router.CreateGroup("role")
	group.AddEventHandler("enterServer", roleController.enterServer)
}

func (roleController *RoleControllerStruct) enterServer(request *net.WsMsgRequestStruct, response *net.WsMsgResponseStruct) {
	//进入游戏的逻辑

	//验证Session 是否合法? 合法情况下从中取出登录的用户id

	//根据用户id 去查询对应的游戏角色,如果有 就继续 没有则提示无角色
	//根据角色id 查询橘色已经拥有的资源 roleRes,没有资源则初始化资源

	enterServerRequestPacket := &model.EnterServerReqPacket{}
	enterServerResponsePacket := &model.EnterServerRspPacket{}

	err := mapstructure.Decode(request.Body.MsgContent, enterServerRequestPacket)

	response.Body.Seq = request.Body.Seq
	response.Body.Name = request.Body.Name
	if err != nil {
		response.Body.Code = constant.InvalidParam
		return
	}
	session := enterServerRequestPacket.Session
	fmt.Println(session)
	_, claims, err := utils.ParseToken(session)
	if err != nil {
		response.Body.Code = constant.SessionInvalid
		return
	}

	//获取uuid
	uid := claims.Uid

	err = logic.DefaultRoleService.EnterServer(uid, enterServerResponsePacket, request)
	if err != nil {
		response.Body.Code = err.(*common.MyError).Code()
		return
	}
	//执行对应的处理逻辑
	response.Body.Code = constant.OK
	response.Body.MsgContent = enterServerResponsePacket

	//执行进入服务逻辑

	roleDo := &data.Role{}
	ok, err := gamedatabase.Engine.Table(roleDo).Where("uid = ?", uid).Get(roleDo)
	if err != nil {
		log.Println("数据库查询出错", err)
		response.Body.Code = constant.DBError
		return
	}
	if ok {
		//角色存在
		response.Body.Code = constant.OK
		response.Body.MsgContent = enterServerResponsePacket
		rid := roleDo.RId
		roleRes := &data.RoleResDo{}
		ok, err := gamedatabase.Engine.Table(roleRes).Where("rid = ?", rid).Get(roleRes)
		if err != nil {
			//资源出错
			log.Println("查询角色资源出错", err)
			response.Body.Code = constant.DBError
			return
		}
		if !ok {
			//资源不存在 则需要新建
			roleRes.RId = rid
			roleRes.Gold = gameconfig.Basic.Role.Gold
			roleRes.Decree = gameconfig.Basic.Role.Decree
			roleRes.Iron = gameconfig.Basic.Role.Iron
			roleRes.Stone = gameconfig.Basic.Role.Stone
			roleRes.Wood = gameconfig.Basic.Role.Wood
			roleRes.Grain = gameconfig.Basic.Role.Grain
			_, err := gamedatabase.Engine.Table(roleRes).Insert(roleRes)
			if err != nil {
				log.Println("创建角色资源失败", err)
				response.Body.Code = constant.DBError
				return
			}
			enterServerResponsePacket.RoleRes = roleRes.ToModel().(model.RoleRes)
			enterServerResponsePacket.Role = roleDo.ToModel().(model.Role)
			enterServerResponsePacket.Time = time.Now().UnixMilli()
			token, _ := utils.Award(rid)

			enterServerResponsePacket.Token = token
		}

	} else {
		//角色不存在 创建角色
		roleDo.UId = uid
		roleDo.HeadId = 1
		roleDo.CreatedAt = time.Now()
		ok, err := gamedatabase.Engine.Table(roleDo).Insert(roleDo)
		if err != nil {
			log.Println("创建角色失败", err)
			response.Body.Code = constant.DBError
			return
		}

		//创建角色成功
		if ok > 0 {
			rid := roleDo.RId
			roleRes := &data.RoleResDo{}

			//给新角色创建资源
			err := CreateNewRoleResource(roleRes, rid)
			if err != nil {
				log.Println("创建新角色的资源失败", err)
				response.Body.Code = constant.DBError
				return
			}
			response.Body.Code = constant.OK

			response.Body.MsgContent = enterServerResponsePacket
			enterServerResponsePacket.RoleRes = roleRes.ToModel().(model.RoleRes)
			enterServerResponsePacket.Role = roleDo.ToModel().(model.Role)
			enterServerResponsePacket.Time = time.Now().UnixMilli()
			token, _ := utils.Award(rid)
			enterServerResponsePacket.Token = token
			fmt.Println("发送数据包")
		} else {
			log.Println("创建角色失败", err)
			response.Body.Code = constant.DBError
			return
		}

	}
	fmt.Println(response)
}
func CreateNewRoleResource(roleRes *data.RoleResDo, rid int) error {

	roleRes.RId = rid
	roleRes.Gold = gameconfig.Basic.Role.Gold
	fmt.Println(gameconfig.Basic.Role)
	roleRes.Decree = gameconfig.Basic.Role.Decree
	roleRes.Iron = gameconfig.Basic.Role.Iron
	roleRes.Stone = gameconfig.Basic.Role.Stone
	roleRes.Wood = gameconfig.Basic.Role.Wood
	roleRes.Grain = gameconfig.Basic.Role.Grain
	_, err := gamedatabase.Engine.Table(roleRes).Insert(roleRes)
	if err != nil {
		log.Println("创建角色资源失败", err)
		return err
	}
	return nil
}
