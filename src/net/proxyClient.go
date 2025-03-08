package net

import (
	"errors"
	"github.com/gorilla/websocket"
	"time"
)

type ProxyClientStruct struct {
	//代理地址
	proxy   string
	channel *wsClientChannelStruct
}

func CreateNewProxyClient(proxy string) *ProxyClientStruct {

	return &ProxyClientStruct{proxy: proxy}
}

func (p *ProxyClientStruct) Connect() error {
	//去连接 websocket服务端
	//通过Dialer连接websocket服务器
	var dialer = websocket.Dialer{
		Subprotocols:     []string{"p1", "p2"},
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}
	ws, _, err := dialer.Dial(p.proxy, nil)
	if err == nil {
		p.channel = NewWsClientChannel(ws)
		if !p.channel.Start() {
			return errors.New("握手失败")
		}
	}
	return err
}
