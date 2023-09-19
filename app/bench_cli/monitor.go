package main

import (
	"context"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/svc"
	"pomelo_bench/pb/bench"
	"strings"
	"time"
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

				leveledList = append(leveledList, pterm.LeveledListItem{Level: 3, Text: totalString(res.Plans[i].Total, res.Plans[i].Plan.RoomNumber*res.Plans[i].Plan.RoomSize)})

				leveledList = append(leveledList, pterm.LeveledListItem{Level: 2, Text: planString(res.Plans[i].Plan)})

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

func totalString(total *bench.Statistics, totalStudentNumber uint64) string {

	if totalStudentNumber != 0 {

		duration := total.OnChatDuration / int64(totalStudentNumber)

		d := time.Duration(duration)

		return fmt.Sprintf("detail: on chat: %d (duration %s) | on server: %d | on add: %d | on leave: %d",
			total.OnChatReceiveCount, d.String(), total.OnServerReceiveCount, total.OnAddReceiveCount, total.OnLeaveReceiveCount)

	} else {
		return fmt.Sprintf("detail: on chat: %d (duration null) | on server: %d | on add: %d | on leave: %d",
			total.OnChatReceiveCount, total.OnServerReceiveCount, total.OnAddReceiveCount, total.OnLeaveReceiveCount)
	}

}

func planString(plan *bench.Plan) string {

	if len(plan.RoomIds) == int(plan.RoomNumber) {

		return fmt.Sprintf("address: %s | room number: %d | room size: %d | channel id: %d | room ids: %s",
			plan.Address,
			plan.RoomNumber,
			plan.RoomSize,
			plan.ChannelId,
			strings.Join(plan.RoomIds, ","))

	} else {

		return fmt.Sprintf("address: %s | room number: %d |  room size: %d | channel id: %d | room id pre: %s",
			plan.Address,
			plan.RoomNumber,
			plan.RoomSize,
			plan.ChannelId,
			pString(plan.RoomIdPre),
		)
	}

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

func pString(a *string) string {
	if a == nil {
		return ""
	}

	return *a
}
