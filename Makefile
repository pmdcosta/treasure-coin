.PHONY: start-db

start-db:
	docker run -d --name treasure-coin-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=treasure_coin -p 5432:5432 postgres:11-alpine

run:
	go run cmd/main.go
