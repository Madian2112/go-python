# API REST de Microservicios en Go

Este proyecto implementa una API REST de microservicios utilizando Go, siguiendo las mejores prácticas de arquitectura de microservicios.

## Estructura del Proyecto

```
api_microservicios/
├── api-gateway/        # API Gateway que enruta las solicitudes a los microservicios
├── product-service/    # Servicio de productos
├── order-service/      # Servicio de pedidos
├── user-service/       # Servicio de usuarios
└── docker-compose.yml  # Configuración para ejecutar todos los servicios
```

## Tecnologías Utilizadas

- **Go**: Lenguaje de programación principal
- **Gin**: Framework web para crear APIs REST
- **gRPC**: Para comunicación entre servicios
- **PostgreSQL**: Base de datos relacional
- **Redis**: Para caché y mensajería
- **Docker**: Para contenerización
- **Prometheus/Grafana**: Para monitoreo
- **Jaeger**: Para tracing distribuido

## Características

- Arquitectura de microservicios
- API Gateway para enrutamiento y autenticación centralizada
- Comunicación síncrona (REST, gRPC) y asíncrona (mensajería)
- Persistencia de datos con PostgreSQL
- Caché con Redis
- Autenticación con JWT
- Documentación de API con Swagger
- Monitoreo y observabilidad
- Despliegue con Docker

## Requisitos

- Go 1.17 o superior
- Docker y Docker Compose
- PostgreSQL
- Redis

## Instalación y Ejecución

1. Clonar el repositorio
2. Ejecutar `docker-compose up` para iniciar todos los servicios
3. Acceder a la API a través de `http://localhost:8080`
4. Acceder a la documentación Swagger en `http://localhost:8080/swagger/index.html`

## Desarrollo

Cada microservicio sigue una estructura similar:

```
service/
├── cmd/                # Punto de entrada de la aplicación
├── internal/
│   ├── config/         # Configuración
│   ├── domain/         # Modelos de dominio
│   ├── repository/     # Acceso a datos
│   ├── service/        # Lógica de negocio
│   ├── transport/      # Controladores HTTP/gRPC
│   └── middleware/     # Middleware
├── pkg/                # Código compartido
└── Dockerfile          # Configuración de Docker
```

## Endpoints de la API

### API Gateway (puerto 8080)

- `GET /health`: Verificación de salud
- `POST /login`: Autenticación
- `GET /docs/*`: Documentación Swagger

### Servicio de Productos (puerto 8081)

- `GET /products`: Listar productos
- `GET /products/:id`: Obtener un producto
- `POST /products`: Crear un producto
- `PUT /products/:id`: Actualizar un producto
- `DELETE /products/:id`: Eliminar un producto

### Servicio de Pedidos (puerto 8082)

- `GET /orders`: Listar pedidos
- `GET /orders/:id`: Obtener un pedido
- `POST /orders`: Crear un pedido
- `PUT /orders/:id`: Actualizar un pedido
- `DELETE /orders/:id`: Cancelar un pedido

### Servicio de Usuarios (puerto 8083)

- `GET /users`: Listar usuarios
- `GET /users/:id`: Obtener un usuario
- `POST /users`: Registrar un usuario
- `PUT /users/:id`: Actualizar un usuario
- `DELETE /users/:id`: Eliminar un usuario

## Monitoreo y Observabilidad

- Métricas de Prometheus: `http://localhost:9090`
- Dashboard de Grafana: `http://localhost:3000`
- UI de Jaeger: `http://localhost:16686`

## Pruebas

Ejecutar pruebas unitarias:

```bash
go test ./...
```

Ejecutar pruebas de integración:

```bash
go test -tags=integration ./...
```

## Contribución

1. Hacer fork del repositorio
2. Crear una rama para la nueva característica
3. Implementar la característica
4. Enviar un pull request

## Licencia

MIT