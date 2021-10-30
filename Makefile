# demo
.PHONY: default build tidy install clean

default: build

build:
	go build -o icurl main.go

tidy:
	mygomod.sh tidy

install:
	sh install.sh

clean:
	rm -rf icurl