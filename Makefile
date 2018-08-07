.PHONY: run open

run:
	go run cmd/main.go

open:
	google-chrome localhost:8080
