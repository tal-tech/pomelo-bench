package main

import (
	"context"
	"github.com/pterm/pterm"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/svc"
	"strconv"
)

func works() {

	tabel := pterm.TableData{
		{
			"Address",
			"Uuid",
			"Status",
			"RoomNumber",
			"Address",
			"RoomIdPre",
			"RoomSize",
			"ChannelId",
		},
	}

	serviceCtx.Manager.Each(func(woker svc.Woker) {
		res, err := woker.Client.ListPlan(context.Background(), &benchclient.ListPlanRequest{})
		if err != nil {
			tabel = append(tabel, []string{
				woker.Address,
				err.Error(),
			})
		} else {

			for i := 0; i < len(res.Plans); i++ {
				tabel = append(tabel, []string{
					woker.Address,
					res.Plans[i].Uuid,
					res.Plans[i].Status.String(),
					strconv.FormatUint(res.Plans[i].Plan.RoomNumber, 10),
					res.Plans[i].Plan.Address,
					res.Plans[i].Plan.RoomIdPre,
					strconv.FormatUint(res.Plans[i].Plan.RoomSize, 10),
					strconv.FormatUint(res.Plans[i].Plan.ChannelId, 10),
				})
			}
		}
	})

	pterm.DefaultTable.WithHasHeader().WithData(tabel).Render()

}
