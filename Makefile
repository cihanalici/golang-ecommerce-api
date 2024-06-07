server:
	go run main.go

test:
	go test -v -cover ./...

sqlc:
	sqlc generate

postgres:
	docker run --name postgres16 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.3-alpine3.19

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_store

dropdb:
	docker exec -it postgres16 dropdb simple_store

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_store?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_store?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_store?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_store?sslmode=disable" -verbose down 1

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/cihanalici/api/db/sqlc Store

mailHog:
	docker run -d -p 1025:1025 -p 8025:8025 mailhog/mailhog

.PHONY: server test sqlc postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 mock mailHog