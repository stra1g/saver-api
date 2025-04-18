generate_migration:
	migrate create -ext sql -dir internal/infra/database/migrations -seq $(name)
	@echo "Migration file created in internal/infra/database/migrations"

migrateup:
	migrate -path internal/infra/database/migrations -database "postgresql://postgres:postgres@127.0.0.1:5432/saver_api?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/infra/database/migrations -database "postgresql://postgres:postgres@postgres:5432/saver_api?sslmode=disable" -verbose down

.PHONY: generate_migration migrateup migratedown

