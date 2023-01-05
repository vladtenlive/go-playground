run:
	go run main.go

build:
	go build

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto

vendor:
	go mod vendor

up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans
