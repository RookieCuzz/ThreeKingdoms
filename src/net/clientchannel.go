package net

import (
	"ThreeKingdoms/src/services/login/proto"
	"ThreeKingdoms/src/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"sync"
	"time"
)

type WsClientChannelStruct struct {
	wsConnection *websocket.Conn
	//是否握手
	isClosed         bool
	properties       map[string]interface{}
	propertyLock     sync.RWMutex
	Seq              int64
	handshake        bool
	handshakeChannel chan bool
	OnPush           func(conn *WsClientChannelStruct, body *ResponseStruct)
	OnClose          func(conn *WsClientChannelStruct)
	syncCtxMap       map[int64]*syncCtx
	syncMutex        sync.RWMutex
}

func NewSyncCtxX() *syncCtx {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*500)

	return &syncCtx{
		ctx:     ctx,
		cancel:  cancelFunc,
		outChan: make(chan interface{}),
	}
}

type syncCtx struct {
	//Goroutine的上下文,包含Goroutine的运行状态,环境,现场等信息
	ctx     context.Context
	cancel  context.CancelFunc
	outChan chan interface{}
}

func (s syncCtx) wait() *ResponseStruct {

	select {
	case msg := <-s.outChan:
		return msg.(*ResponseStruct)
	case <-s.ctx.Done():
		log.Println("代理服务器响应超时")
		return nil
	}
}

func (p *WsClientChannelStruct) Start() bool {

	//做的事情 就是一直不停的接受消息
	//等待握手消息返回
	p.handshake = false
	go p.wsReadLoop()
	return p.waitHandShake()
}
func (p *WsClientChannelStruct) waitHandShake() bool {

	//握手超时 应该关闭
	// 设置超时时间，例如 5 秒
	const timeout = 50 * time.Second
	select {
	case _ = <-p.handshakeChannel:
		log.Println("代理额客户端握手成功啦")
		return true
	case <-time.After(timeout):
		log.Println("握手超时，关闭连接")
		// 这里可以添加关闭连接的逻辑，例如：
		// p.closeConnection() // 假设有这个方法
		return false
	}

	//等待握手成功  等待握手的消息
}
func (p *WsClientChannelStruct) Close() {
	p.wsConnection.Close()
}
func (p *WsClientChannelStruct) wsReadLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("wsReadLoop err:", err)
			p.Close()
		}
	}()

	for {
		_, data, err := p.wsConnection.ReadMessage()
		fmt.Println("读取消息", data)
		if err != nil {
			log.Println("接受消息出现异常:", err)
		}
		//读取消息 可能会有很多的消息,需要进行路由
		message := processProxyClientMessage(data, p)
		fmt.Println(message)
		//得到
		response := &ResponseStruct{}
		err = json.Unmarshal(message, response)
		fmt.Println(response)
		if response.Name == "handshake" {

			//假设收到握手消息,则保存服务端发来的秘钥,并设置握手成功
			pack := &proto.LoginHandshakeServerPacketStruct{}
			mapstructure.Decode(response.MsgContent, pack)
			if pack.Key != "" {
				//保存秘钥
				p.SetProperty("secretKey", pack.Key)
			} else {
				p.RemoveProperty("secretKey")
			}
			p.handshake = true
			p.handshakeChannel <- true

		} else {
			fmt.Println("收到其他消息")
			p.syncMutex.Lock()
			//直接转发到玩家客户端
			ctx, ok := p.syncCtxMap[response.Seq]
			p.syncMutex.Unlock()
			if ok {
				ctx.outChan <- response
			} else {
				log.Println("Seq 未发现: %d,%s", response.Seq, response.Name)
			}
		}

	}

}

func CryptAndZip(data []byte, wsServerBean *WsClientChannelStruct) []byte {

	//加密
	//1获取秘钥
	secretKey, err := wsServerBean.GetProperty("secretKey")
	if err != nil && secretKey == nil {
		log.Fatal(err)
	}
	key := secretKey.(string)
	encrypt, err := utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	if err != nil {
		log.Fatal(err)
	}
	//2压缩
	zip, err2 := utils.Zip(encrypt)
	if err2 != nil {
		log.Println(err2)
	}
	return zip
}

func processProxyClientMessage(data []byte, wsServerBean *WsClientChannelStruct) []byte {
	var out_data []byte
	//收到消息 解析消息,前端发来的消息为json消息
	unzipData, err2 := utils.UnZip(data)
	if err2 != nil {
		log.Println("解析请求发生错误", err2)
	}
	//进行解密
	secretKey, err := wsServerBean.GetProperty("secretKey")

	if err == nil {
		if secretKey == nil {
			log.Println("未握手？or 不需要解密？")
			return unzipData

		}
		//转为字符串
		key := secretKey.(string)

		decrypt, err := utils.AesCBCDecrypt(unzipData, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			//这里是作为代理客户端,所以无法主动发起握手
			log.Println("数据格式有误", err)

		} else {
			//获取解密数据
			out_data = decrypt
		}
	}

	return out_data

}

func NewWsClientChannel(wsConnection *websocket.Conn) *WsClientChannelStruct {
	return &WsClientChannelStruct{
		wsConnection:     wsConnection,
		handshake:        false,
		handshakeChannel: make(chan bool),
		properties:       make(map[string]interface{}),
		syncCtxMap:       make(map[int64]*syncCtx),
	}

}

func (wsServerBean WsClientChannelStruct) SetProperty(key string, value interface{}) {
	//添加属性
	//上隔离写 同时不共享读
	wsServerBean.propertyLock.Lock()

	//保证方法结束回收锁
	defer wsServerBean.propertyLock.Unlock()
	//设置属性
	wsServerBean.properties[key] = value

}

func (wsServerBean WsClientChannelStruct) GetProperty(key string) (interface{}, error) {
	//获取属性 上读锁 共享读
	wsServerBean.propertyLock.RLock()
	defer wsServerBean.propertyLock.RUnlock()
	//返回数据
	return wsServerBean.properties[key], nil

}

func (wsServerBean WsClientChannelStruct) RemoveProperty(key string) {
	//删除数据
	wsServerBean.propertyLock.Lock()
	defer wsServerBean.propertyLock.Unlock()

	//1.
	delete(wsServerBean.properties, key)

}
func (wsServerBean WsClientChannelStruct) Addr() string {
	//客户端的ip:端口(temp)
	return wsServerBean.wsConnection.RemoteAddr().String()
}

func (wsServerBean WsClientChannelStruct) Push(name string, data interface{}) {

	rsp := &WsMsgResponseStruct{
		Body: &ResponseStruct{
			Name:       name,
			MsgContent: data,
			Seq:        0,
		},
	}
	fmt.Println(rsp)
	//写入缓存
	//wsServerBean. <- rsp
}

func (channel *WsClientChannelStruct) Send(name string, content interface{}) *ResponseStruct {

	//把请求发送给代理服务器,登录服务器  等待返回
	channel.Seq += 1
	seq := channel.Seq
	sc := NewSyncCtxX()
	channel.syncMutex.Lock()
	channel.syncCtxMap[seq] = sc
	channel.syncMutex.Unlock()

	//构建代理请求
	proxyRequest := &RequestStruct{
		Seq:        seq,
		Name:       name,
		MsgContent: content,
	}

	err := channel.Write(proxyRequest)
	fmt.Println("@@@", err)
	if err != nil {
		sc.cancel()
		fmt.Println("转发出错", err)
	} else {
		response := sc.wait()
		fmt.Println(response)
		return response
	}

	//proxyResponse := &WsMsgResponseStruct{
	//	Body: &ResponseStruct{
	//
	//	}
	//}
	return nil
}

func (channel *WsClientChannelStruct) Write(request *RequestStruct) error {
	data, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
	}
	zipRequest := CryptAndZip(data, channel)
	err = channel.wsConnection.WriteMessage(websocket.BinaryMessage, zipRequest)
	return err
}
