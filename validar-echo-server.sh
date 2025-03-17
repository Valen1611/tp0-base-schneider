#!/bin/bash

server_address="server"
server_port="12345"
msg="Hello World"

docker run --rm --network tp0_testing_net \
    -e server_port="$server_port" \
    -e msg="$msg" \
    alpine:latest /bin/sh -c '
    apk add --no-cache netcat-openbsd && \
    response=$(echo $msg | nc -w 3 server $server_port) && \
        if [ "$response" == "$msg" ]; then
        echo "action: test_echo_server | result: success"
    else
        echo "action: test_echo_server | result: fail"
    fi
'
