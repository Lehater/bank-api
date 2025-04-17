.PHONY: test cover run lint

# Запустить все тесты
test:
	go test ./... -v

# Запустить тесты с покрытием и сгенерировать HTML-отчет
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Запустить приложение
run:
	go run cmd/api/main.go

# Проверка форматирования (можно заменить на golangci-lint)
lint:
	gofmt -l ./
