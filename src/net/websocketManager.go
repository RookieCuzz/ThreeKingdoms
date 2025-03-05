package net

import "sync"

var WebSocketManager = &WebsocketManagerStruct{
	userCache: make(map[int64]WsConnection),
}

type WebsocketManagerStruct struct {
	uc        sync.RWMutex
	userCache map[int64]WsConnection
}

func (WebSocketManager *WebsocketManagerStruct) UserLogin(wsServerChannel WsConnection, userId int64, token string) {
	WebSocketManager.uc.Lock()
	defer WebSocketManager.uc.Unlock()
	//看看是不是已经在线  有人要挤号
	oldChannel, _ := WebSocketManager.userCache[userId]
	if oldChannel != nil {
		if oldChannel != wsServerChannel {
			//通知旧客户端 有人要挤他
			oldChannel.Push("robLogin", nil)
		}
	}
	WebSocketManager.userCache[userId] = wsServerChannel
}
