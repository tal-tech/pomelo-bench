package main

import (
	"context"
	"github.com/pterm/pterm"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/svc"
)

func sendRecover() {

	options := []string{
		"cancle",
	}

	for router, _ := range serviceCtx.RecoverDataPool {
		options = append(options, router)
	}

	printer := pterm.DefaultInteractiveSelect.WithOptions(options)

	selectedOption, _ := printer.Show()

	data, ok := serviceCtx.RecoverDataPool[selectedOption]
	if !ok {
		return
	}

	// Create and start a fork of the default spinner.
	spinner, _ := pterm.DefaultSpinner.Start("recovering ...")

	serviceCtx.Manager.EachAsync(func(woker svc.Woker) {

		pterm.Info.Println(woker.Address, "recovering ...")

		_, err := woker.Client.CustomSend(context.Background(), &benchclient.CustomSendRequest{
			Uuid: nil, // nil 代表所有的任务都发送消息
			Pool: &benchclient.CustomMessagePool{
				Router: data.Router,
				Data:   data.Data, // 这里发送的任务池 会被一个任务下的所有connector平分
			},
			Number:   1, // 发一次得了
			Limit:    0, // 默认限制一个发送量
			Duration: 0, // 恢复需要控制间隔吗？
		})

		if err != nil {
			pterm.Error.Println(woker.Address, "recovering failed. ", err.Error())
		} else {
			pterm.Success.Println(woker.Address, "recovering success.")
		}
	})

	spinner.Success("all recovering over!")
}
