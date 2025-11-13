.PHONY: help test test-acceptance test-unit test-coverage fmt lint clean deps

help:
	@echo "API Gateway - Comandos disponibles"
	@echo "=================================="
	@echo "make test                - Ejecutar todas las pruebas"
	@echo "make test-acceptance     - Ejecutar pruebas de aceptación (BDD)"
	@echo "make test-coverage       - Ejecutar pruebas con reporte de cobertura"
	@echo "make fmt                 - Formatear código Go"
	@echo "make lint                - Analizar código con golangci-lint"
	@echo "make deps                - Descargar dependencias"
	@echo "make clean               - Limpiar archivos generados"
	@echo "make run                 - Ejecutar el API Gateway"
	@echo "make docker-build        - Construir imagen Docker"

# Ejecutar todas las pruebas
test: deps test-unit test-acceptance

# Ejecutar pruebas de aceptación (BDD con Godog)
test-acceptance:
	@echo "Ejecutando pruebas de aceptación (BDD)..."
	godog run tests/acceptance/features

# Ejecutar pruebas de aceptación con formato pretty
test-acceptance-pretty:
	@echo "Ejecutando pruebas de aceptación con formato pretty..."
	godog run -f pretty tests/acceptance/features

# Ejecutar pruebas unitarias
test-unit:
	@echo "Ejecutando pruebas unitarias..."
	go test -v ./...

# Ejecutar pruebas con cobertura
test-coverage:
	@echo "Ejecutando pruebas con cobertura..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Reporte de cobertura generado en coverage.html"

# Descargar dependencias
deps:
	@echo "Descargando dependencias..."
	go mod download
	go mod tidy

# Formatear código Go
fmt:
	@echo "Formateando código Go..."
	go fmt ./...
	gofmt -s -w .

# Linter
lint:
	@echo "Analizando código..."
	golangci-lint run ./...

# Limpiar archivos generados
clean:
	@echo "Limpiando archivos generados..."
	rm -f coverage.out coverage.html
	go clean -testcache

# Ejecutar el gateway
run:
	@echo "Ejecutando API Gateway..."
	go run app/main.go

# Construir imagen Docker
docker-build:
	@echo "Construyendo imagen Docker..."
	docker build -t api-gateway:latest .

# Ejecutar pruebas en Docker
docker-test:
	@echo "Ejecutando pruebas en Docker..."
	docker-compose -f docker-compose.dev.yml run api-gateway make test

# Instalar herramientas de desarrollo
install-tools:
	@echo "Instalando herramientas de desarrollo..."
	go install github.com/cucumber/godog/cmd/godog@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Ejecutar pruebas en modo watch (requiere entr o similar)
watch-test:
	@echo "Monitoreando cambios para ejecutar pruebas..."
	find . -name "*.go" -o -name "*.feature" | entr make test

# Ejecutar con bandera de tags específicos
test-tag:
	@echo "Ejecutando pruebas con tag: ${TAG}"
	godog run --tags $(TAG) tests/acceptance/features
