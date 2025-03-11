package net

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type ProxyClientStruct struct {
	//代理地址
	proxy   string
	Channel *WsClientChannelStruct
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
		fmt.Println("代理通道创建成功")
		p.Channel = NewWsClientChannel(ws)
		if !p.Channel.Start() {
			return errors.New("握手失败")
		}
	}

	return err
}

func (p *ProxyClientStruct) Send(name string, content interface{}) (*ResponseStruct, error) {
	fmt.Println(name, content)
	if p.Channel != nil {
		return p.Channel.Send(name, content), nil

	}
	return nil, errors.New("未找到连接")

}
