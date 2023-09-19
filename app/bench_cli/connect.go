package main

import (
	"github.com/pterm/pterm"
	"pomelo_bench/app/bench/benchclient"
	"strconv"
)

func connect() {

	room_number, _ := pterm.DefaultInteractiveTextInput.Show("Enter student room number (advise 2)")

	iroom_number, err := strconv.Atoi(room_number)
	if err != nil {
		pterm.Error.Println("invalid room number. must int value")
		return
	}

	room_size, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("35").Show("Enter room size (advise 35)")

	iroom_size, err := strconv.Atoi(room_size)
	if err != nil {
		pterm.Error.Println("invalid room size. must int value")
		return
	}

	channel, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("2").Show("Enter channel (advise 2)")

	ichannel, err := strconv.Atoi(channel)
	if err != nil {
		pterm.Error.Println("invalid channel. must int value")
		return
	}

	pterm.Info.Println("正在尝试连接")

	serviceCtx.Manager.Connect(iroom_number, iroom_size, ichannel, func(address string, response *benchclient.StartPlanResponse, err error) {

		if err != nil {
			pterm.Error.Println(address, "连接失败, err:", err)
		} else {
			pterm.Info.Println(address, "连接成功, uuid:", response.Uuid)
		}

	})

	pterm.Info.Println("任务连接成功")
}
