.PHONY: run test lint db-up db-down db-reset tidy up down

API_DIR  := services/go-api
DB_URL   := postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable

# Starta API:et lokalt (kräver att Postgres körs, se db-up)
run:
	cd $(API_DIR) && DATABASE_URL=$(DB_URL) go run ./cmd/api

# Kör alla tester
test:
	cd $(API_DIR) && go test ./...

# Statisk analys
lint:
	cd $(API_DIR) && go vet ./...

# Ladda ned beroenden och uppdatera go.sum
tidy:
	cd $(API_DIR) && go mod tidy

# Starta Postgres + Redis i bakgrunden
db-up:
	docker compose up -d postgres redis

# Stoppa containers (behåller data)
db-down:
	docker compose down

# Stoppa och radera all data (börja om från scratch)
db-reset:
	docker compose down -v
	docker compose up -d postgres redis

# Bygg och starta hela stacken (go-api + Postgres + Redis) i containers
up:
	docker compose up -d --build

# Stoppa hela stacken (behåller data)
down:
	docker compose down
