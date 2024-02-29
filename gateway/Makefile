proto:
	rm -rf ./pkg/**/pb/*.go
	rm -rf ./docs/swagger/*
	protoc  \
	--proto_path=pkg --go_out=pkg --go_opt=paths=source_relative --go-grpc_out=pkg --go-grpc_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=foedie \
	pkg/**/pb/*.proto

server:
	go run cmd/main.go