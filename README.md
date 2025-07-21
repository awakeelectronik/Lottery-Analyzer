# Lottery Analyzer - Conversor de Grails 2.4.4 a Go 1.21+

Proyecto en Groovy on Grails 2.4.4: https://github.com/awakeelectronik/chance

Proyecto optimizado de análisis de lotería convertido de Grails a Go, manteniendo toda la funcionalidad original mientras implementa las mejores prácticas de desarrollo en Go.

## 🏗️ Arquitectura

Implementación de **Clean Architecture** con separación clara de responsabilidades:

```
lottery-analyzer/
├── cmd/                    # Puntos de entrada
│   ├── main.go            # Aplicación principal (análisis)
│   └── api/main.go        # Servidor API REST
├── internal/              # Código privado de la aplicación
│   ├── config/            # Configuración
│   ├── service/           # Lógica de negocio
│   ├── repository/        # Acceso a datos
│   └── model/             # Modelos de dominio
├── pkg/                   # Código reutilizable
│   └── database/          # Conexión a MySQL
├── bin/                   # Binarios compilados
├── go.mod                 # Dependencias
├── Makefile              # Automatización
├── .env                  # Variables de entorno
└── README.md             # Documentación
```

## 🚀 Configuración Rápida

### 1. Configurar Proyecto

```bash
# Clonar o crear directorio
mkdir lottery-analyzer && cd lottery-analyzer

# Inicializar módulo Go
go mod init lottery-analyzer

# Configurar con Make
make setup
```

### 2. Instalar Dependencias

```bash
make deps
```

## 🎯 Cómo Ejecutar

### Opción 1: Ejecución Rápida (Desarrollo)

```bash
# Ejecutar análisis principal (equivalente al ProcessorController.index())
make run-dev

# Ejecutar servidor API
make run-api-dev
```

### Opción 2: Compilar y Ejecutar

```bash
# Compilar y ejecutar análisis completo
make run

# Compilar y ejecutar API
make run-api
```

### Opción 3: Comandos Específicos

```bash
# Solo scrapping
go run cmd/main.go  # (incluye scrapping automático)

# Solo servidor API
go run cmd/api/main.go
```

### 🔧 Algoritmo Principal

El algoritmo mantiene la **misma lógica exacta** que el ProcessorController original:

1. **Scrapping automático** desde última fecha
2. **Secuencia Fibonacci** para análisis temporal (1, 1, 2, 3, 5, 8...)
3. **15 consultas SQL** para patrones de frecuencia
4. **Cálculo de probabilidades** con factores específicos
5. **Selección de mejores 100 números** con scores ordenados

### 🌐 API REST Endpoints

```bash
# Análisis completo (equivalente al controller original)
curl -X POST http://localhost:8080/api/v1/analysis/process

# Mejores números
curl http://localhost:8080/api/v1/analysis/best-numbers?limit=50

# Health check
curl http://localhost:8080/health

## 🛠️ Comandos Make Disponibles

### Desarrollo

```bash
make help           # Ver todos los comandos
make setup          # Configurar proyecto inicial
make deps           # Instalar dependencias
make run-dev        # Ejecutar en desarrollo
make run-api-dev    # Ejecutar API en desarrollo
```

### Compilación

```bash
make build          # Compilar aplicación principal
make build-api      # Compilar API
make build-all      # Compilar todo
make cross-compile  # Compilar para múltiples OS
```

### Testing y Calidad

```bash
make test           # Ejecutar tests
make test-cover     # Tests con cobertura
make lint           # Linter de código
make format         # Formatear código
make check          # Verificación completa
```

### Utilidades

```bash
make clean          # Limpiar archivos generados
make stats          # Estadísticas del proyecto
make update         # Actualizar dependencias
```

## 📈 Ventajas de la Conversión

### Rendimiento
- **5-10x más rápido** que Grails (compilación nativa)
- **Menor uso de memoria** (sin JVM overhead)
- **Concurrencia nativa** con goroutines

### Mantenibilidad
- **Tipado estático** previene errores
- **Interfaces claras** para testing y mocking
- **Separación de responsabilidades** por capas

### Escalabilidad
- **Arquitectura modular** fácil de extender
- **Repository pattern** para flexibilidad de BD
- **Dependency injection** manual y controlado

## 🔍 Verificación de Funcionamiento

### 1. Probar Conexión a BD

```bash
# El programa verificará automáticamente la conexión
make run-dev
```

### 2. Verificar Scrapping

```bash
# Los logs mostrarán el progreso del scrapping
make run-dev
# Salida esperada:
# "Starting lottery analysis..."
# "Scraping from last date..."
# "Analysis completed in Xs"
```

### 3. Probar API

```bash
# Terminal 1: Iniciar API
make run-api-dev

# Terminal 2: Probar endpoints
curl http://localhost:8080/health
curl -X POST http://localhost:8080/api/v1/analysis/process
```

## ⚡ Optimizaciones Implementadas

### Código Optimizado
- **Reutilización de conexiones** con pool de BD
- **Context cancellation** para operaciones costosas
- **Preparación de statements** SQL
- **Batch processing** para inserciones

### Gestión de Memoria
- **Pre-allocación** de slices con capacidad conocida
- **Reutilización de estructuras** donde es posible
- **Garbage collection** optimizado

### Concurrencia
- **Graceful shutdown** con señales del sistema
- **Context timeouts** para operaciones HTTP
- **Goroutines controladas** para scrapping

### Dependencias Faltantes

```bash
# Reinstalar dependencias
make clean
make deps
```
