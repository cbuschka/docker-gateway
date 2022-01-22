PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SHELL = /bin/bash

run:
	cd ${PROJECT_DIR}
	go run ./cmd/container-watch

build:
	cd ${PROJECT_DIR}
	mkdir -p dist
	go build -ldflags '-extldflags "-static"' -o dist/container-watch ./cmd/container-watch

build-with-docker:
	cd ${PROJECT_DIR}
	mkdir -p dist
	docker run -ti -v ${PWD}:/work -w /work -u $(shell id -u):$(shell id -g) -e HOME=/work golang:1.17-bullseye \
		go build -ldflags '-extldflags "-static"' -o dist/container-watch ./cmd/container-watch


build-container:
	docker build -t gateway:local .

run-container:
	docker run --name gw --rm -v /var/run/docker.sock:/var/run/docker.sock -p 80:80 -p 443:443 gateway:local

start-example-container:
	docker run -d --rm --name nginx --label=ingress-domain=www.example.com nginx

rm-example-container:
	docker rm -f nginx

start-example2-container:
	docker run -d --rm --name nginx2 --label=ingress-domain=www.example.com nginx

rm-example2-container:
	docker rm -f nginx2
