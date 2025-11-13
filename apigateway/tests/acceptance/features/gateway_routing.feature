# language: es
Característica: Enrutamiento de API Gateway
  Como cliente del sistema
  Quiero que el API Gateway enrute correctamente mis requests a los servicios upstream
  Para acceder a través de una interfaz unificada

  Antecedentes:
    Dado que el gateway está disponible en "http://localhost:8000"
    Y que el servicio de autenticación está disponible

  Escenario: Verificar health check del gateway
    Cuando hago una solicitud GET a "/health"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener status "UP"
    Y la respuesta debe contener service "api-gateway"
    Y la estructura de la respuesta debe ser válida según el esquema de health

  Escenario: Acceder a documentación del gateway
    Cuando hago una solicitud GET a "/docs"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener HTML
    Y el contenido debe incluir "API Gateway"

  Escenario: Obtener especificación OpenAPI
    Cuando hago una solicitud GET a "/docs/openapi.yaml"
    Entonces la respuesta debe tener estado 200
    Y el header Content-Type debe ser "application/yaml"

  Escenario: Enrutar login al servicio de autenticación
    Cuando hago una solicitud POST a "/api/v1/auth/login" con datos:
      | field      | value                  |
      | identifier | testuser               |
      | password   | SecurePass123!         |
    Entonces la respuesta debe tener estado 200 o 401
    Y la respuesta debe ser JSON válido

  Escenario: Enrutar registro al servicio de autenticación
    Cuando hago una solicitud POST a "/api/v1/auth/register" con datos:
      | field      | value                  |
      | username   | gatewayuser            |
      | email      | gateway@example.com    |
      | password   | SecurePass123!         |
      | firstName  | Gateway                |
      | lastName   | Test                   |
      | phone      | +34912345678           |
    Entonces la respuesta debe tener estado 201 o 409
    Y la respuesta debe ser JSON válido

  Escenario: El gateway agrega header de autorización
    Cuando hago una solicitud GET a "/api/v1/users/testuser/profile" con token válido
    Entonces la solicitud debe haber incluido el header Authorization
    Y la respuesta debe tener estado 200 o 401

  Escenario: CORS headers están presentes
    Cuando hago una solicitud GET a "/health"
    Entonces el header Access-Control-Allow-Origin debe existir
    Y el header Access-Control-Allow-Methods debe incluir GET

  Escenario: Middleware de logging está activo
    Cuando hago una solicitud GET a "/health"
    Entonces debe haber un log de entrada de request
    Y debe haber un log de salida de request

  Escenario: Manejo de errores de servicio upstream
    Cuando hago una solicitud GET a "/api/v1/users/nonexistent/profile" con token inválido
    Entonces la respuesta debe tener estado 401 o 404
    Y la respuesta debe contener error

  Escenario: Timeout a servicio unavailable
    Dado que el servicio de autenticación no está disponible
    Cuando hago una solicitud POST a "/api/v1/auth/login" con datos:
      | field      | value          |
      | identifier | test           |
      | password   | test123        |
    Entonces la respuesta debe tener estado 503
    Y la respuesta debe contener "Service unavailable"
