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

	// 并发量
	concurrency, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("1").Show("Enter send concurrency (advise 1~100)")

	iconcurrency, err := strconv.Atoi(concurrency)
	if err != nil {
		pterm.Error.Println("invalid send concurrency. must int value")
		return
	}

	customMessagePool := analysisCustomMessagePool(inumber, iconcurrency, data)

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

func analysisCustomMessagePool(number int, concurrency int, data svc.RecoverData) (res *benchclient.CustomSendRequest) {

	var (
		dataLen = concurrency * number // 消息总量 = 消息量 * 并发量
	)

	if len(data.Data) >= dataLen { // 消息量够用

		return &benchclient.CustomSendRequest{
			Uuid: nil, // nil 代表所有的任务都发送消息
			Pool: &benchclient.CustomMessagePool{
				Router: data.Router,
				Data:   data.Data[:dataLen], // 这里发送的任务池 会被一个任务下的所有connector平分
			},
			Number:   1,                   // 发一次得了
			Limit:    uint64(concurrency), // 默认限制一个发送量
			Duration: 1000,                // 恢复需要控制间隔吗？ 1000ms
		}
	}

	factor, number2 := findMaxNumber(dataLen, len(data.Data))

	// 消息量不够用 需要靠 number 找补
	return &benchclient.CustomSendRequest{
		Uuid: nil, // nil 代表所有的任务都发送消息
		Pool: &benchclient.CustomMessagePool{
			Router: data.Router,
			Data:   data.Data[:factor], // 这里发送的任务池 会被一个任务下的所有connector平分
		},
		Number:   uint64(number2),     // 发一次得了
		Limit:    uint64(concurrency), // 默认限制一个发送量
		Duration: 1000,                // 恢复需要控制间隔吗？ 1000ms
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
