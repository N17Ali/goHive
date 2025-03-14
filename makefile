proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/pb/task.proto

build:
	@go build -o ./bin/goHive ./cmd/server/main.go 

run: build
	@./bin/goHive