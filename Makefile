# ================================
# Project Configuration
# ================================
BINARY_NAME := spektr-api
GO_FLAGS :=
GOOS := linux
GOARCH := amd64
USER := user
HOST := 90.188.22.88
DIR := /home/user/friendly-basketball-media/
PASS := `u9F2DCWg`
MAIN_FILE := app/main.go

# ================================
# Build Process
# ================================

# Сборка проекта для текущей ОС и архитектуры
.PHONY: build
build:
	@echo "Building project..."
	go build $(GO_FLAGS) -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build completed: $(BINARY_NAME)"

# Сборка проекта для Linux (GOOS=linux, GOARCH=amd64)
.PHONY: ubuntu
ubuntu:
	@echo "Building for Linux (GOOS=linux, GOARCH=amd64)..."
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build completed for Linux: $(BINARY_NAME)"

# Деплой проекта на удаленный сервер
.PHONY: deploy
deploy:
	@echo "Deploying $(BINARY_NAME) to $(USER)@$(HOST):$(DIR)..."
	# Проверка существования директории на сервере
	@echo "Deploying using SCP..."
	# Команда SCP без пробелов
	scp $(BINARY_NAME) $(USER)@$(HOST):$(DIR)
	@echo "Deployment completed."

# Запуск проекта (после сборки)
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Удаление скомпилированного бинарного файла
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	@echo "Cleaned."

# Полный цикл: очистка, сборка и запуск
.PHONY: all
all: clean run
