package svc

import (
	"bufio"
	"os"
	"pomelo_bench/app/bench_cli/internal/config"
)

type RecoverData struct {
	Describe string
	Router   string   // 消息对应的路由
	Data     [][]byte // 发送的消息信息
}

type ServiceContext struct {
	Config  config.Config
	Manager *WorkManager

	RecoverDataPool map[string]RecoverData // key 是路由 value 是 RecoverData 数据

}

func NewServiceContext(c config.Config) *ServiceContext {
	recoverDataPool, err := getRecoverDataPool(c.CustomSendFiles)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:          c,
		Manager:         NewWorkManager(c),
		RecoverDataPool: recoverDataPool,
	}
}

func getRecoverDataPool(fs []string) (res map[string]RecoverData, err error) {

	res = make(map[string]RecoverData, len(fs))

	for i := 0; i < len(fs); i++ {

		if len(fs[i]) == 0 {
			continue
		}

		data, err := getRecoverData(fs[i])
		if err != nil {
			return nil, err
		}

		res[data.Router] = data
	}

	return res, nil
}

func getRecoverData(filePath string) (res RecoverData, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return RecoverData{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var (
		index int
	)

	res.Data = make([][]byte, 0, 10000)

	for scanner.Scan() {
		if index == 0 {

			res.Describe = scanner.Text()
		} else if index == 1 {

			res.Router = scanner.Text()
		} else {

			res.Data = append(res.Data, []byte(scanner.Text()))
		}

		index++
	}

	if err := scanner.Err(); err != nil {
		return RecoverData{}, err
	}

	return res, nil
}
