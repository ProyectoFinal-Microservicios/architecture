# language: es
Característica: Autenticación a través del API Gateway
  Como usuario del sistema
  Quiero poder autenticarme a través del API Gateway
  Para obtener acceso a recursos protegidos

  Antecedentes:
    Dado que el gateway está disponible en "http://localhost:8000"

  Escenario: Login exitoso a través del gateway
    Cuando hago una solicitud POST a "/api/v1/auth/login" con datos:
      | field      | value                  |
      | identifier | admin                  |
      | password   | admin123               |
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener access_token
    Y la respuesta debe contener token_type "Bearer"
    Y la estructura de la respuesta debe ser válida según el esquema de login

  Escenario: Login fallido - credenciales inválidas
    Cuando hago una solicitud POST a "/api/v1/auth/login" con datos:
      | field      | value                  |
      | identifier | wronguser              |
      | password   | wrongpass              |
    Entonces la respuesta debe tener estado 401
    Y la respuesta debe contener error

  Escenario: Registro exitoso a través del gateway
    Cuando hago una solicitud POST a "/api/v1/auth/register" con datos:
      | field      | value                     |
      | username   | gatewayreguser            |
      | email      | anicu2314@gmail.com       |
      | password   | SecurePass123!            |
      | firstName  | Gateway                   |
      | lastName   | User                      |
      | phone      | +573234030048             |
    Entonces la respuesta debe tener estado 201
    Y la respuesta debe contener message "Usuario registrado exitosamente"
    Y la respuesta debe contener access_token
    Y la estructura de la respuesta debe ser válida según el esquema de registro

  Escenario: Registro fallido - usuario duplicado
    Dado que existe un usuario con username "existinguser"
    Cuando hago una solicitud POST a "/api/v1/auth/register" con datos:
      | field      | value                   |
      | username   | existinguser            |
      | email      | newuser@example.com     |
      | password   | SecurePass123!          |
      | firstName  | New                     |
      | lastName   | User                    |
      | phone      | +34912345678            |
    Entonces la respuesta debe tener estado 409
    Y la respuesta debe contener error

  Escenario: Request sin credentials retorna 401
    Cuando hago una solicitud GET a "/api/v1/users/testuser/profile" sin token
    Entonces la respuesta debe tener estado 401
    Y la respuesta debe contener error "Authorization header required"

  Escenario: Token inválido retorna 401
    Cuando hago una solicitud GET a "/api/v1/users/testuser/profile" con token "invalid.token.here"
    Entonces la respuesta debe tener estado 401
    Y la respuesta debe contener error

  Escenario: Token expirado retorna 401
    Dado que tengo un token expirado
    Cuando hago una solicitud GET a "/api/v1/users/testuser/profile" con token expirado
    Entonces la respuesta debe tener estado 401
    Y la respuesta debe contener error
