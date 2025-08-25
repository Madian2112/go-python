# 17. Arquitecturas de Software en Proyectos Go

## Introducción

Go fue diseñado para construir sistemas simples, fiables y eficientes. Su ecosistema favorece arquitecturas que promueven la claridad, el bajo acoplamiento y la escalabilidad, especialmente en el ámbito de los microservicios y las aplicaciones de red.

A continuación, exploraremos las arquitecturas y estructuras de proyecto más comunes en el mundo de Go.

---

## 1. Standard Go Project Layout (Disposición Estándar de Proyectos de Go)

Más que una arquitectura de software, es una convención de estructura de directorios muy popular para aplicaciones complejas. No es un estándar oficial, pero es ampliamente adoptado.

- **Explicación:** Propone una organización de carpetas para separar el código de la aplicación, el código público, los scripts, la configuración, etc.
- **Cuándo usarla:** Es un excelente punto de partida para cualquier proyecto de Go que vaya más allá de un simple script o herramienta. No es necesario usar todos los directorios, solo los que necesites.
- **Estructura de Directorios Clave:**
  ```
  /mi_proyecto
  ├── /cmd                # Puntos de entrada de la aplicación
  │   └── /mi-app/        # Directorio para cada binario (ej. servidor web)
  │       └── main.go
  ├── /internal           # Código privado de la aplicación (no importable por otros)
  │   ├── /api            # Handlers de la API (ej. HTTP)
  │   ├── /repository     # Lógica de acceso a datos
  │   └── /service        # Lógica de negocio
  ├── /pkg                # Código público que puede ser importado por otros proyectos
  ├── /api                # Especificaciones de la API (ej. OpenAPI, Protobuf)
  ├── /configs            # Archivos de configuración
  ├── /scripts            # Scripts para automatización (build, deploy)
  ├── go.mod              # Módulo de Go
  └── main.go             # (Para proyectos muy simples, en lugar de /cmd)
  ```
- **Recursos:**
  - **Repositorio de Referencia:** [golang-standards/project-layout en GitHub](https://github.com/golang-standards/project-layout)

---

## 2. Arquitectura por Capas (Layered Architecture)

Un enfoque simple y efectivo para separar responsabilidades. Es una base común para arquitecturas más complejas como la Limpia o la Hexagonal.

- **Explicación:** El código se organiza en capas horizontales, donde cada capa solo puede comunicarse con la capa directamente inferior.
  - **Capa de Presentación/API:** Recibe las peticiones (ej. handlers HTTP).
  - **Capa de Servicio/Lógica de Negocio:** Orquesta las operaciones.
  - **Capa de Repositorio/Acceso a Datos:** Interactúa con la base de datos.
- **Cuándo usarla:** Es ideal para la mayoría de las aplicaciones monolíticas de tamaño pequeño a mediano. Proporciona una buena separación de responsabilidades sin ser demasiado compleja.
- **Estructura de Directorios (dentro de `/internal`):**
  ```
  /internal
  ├── /server             # Capa de presentación
  │   ├── handler.go      # Define los handlers HTTP
  │   └── routes.go       # Registra las rutas
  ├── /service            # Capa de negocio
  │   └── product.go      # Lógica de negocio para productos
  └── /repository         # Capa de datos
      └── product_db.go   # Implementación del repositorio con la BBDD
  ```
- **Recursos:**
  - [Building Go Apps with a Layered Architecture](https://www.cosmicdevelopment.com/blog/building-go-apps-with-a-layered-architecture-and-ddd)

---

## 3. Arquitectura Limpia (Clean Architecture)

Al igual que en otros lenguajes, su objetivo en Go es producir sistemas con bajo acoplamiento, alta cohesión y total independencia de la infraestructura.

- **Explicación:** Se basa en la regla de dependencia: las capas externas dependen de las internas, nunca al revés. El dominio (entidades y casos de uso) no sabe nada sobre la base de datos o la API web. La comunicación entre capas se logra mediante interfaces definidas en las capas internas.
- **Cuándo usarla:** Para aplicaciones complejas y de larga duración donde la mantenibilidad y la capacidad de prueba son críticas. Cuando se prevé que la base de datos o el framework web puedan cambiar en el futuro.
- **Estructura de Directorios (conceptual):**
  ```
  /mi_proyecto
  ├── /domain             # Entidades puras y sus interfaces
  │   ├── user.go
  │   └── user_repository.go (interface)
  ├── /usecase            # Lógica de negocio que implementa los casos de uso
  │   └── user_interactor.go
  ├── /interfaces         # Adaptadores (implementaciones de BBDD, etc.)
  │   ├── /repository
  │   │   └── user_repository_mysql.go (implementation)
  │   └── /controller
  │       └── user_controller.go (handlers)
  └── /main.go            # Inyección de dependencias y arranque
  ```
- **Recursos:**
  - [Ejemplo de Arquitectura Limpia en Go por Manuel Zapf](https://github.com/manuelzapf/go-clean-architecture)
  - [Clean Architecture for Go developers](https://medium.com/@hugo.nobre.leite/clean-architecture-for-go-developers-a-practical-approach-1-3-7b799981a8b)

---

## 4. Arquitectura de Microservicios

Go es uno de los lenguajes preferidos para construir microservicios debido a su rendimiento, concurrencia nativa, bajo consumo de memoria y facilidad de despliegue (binarios estáticos).

- **Explicación:** La aplicación se divide en servicios pequeños y autónomos que se comunican a través de la red, usualmente con gRPC/Protobuf o una API REST.
- **Cuándo usarla:** En sistemas grandes y complejos que necesitan escalar de forma independiente, o cuando diferentes equipos trabajan en diferentes áreas de negocio.
- **Estructura de Directorios (Monorepo con gRPC):**
  ```
  /mi_plataforma
  ├── /api                # Archivos .proto para la definición de la API gRPC
  │   └── /proto
  │       └── /v1
  │           └── user.proto
  ├── /gen                # Código Go generado a partir de los .proto
  │   └── /go
  │       └── /v1
  │           └── user.pb.go
  ├── /services           # Un directorio por cada microservicio
  │   ├── /user-service
  │   │   └── main.go
  │   └── /order-service
  │       └── main.go
  └── /pkg                # Código compartido entre servicios
  ```
- **Recursos:**
  - [gRPC.io](https://grpc.io/) (Framework RPC de Google, muy usado en Go)
  - [Go Kit](https://gokit.io/) (Un toolkit para construir microservicios en Go)
  - [Building Microservices with Go](https://www.oreilly.com/library/view/building-microservices-with/9781491986441/) (Libro)
