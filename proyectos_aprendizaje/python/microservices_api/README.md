# API de Microservicios con Python (FastAPI)

## Descripción

Este proyecto implementa una arquitectura de microservicios utilizando Python y FastAPI. La aplicación está compuesta por múltiples servicios independientes que se comunican entre sí, cada uno con su propia responsabilidad y base de datos.

## Características

- **Arquitectura de Microservicios**: Servicios independientes y desacoplados
- **API Gateway**: Punto de entrada único que enruta las solicitudes a los servicios correspondientes
- **Autenticación JWT**: Servicio de autenticación centralizado
- **Documentación Automática**: Cada servicio incluye documentación Swagger/OpenAPI
- **Comunicación entre Servicios**: Comunicación HTTP entre microservicios
- **Manejo de Errores**: Respuestas de error consistentes en todos los servicios
- **Logging**: Registro de solicitudes y respuestas

## Estructura del Proyecto

```
microservices_api/
├── requirements.txt           # Dependencias del proyecto
├── auth_service/             # Servicio de autenticación
│   └── main.py               # Implementación del servicio de autenticación
├── product_service/          # Servicio de productos
│   └── main.py               # Implementación del servicio de productos
├── order_service/            # Servicio de órdenes
│   └── main.py               # Implementación del servicio de órdenes
└── gateway/                  # API Gateway
    └── main.py               # Implementación del API Gateway
```

## Servicios Implementados

### 1. Servicio de Autenticación (Puerto 8000)

Maneja la autenticación de usuarios y la generación de tokens JWT.

**Endpoints:**
- `POST /token`: Genera un token JWT para un usuario válido
- `GET /users/me/`: Obtiene información del usuario actual
- `GET /health`: Verifica el estado del servicio

### 2. Servicio de Productos (Puerto 8001)

Gestiona el catálogo de productos.

**Endpoints:**
- `GET /products/`: Lista todos los productos
- `GET /products/{id}`: Obtiene un producto específico
- `POST /products/`: Crea un nuevo producto
- `PUT /products/{id}`: Actualiza un producto existente
- `DELETE /products/{id}`: Elimina un producto
- `GET /products/category/{category}`: Lista productos por categoría
- `GET /health`: Verifica el estado del servicio

### 3. Servicio de Órdenes (Puerto 8002)

Gestiona las órdenes de compra.

**Endpoints:**
- `GET /orders/`: Lista todas las órdenes
- `GET /orders/{id}`: Obtiene una orden específica
- `POST /orders/`: Crea una nueva orden
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

- Python 3.8+
- FastAPI
- Uvicorn
- Pydantic
- Python-jose (JWT)
- Passlib (hashing de contraseñas)
- Httpx (cliente HTTP asíncrono)

## Instalación

1. Instalar las dependencias:

```bash
pip install -r requirements.txt
```

## Ejecución

Para ejecutar todos los servicios, necesitarás abrir 4 terminales diferentes:

1. **Servicio de Autenticación:**

```bash
cd auth_service
uvicorn main:app --reload --port 8000
```

2. **Servicio de Productos:**

```bash
cd product_service
uvicorn main:app --reload --port 8001
```

3. **Servicio de Órdenes:**

```bash
cd order_service
uvicorn main:app --reload --port 8002
```

4. **API Gateway:**

```bash
cd gateway
uvicorn main:app --reload --port 8080
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
curl -X POST "http://localhost:8080/api/orders" -H "Content-Type: application/json" -H "Authorization: Bearer {token}" -d '{"customer_id":"user123", "items":[{"product_id":"1", "quantity":2}], "total":1999.98}'
```

2. Obtener órdenes de un cliente:

```bash
curl -X GET "http://localhost:8080/api/orders/customer/user123" -H "Authorization: Bearer {token}"
```

## Conceptos Aplicados

### Patrones de Diseño

1. **Patrón Gateway**: Implementado en el API Gateway para centralizar el acceso a los microservicios.
2. **Patrón Repository**: Utilizado para abstraer el acceso a datos en cada servicio.
3. **Patrón Dependency Injection**: Implementado a través del sistema de dependencias de FastAPI.

### Principios SOLID

1. **Principio de Responsabilidad Única (SRP)**: Cada servicio tiene una única responsabilidad.
2. **Principio de Abierto/Cerrado (OCP)**: Los servicios están abiertos para extensión pero cerrados para modificación.
3. **Principio de Sustitución de Liskov (LSP)**: Las implementaciones concretas pueden sustituir a las interfaces.
4. **Principio de Segregación de Interfaces (ISP)**: Cada servicio expone solo los métodos necesarios.
5. **Principio de Inversión de Dependencias (DIP)**: Dependencia de abstracciones, no de implementaciones concretas.

### Conceptos de Python Avanzados

1. **Programación Asíncrona**: Uso de `async/await` para operaciones no bloqueantes.
2. **Type Hints**: Uso de anotaciones de tipo para mejorar la legibilidad y el mantenimiento.
3. **Pydantic Models**: Validación de datos y serialización/deserialización.
4. **Dependency Injection**: Inyección de dependencias para facilitar las pruebas y la modularidad.

### Conceptos de Microservicios

1. **Desacoplamiento**: Servicios independientes con responsabilidades específicas.
2. **Escalabilidad**: Cada servicio puede escalar de forma independiente.
3. **Resiliencia**: Fallos en un servicio no afectan a otros servicios.
4. **API Gateway**: Punto de entrada único para todos los servicios.

## Problemas del Mundo Real que Resuelve

1. **Escalabilidad**: Permite escalar servicios individuales según la demanda.
2. **Mantenibilidad**: Facilita el mantenimiento y la evolución de cada servicio de forma independiente.
3. **Despliegue Continuo**: Permite desplegar servicios de forma independiente sin afectar a otros.
4. **Resiliencia**: Mejora la tolerancia a fallos al aislar los servicios.
5. **Equipos Autónomos**: Permite que diferentes equipos trabajen en diferentes servicios.

## Mejoras Potenciales

1. **Implementar Circuit Breaker**: Para manejar fallos en la comunicación entre servicios.
2. **Añadir Service Discovery**: Para descubrir dinámicamente la ubicación de los servicios.
3. **Implementar Caching**: Para mejorar el rendimiento.
4. **Añadir Monitorización**: Para supervisar el estado y el rendimiento de los servicios.
5. **Implementar Bases de Datos Reales**: Reemplazar las bases de datos simuladas por bases de datos reales.
6. **Añadir Tests**: Implementar pruebas unitarias, de integración y de carga.
7. **Implementar Contenedores**: Dockerizar los servicios para facilitar el despliegue.
8. **Añadir Mensajería Asíncrona**: Implementar comunicación asíncrona entre servicios mediante colas de mensajes.