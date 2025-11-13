# 🎯 API Gateway BDD Tests - Quick Reference Card

## ⚡ Comandos Más Usados

| Comando | Resultado |
|---------|-----------|
| `make test-acceptance` | Ejecutar todas las pruebas BDD |
| `make test-acceptance-pretty` | Pruebas en formato legible |
| `godog run -v tests/acceptance/features` | Pruebas con detalles |
| `make help` | Ver todos los comandos |
| `make clean` | Limpiar archivos generados |

---

## 📂 Archivos Importantes

```
tests/
├── features/             ← Donde escribes nuevos tests
├── steps/                ← Implementación de steps
├── schemas/              ← Validación JSON
├── support/              ← Cliente HTTP + Validator
├── README.md             ← Referencia general
├── QUICKSTART.md         ← 5 minutos para empezar
├── TESTING.md            ← Documentación completa
├── ARCHITECTURE.md       ← Diseño técnico
└── EXAMPLES.md           ← 9 ejemplos prácticos
```

---

## 🧪 Estructura de un Test

### 1. Feature (Gherkin - Español)
```gherkin
# language: es
Escenario: Descripción del test
  Cuando realizo una solicitud GET a "/endpoint"
  Entonces el código de estado debe ser 200
```

### 2. Schema (JSON-Schema)
```json
{
  "mySchema": {
    "type": "object",
    "properties": {"field": {"type": "string"}},
    "required": ["field"]
  }
}
```

### 3. Step (Go)
```go
func (ctx *APIGatewayContext) MiStep() error {
    // Tu código aquí
    return nil
}
```

---

## 🔑 Steps Principales

### Solicitudes HTTP
```gherkin
Cuando realizo una solicitud GET a "/path"
Cuando realizo una solicitud POST a "/path"
Cuando realizo una solicitud PATCH a "/path"
Cuando realizo una solicitud DELETE a "/path"
```

### Validaciones
```gherkin
Entonces el código de estado debe ser 200
Y el código de estado debe ser 401 o 404 o 400
Y el código de estado está entre 200 y 299
Y la respuesta debe contener el campo "status"
Y la respuesta debe validarse contra el esquema "schemaName"
Y la respuesta debe tener el header "Content-Type"
```

### Data Tables
```gherkin
| field    | value       |
| username | john_doe    |
| password | MyPass123   |
```

---

## 🎓 Ejemplo Completo

```gherkin
# language: es
Funcionalidad: Mi test

Antecedentes:
  Dado que el gateway está accesible en "http://localhost:8080"

Escenario: El gateway responde con estado 200
  Cuando realizo una solicitud GET a "/health"
  Entonces el código de estado debe ser 200
  Y la respuesta debe validarse contra el esquema "healthResponse"
```

---

## 🚀 Flujo de Inicio

1. **Leer QUICKSTART.md** (5 minutos)
   ```bash
   cat tests/QUICKSTART.md
   ```

2. **Instalar herramientas** (si es la primera vez)
   ```bash
   make install-tools
   ```

3. **Ejecutar tests**
   ```bash
   make test-acceptance
   ```

4. **Ver ejemplos prácticos**
   ```bash
   cat tests/EXAMPLES.md
   ```

5. **Escribir nuevos tests**
   - Copia un ejemplo
   - Crea nuevo .feature file
   - Implementa steps necesarios
   - Ejecuta y verifica

---

## 📊 Cobertura Actual

| Área | Escenarios | Estado |
|------|-----------|--------|
| Routing | 9 | ✅ Completo |
| Authentication | 8 | ✅ Completo |
| Users | 9 | ✅ Completo |
| **Total** | **26** | **✅ Completo** |

---

## 🔧 Herramientas Necesarias

```bash
# Go 1.21+
go version

# Godog CLI
godog --version

# Verificar dependencies
go mod tidy
go mod download
```

---

## 💾 Estructura de Carpetas

```
tests/
├── acceptance/
│   ├── features/          ← Archivos .feature
│   ├── steps/             ← gateway_steps.go
│   ├── schemas/           ← gateway-schemas.json
│   └── main_test.go
└── support/
    ├── http_client.go
    └── schema_validator.go
```

---

## 🎯 Patrones de Steps

### Setup/Precondición
```go
func (ctx *APIGatewayContext) GatewayIsAccessible(baseURL string) error {
    ctx.client = support.NewHTTPClient(baseURL)
    return nil
}
```

### Acción
```go
func (ctx *APIGatewayContext) MakeGetRequest(endpoint string) error {
    ctx.lastResponse, ctx.lastError = ctx.client.GET(endpoint, ctx.customHeaders)
    return ctx.lastError
}
```

### Aserción
```go
func (ctx *APIGatewayContext) StatusCodeShouldBe(expected int) error {
    if ctx.lastResponse.StatusCode != expected {
        return fmt.Errorf("expected %d but got %d", expected, ctx.lastResponse.StatusCode)
    }
    return nil
}
```

---

## 📌 Headers Comunes

```gherkin
Cuando incluyo el header "Authorization" con valor "Bearer TOKEN"
Cuando incluyo el header "Content-Type" con valor "application/json"
Cuando incluyo el header "X-Custom-Header" con valor "value"
```

---

## 🛠️ Debug

### Ver detalles de ejecución
```bash
godog run -v tests/acceptance/features
```

### Ejecutar escenario específico
```bash
godog run -n "Mi escenario" tests/acceptance/features
```

### Generar reporte JSON
```bash
godog run -f json tests/acceptance/features > report.json
```

---

## ⚠️ Errores Comunes

| Error | Solución |
|-------|----------|
| `step is undefined` | Implementar en `gateway_steps.go` |
| `connection refused` | Verificar gateway en `:8080` |
| `schema not found` | Agregar a `gateway-schemas.json` |
| `godog: command not found` | Ejecutar `make install-tools` |

---

## 📖 Documentación por Necesidad

**Necesito empezar rápido:**
→ Lee `tests/QUICKSTART.md`

**Necesito entender cómo funciona:**
→ Lee `tests/ARCHITECTURE.md`

**Necesito ver ejemplos:**
→ Lee `tests/EXAMPLES.md`

**Necesito documentación completa:**
→ Lee `tests/TESTING.md`

**Necesito referencia general:**
→ Lee `tests/README.md`

---

## 🎨 Formatos de Salida

```bash
# Progress (por defecto)
godog run tests/acceptance/features

# Pretty (legible)
godog run -f pretty tests/acceptance/features

# JSON (parseable)
godog run -f json tests/acceptance/features

# XML (Jenkins)
godog run -f junit tests/acceptance/features > results.xml
```

---

## 🔐 Autenticación

### Usar token JWT
```gherkin
Dado que he realizado login con usuario "alice@mail.com"
Cuando realizo una solicitud GET a "/users/alice"
```

### Establecer token manualmente
```go
ctx.client.SetToken("eyJhbGc...")
```

---

## 💡 Tips Útiles

✅ Usa data tables para datos complejos
✅ Reutiliza steps existentes
✅ Escribe features primero (BDD)
✅ Mantén features independientes
✅ Usa nombres descriptivos
✅ Agrupa tests con tags
✅ Documenta casos complejos

---

## 🚀 Próximos Pasos

1. ✅ Lee QUICKSTART
2. ✅ Ejecuta `make test-acceptance`
3. ✅ Revisa EXAMPLES
4. ✅ Escribe tu primer test
5. ✅ Integra en CI/CD

---

## 📞 Referencias Rápidas

| Recurso | URL/Comando |
|---------|-----------|
| Godog | https://github.com/cucumber/godog |
| Gherkin | https://cucumber.io/docs/gherkin |
| JSON-Schema | https://json-schema.org |
| Help | `make help` |

---

**¡Listo para empezar! 🎉**

Ejecuta: `make test-acceptance`
