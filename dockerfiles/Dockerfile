FROM golang:1.13.6
COPY ./ .
ENV GOPROXY=https://goproxy.cn,direct
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && make build
EXPOSE 7890
CMD ["./bin/extractor-server","run","-c","./examples/extractor-server.yaml"]