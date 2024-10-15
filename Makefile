#include ./src/.env

PROJECT = jachobs-mind

all: up

init:
	-mkdir ./libs
	GOPATH=$$(echo "$${PWD}/libs") go build ./src/main.go

up: build
	PROJECT=$(PROJECT) docker compose -f ./deploy/docker/docker-compose.yml -p $(PROJECT) up -d

build:
	PROJECT=$(PROJECT) docker build -t $(PROJECT)-image .

down:
	PROJECT=$(PROJECT) docker compose -f ./deploy/docker/docker-compose.yml -p $(PROJECT) down

deploy:
	ansible-playbook ./deploy/docker/playbook-up.yaml

destroy:
	ansible-playbook ./deploy/docker/playbook-down.yaml

db-term:
	docker exec -it $(PROJECT)-db-1 bash

redis-cli:
	docker exec -it $(PROJECT)-db-1 redis-cli

restart: down
	PROJECT=$(PROJECT) docker compose -f ./deploy/docker/docker-compose.yml -p $(PROJECT) up -d

k8s:
	kubectl apply -f ./deploy/k8s/

