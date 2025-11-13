# API Gateway - BDD Test Suite

## Descripción General

Suite completa de pruebas de aceptación basadas en BDD (Behavior-Driven Development) para el API Gateway usando:

- **Lenguaje:** Gherkin (español)
- **Framework:** Godog (Go + Cucumber)
- **Validación:** JSON-Schema
- **HTTP:** Cliente personalizado con soporte JWT

## Características

✅ **26 Escenarios** basados en comportamiento esperado
✅ **Gherkin en Español** - Feature files mantenibles por todo el equipo
✅ **Validación de Esquemas** - Asegura formato correcto de respuestas
✅ **Autenticación JWT** - Pruebas con tokens
✅ **Cobertura Completa** - Routing, autenticación, operaciones CRUD, manejo de errores

## Inicio Rápido

### 1. Instalación

```bash
# Descargar dependencias
make deps

# Instalar herramientas (si es la primera vez)
make install-tools
```

### 2. Ejecutar Pruebas

```bash
# Ejecución simple
make test-acceptance

# Formato más legible
make test-acceptance-pretty

# Con detalles verbosos
godog run -v tests/acceptance/features
```

### 3. Ver Resultados

```bash
# En formato JSON (para CI/CD)
godog run -f json tests/acceptance/features

# En formato XML (para Jenkins)
godog run -f junit tests/acceptance/features
```

## Archivos Clave

| Archivo | Propósito |
|---------|-----------|
| `features/gateway_routing.feature` | Tests de enrutamiento y middleware |
| `features/gateway_authentication.feature` | Tests de autenticación |
| `features/gateway_user_operations.feature` | Tests de operaciones de usuarios |
| `features/schemas/gateway-schemas.json` | Esquemas de validación |
| `steps/gateway_steps.go` | Implementación de steps |
| `support/http_client.go` | Cliente HTTP especializado |
| `support/schema_validator.go` | Validador de esquemas |

## Ejemplos de Uso

### Ejecutar un feature específico

```bash
godog run tests/acceptance/features/gateway_routing.feature
```

### Ejecutar escenarios con un tag

```bash
godog run --tags @authentication tests/acceptance/features
```

### Ejecutar en paralelo

```bash
godog run -c 4 tests/acceptance/features
```

## Estructura de un Test BDD

### Feature File (Gherkin)
```gherkin
# language: es
Funcionalidad: Enrutamiento del API Gateway
  
  Antecedentes:
    Dado que el gateway está accesible en "http://localhost:8080"
  
  Escenario: El gateway responde con estado 200 en /health
    Cuando realizo una solicitud GET a "/health"
    Entonces el código de estado debe ser 200
    Y la respuesta debe validarse contra el esquema "healthResponse"
```

### Step Definition (Go)
```go
func (ctx *APIGatewayContext) MakeGetRequest(endpoint string) error {
    response, err := ctx.client.GET(endpoint, ctx.customHeaders)
    ctx.lastResponse = response
    return err
}
```

## Data Tables (Tablas de Datos)

Los tests usan tablas para pasar datos estructurados:

```gherkin
Cuando realizo una solicitud POST a "/auth/login"
| field    | value       |
| username | john_doe    |
| password | SecurePass1 |
```

**Importante:** Los headers siempre están en inglés, pero los datos pueden estar en cualquier idioma.

## Validaciones Disponibles

### Código de Estado
```gherkin
Entonces el código de estado debe ser 200
Y el código de estado debe ser 401 o 404 o 400
Y el código de estado está entre 200 y 299
```

### Contenido de Respuesta
```gherkin
Y la respuesta debe validarse contra el esquema "healthResponse"
Y la respuesta debe contener el campo "status"
Y la respuesta debe tener el header "Content-Type"
```

### Headers
```gherkin
Y la respuesta debe tener un header con patrón "X-Request-ID|request-id"
```

## Schemas (Esquemas JSON)

Los esquemas validan que las respuestas tengan la estructura correcta:

```json
{
  "healthResponse": {
    "type": "object",
    "properties": {
      "status": {"type": "string"},
      "timestamp": {"type": "string"}
    },
    "required": ["status"],
    "additionalProperties": false
  }
}
```

## Contexto de Prueba

Cada test mantiene un contexto compartido:

```go
type APIGatewayContext struct {
    client           *support.HTTPClient        // Cliente HTTP
    validator        *support.SchemaValidator   // Validador
    lastResponse     *support.Response          // Última respuesta
    lastError        error                      // Último error
    customHeaders    map[string]string          // Headers personalizados
    requestDataTable map[string]interface{}     // Datos de tabla
}
```

## Integración con CI/CD

### GitHub Actions
```yaml
- name: Run BDD Tests
  run: godog run tests/acceptance/features
```

### Jenkins
```groovy
stage('BDD Tests') {
    steps {
        sh 'godog run -f junit tests/acceptance/features > test-results.xml'
    }
}
```

### GitLab CI
```yaml
test:
  script:
    - godog run tests/acceptance/features
```

## Troubleshooting

### El gateway no responde
1. Verifica que el gateway está ejecutándose: `http://localhost:8080`
2. Comprueba los logs del gateway
3. Ajusta la URL base en `gateway_steps.go`

### Tests fallan por validación de schema
1. Revisa la respuesta real comparándola con el schema
2. Actualiza `gateway-schemas.json` si es necesario
3. Usa `godog run -v` para ver detalles

### Godog no se encuentra
```bash
go install github.com/cucumber/godog/cmd/godog@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

## Comandos Disponibles

```bash
make test                   # Ejecutar todas las pruebas
make test-acceptance        # Solo pruebas BDD
make test-acceptance-pretty # BDD con formato legible
make test-unit              # Pruebas unitarias
make test-coverage          # Pruebas con cobertura
make fmt                    # Formatear código
make lint                   # Análisis estático
make clean                  # Limpiar archivos generados
make deps                   # Descargar dependencias
```

## Cobertura de Pruebas

| Área | Escenarios | Estado |
|------|-----------|--------|
| Enrutamiento | 9 | ✅ Completo |
| Autenticación | 8 | ✅ Completo |
| Usuarios | 9 | ✅ Completo |
| **Total** | **26** | **✅ Completo** |

## Próximos Pasos

- [ ] Agregar tests de rendimiento
- [ ] Integrar con monitoring
- [ ] Agregar tests de seguridad
- [ ] Documentar casos edge
- [ ] Implementar retry logic

## Referencias

- [Godog GitHub](https://github.com/cucumber/godog)
- [Gherkin Specification](https://cucumber.io/docs/gherkin/)
- [JSON-Schema](https://json-schema.org/)
- [BDD Best Practices](https://cucumber.io/docs/bdd/)

---

**¿Necesitas ayuda?**
- Consulta `TESTING.md` para documentación completa
- Revisa los ejemplos en los archivos `.feature`
- Ejecuta `make help` para ver todos los comandos
