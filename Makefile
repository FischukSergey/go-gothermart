ipAddr:=localhost:8080
envRunAddr:=RUN_ADDRESS=$(ipAddr)
envDatabaseDSN:=DATABASE_URI="user=postgres password=postgres host=localhost port=5432 dbname=urlshortdb sslmode=disable"

server:
				@echo "Running server"
				$(envRunAddr) $(envFlagFileStoragePath) go run ./cmd/gophermart/main.go ./cmd/gophermart/config.go
.PHONY: server
