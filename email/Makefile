proto:
	rm -rf domain/pb/*.go
	protoc --proto_path=domain/proto --go_out=domain/pb --go_opt=paths=source_relative \
    --go-grpc_out=domain/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=domain/pb --grpc-gateway_opt=paths=source_relative \
    domain/proto/*.proto

server:
	go run cmd/main.go