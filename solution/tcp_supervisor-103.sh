#!/bin/sh

go build tcp_client_middleware1.go

while true; do
	server=`ps -ef | grep tcp_client_middleware1 | grep "192.168.1.103:2102" | grep -v grep`
	if [ ! "$server" ]; then
		./tcp_client_middleware1 192.168.1.103:2102 3000  coinpush.pokekara.com:10112 3000 >>log-103.log 2>&1 &
		sleep 6
	fi
	sleep 1
done