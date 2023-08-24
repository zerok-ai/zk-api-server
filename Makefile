fmt:
	gofmt -s -w .

test:
	go test ./... -cover

mock:
	mockery --all

coverage_cli:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage_html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

create-migration-file:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	migrate -path db/migrations -database "postgres://$$PL_POSTGRES_USERNAME:$$PL_POSTGRES_PASSWORD=@localhost:5432/zk?sslmode=disable&x-migrations-table=$$ZK_SCHEMA_MIGRATIONS_TABLE_NAME" -verbose up $(count)

migrate-down:
	migrate -path db/migrations -database "postgres://$$PL_POSTGRES_USERNAME:$$PL_POSTGRES_PASSWORD@localhost:5432/zk?sslmode=disable&x-migrations-table=$$ZK_SCHEMA_MIGRATIONS_TABLE_NAME" -verbose down $(count)

fix-migration:
	migrate -path db/migrations -database "postgres://$$PL_POSTGRES_USERNAME:$$PL_POSTGRES_PASSWORD@localhost:5432/zk?sslmode=disable&x-migrations-table=$$ZK_SCHEMA_MIGRATIONS_TABLE_NAME" force $(version)
