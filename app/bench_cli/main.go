package main

import (
	"flag"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/config"
	"pomelo_bench/app/bench_cli/internal/svc"
)

var (
	serviceCtx *svc.ServiceContext
	exit       = false
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {

	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := logx.SetUp(logx.LogConf{
		ServiceName: "bench_cli",
		Mode:        "file",
		Path:        "logs",
	}); err != nil {
		panic(err)
	}

	serviceCtx = svc.NewServiceContext(c)

	// 连接上woker机器
	connectWorker(c.WorksAddr)

	for !exit {

		menu()

	}

}

// 主菜单
func menu() {

	options := []string{
		"connect",
		"send",
		"recover",
		"close",
		"tree",
		"quit",
	}

	printer := pterm.DefaultInteractiveSelect.WithOptions(options)

	selectedOption, _ := printer.Show()
	if selectedOption == "connect" {

		connect()

	} else if selectedOption == "send" {

		send()

	} else if selectedOption == "recover" {

		sendRecover()

	} else if selectedOption == "close" {

		closePlan()

	} else if selectedOption == "tree" {

		monitorTree()

	} else if selectedOption == "quit" {

		exit = true
		pterm.Info.Println("good bye!")
	}

}

func connectWorker(worksAddr []string) {
	pterm.DefaultSection.Println("pomelo bench cli!")
	pterm.Info.Println("尝试连接压测woker机器")

	for i := 0; i < len(worksAddr); i++ {
		cli, err := zrpc.NewClient(zrpc.RpcClientConf{
			Target:  worksAddr[i],
			Timeout: 10 * 1000,
		})

		if err != nil {

			pterm.Error.Println(fmt.Sprintf("连接 %s 失败,err:%s", worksAddr[i], err))

		} else {

			pterm.Info.Println(fmt.Sprintf("连接 %s 成功", worksAddr[i]))

			c := benchclient.NewBench(cli)
			serviceCtx.Manager.Add(svc.Woker{
				Address: worksAddr[i],
				Client:  c,
			})
		}
	}

	pterm.Info.Println("连接压测woker机器完成")
}
