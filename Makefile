postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres17 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
migrateversion:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose version
sqlc:
	sqlc generate

# Run all tests in the db/sqlc package (TestMain runs automatically)
test:
	go test -v -cover ./...

# Run a single test by name: make test-one t=TestCreateAccount
test-one:
	@if [ -z "$(t)" ]; then echo "Usage: make test-one t=TestNameRegex"; exit 1; fi
	go test ./db/sqlc -v -run $(t)

# Run the server (convenience)
runserver:
	@echo "Starting server..."
	@go run main.go

mock:
	mockgen -package=mockdb -destination=db/mock/store.go github.com/ShubhKanodia/GoBank/db/sqlc Store
# Start DB and run migrations, then execute tests (convenience)
runtest: postgres createdb migrateup test

.PHONY: createdb postgres dropdb migrateup migrateup1 migratedown migratedown1 migrateversion sqlc test test-one runserver mock runtest
 