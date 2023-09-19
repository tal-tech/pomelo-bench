package transform

import (
	"pomelo_bench/app/bench/internal/metrics"
	"pomelo_bench/app/bench/internal/service/lcpomelo"
	"pomelo_bench/pb/bench"
)

// Connectors 转换Connector数据
func Connectors(connectors []lcpomelo.ClientDetail) []*bench.Connector {
	res := make([]*bench.Connector, 0, len(connectors))

	for i := 0; i < len(connectors); i++ {
		c := &bench.Connector{
			Uid:       int64(connectors[i].Uid),
			ChannelId: int64(connectors[i].ChannelId),
			RoomId:    connectors[i].RoomId,
			Total:     Statistics(connectors[i].Statistics),
			PomeloGate: &bench.PomeloConnector{
				Connected: connectors[i].PomeloGate.ConnectorConnected,
				Address:   connectors[i].PomeloGate.Address,
				ReqId:     connectors[i].PomeloGate.ReqId,
			},
			PomeloChat: &bench.PomeloConnector{
				Connected: connectors[i].PomeloChat.ConnectorConnected,
				Address:   connectors[i].PomeloChat.Address,
				ReqId:     connectors[i].PomeloChat.ReqId,
			},
		}

		res = append(res, c)
	}

	return res
}

func Statistics(a lcpomelo.Statistics) *bench.Statistics {
	return &bench.Statistics{
		SendCount:            a.SendCount,
		CustomSendCount:      a.CustomSendCount,
		OnServerReceiveCount: a.OnServerReceiveCount,
		OnAddReceiveCount:    a.OnAddReceiveCount,
		OnLeaveReceiveCount:  a.OnLeaveReceiveCount,
		OnChatReceiveCount:   a.OnChatReceiveCount,
		OnChatDuration:       a.OnChatDuration,
		OnlineNum:            a.OnlineNum,
	}
}

func Stat(a metrics.StatReport) *bench.Metrics {
	return &bench.Metrics{
		Drops:     uint64(a.Drops),
		Average:   a.Average,
		Median:    a.Median,
		Top90Th:   a.Top90th,
		Top99Th:   a.Top99th,
		Top99P9Th: a.Top99p9th,
	}
}
