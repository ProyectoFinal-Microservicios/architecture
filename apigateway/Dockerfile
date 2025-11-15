FROM golang:1.21-alpine3.19 AS builder

WORKDIR /app

# Instalar dependencias
# RUN apk add --no-cache ca-certificates && update-ca-certificates && apk add --no-cache git
RUN apk add --no-cache git

# Copiar go.mod y go.sum
COPY go.mod go.sum ./
RUN go mod download


# Copiar código fuente
COPY . .

# Compilar
RUN go get github.com/gorilla/mux
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./app/main.go

# Stage final
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binario
COPY --from=builder /app/gateway .

# Copiar documentación
COPY --from=builder /app/docs ./docs

EXPOSE 8000

CMD ["./gateway"]