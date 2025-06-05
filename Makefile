init: up\
packages-install \
build \
run

up:
	docker-compose -f docker-compose.yaml up -d --build --remove-orphans

packages-install:
	go mod tidy

build:
	go build -a -o main cmd/api/main.go

run:
	./main