# Multi-stage build para Java/Spring Boot
FROM gradle:8.4-jdk17 AS build

WORKDIR /app

# Copiar el proyecto
COPY . .

# Compilar el proyecto (sin pruebas para rapidez)
RUN gradle clean build -x test --no-daemon

# Stage de runtime
FROM eclipse-temurin:21-jre-alpine

WORKDIR /app

# Crear usuario no-root
RUN addgroup -g 1001 -S orchestrator && \
    adduser -u 1001 -S orchestrator -G orchestrator

# Copiar el JAR desde el build anterior
COPY --from=build /app/build/libs/*.jar app.jar

# Cambiar propietario
RUN chown orchestrator:orchestrator app.jar

USER orchestrator

# Exponer puerto
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/actuator/health || exit 1

# Comando por defecto
CMD ["java", "-jar", "app.jar"]
