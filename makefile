proto:
	@mkdir -p api/taskpb
	@protoc --go_out=. --go_opt=module=github.com/n17ali/gohive \
		--go-grpc_out=. --go-grpc_opt=module=github.com/n17ali/gohive \
		api/proto/task.proto

build:
	@go build -o ./bin/goHive ./cmd/server/main.go 

run: build
	@./bin/goHive