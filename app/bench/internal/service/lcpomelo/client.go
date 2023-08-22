// Package lcpomelo pomelo客户端
package lcpomelo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"pomelo_bench/pomelosdk"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// pomelo 路由
const (
	ROUTE_GATE                         = "gate.gateHandler.queryEntry"
	ROUTE_Connector_EntryHandler_Enter = "connector.entryHandler.enter"
	ROUTE_Chat_ChatHandler_Send        = "chat.chatHandler.send"

	ROUTE_ACK = "ack"
)

const (
	Event_OnServer = "onServer"
	Event_OnAdd    = "onAdd"
	Event_OnLeave  = "onLeave"
	Event_OnChat   = "onChat"
)

type ClientConnector struct {
	uid       int
	channelId int
	roomId    string

	pomeloGateConnectorConnected bool
	pomeloGateAddress            string
	pomeloGateConnector          *pomelosdk.Connector
	gateReqId                    *uint64

	pomeloChatConnectorConnected bool
	pomeloChatAddress            string
	pomeloChatConnector          *pomelosdk.Connector
	chatReqId                    *uint64
	chatAckReqId                 *uint64

	sendCount            *uint64
	customSendCount      *uint64
	onServerReceiveCount *uint64 // 接收到包的数量
	onAddReceiveCount    *uint64 // 接收到包的数量
	onLeaveReceiveCount  *uint64 // 接收到包的数量
	onChatReceiveCount   *uint64 // 接收到包的数量
	onlineNum            uint64  // 房间在线人数
}

func NewClientConnector(addr string, uid int, channelId int, roomId string) *ClientConnector {
	var (
		gateReqId       uint64
		chatReqId       uint64
		chatAckReqId    uint64
		sendCount       uint64
		customSendCount uint64

		onServerReceiveCount uint64
		onAddReceiveCount    uint64
		onLeaveReceiveCount  uint64
		onChatReceiveCount   uint64
	)

	return &ClientConnector{
		uid:       uid,
		channelId: channelId,
		roomId:    roomId,

		pomeloGateConnectorConnected: false,
		pomeloGateAddress:            addr,
		pomeloGateConnector:          pomelosdk.NewConnector(),
		gateReqId:                    &gateReqId,

		pomeloChatConnectorConnected: false,
		pomeloChatAddress:            "", // gate 请求后获得
		pomeloChatConnector:          pomelosdk.NewConnector(),
		chatReqId:                    &chatReqId,
		chatAckReqId:                 &chatAckReqId,

		sendCount:            &sendCount,
		customSendCount:      &customSendCount,
		onServerReceiveCount: &onServerReceiveCount, // 接收到包的数量
		onAddReceiveCount:    &onAddReceiveCount,    // 接收到包的数量
		onLeaveReceiveCount:  &onLeaveReceiveCount,  // 接收到包的数量
		onChatReceiveCount:   &onChatReceiveCount,   // 接收到包的数量

	}
}

// RunGateConnectorAndWaitConnect 初始化GateConnector握手
func (c *ClientConnector) RunGateConnectorAndWaitConnect(ctx context.Context, timeout time.Duration) error {
	if !c.pomeloGateConnectorConnected {

		err := c.runAndWaitConnect(ctx, c.pomeloGateConnector, c.pomeloGateAddress, timeout, nil)
		if err != nil {
			return errors.New(fmt.Sprintf("GateConnector tailed, %s", err.Error()))
		}

		c.pomeloGateConnectorConnected = true
	}

	return nil
}

// RunChatConnectorAndWaitConnect 初始化GateConnector握手
func (c *ClientConnector) RunChatConnectorAndWaitConnect(ctx context.Context, timeout time.Duration) error {
	if c.pomeloChatAddress == "" {
		return errors.New("invalid pomelo chat address")
	}

	if !c.pomeloChatConnectorConnected {

		err := c.runAndWaitConnect(ctx, c.pomeloChatConnector, c.pomeloChatAddress, timeout, func(message string) {

			c.pomeloChatConnectorConnected = false

			logx.Error("chat connector failed, err:", message)
		})
		if err != nil {
			return errors.New(fmt.Sprintf("ChatConnector tailed, %s", err.Error()))
		}

		c.pomeloChatConnectorConnected = true
	}
	return nil
}

// SyncGateRequest 通过Gate请求Chat地址 并关闭Gate
func (c *ClientConnector) SyncGateRequest(ctx context.Context) error {

	if c.pomeloGateConnectorConnected == false {
		return errors.New("GateConnector not connected")
	}

	type GateRequest struct {
		Rid       string `json:"rid"` // room id
		Uid       int    `json:"uid"`
		RType     int    `json:"rtype"`
		UType     int    `json:"utype"`
		RetryTime int64  `json:"retrytime"`
	}

	type GateResponse struct {
		Code    int    `json:"code"`
		Host    string `json:"host"`
		Port    int    `json:"port"`
		Message string `json:"message"`
	}

	var (
		gateRequest = GateRequest{
			Rid:       c.roomId,
			Uid:       c.uid,
			RType:     c.channelId,
			UType:     0,
			RetryTime: time.Now().UnixNano() / 1e6,
		}

		gateResponse GateResponse
	)

	err := c.syncRequest(ctx, c.pomeloGateConnector, 10*time.Second, c.gateReqId, ROUTE_GATE, gateRequest, &gateResponse)
	if err != nil {
		return err
	}

	if strings.Contains(gateResponse.Host, "com") {
		c.pomeloChatAddress = fmt.Sprintf("wss://%s:%d", gateResponse.Host, gateResponse.Port)
	} else {
		c.pomeloChatAddress = fmt.Sprintf("ws://%s:%d", gateResponse.Host, gateResponse.Port)
	}

	return nil
}

// CloseGate 关闭Gate
func (c *ClientConnector) CloseGate(ctx context.Context) error {
	if c.pomeloGateConnectorConnected {
		c.pomeloGateConnector.Close()
		c.pomeloGateConnectorConnected = false
	}
	return nil
}

// CloseChat 关闭Chat
func (c *ClientConnector) CloseChat(ctx context.Context) error {
	if c.pomeloChatConnectorConnected {
		c.pomeloChatConnector.Close()
		c.pomeloChatConnectorConnected = false
	}
	return nil
}

// SyncChatEnterConnectorRequest Chat 进入房间
func (c *ClientConnector) SyncChatEnterConnectorRequest(ctx context.Context) error {

	if c.pomeloChatConnectorConnected == false {
		return errors.New("ChatConnector not connected")
	}

	type ConnectorRequest struct {
		Uid          int    `json:"uid"`
		Username     int    `json:"username"`
		Rtype        int    `json:"rtype"`
		Rid          string `json:"rid"`
		Role         int    `json:"role"`
		Ulevel       int    `json:"ulevel"`
		Uname        int    `json:"uname"`
		Classid      string `json:"classid"`
		Mtcv         string `json:"mainTeacherClientVer"`
		Pv           string `json:"protocolVersion"`
		UniqId       string `json:"uniqId"`
		InteractMode int    `json:"interactMode"`
		LiveType     int    `json:"liveType"`
		Route        string `json:"route"`
		ReqId        int    `json:"reqId"`
	}

	uniqId := rand.Int()

	request := ConnectorRequest{
		Uid:      c.uid,
		Username: c.uid,
		Uname:    c.uid,

		Rtype:        c.channelId,
		Rid:          c.roomId,
		Role:         1,
		Ulevel:       1,
		Classid:      c.roomId,
		Mtcv:         "0.0.1",
		Pv:           "1.1",
		UniqId:       strconv.Itoa(uniqId),
		InteractMode: 1,
		LiveType:     1,
		Route:        ROUTE_Connector_EntryHandler_Enter,
		ReqId:        int(atomic.LoadUint64(c.chatReqId)),
	}

	c.onEvent()

	err := c.syncRequest(ctx, c.pomeloChatConnector, 30*time.Second, c.chatReqId, ROUTE_Connector_EntryHandler_Enter, request, nil)
	if err != nil {
		return err
	}

	return nil
}

// 事件监听
func (c *ClientConnector) onEvent() {

	type callbackMessage struct {
		MsgId     int    `json:"msgId"`
		OnlineNum uint64 `json:"onlineNum"`
	}

	type eventAck struct {
		Ack   int `json:"ack"`
		MsgId int `json:"msgId"`
	}

	ack := func(data []byte) {
		var cme callbackMessage
		if err := json.Unmarshal(data, &cme); err != nil {
			logx.Errorf("[%d] OnServer,data: %s", c.uid, string(data))
		}

		ack := eventAck{
			Ack:   1,
			MsgId: cme.MsgId,
		}

		requestBytes, err := json.Marshal(ack)
		if err != nil {
			return
		}

		err = c.asyncRequest(context.Background(), c.pomeloChatConnector, c.chatAckReqId, ROUTE_ACK, requestBytes, nil)
		if err != nil {
			logx.Errorf("[%d] asyncRequest failed: %s", c.uid, err)
		}

		if cme.OnlineNum != 0 {
			c.onlineNum = cme.OnlineNum
		}

	}

	c.pomeloChatConnector.On(Event_OnServer, func(data []byte) {
		atomic.AddUint64(c.onServerReceiveCount, 1)
		logx.Infof("[%d] onServer,data: %s", c.uid, string(data))

		ack(data)
	})

	c.pomeloChatConnector.On(Event_OnAdd, func(data []byte) {
		atomic.AddUint64(c.onAddReceiveCount, 1)

		logx.Infof("[%d] onAdd,data: %s", c.uid, string(data))

		ack(data)
	})

	c.pomeloChatConnector.On(Event_OnLeave, func(data []byte) {
		atomic.AddUint64(c.onLeaveReceiveCount, 1)

		logx.Infof("[%d] onLeave,data: %s", c.uid, string(data))

		ack(data)
	})

	c.pomeloChatConnector.On(Event_OnChat, func(data []byte) {
		atomic.AddUint64(c.onChatReceiveCount, 1)

		logx.Infof("[%d] onChat,data: %s", c.uid, string(data))

		ack(data)
	})

}

func (c *ClientConnector) AsyncCustomSend(ctx context.Context, route string, data []byte, cb pomelosdk.Callback) error {
	if c.pomeloChatConnectorConnected == false {
		return errors.New("ChatConnector not connected")
	}

	err := c.asyncRequest(ctx, c.pomeloChatConnector, c.chatReqId, route, data, cb)
	if err != nil {
		return err
	}

	atomic.AddUint64(c.customSendCount, 1)

	return nil
}

func (c *ClientConnector) AsyncChatSendMessage(ctx context.Context, message string, cb pomelosdk.Callback) error {
	if c.pomeloChatConnectorConnected == false {
		return errors.New("ChatConnector not connected")
	}

	type ctimeContent struct {
		Ctime string `json:"ctime"`
	}

	type chatMessage struct {
		Route   string `json:"route"`
		Target  string `json:"target"`
		Content string `json:"content"`
	}

	content, _ := jsonx.MarshalToString(ctimeContent{
		Ctime: strconv.Itoa(int(time.Now().UnixMilli())),
	})

	// {"route":"onChat","target":"*","content":"{\"ctime\":${time}}"}
	msg := chatMessage{
		Route:   "onChat",
		Target:  "*",
		Content: content,
	}

	requestBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	//err := c.syncRequest(ctx, c.pomeloChatConnector, c.chatReqId, ROUTE_Chat_ChatHandler_Send, msg, nil)
	err = c.asyncRequest(ctx, c.pomeloChatConnector, c.chatReqId, ROUTE_Chat_ChatHandler_Send, requestBytes, cb)
	//err := c.notify(ctx, c.pomeloChatConnector, c.chatReqId, ROUTE_Chat_ChatHandler_Send, request)
	if err != nil {
		return err
	}

	atomic.AddUint64(c.sendCount, 1)

	return nil
}
