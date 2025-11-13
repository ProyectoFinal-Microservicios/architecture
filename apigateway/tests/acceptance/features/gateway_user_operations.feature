# language: es
Característica: Operaciones de Usuarios a través del API Gateway
  Como usuario autenticado
  Quiero poder obtener, actualizar y eliminar mi perfil a través del gateway
  Para mantener mis datos a través de una interfaz unificada

  Antecedentes:
    Dado que el gateway está disponible en "http://localhost:8000"
    Y existe un usuario autenticado con username "gatewayuser"

  Escenario: Obtener perfil de usuario a través del gateway
    Cuando hago una solicitud GET a "/api/v1/users/gatewayuser/profile" con token válido
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener user data
    Y la respuesta debe contener username "gatewayuser"
    Y la estructura de la respuesta debe ser válida según el esquema de perfil

  Escenario: Obtener perfil sin autenticación retorna 401
    Cuando hago una solicitud GET a "/api/v1/users/gatewayuser/profile" sin token
    Entonces la respuesta debe tener estado 401
    Y la respuesta debe contener error "Authorization header required"

  Escenario: Actualizar perfil a través del gateway
    Cuando hago una solicitud PATCH a "/api/v1/users/gatewayuser/profile" con datos:
      | field      | value                |
      | firstName  | UpdatedFirstName     |
      | lastName   | UpdatedLastName      |
      | phone      | +34987654321         |
    Y incluyo token válido
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener firstName "UpdatedFirstName"
    Y la estructura de la respuesta debe ser válida según el esquema de perfil

  Escenario: Actualizar perfil sin autenticación retorna 401
    Cuando hago una solicitud PATCH a "/api/v1/users/gatewayuser/profile" con datos:
      | field      | value            |
      | firstName  | NewName          |
    Y no incluyo token
    Entonces la respuesta debe tener estado 401

  Escenario: Eliminar usuario a través del gateway
    Dado que existe un usuario con username "userdelete"
    Cuando hago una solicitud DELETE a "/api/v1/users/userdelete" con token válido
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener message "Cuenta eliminada exitosamente"

  Escenario: Eliminar usuario sin autenticación retorna 401
    Cuando hago una solicitud DELETE a "/api/v1/users/gatewayuser" sin token
    Entonces la respuesta debe tener estado 401
    Y la respuesta debe contener error "Authorization header required"

  Escenario: El gateway enruta DELETE correctamente
    Dado que existe un usuario con username "testdeleteuser"
    Cuando hago una solicitud DELETE a "/api/v1/users/testdeleteuser" con token válido
    Entonces el servicio de autenticación debe haber recibido DELETE request
    Y la URL debe incluir el username correcto

  Escenario: El gateway enruta GET correctamente
    Cuando hago una solicitud GET a "/api/v1/users/gatewayuser/profile" con token válido
    Entonces el servicio de autenticación debe haber recibido GET request
    Y la URL debe incluir "/accounts/gatewayuser"
    Y la respuesta debe tener estado 200

  Escenario: El gateway enruta PATCH correctamente
    Cuando hago una solicitud PATCH a "/api/v1/users/gatewayuser/profile" con datos:
      | field | value    |
      | phone | +123456  |
    Y incluyo token válido
    Entonces el servicio de autenticación debe haber recibido PATCH request
    Y la URL debe incluir "/accounts/gatewayuser"
    Y la respuesta debe tener estado 200
