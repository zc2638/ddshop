FROM golang:1.17 as builder

RUN apt-get update && apt-get install -y libasound2-dev

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /work
ADD . .
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-w -s" -o /usr/local/bin/ddshop github.com/zc2638/ddshop/cmd/ddshop

FROM alpine:3.6
MAINTAINER zc
LABEL maintainer="zc" \
    email="zc2638@qq.com"

ENV TZ="Asia/Shanghai"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++ alsa-lib-dev

COPY --from=builder /usr/local/bin/ddshop /usr/local/bin/ddshop
COPY --from=builder /work/config/config.yaml /work/config/config.yaml

WORKDIR /work
CMD ["ddshop", "-c", "config/config.yaml"]