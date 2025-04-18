SHELL := /bin/bash
PWD := $(shell pwd)

GIT_REMOTE = github.com/7574-sistemas-distribuidos/docker-compose-init

default: build

all:

deps:
	go mod tidy
	go mod vendor

build: deps
	GOOS=linux go build -o bin/client github.com/7574-sistemas-distribuidos/docker-compose-init/client
.PHONY: build

docker-image:
	docker build -f ./server/Dockerfile -t "server:latest" .
	docker build -f ./client/Dockerfile -t "client:latest" .
	# Execute this command from time to time to clean up intermediate stages generated 
	# during client build (your hard drive will like this :) ). Don't left uncommented if you 
	# want to avoid rebuilding client image every time the docker-compose-up command 
	# is executed, even when client code has not changed
	# docker rmi `docker images --filter label=intermediateStageToBeDeleted=true -q`
.PHONY: docker-image

docker-compose-up: docker-image
	docker compose -f docker-compose-dev.yaml up -d --build
.PHONY: docker-compose-up

docker-compose-down:
	docker compose -f docker-compose-dev.yaml stop -t 1
	docker compose -f docker-compose-dev.yaml down
.PHONY: docker-compose-down

docker-compose-logs:
	docker compose -f docker-compose-dev.yaml logs -f
.PHONY: docker-compose-logs


# Para modificar el codigo y verlo mas rapido
reboot-client:
	docker rm client1 --force
	docker image remove client:latest
	docker build -f ./client/Dockerfile -t "client:latest" .
	docker compose -f docker-compose-dev.yaml up --build client1
.PHONY: reboot-client

reboot-server:
	docker rm server --force
	docker image remove server:latest
	docker build -f ./server/Dockerfile -t "server:latest" .
	docker compose -f docker-compose-dev.yaml up --build server
.PHONY: reboot-server

quick-run-client:
	docker compose -f docker-compose-dev.yaml up client1
.PHONY: quick-run-client

quick-run-server:
	docker compose -f docker-compose-dev.yaml up server
.PHONY: quick-run-server
	