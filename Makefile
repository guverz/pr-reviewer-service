.PHONY: build run test clean docker-build docker-up docker-down

# Определяем расширение для исполняемых файлов
ifeq ($(OS),Windows_NT)
    EXE := .exe
else
    EXE :=
endif

# Сборка приложения
build:
	go build -o bin/server$(EXE) ./cmd/server

# Запуск приложения
run: build
ifeq ($(OS),Windows_NT)
	bin\server$(EXE)
else
	./bin/server$(EXE)
endif

# Интеграционные тесты (требуют запущенный сервер)
test-integration:
	go test ./test/integration_test.go -v

# Очистка артефактов сборки
clean:
ifeq ($(OS),Windows_NT)
	if exist bin rmdir /s /q bin
else
	rm -rf bin/
endif

# Сборка Docker образа
docker-build:
	docker-compose build

# Запуск через docker-compose
docker-up:
	docker-compose up -d

# Остановка docker-compose
docker-down:
	docker-compose down

# Полная пересборка и запуск
docker-restart: docker-down docker-build docker-up


