package main

import (
	"context"
	"github.com/pterm/pterm"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/svc"
)

func clear() {

	// Create and start a fork of the default spinner.
	spinner, _ := pterm.DefaultSpinner.Start("clearing ...")

	serviceCtx.Manager.EachAsync(func(woker svc.Woker) {

		pterm.Info.Println(woker.Address, "clearing ...")

		_, err := woker.Client.ClearStatistics(context.Background(), &benchclient.ClearStatisticsRequest{
			Uuid: nil, // nil 代表所有的任务都发送消息
		})

		if err != nil {
			pterm.Error.Println(woker.Address, "send failed. ", err.Error())
		} else {
			pterm.Success.Println(woker.Address, "send success.")
		}
	})

	spinner.Success("all clear over!")

}
