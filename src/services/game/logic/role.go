package logic

import (
	"ThreeKingdoms/src/constant"
	"ThreeKingdoms/src/gamedatabase"
	"ThreeKingdoms/src/net"
	"ThreeKingdoms/src/services/common"
	"ThreeKingdoms/src/services/game/gameconfig"
	"ThreeKingdoms/src/services/game/model"
	"ThreeKingdoms/src/services/game/model/data"
	"ThreeKingdoms/src/utils"
	"fmt"
	"log"
	"time"
)

var DefaultRoleService = &RoleService{}

type RoleService struct {
}

func (r *RoleService) EnterServer(uid int, enterServerResponsePacket *model.EnterServerRspPacket, request *net.WsMsgRequestStruct) error {

	role := &data.Role{}

	ok, err := gamedatabase.Engine.Table(role).Where("uid = ?", uid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return common.New(constant.DBError, "查询角色出错")
	}

	//找到了对应角色则说明 玩家之前注册过角色
	if ok {

		//查询游戏资源
		rid := role.RId
		roleRes := &data.RoleResDo{}
		ok, err := gamedatabase.Engine.Table(roleRes).Where("rid = ?", rid).Get(roleRes)
		if err != nil {
			//资源出错
			log.Println("查询角色资源出错", err)
			return common.New(constant.DBError, "查询角色资源出错")
		}
		if !ok {
			//资源不存在 则需要新建
			roleRes := CreateNewRoleResource(rid)

			_, err := gamedatabase.Engine.Table(roleRes).Insert(roleRes)
			if err != nil {
				log.Println("创建角色资源失败", err)
				return common.New(constant.DBError, "创建角色资源失败")
			}

			enterServerResponsePacket.RoleRes = roleRes.ToModel().(model.RoleRes)
			enterServerResponsePacket.Role = role.ToModel().(model.Role)
			enterServerResponsePacket.Time = time.Now().UnixMilli()
			token, _ := utils.Award(rid)

			enterServerResponsePacket.Token = token
		}

	} else {
		//说明玩家还未创建角色自然也未创建资源
		role = CreateNewRole(uid)
		if role == nil {
			log.Println("创建新角色出错", err)
			return common.New(constant.DBError, "创建新角色出错")
		}
		resource := CreateNewRoleResource(role.RId)
		if resource == nil {
			log.Println("创建新角色的资源出错", err)
			return common.New(constant.DBError, "创建新角色的资源出错")
		}
		enterServerResponsePacket.Role = role.ToModel().(model.Role)
		enterServerResponsePacket.RoleRes = resource.ToModel().(model.RoleRes)
		enterServerResponsePacket.Time = time.Now().UnixMilli()
		token, _ := utils.Award(role.RId)
		enterServerResponsePacket.Token = token
	}
	return nil
}

func CreateNewRoleResource(rid int) *data.RoleResDo {

	roleRes := &data.RoleResDo{}
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
		return nil
	}
	return roleRes
}
func CreateNewRole(uid int) *data.Role {

	roleDo := &data.Role{}

	//角色不存在 创建角色
	roleDo.UId = uid
	roleDo.HeadId = 1
	roleDo.CreatedAt = time.Now()
	ok, err := gamedatabase.Engine.Table(roleDo).Insert(roleDo)
	if err != nil {
		log.Println("创建角色失败", err)
		return nil
	}
	if ok > 0 {
		log.Println("创建角色失败", err)
		return nil
	}
	return roleDo

}
