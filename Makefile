#include ./bot/.env

PROJECT = jachobs-mind

all: up

init:
	touch .env
	GOPATH=$$(echo "$${PWD}/libs") go build -o ./jachobs-mind ./bot/
	-rm jachobs-mind


up: build
	PROJECT=$(PROJECT) docker compose -f ./deploy/docker/docker-compose-run.yml -p $(PROJECT) up -d

build:
	PROJECT=$(PROJECT) docker build -f Dockerfile-Build -t jdolakk/$(PROJECT)-build-image .
	PROJECT=$(PROJECT) docker compose -f ./deploy/docker/docker-compose-build.yml -p $(PROJECT) create
	docker cp $(PROJECT)-bot-build-1:/$(PROJECT)/bin/$(PROJECT) ./bin/

	PROJECT=$(PROJECT) docker build -f Dockerfile-Run -t jdolakk/$(PROJECT)-image .

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

test-up:
	PROJECT=$(PROJECT) docker compose -f ./deploy/docker/docker-compose.yml -p $(PROJECT)-test up -d
	bash ./ci/tests.sh

test-down:
	PROJECT=$(PROJECT) docker compose -f ./deploy/docker/docker-compose.yml -p $(PROJECT)-test down



