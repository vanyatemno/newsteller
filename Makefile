.PHONY: run

run:
	docker compose up --build

test:
	go test ./internal/*