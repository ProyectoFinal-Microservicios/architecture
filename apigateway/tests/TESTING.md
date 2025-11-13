# API Gateway - Pruebas de Aceptación (BDD)

Esta documentación describe cómo ejecutar y trabajar con las pruebas de aceptación basadas en BDD para el API Gateway.

## Estructura del Proyecto

```
tests/
├── acceptance/
│   ├── features/                    # Archivos Gherkin con escenarios
│   │   ├── gateway_routing.feature
│   │   ├── gateway_authentication.feature
│   │   └── gateway_user_operations.feature
│   ├── steps/
│   │   └── gateway_steps.go         # Implementación de steps en Go
│   ├── main_test.go                 # Entry point para Godog
│   └── schemas/
│       └── gateway-schemas.json     # Esquemas JSON para validación
└── support/
    ├── http_client.go               # Cliente HTTP para solicitudes
    └── schema_validator.go          # Validador de esquemas JSON
```

## Requisitos Previos

- Go 1.21 o superior
- Godog CLI: `go install github.com/cucumber/godog/cmd/godog@latest`
- API Gateway ejecutándose en `http://localhost:8080`

## Instalación

1. **Descargar dependencias:**
```bash
make deps
# o manualmente:
go mod download
go mod tidy
```

2. **Instalar Godog (si no está instalado):**
```bash
go install github.com/cucumber/godog/cmd/godog@latest
```

## Ejecutar Pruebas

### Ejecución Rápida

```bash
# Ejecutar todas las pruebas
make test-acceptance

# Ejecutar con formato pretty (más legible)
make test-acceptance-pretty
```

### Ejecución Avanzada

```bash
# Ejecutar un feature específico
godog run tests/acceptance/features/gateway_routing.feature

# Ejecutar con tags específicos
godog run --tags @routing tests/acceptance/features

# Ejecutar con formato JSON (para CI/CD)
godog run -f json tests/acceptance/features > report.json

# Ejecutar con nivel de detalle aumentado
godog run -v tests/acceptance/features

# Ejecutar paralelo (experimental)
godog run -c 3 tests/acceptance/features
```

### Formatos de Salida Disponibles

```bash
# Progress (por defecto)
godog run -f progress tests/acceptance/features

# Pretty (legible)
godog run -f pretty tests/acceptance/features

# JSON (para parseo)
godog run -f json tests/acceptance/features

# JUnit XML (para CI/CD)
godog run -f junit tests/acceptance/features > report.xml
```

## Archivos de Features

### gateway_routing.feature
Verifica que el gateway enrute correctamente las solicitudes a los servicios upstream:
- Health checks (`/health`)
- Documentación OpenAPI (`/docs/swagger.json`)
- CORS (Cross-Origin Resource Sharing)
- Enrutamiento a servicios de autenticación y perfiles
- Headers de traza
- Logging de solicitudes
- Manejo de errores de conexión

**Escenarios:** 9
**Tags:** `@routing`, `@health`, `@cors`, `@middleware`

### gateway_authentication.feature
Valida el flujo de autenticación a través del gateway:
- Login (POST `/auth/login`)
- Registro de usuarios (POST `/auth/register`)
- Validación de credenciales
- Manejo de errores de autenticación
- Persistencia de tokens

**Escenarios:** 8
**Tags:** `@authentication`, `@login`, `@register`

### gateway_user_operations.feature
Prueba operaciones CRUD de usuarios a través del gateway:
- Obtener perfil de usuario (GET `/users/{username}`)
- Actualizar perfil (PATCH `/users/{username}/profile`)
- Eliminar usuario (DELETE `/users/{username}`)
- Validación de autorización
- Verificación de enrutamiento correcto

**Escenarios:** 9
**Tags:** `@users`, `@crud`, `@authorization`

## Estructura del Código Go

### HTTPClient (`tests/support/http_client.go`)

Cliente HTTP especializado para pruebas que maneja:
- Autenticación JWT
- Headers personalizados
- Timeouts
- Parseo de respuestas JSON

**Métodos principales:**
- `NewHTTPClient(baseURL)` - Crear cliente
- `SetToken(token)` - Establecer token JWT
- `GET(endpoint, headers)` - Solicitud GET
- `POST(endpoint, body, headers)` - Solicitud POST
- `PATCH(endpoint, body, headers)` - Solicitud PATCH
- `DELETE(endpoint, headers)` - Solicitud DELETE

**Ejemplo:**
```go
client := support.NewHTTPClient("http://localhost:8080")
response, err := client.GET("/health", make(map[string]string))
if err != nil {
    t.Fatalf("Request failed: %v", err)
}
```

### SchemaValidator (`tests/support/schema_validator.go`)

Validador de esquemas JSON para verificar respuestas:
- Carga esquemas desde archivo JSON
- Compila esquemas al inicializar
- Valida datos contra esquemas
- Reporta errores de validación

**Métodos principales:**
- `NewSchemaValidator(schemasPath)` - Crear validador
- `Validate(data, schemaName)` - Validar datos
- `ValidateJSON(jsonStr, schemaName)` - Validar JSON string

**Ejemplo:**
```go
validator, _ := support.NewSchemaValidator("./features/schemas/gateway-schemas.json")
data := map[string]interface{}{"status": "UP"}
result := validator.Validate(data, "healthResponse")
if !result.IsValid {
    t.Fatalf("Validation failed: %v", result.Errors)
}
```

### Steps (`tests/acceptance/steps/gateway_steps.go`)

Implementación de los steps del Gherkin en Go usando Godog.

**Estructura:**
- `APIGatewayContext` - Contexto compartido entre steps
- Métodos para cada step del Gherkin
- `InitializeScenario` - Registro de todos los steps

**Patrones de implementación:**

**1. Solicitud HTTP:**
```gherkin
Cuando realizo una solicitud GET a "/health"
```
```go
func (ctx *APIGatewayContext) MakeGetRequest(endpoint string) error {
    response, err := ctx.client.GET(endpoint, ctx.customHeaders)
    ctx.lastResponse = response
    return err
}
```

**2. Validación de código de estado:**
```gherkin
Entonces el código de estado debe ser 200
```
```go
func (ctx *APIGatewayContext) StatusCodeShouldBe(expectedStatus int) error {
    if ctx.lastResponse.StatusCode != expectedStatus {
        return fmt.Errorf("expected %d but got %d", expectedStatus, ctx.lastResponse.StatusCode)
    }
    return nil
}
```

**3. Validación con schema:**
```gherkin
Y la respuesta debe validarse contra el esquema "healthResponse"
```
```go
func (ctx *APIGatewayContext) ResponseValidatesAgainstSchema(schemaName string) error {
    result := ctx.validator.Validate(data, schemaName)
    if !result.IsValid {
        return fmt.Errorf("validation failed: %v", result.Errors)
    }
    return nil
}
```

## Data Tables (Tablas de Datos)

Las pruebas usan tablas de datos para pasar información estructurada:

```gherkin
Cuando realizo una solicitud POST a "/auth/login"
| field    | value          |
| username | testuser       |
| password | password123    |
```

**Notas importantes:**
- Headers en inglés (field, value)
- Datos pueden estar en cualquier idioma
- Los headers se convierten automáticamente a un mapa Go

## Esquemas JSON

Los esquemas se definen en `tests/acceptance/schemas/gateway-schemas.json`:

```json
{
  "healthResponse": {
    "type": "object",
    "properties": {
      "status": {"type": "string"},
      "timestamp": {"type": "string"}
    },
    "required": ["status"]
  }
}
```

### Esquemas Disponibles

1. **healthResponse** - Respuesta del endpoint `/health`
2. **loginResponse** - Respuesta de login con token
3. **registerResponse** - Respuesta de registro de usuario
4. **profileResponse** - Datos de perfil de usuario
5. **errorResponse** - Respuesta de error estándar
6. **deleteResponse** - Respuesta de eliminación de usuario

## Integración Continua

### GitHub Actions Ejemplo

```yaml
name: API Gateway Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      gateway:
        image: api-gateway:latest
        ports:
          - 8080:8080
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run BDD Tests
        run: make test-acceptance
      
      - name: Upload Coverage
        uses: codecov/codecov-action@v3
```

### Docker Compose

```yaml
version: '3.8'
services:
  api-gateway:
    build: .
    ports:
      - 8080:8080
    
  tests:
    build:
      context: .
      dockerfile: Dockerfile.test
    depends_on:
      - api-gateway
    command: make test-acceptance
```

## Troubleshooting

### Problema: "godog: command not found"

**Solución:**
```bash
go install github.com/cucumber/godog/cmd/godog@latest
# Asegúrate de que $GOPATH/bin está en tu PATH
```

### Problema: "connection refused" en localhost:8080

**Solución:**
- Verifica que el API Gateway está ejecutándose
- Comprueba el puerto en la configuración
- Ajusta la URL base si es necesario en `gateway_steps.go`

### Problema: "schema not found"

**Solución:**
- Verifica la ruta en `NewSchemaValidator("./features/schemas/gateway-schemas.json")`
- Asegúrate de que el archivo existe y es JSON válido
- Revisa que el nombre del schema coincide exactamente

### Problema: "response does not match schema"

**Solución:**
- Verifica la estructura JSON de la respuesta real
- Actualiza el esquema en `gateway-schemas.json`
- Agrega campos requeridos al esquema
- Usa `godog run -v` para ver detalles

## Mejores Prácticas

1. **Mantén tests independientes** - Cada escenario debe funcionar solo
2. **Usa datos realistas** - Los datos de prueba deben ser similares a los reales
3. **Evita acoplamiento** - No dependa del orden de ejecución
4. **Documentación clara** - Escribe steps entendibles en español
5. **Schemas actualizados** - Sincroniza schemas con API real
6. **Reutiliza contexto** - Aprovecha `APIGatewayContext` para compartir estado
7. **Manejo de errores** - Verifica casos de error en los tests

## Recursos Útiles

- [Godog Documentation](https://github.com/cucumber/godog)
- [Gherkin Syntax](https://cucumber.io/docs/gherkin/)
- [JSON Schema Specification](https://json-schema.org/)
- [Go Testing Documentation](https://golang.org/pkg/testing/)

## Contribuciones

Al agregar nuevos tests:
1. Escribe el feature primero (BDD)
2. Implementa los steps en Go
3. Actualiza esquemas si es necesario
4. Verifica que todos los tests pasen
5. Actualiza esta documentación

---

**Última actualización:** 2024
**Versión:** 1.0
