#!/bin/bash

server_address="server"
server_port="12345"
msg="Hello World"

docker run --rm --network tp0_testing_net alpine:latest /bin/sh -c '
  apk add --no-cache netcat-openbsd && \
  response=\$(echo $msg | nc $server_address $server_port) && \
    if [ "$response" == "$msg" ]; then
    echo "action: test_echo_server | result: success"
  else
    echo "action: test_echo_server | result: fail"
  fi
'
