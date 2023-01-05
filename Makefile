run:
	go run main.go

build:
	go build

vendor:
	go mod vendor

up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans
