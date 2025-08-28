# 18. Guía Completa: API RESTful en Go y Conexión a Bases de Datos

## Introducción

Esta guía es un tutorial integral que te llevará desde un proyecto vacío hasta una API RESTful robusta, segura y lista para producción en Go. A diferencia de otras lecciones que se enfocan en un solo aspecto, aquí uniremos todo: arquitectura, frameworks, operaciones CRUD, autenticación y, fundamentalmente, la **interacción con diferentes tipos de bases de datos y distintas filosofías de conexión**.

**Proyecto de Ejemplo:** Construiremos una API para gestionar un inventario de productos (`/products`).

---

## 1. Arquitectura y Configuración del Proyecto

Antes de escribir una sola línea de código de la API, definimos una estructura sólida que nos permita crecer y mantener el proyecto a largo plazo.

### 1.1. Estructura de Directorios (Standard Go Project Layout)

Usaremos la disposición estándar como base para organizar nuestro código de manera lógica.

```
/inventario-api
├── /cmd
│   └── /server
│       └── main.go         # Punto de entrada y arranque del servidor
├── /internal
│   ├── /api                # Handlers HTTP, middlewares, enrutamiento
│   ├── /service            # Lógica de negocio (Casos de Uso)
│   └── /repository         # Acceso a datos (Nuestra capa de BBDD)
├── /pkg                    # (Opcional) Código seguro para compartir
├── /configs                # Archivos de configuración (ej. config.yaml)
├── go.mod
└── go.sum
```

- **/cmd/server:** Contiene la función `main`. Su única responsabilidad es "ensamblar" la aplicación: inicializar la base de datos, inyectar las dependencias (repositorio -> servicio -> handler) y arrancar el servidor.
- **/internal/api:** La capa más externa. Maneja las peticiones HTTP, valida los datos de entrada, y llama a la capa de servicio. No sabe nada de la base de datos.
- **/internal/service:** El cerebro de la aplicación. Contiene la lógica de negocio pura. Orquesta las operaciones y depende de una *interfaz* de repositorio, no de una implementación concreta.
- **/internal/repository:** La capa de datos. Implementa la lógica para comunicarse con la base de datos. Aquí es donde pondremos todo el código relacionado con SQL, GORM, Mongo, etc.

### 1.2. El Patrón de Repositorio

Para separar la lógica de negocio de la base de datos, usaremos el Patrón de Repositorio. Definimos una interfaz en nuestra capa de servicio que dicta "qué" se puede hacer con los datos, y luego creamos implementaciones concretas en la capa de repositorio que definen "cómo" se hace.

**Ejemplo de Interfaz:**
```go
// internal/service/product_service.go
package service

type ProductRepository interface {
    GetByID(id string) (*Product, error)
    Create(product *Product) error
    // ... otros métodos CRUD
}
```

---

## 2. La Capa de Persistencia: Conexión a Bases de Datos

Esta es la sección central. Aquí exploraremos cómo conectar nuestra API a diferentes bases de datos usando varias técnicas.

### 2.1. Fundamentos: El paquete `database/sql`

Go provee una interfaz estándar para bases de datos SQL. No es un driver, sino una capa de abstracción. Para cada base de datos, necesitarás un driver específico que "implemente" esta interfaz.

- **Recurso Principal:** [Documentación de `database/sql`](https://go.dev/doc/database/sql-tutorial)

### 2.2. Bases de Datos Relacionales (SQL)

#### **PostgreSQL / MySQL**

**Método 1: `database/sql` + Driver (El enfoque "puro")**

Es el equivalente a usar ADO.NET con `SqlCommand` en C#. Tienes control total sobre el SQL, pero es más verboso.

- **Implementación (`internal/repository/product_postgres.go`):**
  ```go
  import (
      "database/sql"
      _ "github.com/jackc/pgx/v5/stdlib" // Driver de PostgreSQL
  )

  type postgresRepo struct {
      db *sql.DB
  }

  func (r *postgresRepo) GetByID(id string) (*service.Product, error) {
      // Escribir la consulta SQL cruda
      // Escanear los resultados en los campos del struct manualmente
  }
  ```
- **Cuándo usarlo:** Cuando necesitas optimizar consultas SQL al máximo o quieres cero dependencias externas (aparte del driver).
- **Recursos:**
  - **Driver `pgx` para PostgreSQL:** [github.com/jackc/pgx](https://github.com/jackc/pgx)
  - **Driver para MySQL:** [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)

**Método 2: `sqlx` (El ayudante)**

`sqlx` es una extensión ligera sobre `database/sql` que simplifica enormemente el escaneo de resultados a structs y el manejo de sentencias `IN`.

- **Implementación:**
  ```go
  import "github.com/jmoiron/sqlx"

  // ...
  func (r *postgresRepo) GetByID(id string) (*service.Product, error) {
      // La consulta es la misma, pero el escaneo es automático
      err := r.db.Get(&product, "SELECT * FROM products WHERE id=$1", id)
  }
  ```
- **Cuándo usarlo:** En la mayoría de los casos. Es el punto medio perfecto entre control y conveniencia.
- **Recursos:**
  - **Repositorio `sqlx`:** [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx)

**Método 3: `GORM` (El ORM completo)**

Es el equivalente a Entity Framework. Abstrae completamente el SQL y te permite trabajar con objetos de Go.

- **Implementación:**
  ```go
  import "gorm.io/gorm"

  // ...
  func (r *gormRepo) GetByID(id string) (*service.Product, error) {
      var product ProductModel
      result := r.db.First(&product, "id = ?", id)
      // Convertir ProductModel a service.Product
      return &product, result.Error
  }
  ```
- **Cuándo usarlo:** Para prototipado rápido, CRUDs sencillos, o si el equipo prefiere la abstracción de un ORM.
- **Contras:** Puede generar consultas ineficientes si no se usa con cuidado y añade una capa de "magia".
- **Recursos:**
  - **Sitio oficial de GORM:** [gorm.io](https://gorm.io/)

### 2.3. Bases de Datos NoSQL

#### **MongoDB**

La interacción con MongoDB es completamente diferente, ya que es una base de datos de documentos (BSON) y no relacional.

- **Método: Driver Oficial de MongoDB**
  No hay un equivalente a `database/sql`. Se usa directamente el driver oficial.
- **Implementación (`internal/repository/product_mongo.go`):**
  ```go
  import "go.mongodb.org/mongo-driver/mongo"

  // ...
  func (r *mongoRepo) GetByID(id string) (*service.Product, error) {
      var product BSONProduct
      // Usar filter, collection.FindOne(), etc.
      // Trabajar con BSON y contextos
  }
  ```
- **Cuándo usarlo:** Cuando tu modelo de datos es flexible, basado en documentos, y necesitas escalar horizontalmente.
- **Recursos:**
  - **Driver Oficial de Go para MongoDB:** [mongodb.com/docs/drivers/go/](https://www.mongodb.com/docs/drivers/go/)

---

## 3. Implementación de la API (CRUD Completo)

Aquí uniremos todo. Usaremos **Gin** como framework por su popularidad y rendimiento, y **sqlx** para la conexión a PostgreSQL.

### 3.1. Handlers (`/internal/api/handlers.go`)
El handler es responsable de manejar el contexto HTTP, decodificar el JSON de entrada, llamar al servicio y devolver una respuesta HTTP.

```go
package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "tu-proyecto/internal/service"
)

type ProductHandler struct {
    service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
    return &ProductHandler{service: s}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var req service.CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product, err := h.service.Create(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
    id := c.Param("id")
    product, err := h.service.GetByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
        return
    }
    c.JSON(http.StatusOK, product)
}

// ... Aquí irían UpdateProduct, DeleteProduct, ListProducts (con paginación) ...
```

### 3.2. Router (`/internal/api/router.go`)
Aquí definimos los endpoints y los agrupamos.

```go
package api

import "github.com/gin-gonic/gin"

func SetupRouter(handler *ProductHandler) *gin.Engine {
    r := gin.Default()

    // Agrupar rutas bajo /api/v1
    v1 := r.Group("/api/v1")
    {
        products := v1.Group("/products")
        {
            products.POST("/", handler.CreateProduct)
            products.GET("/:id", handler.GetProductByID)
            // ... PUT, DELETE, GET ...
        }
    }

    return r
}
```

### 3.3. Ensamblaje (`/cmd/server/main.go`)
La función `main` es la "maestra de orquesta". Inicializa todo y lo pone en marcha.

```go
package main

import (
    "log"
    "tu-proyecto/internal/api"
    "tu-proyecto/internal/repository"
    "tu-proyecto/internal/service"
    // ... imports de config y BBDD
)

func main() {
    // 1. Cargar configuración (ver sección 4)
    // 2. Conectar a la base de datos (ver sección 2)
    db := repository.ConnectToDB() // Esta función devuelve *sqlx.DB

    // 3. Inyección de Dependencias
    productRepo := repository.NewProductPostgresRepo(db)
    productService := service.NewProductService(productRepo)
    productHandler := api.NewProductHandler(productService)

    // 4. Configurar y arrancar el router
    router := api.SetupRouter(productHandler)
    log.Fatal(router.Run(":8080"))
}
```

---

## 4. Tópicos Avanzados para Producción

Una API funcional es solo el principio. Para que esté lista para producción, necesita ser robusta, segura y observable.

### 4.1. Validación de Datos de Entrada
Nunca confíes en los datos del cliente. Usa `struct tags` para definir reglas de validación.

- **Recurso:** [github.com/go-playground/validator](https://github.com/go-playground/validator)

```go
// En el struct de la petición del servicio
type CreateProductRequest struct {
    Name  string  `json:"name" binding:"required,min=3"`
    Price float64 `json:"price" binding:"required,gt=0"`
    SKU   string  `json:"sku" binding:"required,alphanum"`
}

// Gin usa esta librería por debajo. Si `ShouldBindJSON` falla,
// es porque la validación no pasó.
```

### 4.2. Gestión de Configuración con Viper
Maneja la configuración (puertos, credenciales de BBDD, secretos) de forma segura y flexible.

- **Recurso:** [github.com/spf13/viper](https://github.com/spf13/viper)

```go
// configs/config.go
import "github.com/spf13/viper"

type Config struct {
    DBHost string `mapstructure:"DB_HOST"`
    // ... otros campos
}

func LoadConfig() (config Config, err error) {
    viper.AddConfigPath(".")
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()
    // ...
    err = viper.Unmarshal(&config)
    return
}
```

### 4.3. Logging Estructurado con `slog`
Desde Go 1.21, `slog` es la librería estándar para logging estructurado (JSON), esencial para herramientas como Datadog o Splunk.

```go
import "log/slog"

// En main.go, configura un logger global
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.SetDefault(logger)

// En cualquier parte del código:
slog.Info("Usuario creado", "user_id", user.ID, "email", user.Email)
// Salida: {"time":"...","level":"INFO","msg":"Usuario creado","user_id":"...","email":"..."}
```

### 4.4. Autenticación JWT (Flujo Completo)
**1. Endpoint de Login (emite el token):**
```go
func (h *AuthHandler) Login(c *gin.Context) {
    // ... verificar usuario y contraseña ...

    // Crear token JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    })
    
    tokenString, err := token.SignedString([]byte("MI_SECRETO_SUPER_SEGURO"))
    // ... manejar error ...

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
```

**2. Middleware de Autenticación (protege rutas):**
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... extraer token del header "Authorization: Bearer <token>" ...
        // ... validar el token (firma y expiración) ...

        // Si es válido, añadir claims al contexto y continuar
        c.Set("userID", claims["user_id"])
        c.Next()

        // Si no, devolver 401 Unauthorized
        c.AbortWithStatus(http.StatusUnauthorized)
    }
}
```
**3. Aplicar el middleware a un grupo de rutas:**
```go
// en router.go
protected := v1.Group("/protected").Use(AuthMiddleware())
{
    protected.GET("/profile", handler.GetProfile)
}
```
---

## 5. Preparación para Producción

### 5.1. Contenerización con Docker
**`Dockerfile` (Optimizado para Go):**
```dockerfile
# ---- Build Stage ----
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Compila la aplicación de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# ---- Final Stage ----
FROM alpine:latest

WORKDIR /root/
# Copia solo el binario compilado desde la etapa de construcción
COPY --from=builder /app/main .

# Expone el puerto y ejecuta la aplicación
EXPOSE 8080
CMD ["./main"]
```
**`docker-compose.yml` (para desarrollo local):**
```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db
  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=mydb
    ports:
      - "5432:5432"
```

### 5.2. Graceful Shutdown
Esto asegura que la API termine de procesar las peticiones en curso antes de apagarse.

```go
// en main.go
import (
    "context"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // ... setup de la app ...
    
    srv := &http.Server{Addr: ":8080", Handler: router}

    // Iniciar servidor en una goroutine
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // Esperar por una señal de interrupción
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    <-ctx.Done() // Bloquea hasta que se recibe la señal

    log.Println("Apagando servidor de forma controlada...")

    // Dar 5 segundos para que las peticiones terminen
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Apagado forzado del servidor:", err)
    }
    log.Println("Servidor apagado.")
}
```
