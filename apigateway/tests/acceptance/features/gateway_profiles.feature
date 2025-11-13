# language: es
Característica: Gestión de perfiles a través del API Gateway
  Como usuario autenticado del sistema
  Quiero gestionar mi perfil y ver perfiles públicos a través del API Gateway
  Para mantener y consultar información extendida de usuarios

  Antecedentes:
    Dado que el gateway está disponible en "http://localhost:8888"
    Y que el servicio de profiles está disponible

  Escenario: Obtener mi perfil autenticado
    Dado que estoy autenticado como "testuser"
    Cuando hago una solicitud GET a "/api/v1/profiles/me"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener "user_id"
    Y la respuesta debe contener "username"
    Y la estructura de la respuesta debe ser válida según el esquema de perfil

  Escenario: Obtener mi perfil sin autenticación
    Cuando hago una solicitud GET a "/api/v1/profiles/me" sin token
    Entonces la respuesta debe tener estado 401
    Y la respuesta debe contener error

  Escenario: Actualizar mi perfil con datos válidos
    Dado que estoy autenticado como "testuser"
    Cuando hago una solicitud PUT a "/api/v1/profiles/me" con datos:
      | field             | value                                           |
      | nickname          | TestNick                                        |
      | bio               | Desarrollador full-stack especializado en Go    |
      | organization      | Tech Corp                                       |
      | country           | España                                          |
      | profile_visibility| public                                          |
      | github_url        | https://github.com/testuser                     |
      | linkedin_url      | https://linkedin.com/in/testuser                |
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener "message"
    Y la respuesta debe contener "profile"
    Y el perfil debe tener nickname "TestNick"
    Y el perfil debe tener bio "Desarrollador full-stack especializado en Go"

  Escenario: Actualizar mi perfil con URL inválida
    Dado que estoy autenticado como "testuser"
    Cuando hago una solicitud PUT a "/api/v1/profiles/me" con datos:
      | field      | value              |
      | github_url | invalid-url        |
    Entonces la respuesta debe tener estado 400 o 500
    Y la respuesta debe contener error

  Escenario: Buscar perfiles públicos sin filtros
    Cuando hago una solicitud GET a "/api/v1/profiles/search"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener "results"
    Y la respuesta debe contener "count"
    Y la respuesta debe contener "limit"
    Y la respuesta debe contener "offset"

  Escenario: Buscar perfiles públicos con término de búsqueda
    Cuando hago una solicitud GET a "/api/v1/profiles/search?q=desarrollador"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener "results"
    Y todos los resultados deben contener "desarrollador" en username, nickname o bio

  Escenario: Buscar perfiles públicos por país
    Cuando hago una solicitud GET a "/api/v1/profiles/search?country=España"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener "results"
    Y todos los resultados deben tener country "España"

  Escenario: Buscar perfiles con paginación
    Cuando hago una solicitud GET a "/api/v1/profiles/search?limit=5&offset=0"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe tener máximo 5 resultados
    Y el limit debe ser 5
    Y el offset debe ser 0

  Escenario: Obtener perfil público por username
    Dado que existe un perfil público con username "publicuser"
    Cuando hago una solicitud GET a "/api/v1/profiles/publicuser"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener "username" con valor "publicuser"
    Y la respuesta debe contener "profile_visibility" con valor "public"

  Escenario: Obtener perfil privado sin autenticación
    Dado que existe un perfil privado con username "privateuser"
    Cuando hago una solicitud GET a "/api/v1/profiles/privateuser" sin token
    Entonces la respuesta debe tener estado 404
    Y la respuesta debe contener error "Perfil no encontrado o privado"

  Escenario: Obtener mi propio perfil privado estando autenticado
    Dado que estoy autenticado como "privateuser"
    Y mi perfil tiene visibilidad "private"
    Cuando hago una solicitud GET a "/api/v1/profiles/privateuser"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener "username" con valor "privateuser"

  Escenario: Obtener perfil inexistente
    Cuando hago una solicitud GET a "/api/v1/profiles/nonexistentuser"
    Entonces la respuesta debe tener estado 404
    Y la respuesta debe contener error

  Escenario: Obtener estadísticas de mi perfil
    Dado que estoy autenticado como "testuser"
    Cuando hago una solicitud GET a "/api/v1/profiles/stats/me"
    Entonces la respuesta debe tener estado 200
    Y la respuesta debe contener "total_views"
    Y la respuesta debe contener "recent_activity"
    Y total_views debe ser un número

  Escenario: Obtener estadísticas sin autenticación
    Cuando hago una solicitud GET a "/api/v1/profiles/stats/me" sin token
    Entonces la respuesta debe tener estado 401
    Y la respuesta debe contener error

  Escenario: Actualizar perfil y verificar cambio en GET unificado
    Dado que estoy autenticado como "testuser"
    Cuando hago una solicitud PUT a "/api/v1/profiles/me" con datos:
      | field        | value                      |
      | bio          | Bio actualizada desde test |
      | organization | Nueva Empresa              |
    Y hago una solicitud GET a "/api/v1/users/testuser/profile"
    Entonces la respuesta debe tener estado 200
    Y el usuario debe tener bio "Bio actualizada desde test"
    Y el usuario debe tener organization "Nueva Empresa"

  Escenario: Cambiar visibilidad de perfil a privado
    Dado que estoy autenticado como "testuser"
    Cuando hago una solicitud PUT a "/api/v1/profiles/me" con datos:
      | field             | value   |
      | profile_visibility| private |
    Entonces la respuesta debe tener estado 200
    Y el perfil debe tener profile_visibility "private"

  Escenario: Actualizar URLs sociales del perfil
    Dado que estoy autenticado como "testuser"
    Cuando hago una solicitud PUT a "/api/v1/profiles/me" con datos:
      | field         | value                              |
      | github_url    | https://github.com/newtestuser     |
      | linkedin_url  | https://linkedin.com/in/newtestuser|
      | twitter_url   | https://twitter.com/newtestuser    |
      | website_url   | https://newtestuser.dev            |
    Entonces la respuesta debe tener estado 200
    Y el perfil debe tener github_url "https://github.com/newtestuser"
    Y el perfil debe tener linkedin_url "https://linkedin.com/in/newtestuser"
    Y el perfil debe tener twitter_url "https://twitter.com/newtestuser"
    Y el perfil debe tener website_url "https://newtestuser.dev"

  Escenario: Proxy correcto de errores del servicio profiles
    Dado que estoy autenticado como "testuser"
    Y el servicio de profiles está caído
    Cuando hago una solicitud GET a "/api/v1/profiles/me"
    Entonces la respuesta debe tener estado 503
    Y la respuesta debe contener error "Service unavailable"
