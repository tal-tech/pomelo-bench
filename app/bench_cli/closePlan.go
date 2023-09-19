package main

import (
	"context"
	"github.com/pterm/pterm"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/svc"
)

func closePlan() {

	result, _ := pterm.DefaultInteractiveConfirm.Show("do you want to close all plan ?")

	if !result {
		return
	}

	// Create and start a fork of the default spinner.
	spinner, _ := pterm.DefaultSpinner.Start("close ...")

	serviceCtx.Manager.EachAsync(func(woker svc.Woker) {

		pterm.Info.Println(woker.Address, "close ...")

		_, err := woker.Client.ClosePlan(context.Background(), &benchclient.ClosePlanRequest{
			Uuid: nil, // nil 代表所有的任务都发送消息
		})

		if err != nil {
			pterm.Error.Println(woker.Address, "close failed. ", err.Error())
		} else {
			pterm.Success.Println(woker.Address, "close success.")
		}
	})

	// 清理学生id和roomid
	serviceCtx.Manager.Clear()

	spinner.Success("all close over!")
}
