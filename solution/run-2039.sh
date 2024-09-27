#!/bin/sh

go build tcp_client_middleware.go

while true; do
        server=`ps -ef | grep tcp_client_middleware | grep "192.168.1.39:2039" | grep -v grep`
        if [ ! "$server" ]; then
            ./tcp_client_middleware 192.168.1.39:2039 3000  406a7637n7.goho.co:11180 10000 >>log-2039.log 2>&1 &
            sleep 10
        fi
        sleep 5
done