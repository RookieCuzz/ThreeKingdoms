package data

import (
	"ThreeKingdoms/src/gamedatabase"
	"ThreeKingdoms/src/services/game/model"
	"log"
)

var RoleResDao = &roleResDao{
	rrChan: make(chan *RoleResDo, 100),
}

type roleResDao struct {
	rrChan chan *RoleResDo
}

func (r *roleResDao) running() {
	for {
		select {
		case rr := <-r.rrChan:
			_, err := gamedatabase.Engine.
				Table(new(RoleResDo)).
				ID(rr.Id).
				Cols("wood", "iron", "stone", "grain", "gold").
				Update(rr)
			if err != nil {
				log.Println("RoleResDao update error", err)
			}
		}
	}
}

func init() {
	go RoleResDao.running()
}

type Yield struct {
	Wood  int
	Iron  int
	Stone int
	Grain int
	Gold  int
}

type RoleResDo struct {
	Id     int `xorm:"id pk autoincr"`
	RId    int `xorm:"rid"`
	Wood   int `xorm:"wood"`
	Iron   int `xorm:"iron"`
	Stone  int `xorm:"stone"`
	Grain  int `xorm:"grain"`
	Gold   int `xorm:"gold"`
	Decree int `xorm:"decree"` //令牌
}

func (r *RoleResDo) TableName() string {
	return "role_res"
}

func (r *RoleResDo) ToModel() interface{} {
	p := model.RoleRes{}
	p.Gold = r.Gold
	p.Grain = r.Grain
	p.Stone = r.Stone
	p.Iron = r.Iron
	p.Wood = r.Wood
	p.Decree = r.Decree

	//yield := GetYield(r.RId)
	//p.GoldYield = yield.Gold
	//p.GrainYield = yield.Grain
	//p.StoneYield = yield.Stone
	//p.IronYield = yield.Iron
	//p.WoodYield = yield.Wood

	p.GoldYield = 100
	p.GrainYield = 100
	p.StoneYield = 100
	p.IronYield = 100
	p.WoodYield = 100
	p.DepotCapacity = 10000
	return p
}

func (r *RoleResDo) SyncExecute() {
	RoleResDao.rrChan <- r
	//r.Push()
}

/* 推送同步 begin */
func (r *RoleResDo) IsCellView() bool {
	return false
}

func (r *RoleResDo) IsCanView(rid, x, y int) bool {
	return false
}

func (r *RoleResDo) BelongToRId() []int {
	return []int{r.RId}
}

func (r *RoleResDo) PushMsgName() string {
	return "roleRes.push"
}

func (r *RoleResDo) Position() (int, int) {
	return -1, -1
}

func (r *RoleResDo) TPosition() (int, int) {
	return -1, -1
}

//func (r *RoleResDo) Push() {
//	net.Mgr.Push(r)
//}
