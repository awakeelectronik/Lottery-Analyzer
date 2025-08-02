# Makefile
.PHONY: help build run test clean deps setup run-api docker-build docker-run

# Variables
BINARY_NAME=lottery-analyzer
API_BINARY_NAME=lottery-api
MAIN_PATH=./cmd/main.go
API_PATH=./cmd/api/main.go
BUILD_DIR=./bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags="-s -w"

## help: Mostrar este mensaje de ayuda
help:
	@echo "Comandos disponibles:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## setup: Configurar proyecto por primera vez
setup:
	@echo "🚀 Configurando proyecto..."
	@mkdir -p $(BUILD_DIR)
	@$(GOMOD) tidy
	@echo "✅ Proyecto configurado"

## deps: Instalar dependencias
deps:
	@echo "📦 Instalando dependencias..."
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "✅ Dependencias instaladas"

## build: Compilar aplicación principal
build: deps
	@echo "🔨 Compilando aplicación principal..."
	@$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Compilación completada: $(BUILD_DIR)/$(BINARY_NAME)"

## build-api: Compilar API server
build-api: deps
	@echo "🔨 Compilando API server..."
	@$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(API_BINARY_NAME) $(API_PATH)
	@echo "✅ API compilada: $(BUILD_DIR)/$(API_BINARY_NAME)"

## build-all: Compilar todo
build-all: build build-api
	@echo "✅ Todas las aplicaciones compiladas"

## run: Ejecutar aplicación principal
run: build
	@echo "🚀 Ejecutando análisis de lotería..."
	@$(BUILD_DIR)/$(BINARY_NAME)

## run-api: Ejecutar API server
run-api: build-api
	@echo "🌐 Iniciando API server..."
	@$(BUILD_DIR)/$(API_BINARY_NAME)

## run-dev: Ejecutar en modo desarrollo
run-dev:
	@echo "🔧 Ejecutando en modo desarrollo..."
	@$(GOCMD) run $(MAIN_PATH)

## run-api-dev: Ejecutar API en modo desarrollo
run-api-dev:
	@echo "🔧 Ejecutando API en modo desarrollo..."
	@$(GOCMD) run $(API_PATH)

## test: Ejecutar tests
test:
	@echo "🧪 Ejecutando tests..."
	@$(GOTEST) -v ./...

## test-cover: Ejecutar tests con cobertura
test-cover:
	@echo "🧪 Ejecutando tests con cobertura..."
	@$(GOTEST) -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "📊 Reporte de cobertura: coverage.html"

## bench: Ejecutar benchmarks
bench:
	@echo "⚡ Ejecutando benchmarks..."
	@$(GOTEST) -bench=. -benchmem ./...

## lint: Ejecutar linter
lint:
	@echo "🔍 Ejecutando linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint no instalado. Instalando..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
        export PATH="$$PATH:$$($(GOCMD) env GOPATH)/bin"; golangci-lint run; \
	fi

## format: Formatear código
format:
	@echo "✨ Formateando código..."
	@$(GOCMD) fmt ./...
	@$(GOCMD) mod tidy

## clean: Limpiar archivos generados
clean:
	@echo "🧹 Limpiando archivos generados..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Limpieza completada"

## cross-compile: Compilar para múltiples plataformas
cross-compile: deps
	@echo "🌍 Compilando para múltiples plataformas..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "✅ Compilación multiplataforma completada"

## docker-build: Construir imagen Docker
docker-build:
	@echo "🐳 Construyendo imagen Docker..."
	@docker build -t $(BINARY_NAME):latest .
	@echo "✅ Imagen Docker construida"

## docker-run: Ejecutar en Docker
docker-run: docker-build
	@echo "🐳 Ejecutando en Docker..."
	@docker run --rm $(BINARY_NAME):latest

## install: Instalar binario en GOPATH
install: build
	@echo "📦 Instalando binario..."
	@$(GOCMD) install $(MAIN_PATH)
	@echo "✅ Binario instalado en GOPATH"

## check: Verificar código completo
check: format lint test
	@echo "✅ Verificación completa finalizada"

## quick: Compilar y ejecutar rápido
quick:
	@$(GOCMD) run $(MAIN_PATH)

## api-quick: Compilar y ejecutar API rápido
api-quick:
	@$(GOCMD) run $(API_PATH)

## stats: Mostrar estadísticas del proyecto
stats:
	@echo "📊 Estadísticas del proyecto:"
	@echo "   Líneas de código Go:"
	@find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1
	@echo "   Archivos Go:"
	@find . -name "*.go" -not -path "./vendor/*" | wc -l
	@echo "   Dependencias:"
	@$(GOMOD) graph | wc -l

## update: Actualizar dependencias
update:
	@echo "⬆️  Actualizando dependencias..."
	@$(GOMOD) get -u ./...
	@$(GOMOD) tidy
	@echo "✅ Dependencias actualizadas"

# Definir objetivo por defecto
.DEFAULT_GOAL := help