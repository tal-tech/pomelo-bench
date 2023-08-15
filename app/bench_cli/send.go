package main

import (
	"context"
	"github.com/pterm/pterm"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/svc"
	"strconv"
)

func send() {
	// 目前没用
	//message, _ := pterm.DefaultInteractiveTextInput.Show("Enter send message (advise hello)")

	number, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("1").Show("Enter send number (advise 1)")

	inumber, err := strconv.Atoi(number)
	if err != nil {
		pterm.Error.Println("invalid send number. must int value")
		return
	}

	// Create and start a fork of the default spinner.
	spinner, _ := pterm.DefaultSpinner.Start("sending ...")

	serviceCtx.Manager.EachAsync(func(woker svc.Woker) {

		pterm.Info.Println(woker.Address, "sending ...")

		_, err := woker.Client.SendChat(context.Background(), &benchclient.SendChatRequest{
			Message:  "", // 未使用
			Number:   uint64(inumber),
			Limit:    1,    // 默认限制一个发送量
			Duration: 1000, // 1000毫秒
			Uuid:     nil,  // nil 代表所有的任务都发送消息
		})

		if err != nil {
			pterm.Error.Println(woker.Address, "send failed. ", err.Error())
		} else {
			pterm.Success.Println(woker.Address, "send success.")
		}
	})

	spinner.Success("all send over!")
}
