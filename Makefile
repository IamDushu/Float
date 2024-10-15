postgres: 
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root floatdb

dropdb:
	docker exec -it postgres12 dropdb floatdb

migrateup: 
	migrate -path internal/db/migration -database="postgresql://root:secret@localhost:5432/floatdb?sslmode=disable" -verbose up

migratedown: 
	migrate -path internal/db/migration -database="postgresql://root:secret@localhost:5432/floatdb?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc