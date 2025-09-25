wire-gen:
	@echo "Generating code..."

	go fmt ./cmd/api/di

	wire ./...

	@echo "Code generation complete."

start-app:
	@echo "Starting the application..."

	go run ./cmd/api/main.go

migrate-schema:
	@echo "Strating the migrate schema"

	go run ./cmd/migration/main.go
	