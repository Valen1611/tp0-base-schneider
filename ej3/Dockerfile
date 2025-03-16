FROM alpine:latest

RUN apk add --no-cache netcat

COPY ../validar-echo-server.sh /validar-echo-server.sh

CMD ["/bin/sh", "/validar-echo-server.sh"]
