package net

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type wsClientChannelStruct struct {
	wsConnection *websocket.Conn
	//是否握手
	handshake        bool
	handshakeChannel chan bool
}

func (p *wsClientChannelStruct) Start() bool {

	//做的事情 就是一直不停的接受消息
	//等待握手消息返回
	p.handshake = false
	go p.wsReadLoop()
	return p.waitHandShake()
}
func (p *wsClientChannelStruct) waitHandShake() bool {

	//握手超时 应该关闭
	// 设置超时时间，例如 5 秒
	const timeout = 5 * time.Second
	select {
	case _ = <-p.handshakeChannel:
		log.Println("握手成功啦")
		return true
	case <-time.After(timeout):
		log.Println("握手超时，关闭连接")
		// 这里可以添加关闭连接的逻辑，例如：
		// p.closeConnection() // 假设有这个方法
		return false
	}

	//等待握手成功  等待握手的消息
}

func (p *wsClientChannelStruct) wsReadLoop() {
	for {
		_, data, err := p.wsConnection.ReadMessage()
		fmt.Println(string(data))
		fmt.Println(err)
		//假设收到握手消息 则通知
		p.handshake = true
		p.handshakeChannel <- true
	}

}

func NewWsClientChannel(wsConnection *websocket.Conn) *wsClientChannelStruct {
	return &wsClientChannelStruct{
		wsConnection:     wsConnection,
		handshake:        false,
		handshakeChannel: make(chan bool),
	}

}
