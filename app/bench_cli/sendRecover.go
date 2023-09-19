package main

import (
	"context"
	"github.com/pterm/pterm"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/svc"
	"strconv"
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

	// 消息量
	number, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("1").Show("Enter send number (advise 1)")

	inumber, err := strconv.Atoi(number)
	if err != nil {
		pterm.Error.Println("invalid send number. must int value")
		return
	}

	limit, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("1").Show("Enter send room limit (advise 1)")

	ilimit, err := strconv.Atoi(limit)
	if err != nil {
		pterm.Error.Println("invalid send number. must int value")
		return
	}

	customMessagePool := analysisCustomMessagePool(inumber, ilimit, data)

	// Create and start a fork of the default spinner.
	spinner, _ := pterm.DefaultSpinner.Start("recovering ...")

	serviceCtx.Manager.EachAsync(func(woker svc.Woker) {

		pterm.Info.Println(woker.Address, "recovering ...")

		_, err := woker.Client.CustomSend(context.Background(), customMessagePool)

		if err != nil {
			pterm.Error.Println(woker.Address, "recovering failed. ", err.Error())
		} else {
			pterm.Success.Println(woker.Address, "recovering success.")
		}
	})

	spinner.Success("all recovering over!")
}

func analysisCustomMessagePool(number int, limit int, data svc.RecoverData) (res *benchclient.CustomSendRequest) {

	// 消息量不够用 需要靠 number 找补
	return &benchclient.CustomSendRequest{
		Uuid: nil, // nil 代表所有的任务都发送消息
		Pool: &benchclient.CustomMessagePool{
			Router: data.Router,
			Data:   data.Data, // 这里发送的任务池 会被一个任务下的所有connector平分
		},
		Number:   uint64(number), // 发一次得了
		Limit:    uint64(limit),  // 一个房间的限制发送量
		Duration: 1000,           // 恢复需要控制间隔吗？ 1000ms
	}
}

func findMaxNumber(target int, maxFactor int) (factor int, number int) {

	for i := maxFactor; i >= 0; i-- {
		if target%i == 0 { // 说明正合适

			return i, target / i

		}
	}

	return 1, target
}
