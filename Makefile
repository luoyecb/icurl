# demo
.PHONY: default build tidy install clean

default: build

build:
	go build -o icurl

tidy:
	mygomod.sh tidy

install:
	sh install.sh

clean:
	rm -rf icurl
