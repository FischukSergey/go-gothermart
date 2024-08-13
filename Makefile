ipAddr:=localhost:8080
envRunAddr:=RUN_ADDRESS=$(ipAddr)
envDatabaseDSN:=DATABASE_URI="user=postgres password=postgres host=localhost port=5432 dbname=gophermartdb sslmode=disable"

server:
				@echo "Running server"
				$(envRunAddr) $(envDatabaseDSN) go run ./cmd/gophermart/
.PHONY: server

migration:
	@echo "Running migration"
	go run ./cmd/migrator \
 		--storage-path="postgresql://postgres:postgres@localhost:5432/gophermartdb?sslmode=disable" \
 		--migrations-path=./migrations
.PHONY: migration

accrual:
	./cmd/accrual/accrual_darwin_arm64 -a=localhost:8082
.PHONY: accrual

# pg_ctl -D /usr/local/pgsql/data stop/start