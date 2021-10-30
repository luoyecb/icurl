#!/bin/bash
#
# install to GOPATH

if [ -z "$GOPATH" ]; then
	exit 0
fi

first_part=`echo $GOPATH | awk -F":" '{print $1}'`
if [ -z "$first_part" ]; then
	exit 0
fi

mv icurl "$first_part/bin/"
