package main

import (
	"context"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/svc"
)

func monitorTree() {

	// You can use a LeveledList here, for easy generation.
	leveledList := pterm.LeveledList{}

	serviceCtx.Manager.Each(func(woker svc.Woker) {
		res, err := woker.Client.ListPlan(context.Background(), &benchclient.ListPlanRequest{})

		leveledList = append(leveledList, pterm.LeveledListItem{Level: 0, Text: woker.Address})

		if err != nil {
			leveledList = append(leveledList, pterm.LeveledListItem{Level: 1, Text: err.Error()})

		} else {

			for i := 0; i < len(res.Plans); i++ {
				leveledList = append(leveledList, pterm.LeveledListItem{Level: 1, Text: res.Plans[i].Uuid})
				leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("status: %s", res.Plans[i].Status)})
				leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("total student number: %d", res.Plans[i].Plan.RoomNumber*res.Plans[i].Plan.RoomSize)})

				leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("send: %d", res.Plans[i].Total.SendCount)})
				leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("custom send: %d", res.Plans[i].Total.CustomSendCount)})

				leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("gate connector : %d", res.Plans[i].Connector.GateConnector)})
				if res.Plans[i].Connector.ChatConnector != res.Plans[i].Plan.RoomNumber*res.Plans[i].Plan.RoomSize {
					leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("[WARNING] online: %d. should be %d", res.Plans[i].Connector.ChatConnector,
						res.Plans[i].Plan.RoomNumber*res.Plans[i].Plan.RoomSize)})
				} else {
					leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("[OK] online: %d", res.Plans[i].Connector.ChatConnector)})
				}

				if res.Plans[i].Total.OnChatReceiveCount != res.Plans[i].Plan.RoomSize*res.Plans[i].Total.SendCount {
					leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("[WARNING] on event. on chat receive should be %d", res.Plans[i].Plan.RoomSize*res.Plans[i].Total.SendCount)})
				} else {
					leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: "[OK] on event"})
				}

				leveledList = append(leveledList, pterm.LeveledListItem{Level: 3, Text: fmt.Sprintf("detail: on chat: %d | on server: %d | on add: %d | on leave: %d",
					res.Plans[i].Total.OnChatReceiveCount, res.Plans[i].Total.OnServerReceiveCount, res.Plans[i].Total.OnAddReceiveCount, res.Plans[i].Total.OnLeaveReceiveCount)})

				leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("address: %s | room number: %d | room id pre: %s | room size: %d | channel id: %d",
					res.Plans[i].Plan.Address, res.Plans[i].Plan.RoomNumber, res.Plans[i].Plan.RoomIdPre, res.Plans[i].Plan.RoomSize, res.Plans[i].Plan.ChannelId)})

				// bench 指标统计
				leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: statString(res.Plans[i].Stat)})

			}
		}
	})

	// Generate tree from LeveledList.
	root := putils.TreeFromLeveledList(leveledList)
	root.Text = "monitor_tree"

	// Render TreePrinter
	_ = pterm.DefaultTree.WithRoot(root).Render()
}

func statString(s *benchclient.Metrics) string {
	return fmt.Sprintf("drops: %d | average: %.2f | median: %.2f | top90th: %.2f | top99th: %.2f | top99p9th: %.2f ",
		s.Drops,
		s.Average,
		s.Median,
		s.Top90Th,
		s.Top99Th,
		s.Top99P9Th)

}
