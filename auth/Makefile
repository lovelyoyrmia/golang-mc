DB_URL=postgresql://root:root@localhost:5432/foedie_auth?sslmode=disable

network:
	docker network create foedie-network

redis:
	docker run --name redis -p 6379:6379 -d redis

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root -d postgres:12-alpine

server:
	go run cmd/main.go

proto:
	rm -rf domain/pb/*.go
	protoc --proto_path=domain/proto --go_out=domain/pb --go_opt=paths=source_relative \
    --go-grpc_out=domain/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=domain/pb --grpc-gateway_opt=paths=source_relative \
    domain/proto/*.proto

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root foedie_auth

dropdb:
	docker exec -it postgres12 dropdb foedie_auth

sqlc:
	sqlc generate

migrate_up:
	migrate -path migrations -database "$(DB_URL)" --verbose up

create_migrate:
	migrate create -ext sql -dir migrations -seq init_schema

migrate_down:
	migrate -path migrations -database "$(DB_URL)" --verbose down

migrate_up1:
	migrate -path migrations -database "$(DB_URL)" --verbose up 1

migrate_down1:
	migrate -path migrations -database "$(DB_URL)" --verbose down 1

new_migration:
	migrate create -ext sql -dir migrations -seq $(name)
	
db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

test:
	go test -v -cover ./...

mock:
	mockgen -package mock -destination internal/mock/store.go github.com/Foedie/foedie-server-v2/auth/internal/db Store && \
	mockgen -package mock -destination internal/mock/user.go github.com/Foedie/foedie-server-v2/auth/internal/worker TaskUser && \
	mockgen -package mock -destination internal/mock/token.go github.com/Foedie/foedie-server-v2/auth/pkg/token Maker

.PHONY: postgres server createdb dropdb proto_start sqlc create_migrate migrate_up migrate_up1 migrate_down migrate_down1 new_migration db_docs db_schema test mock