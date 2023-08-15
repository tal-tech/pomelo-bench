package lcpomelo

import (
	"context"
	"sync/atomic"
)

// Detail 获取运行详情
func (c *ClientConnector) Detail(ctx context.Context) ClientDetail {

	return ClientDetail{
		Uid:       c.uid,
		ChannelId: c.channelId,
		RoomId:    c.roomId,
		PomeloGate: ConnectorDetail{
			ConnectorConnected: c.pomeloGateConnectorConnected,
			Address:            c.pomeloGateAddress,
			ReqId:              atomic.LoadUint64(c.gateReqId),
		},
		PomeloChat: ConnectorDetail{
			ConnectorConnected: c.pomeloChatConnectorConnected,
			Address:            c.pomeloChatAddress,
			ReqId:              atomic.LoadUint64(c.chatReqId),
		},
		Statistics: Statistics{
			SendCount:            atomic.LoadUint64(c.sendCount),
			CustomSendCount:      atomic.LoadUint64(c.customSendCount),
			OnlineNum:            c.onlineNum,
			OnServerReceiveCount: atomic.LoadUint64(c.onServerReceiveCount),
			OnAddReceiveCount:    atomic.LoadUint64(c.onAddReceiveCount),
			OnLeaveReceiveCount:  atomic.LoadUint64(c.onLeaveReceiveCount),
			OnChatReceiveCount:   atomic.LoadUint64(c.onChatReceiveCount),
		},
	}
}

type ClientDetail struct {
	Uid       int
	ChannelId int
	RoomId    string

	PomeloGate ConnectorDetail
	PomeloChat ConnectorDetail

	Statistics Statistics
}

// Statistics 统计信息
type Statistics struct {
	// 发送量
	SendCount uint64
	// 自定义发送量
	CustomSendCount uint64
	// 总在线人数
	OnlineNum uint64
	// 总接收量
	OnServerReceiveCount uint64
	OnAddReceiveCount    uint64
	OnLeaveReceiveCount  uint64
	OnChatReceiveCount   uint64
}

func (s *Statistics) Add(b Statistics) {
	s.SendCount += b.SendCount
	s.CustomSendCount += b.CustomSendCount
	s.OnlineNum += b.OnlineNum
	s.OnServerReceiveCount += b.OnServerReceiveCount
	s.OnAddReceiveCount += b.OnAddReceiveCount
	s.OnLeaveReceiveCount += b.OnLeaveReceiveCount
	s.OnChatReceiveCount += b.OnChatReceiveCount
}

type ConnectorDetail struct {
	ConnectorConnected bool
	Address            string
	ReqId              uint64
}
