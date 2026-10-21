.PHONY: up down down-v logs test lint swag clean

# Запустить линтер, сгенерировать сваггер, а затем поднять все сервисы в фоне
up: lint swag
	docker compose up --build -d

# Генерация документации Swagger
swag:
	swag init -g cmd/server_api/main.go --parseDependency --parseInternal

# Запуск линтера
lint:
	golangci-lint run ./...

# Остановить контейнеры
down:
	docker compose down

# Остановить контейнеры и полностью удалить данные БД (volumes)
down-v:
	docker compose down -v

# Посмотреть логи всех сервисов
logs:
	docker compose logs -f

# Очистка локальных артефактов сборки
clean:
	rm -rf bin/ docs/