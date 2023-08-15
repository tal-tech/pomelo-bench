pb:
	goctl  rpc protoc bench.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=./app/bench

build:
	cd app/bench && go build
	cd app/bench_cli && go build

run_bench:
	nohup app/bench/bench -f app/bench/etc/bench.yaml &
	nohup app/bench/bench -f app/bench/etc/bench2.yaml &

run_bench_cli:
	cd app/bench_cli && ./bench_cli

kill:
	pkill bench
