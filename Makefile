pb:
	goctl  rpc protoc bench.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=./app/bench

build:
	cd app/bench && go build
	cd app/bench_cli && go build

build_linux:
	cd app/bench && GOOS=linux GOARCH=amd64 go build
	cd app/bench_cli && GOOS=linux GOARCH=amd64 go build

run_bench:
	nohup app/bench/bench -f app/bench/etc/bench.yaml &
	#nohup app/bench/bench -f app/bench/etc/bench2.yaml &

run_bench_cli:
	cd app/bench_cli && ./bench_cli

kill:
	pkill bench

docker_build:
	docker build -t pomelo_bench:v1.3 .

docker_run:
	docker run -d --restart=always --name=pomelo_bench -p 8080:8080 -p 9101:9101 pomelo_bench:v1.3

docker_save:
	docker save pomelo_bench:v1.3 -o pomelo_bench_1_3.tar

docker_load:
	docker load -i pomelo_bench_1_3.tar