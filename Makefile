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
	migrate -path db/migrations -database "postgres://pl:pl=@localhost:5432/pl?sslmode=disable&x-migrations-table=zk_api_server_migrations" -verbose up $(count)

migrate-down:
	migrate -path db/migrations -database "postgres://pl:pl=@localhost:5432/pl?sslmode=disable&x-migrations-table=zk_api_server_migrations" -verbose down $(count)

fix-migration:
	migrate -path db/migrations -database "postgres://pl:pl=@localhost:5432/pl?sslmode=disable&x-migrations-table=zk_api_server_migrations" force $(version)
