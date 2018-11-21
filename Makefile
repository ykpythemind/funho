run:
	go run main.go handler.go presenter.go websocket.go

seeds:
	go run ./cmd/seed

migrate:
	go run ./cmd/migrate

migrate-down:
	DOWN=true go run ./cmd/migrate
