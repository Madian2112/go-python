FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Imagen final
FROM alpine:latest

WORKDIR /app

# Instalar certificados CA
RUN apk --no-cache add ca-certificates

# Copiar el binario compilado
COPY --from=builder /app/main .

# Exponer puerto
EXPOSE 8082

# Ejecutar la aplicación
CMD ["./main"]