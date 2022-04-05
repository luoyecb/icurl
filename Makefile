# demo
.PHONY: default build tidy clean

GOPROXY := https://goproxy.cn,direct
GOPRIVATE :=
GO111MODULE := auto

export GOPROXY
export GOPRIVATE
export GO111MODULE

default: build

build:
	go build -o icurl

tidy:
	go mod tidy

clean:
	rm -rf icurl
