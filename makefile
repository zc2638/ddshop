cookie="DDXQSESSID=XXOO请在这里设置你的cookie，或者conf/conf.yaml填写，然后点击run运行,或者根目录下make run f******24***c***b"

.PHONY: go-proxy
go-proxy:
	export GOPROXY=https://mirrors.aliyun.com/goproxy/

.PHONY: dep
dep: go-proxy
	go mod tidy -compat=1.17

# build之后 执行这个让程序后台运行 nohup ./build/app > info.log 2>&1 &  用 jobs -l 查看
build: dep
	go build -o ./build/app ./cmd/ddshop/main.go

.PHONY: run
run: build
	./build/app --cookie ${cookie}

setup: go-proxy
	go install -v mvdan.cc/gofumpt@latest
	go install -v golang.org/x/tools/cmd/goimports@latest
	go install -v github.com/daixiang0/gci@latest

# 格式化 import 分组 和 代码格式化
format:
	$(shell go env GOPATH)/bin/gofumpt -w -l .
	$(shell go env GOPATH)/bin/gci -w .
	$(shell go env GOPATH)/bin/goimports -w .

