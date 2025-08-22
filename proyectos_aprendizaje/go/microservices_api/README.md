# API de Microservicios con Go (Gin)

## Descripción

Este proyecto implementa una arquitectura de microservicios utilizando Go y el framework Gin. La aplicación está compuesta por múltiples servicios independientes que se comunican entre sí, cada uno con su propia responsabilidad y base de datos.

## Características

- **Arquitectura de Microservicios**: Servicios independientes y desacoplados
- **API Gateway**: Punto de entrada único que enruta las solicitudes a los servicios correspondientes
- **Autenticación JWT**: Servicio de autenticación centralizado
- **Persistencia de Datos**: Almacenamiento en archivos JSON (simulando bases de datos)
- **Comunicación entre Servicios**: Comunicación HTTP entre microservicios
- **Manejo de Errores**: Respuestas de error consistentes en todos los servicios
- **Logging**: Registro de solicitudes y respuestas

## Estructura del Proyecto

```
microservices_api/
├── go.mod                  # Dependencias del proyecto
├── auth_service/           # Servicio de autenticación
│   └── main.go             # Implementación del servicio de autenticación
├── product_service/        # Servicio de productos
│   └── main.go             # Implementación del servicio de productos
├── order_service/          # Servicio de órdenes
│   └── main.go             # Implementación del servicio de órdenes
└── gateway/                # API Gateway
    └── main.go             # Implementación del API Gateway
```

## Servicios Implementados

### 1. Servicio de Autenticación (Puerto 8000)

Maneja la autenticación de usuarios y la generación de tokens JWT.

**Endpoints:**
- `POST /token`: Genera un token JWT para un usuario válido
- `GET /users/me`: Obtiene información del usuario actual
- `GET /health`: Verifica el estado del servicio

### 2. Servicio de Productos (Puerto 8001)

Gestiona el catálogo de productos.

**Endpoints:**
- `GET /products`: Lista todos los productos
- `GET /products/{id}`: Obtiene un producto específico
- `POST /products`: Crea un nuevo producto
- `PUT /products/{id}`: Actualiza un producto existente
- `DELETE /products/{id}`: Elimina un producto
- `GET /products/category/{category}`: Lista productos por categoría
- `GET /health`: Verifica el estado del servicio

### 3. Servicio de Órdenes (Puerto 8002)

Gestiona las órdenes de compra.

**Endpoints:**
- `GET /orders`: Lista todas las órdenes
- `GET /orders/{id}`: Obtiene una orden específica
- `POST /orders`: Crea una nueva orden
- `PUT /orders/{id}`: Actualiza una orden existente
- `DELETE /orders/{id}`: Elimina una orden
- `GET /orders/customer/{customer_id}`: Lista órdenes por cliente
- `GET /health`: Verifica el estado del servicio

### 4. API Gateway (Puerto 8080)

Actúa como punto de entrada único para todos los servicios.

**Endpoints:**
- `GET /health`: Verifica el estado de todos los servicios
- `POST /api/auth/token`: Redirige al servicio de autenticación
- `GET /api/auth/users/me`: Redirige al servicio de autenticación
- `GET /api/products`: Redirige al servicio de productos
- `GET /api/products/{id}`: Redirige al servicio de productos
- `POST /api/products`: Redirige al servicio de productos
- `PUT /api/products/{id}`: Redirige al servicio de productos
- `DELETE /api/products/{id}`: Redirige al servicio de productos
- `GET /api/products/category/{category}`: Redirige al servicio de productos
- `GET /api/orders`: Redirige al servicio de órdenes
- `GET /api/orders/{id}`: Redirige al servicio de órdenes
- `POST /api/orders`: Redirige al servicio de órdenes
- `PUT /api/orders/{id}`: Redirige al servicio de órdenes
- `DELETE /api/orders/{id}`: Redirige al servicio de órdenes
- `GET /api/orders/customer/{customer_id}`: Redirige al servicio de órdenes

## Requisitos

- Go 1.18+
- Gin Framework
- JWT-Go
- UUID

## Compilación y Ejecución

1. Compilar los servicios:

```bash
# Compilar el servicio de autenticación
cd auth_service
go build -o auth_service

# Compilar el servicio de productos
cd ../product_service
go build -o product_service

# Compilar el servicio de órdenes
cd ../order_service
go build -o order_service

# Compilar el API Gateway
cd ../gateway
go build -o gateway
```

2. Ejecutar los servicios:

Para ejecutar todos los servicios, necesitarás abrir 4 terminales diferentes:

```bash
# Terminal 1: Servicio de Autenticación
cd auth_service
./auth_service

# Terminal 2: Servicio de Productos
cd product_service
./product_service

# Terminal 3: Servicio de Órdenes
cd order_service
./order_service

# Terminal 4: API Gateway
cd gateway
./gateway
```

Alternativamente, puedes ejecutar cada servicio directamente con `go run`:

```bash
# Terminal 1: Servicio de Autenticación
cd auth_service
go run main.go

# Terminal 2: Servicio de Productos
cd product_service
go run main.go

# Terminal 3: Servicio de Órdenes
cd order_service
go run main.go

# Terminal 4: API Gateway
cd gateway
go run main.go
```

## Uso

### Autenticación

1. Obtener un token JWT:

```bash
curl -X POST "http://localhost:8080/api/auth/token" -H "Content-Type: application/json" -d '{"username":"testuser", "password":"password123"}'
```

2. Usar el token para acceder a recursos protegidos:

```bash
curl -X GET "http://localhost:8080/api/products" -H "Authorization: Bearer {token}"
```

### Productos

1. Crear un producto:

```bash
curl -X POST "http://localhost:8080/api/products" -H "Content-Type: application/json" -H "Authorization: Bearer {token}" -d '{"name":"Laptop", "description":"Potente laptop para desarrollo", "price":999.99, "category":"Electronics", "stock":10}'
```

2. Obtener todos los productos:

```bash
curl -X GET "http://localhost:8080/api/products"
```

### Órdenes

1. Crear una orden:

```bash
curl -X POST "http://localhost:8080/api/orders" -H "Content-Type: application/json" -H "Authorization: Bearer {token}" -d '{"customer_id":"user123", "items":[{"product_id":"1", "quantity":2}]}'
```

2. Obtener órdenes de un cliente:

```bash
curl -X GET "http://localhost:8080/api/orders/customer/user123" -H "Authorization: Bearer {token}"
```

## Conceptos de Go Aplicados

1. **Goroutines y Concurrencia**: Manejo de múltiples solicitudes concurrentes
2. **Interfaces**: Uso de interfaces para abstraer comportamientos
3. **Structs y Métodos**: Modelado de datos y comportamiento
4. **JSON Marshaling/Unmarshaling**: Serialización y deserialización de datos
5. **Manejo de Errores**: Propagación y manejo de errores
6. **Middleware**: Uso de middleware para funcionalidades transversales
7. **Punteros**: Uso eficiente de memoria y paso por referencia

## Patrones de Diseño

1. **Patrón Gateway**: Implementado en el API Gateway para centralizar el acceso a los microservicios
2. **Patrón Middleware**: Utilizado para funcionalidades transversales como autenticación y logging
3. **Patrón Repository**: Simulado para abstraer el acceso a datos
4. **Patrón Singleton**: Utilizado para instancias únicas como el router

## Problemas del Mundo Real que Resuelve

1. **Escalabilidad**: Permite escalar servicios individuales según la demanda
2. **Mantenibilidad**: Facilita el mantenimiento y la evolución de cada servicio de forma independiente
3. **Despliegue Continuo**: Permite desplegar servicios de forma independiente sin afectar a otros
4. **Resiliencia**: Mejora la tolerancia a fallos al aislar los servicios
5. **Equipos Autónomos**: Permite que diferentes equipos trabajen en diferentes servicios

## Diferencias con la Implementación en Python

1. **Tipado Estático vs Dinámico**: Go utiliza tipado estático, lo que proporciona mayor seguridad en tiempo de compilación
2. **Concurrencia**: Go tiene soporte nativo para concurrencia con goroutines y channels
3. **Rendimiento**: Go generalmente ofrece mejor rendimiento y menor consumo de recursos
4. **Compilación vs Interpretación**: Go es compilado, lo que elimina muchos errores en tiempo de ejecución
5. **Manejo de Errores**: Go utiliza valores de retorno múltiples para el manejo de errores, en lugar de excepciones

## Mejoras Potenciales

1. **Implementar Circuit Breaker**: Para manejar fallos en la comunicación entre servicios
2. **Añadir Service Discovery**: Para descubrir dinámicamente la ubicación de los servicios
3. **Implementar Caching**: Para mejorar el rendimiento
4. **Añadir Monitorización**: Para supervisar el estado y el rendimiento de los servicios
5. **Implementar Bases de Datos Reales**: Reemplazar las bases de datos simuladas por bases de datos reales
6. **Añadir Tests**: Implementar pruebas unitarias, de integración y de carga
7. **Implementar Contenedores**: Dockerizar los servicios para facilitar el despliegue
8. **Añadir Mensajería Asíncrona**: Implementar comunicación asíncrona entre servicios mediante colas de mensajes