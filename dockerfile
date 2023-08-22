FROM alpine:latest

COPY app/bench/bench app/bench/etc/bench.yaml /app/bench/

WORKDIR app/bench/

EXPOSE 8080
EXPOSE 9101

ENTRYPOINT ["/app/bench/bench","-f","/app/bench/bench.yaml"]
