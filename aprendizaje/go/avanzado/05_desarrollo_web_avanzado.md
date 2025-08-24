# Desarrollo Web Avanzado en Go

## Introducción

Go se ha convertido en una opción popular para el desarrollo de aplicaciones web debido a su rendimiento, simplicidad y soporte para concurrencia. En este documento, exploraremos técnicas avanzadas para el desarrollo web en Go, incluyendo frameworks, patrones de diseño, optimización de rendimiento y seguridad.

## Frameworks Web en Go

### Comparativa de Frameworks

Go ofrece varios frameworks web, cada uno con sus propias ventajas y enfoques:

#### 1. Estándar: net/http

La biblioteca estándar `net/http` proporciona una base sólida para construir aplicaciones web:

```go
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hola, %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

**Ventajas:**
- Parte de la biblioteca estándar
- Sin dependencias externas
- Rendimiento excelente

**Desventajas:**
- Funcionalidad básica
- Requiere más código para características avanzadas

#### 2. Gin

Gin es un framework web de alto rendimiento con una sintaxis expresiva:

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run() // escucha en :8080
}
```

**Ventajas:**
- Alto rendimiento
- Middleware integrado
- Validación de formularios
- Enrutamiento de grupos

**Desventajas:**
- Dependencia externa

#### 3. Echo

Echo es un framework web de alto rendimiento y minimalista:

```go
package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":8080"))
}
```

**Ventajas:**
- API limpia y minimalista
- Buen rendimiento
- Extensible con middleware

**Desventajas:**
- Comunidad más pequeña que Gin

#### 4. Fiber

Fiber es un framework web inspirado en Express.js:

```go
package main

import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Listen(":8080")
}
```

**Ventajas:**
- Sintaxis similar a Express.js
- Muy rápido (basado en fasthttp)
- Zero memory allocation en rutas comunes

**Desventajas:**
- No usa net/http estándar
- API menos estable

### Selección de Framework

Factores a considerar al elegir un framework:

1. **Rendimiento**: Si el rendimiento es crítico, considera Gin, Echo o Fiber.
2. **Simplicidad**: Para proyectos pequeños, net/http puede ser suficiente.
3. **Ecosistema**: Gin tiene el ecosistema más grande con muchos plugins.
4. **Familiaridad**: Si vienes de Node.js, Fiber puede ser más familiar.
5. **Mantenimiento**: Todos los frameworks mencionados están bien mantenidos.

## Arquitecturas Web Avanzadas

### Arquitectura Limpia (Clean Architecture)

La arquitectura limpia separa las preocupaciones en capas distintas:

```
┌─────────────────────────────┐
│    Entidades (Dominio)      │
└─────────────┬───────────────┘
              │
┌─────────────▼───────────────┐
│    Casos de Uso (Servicios) │
└─────────────┬───────────────┘
              │
┌─────────────▼───────────────┐
│    Adaptadores (Controllers)│
└─────────────┬───────────────┘
              │
┌─────────────▼───────────────┐
│    Frameworks & Drivers     │
└─────────────────────────────┘
```

Ejemplo de estructura de proyecto:

```
/cmd
  /api
    main.go
/internal
  /domain
    user.go
    product.go
  /usecase
    user_service.go
    product_service.go
  /repository
    user_repository.go
    product_repository.go
  /delivery
    /http
      user_handler.go
      product_handler.go
/pkg
  /middleware
  /logger
  /config
```

Implementación de ejemplo:

```go
// domain/user.go
package domain

type User struct {
    ID    string
    Name  string
    Email string
}

type UserRepository interface {
    GetByID(id string) (*User, error)
    Save(user *User) error
}

type UserService interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
}

// usecase/user_service.go
package usecase

import "myapp/internal/domain"

type userService struct {
    userRepo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
    return &userService{userRepo: repo}
}

func (s *userService) GetUser(id string) (*domain.User, error) {
    return s.userRepo.GetByID(id)
}

func (s *userService) CreateUser(user *domain.User) error {
    return s.userRepo.Save(user)
}

// delivery/http/user_handler.go
package http

import (
    "encoding/json"
    "net/http"
    
    "github.com/gorilla/mux"
    "myapp/internal/domain"
)

type UserHandler struct {
    UserService domain.UserService
}

func NewUserHandler(r *mux.Router, us domain.UserService) {
    handler := &UserHandler{UserService: us}
    r.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
    r.HandleFunc("/users", handler.CreateUser).Methods("POST")
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    user, err := h.UserService.GetUser(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user domain.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := h.UserService.CreateUser(&user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}
```

### Arquitectura Hexagonal (Ports and Adapters)

La arquitectura hexagonal separa la lógica de negocio de los detalles de implementación:

```
┌─────────────────────────────────────────────────┐
│                                                 │
│                  Dominio                        │
│                                                 │
│  ┌─────────────────────────────────────────┐    │
│  │                                         │    │
│  │             Lógica de Negocio          │    │
│  │                                         │    │
│  └───────────┬─────────────────┬───────────┘    │
│              │                 │                 │
└──────────────┼─────────────────┼─────────────────┘
               │                 │
┌──────────────▼─────┐ ┌─────────▼──────────────┐
│                    │ │                         │
│  Puertos Primarios │ │  Puertos Secundarios   │
│    (API, UI)       │ │  (Persistencia, etc.)  │
│                    │ │                         │
└──────────────┬─────┘ └─────────┬──────────────┘
               │                 │
┌──────────────▼─────┐ ┌─────────▼──────────────┐
│                    │ │                         │
│     Adaptadores    │ │      Adaptadores       │
│     Primarios      │ │      Secundarios       │
│                    │ │                         │
└────────────────────┘ └─────────────────────────┘
```

Ejemplo de implementación:

```go
// domain/user.go
package domain

type User struct {
    ID    string
    Name  string
    Email string
}

// ports/user_repository.go
package ports

import "myapp/domain"

type UserRepository interface {
    GetByID(id string) (*domain.User, error)
    Save(user *domain.User) error
}

// ports/user_service.go
package ports

import "myapp/domain"

type UserService interface {
    GetUser(id string) (*domain.User, error)
    CreateUser(name, email string) (*domain.User, error)
}

// core/user_service.go
package core

import (
    "myapp/domain"
    "myapp/ports"
)

type userService struct {
    userRepo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
    return &userService{userRepo: repo}
}

func (s *userService) GetUser(id string) (*domain.User, error) {
    return s.userRepo.GetByID(id)
}

func (s *userService) CreateUser(name, email string) (*domain.User, error) {
    user := &domain.User{
        Name:  name,
        Email: email,
    }
    err := s.userRepo.Save(user)
    return user, err
}

// adapters/postgres/user_repository.go
package postgres

import (
    "database/sql"
    "myapp/domain"
)

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByID(id string) (*domain.User, error) {
    // Implementación con PostgreSQL
    return nil, nil
}

func (r *PostgresUserRepository) Save(user *domain.User) error {
    // Implementación con PostgreSQL
    return nil
}

// adapters/http/user_handler.go
package http

import (
    "encoding/json"
    "net/http"
    
    "github.com/gorilla/mux"
    "myapp/ports"
)

type UserHandler struct {
    userService ports.UserService
}

func NewUserHandler(r *mux.Router, us ports.UserService) {
    handler := &UserHandler{userService: us}
    r.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
    r.HandleFunc("/users", handler.CreateUser).Methods("POST")
}

// Implementación de los handlers...
```

## Middleware Avanzado

### Middleware Personalizado

Los middleware son funciones que procesan una solicitud HTTP antes o después del handler principal:

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Crear un ResponseWriter personalizado para capturar el código de estado
        ww := NewResponseWriter(w)
        
        // Llamar al siguiente handler
        next.ServeHTTP(ww, r)
        
        // Logging después de que el handler ha procesado la solicitud
        duration := time.Since(start)
        log.Printf(
            "%s %s %d %s",
            r.Method,
            r.RequestURI,
            ww.Status(),
            duration,
        )
    })
}

// ResponseWriter personalizado para capturar el código de estado
type ResponseWriter struct {
    http.ResponseWriter
    status int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
    return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Status() int {
    return rw.status
}
```

### Middleware de Autenticación

```go
func JWTMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extraer token del encabezado Authorization
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }
            
            // Formato: "Bearer {token}"
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
                return
            }
            
            tokenString := parts[1]
            
            // Validar token
            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                // Verificar algoritmo de firma
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                }
                return []byte(secret), nil
            })
            
            if err != nil || !token.Valid {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            // Extraer claims y añadir al contexto
            if claims, ok := token.Claims.(jwt.MapClaims); ok {
                ctx := context.WithValue(r.Context(), "user", claims)
                next.ServeHTTP(w, r.WithContext(ctx))
            } else {
                http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            }
        })
    }
}
```

### Middleware de CORS

```go
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Verificar si el origen está permitido
            allowed := false
            for _, allowedOrigin := range allowedOrigins {
                if allowedOrigin == "*" || allowedOrigin == origin {
                    allowed = true
                    break
                }
            }
            
            if allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
            
            // Manejar solicitudes preflight OPTIONS
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### Middleware de Rate Limiting

```go
func RateLimitMiddleware(limit int, windowSec int) func(http.Handler) http.Handler {
    // Usar un mapa para almacenar contadores por IP
    counters := make(map[string]*rateLimiter)
    mu := &sync.Mutex{}
    
    // Limpiar contadores viejos periódicamente
    go func() {
        for {
            time.Sleep(time.Duration(windowSec) * time.Second)
            mu.Lock()
            for ip, limiter := range counters {
                if time.Since(limiter.lastSeen) > time.Duration(windowSec)*time.Second {
                    delete(counters, ip)
                }
            }
            mu.Unlock()
        }
    }()
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            
            mu.Lock()
            if _, exists := counters[ip]; !exists {
                counters[ip] = &rateLimiter{
                    count:    0,
                    lastSeen: time.Now(),
                }
            }
            
            limiter := counters[ip]
            limiter.count++
            limiter.lastSeen = time.Now()
            
            if limiter.count > limit {
                mu.Unlock()
                w.Header().Set("Retry-After", fmt.Sprintf("%d", windowSec))
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            mu.Unlock()
            next.ServeHTTP(w, r)
        })
    }
}

type rateLimiter struct {
    count    int
    lastSeen time.Time
}
```

## Manejo Avanzado de Rutas

### Enrutamiento con Parámetros

```go
func main() {
    r := mux.NewRouter()
    
    // Parámetros en la ruta
    r.HandleFunc("/users/{id:[0-9]+}", getUserHandler).Methods("GET")
    
    // Parámetros opcionales
    r.HandleFunc("/articles/{category}/{id:[0-9]+}", getArticleHandler).Methods("GET")
    r.HandleFunc("/articles/{id:[0-9]+}", getArticleHandler).Methods("GET")
    
    // Subrutas
    apiRouter := r.PathPrefix("/api").Subrouter()
    apiRouter.HandleFunc("/users", getUsersHandler).Methods("GET")
    
    // Middleware específico para subrutas
    apiRouter.Use(JWTMiddleware("secret"))
    
    http.ListenAndServe(":8080", r)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    fmt.Fprintf(w, "User ID: %s", id)
}

func getArticleHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    category, ok := vars["category"]
    if !ok {
        category = "general"
    }
    fmt.Fprintf(w, "Article ID: %s, Category: %s", id, category)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    // Obtener usuario del contexto (establecido por el middleware JWT)
    user := r.Context().Value("user")
    fmt.Fprintf(w, "Users for: %v", user)
}
```

### Enrutamiento Basado en Dominio

```go
func main() {
    r := mux.NewRouter()
    
    // Rutas específicas de dominio
    apiRouter := r.Host("api.example.com").Subrouter()
    apiRouter.HandleFunc("/users", apiUsersHandler)
    
    webRouter := r.Host("www.example.com").Subrouter()
    webRouter.HandleFunc("/users", webUsersHandler)
    
    // Ruta por defecto
    r.HandleFunc("/", defaultHandler)
    
    http.ListenAndServe(":8080", r)
}
```

## Optimización de Rendimiento

### Uso de Caché

```go
type CacheMiddleware struct {
    cache    map[string]cacheEntry
    mu       sync.RWMutex
    duration time.Duration
}

type cacheEntry struct {
    content    []byte
    expiration time.Time
}

func NewCacheMiddleware(duration time.Duration) *CacheMiddleware {
    return &CacheMiddleware{
        cache:    make(map[string]cacheEntry),
        duration: duration,
    }
}

func (m *CacheMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Solo cachear solicitudes GET
        if r.Method != "GET" {
            next.ServeHTTP(w, r)
            return
        }
        
        // Crear clave de caché
        key := r.URL.String()
        
        // Verificar si está en caché
        m.mu.RLock()
        entry, found := m.cache[key]
        m.mu.RUnlock()
        
        if found && time.Now().Before(entry.expiration) {
            // Servir desde caché
            w.Write(entry.content)
            return
        }
        
        // Capturar la respuesta
        rec := httptest.NewRecorder()
        next.ServeHTTP(rec, r)
        
        // Guardar en caché
        result := rec.Result()
        content, _ := io.ReadAll(result.Body)
        result.Body.Close()
        
        m.mu.Lock()
        m.cache[key] = cacheEntry{
            content:    content,
            expiration: time.Now().Add(m.duration),
        }
        m.mu.Unlock()
        
        // Copiar encabezados
        for k, v := range result.Header {
            w.Header()[k] = v
        }
        
        // Escribir el cuerpo y el código de estado
        w.WriteHeader(result.StatusCode)
        w.Write(content)
    })
}
```

### Compresión de Respuestas

```go
func GzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verificar si el cliente acepta gzip
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }
        
        // Crear un ResponseWriter que comprime con gzip
        gz := gzip.NewWriter(w)
        defer gz.Close()
        
        w.Header().Set("Content-Encoding", "gzip")
        gzw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}
        
        next.ServeHTTP(gzw, r)
    })
}

type gzipResponseWriter struct {
    http.ResponseWriter
    Writer *gzip.Writer
}

func (gzw *gzipResponseWriter) Write(data []byte) (int, error) {
    return gzw.Writer.Write(data)
}
```

### Uso de Pools para Reducir la Presión del GC

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Obtener un buffer del pool
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset() // Asegurarse de que está vacío
    defer bufferPool.Put(buf) // Devolver al pool cuando termine
    
    // Usar el buffer
    buf.WriteString("Hello, World!")
    
    // Escribir la respuesta
    w.Write(buf.Bytes())
}
```

## Seguridad Web Avanzada

### Protección contra CSRF

```go
func CSRFMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Solo verificar en métodos que modifican estado
            if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" || r.Method == "TRACE" {
                next.ServeHTTP(w, r)
                return
            }
            
            // Verificar token CSRF
            token := r.Header.Get("X-CSRF-Token")
            if token == "" {
                http.Error(w, "CSRF token required", http.StatusForbidden)
                return
            }
            
            // Verificar cookie
            cookie, err := r.Cookie("csrf_token")
            if err != nil || cookie.Value == "" {
                http.Error(w, "CSRF cookie missing", http.StatusForbidden)
                return
            }
            
            // Comparar token y cookie
            if subtle.ConstantTimeCompare([]byte(token), []byte(cookie.Value)) != 1 {
                http.Error(w, "CSRF token invalid", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Generar y establecer token CSRF
func setCSRFToken(w http.ResponseWriter) string {
    token := make([]byte, 32)
    rand.Read(token)
    tokenStr := base64.StdEncoding.EncodeToString(token)
    
    http.SetCookie(w, &http.Cookie{
        Name:     "csrf_token",
        Value:    tokenStr,
        Path:     "/",
        HttpOnly: false, // Debe ser accesible por JavaScript
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        MaxAge:   3600,
    })
    
    return tokenStr
}
```

### Protección contra XSS

```go
func sanitizeHTML(input string) string {
    p := bluemonday.UGCPolicy()
    return p.Sanitize(input)
}

func renderTemplate(w http.ResponseWriter, templateName string, data map[string]interface{}) {
    // Sanitizar datos antes de renderizar
    for key, value := range data {
        if strValue, ok := value.(string); ok {
            data[key] = sanitizeHTML(strValue)
        }
    }
    
    // Renderizar plantilla
    tmpl, err := template.ParseFiles(templateName + ".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    tmpl.Execute(w, data)
}
```

### Cabeceras de Seguridad

```go
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Protección contra clickjacking
        w.Header().Set("X-Frame-Options", "DENY")
        
        // Protección XSS
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        
        // Prevenir MIME sniffing
        w.Header().Set("X-Content-Type-Options", "nosniff")
        
        // Content Security Policy
        w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none'")
        
        // HTTP Strict Transport Security
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // Referrer Policy
        w.Header().Set("Referrer-Policy", "no-referrer-when-downgrade")
        
        // Feature Policy
        w.Header().Set("Feature-Policy", "camera 'none'; microphone 'none'")
        
        next.ServeHTTP(w, r)
    })
}
```

## Pruebas Avanzadas

### Pruebas de Integración

```go
func TestUserAPI(t *testing.T) {
    // Configurar base de datos de prueba
    db, err := setupTestDB()
    if err != nil {
        t.Fatalf("Error setting up test database: %v", err)
    }
    defer cleanupTestDB(db)
    
    // Configurar servidor
    repo := postgres.NewPostgresUserRepository(db)
    service := core.NewUserService(repo)
    handler := http.NewUserHandler(service)
    
    r := mux.NewRouter()
    r.HandleFunc("/users", handler.CreateUser).Methods("POST")
    r.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
    
    server := httptest.NewServer(r)
    defer server.Close()
    
    // Prueba: Crear usuario
    userData := map[string]string{
        "name":  "Test User",
        "email": "test@example.com",
    }
    userJSON, _ := json.Marshal(userData)
    
    resp, err := http.Post(server.URL+"/users", "application/json", bytes.NewBuffer(userJSON))
    if err != nil {
        t.Fatalf("Error creating user: %v", err)
    }
    
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
    }
    
    // Extraer ID del usuario creado
    var createdUser map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&createdUser)
    resp.Body.Close()
    
    userID := createdUser["id"].(string)
    
    // Prueba: Obtener usuario
    resp, err = http.Get(server.URL + "/users/" + userID)
    if err != nil {
        t.Fatalf("Error getting user: %v", err)
    }
    
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
    }
    
    var fetchedUser map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&fetchedUser)
    resp.Body.Close()
    
    if fetchedUser["name"] != userData["name"] || fetchedUser["email"] != userData["email"] {
        t.Errorf("User data mismatch. Expected %v, got %v", userData, fetchedUser)
    }
}
```

### Pruebas de Carga

```go
func BenchmarkAPI(b *testing.B) {
    // Configurar servidor
    handler := setupHandler()
    server := httptest.NewServer(handler)
    defer server.Close()
    
    b.ResetTimer()
    
    // Ejecutar prueba de carga
    for i := 0; i < b.N; i++ {
        resp, err := http.Get(server.URL + "/users/1")
        if err != nil {
            b.Fatalf("Error: %v", err)
        }
        resp.Body.Close()
    }
}

func BenchmarkConcurrentAPI(b *testing.B) {
    // Configurar servidor
    handler := setupHandler()
    server := httptest.NewServer(handler)
    defer server.Close()
    
    b.ResetTimer()
    
    // Ejecutar prueba de carga concurrente
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            resp, err := http.Get(server.URL + "/users/1")
            if err != nil {
                b.Fatalf("Error: %v", err)
            }
            resp.Body.Close()
        }
    })
}
```

## Despliegue y Operaciones

### Configuración con Variables de Entorno

```go
type Config struct {
    Server struct {
        Port    int
        Timeout time.Duration
    }
    Database struct {
        Host     string
        Port     int
        User     string
        Password string
        Name     string
    }
    JWT struct {
        Secret  string
        Expires time.Duration
    }
}

func LoadConfig() (*Config, error) {
    var cfg Config
    
    // Valores por defecto
    cfg.Server.Port = 8080
    cfg.Server.Timeout = 30 * time.Second
    
    // Cargar desde variables de entorno
    if port := os.Getenv("SERVER_PORT"); port != "" {
        p, err := strconv.Atoi(port)
        if err != nil {
            return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
        }
        cfg.Server.Port = p
    }
    
    if timeout := os.Getenv("SERVER_TIMEOUT"); timeout != "" {
        t, err := time.ParseDuration(timeout)
        if err != nil {
            return nil, fmt.Errorf("invalid SERVER_TIMEOUT: %v", err)
        }
        cfg.Server.Timeout = t
    }
    
    // Database config
    cfg.Database.Host = getEnv("DB_HOST", "localhost")
    cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
    cfg.Database.User = getEnv("DB_USER", "postgres")
    cfg.Database.Password = getEnv("DB_PASSWORD", "")
    cfg.Database.Name = getEnv("DB_NAME", "app")
    
    // JWT config
    cfg.JWT.Secret = getEnv("JWT_SECRET", "")
    if cfg.JWT.Secret == "" {
        return nil, fmt.Errorf("JWT_SECRET is required")
    }
    cfg.JWT.Expires = getEnvAsDuration("JWT_EXPIRES", 24*time.Hour)
    
    return &cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if i, err := strconv.Atoi(value); err == nil {
            return i
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if d, err := time.ParseDuration(value); err == nil {
            return d
        }
    }
    return defaultValue
}
```

### Graceful Shutdown

```go
func main() {
    // Cargar configuración
    cfg, err := LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }
    
    // Configurar router
    r := mux.NewRouter()
    // ... configurar rutas
    
    // Configurar servidor
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
        Handler:      r,
        ReadTimeout:  cfg.Server.Timeout,
        WriteTimeout: cfg.Server.Timeout,
        IdleTimeout:  120 * time.Second,
    }
    
    // Canal para señales de sistema operativo
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    
    // Canal para errores del servidor
    errChan := make(chan error, 1)
    
    // Iniciar servidor en goroutine
    go func() {
        log.Printf("Server listening on %s", srv.Addr)
        errChan <- srv.ListenAndServe()
    }()
    
    // Esperar señal de interrupción o error
    select {
    case <-stop:
        log.Println("Shutdown signal received")
    case err := <-errChan:
        log.Printf("Server error: %v", err)
    }
    
    // Crear contexto con timeout para shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Intentar shutdown graceful
    log.Println("Shutting down server...")
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    }
    
    log.Println("Server stopped")
}
```

### Monitoreo y Métricas con Prometheus

```go
func main() {
    // Configurar métricas Prometheus
    requestCounter := prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestDuration := prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    // Registrar métricas
    prometheus.MustRegister(requestCounter, requestDuration)
    
    // Middleware para métricas
    metricsMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Capturar código de estado
            ww := NewResponseWriter(w)
            
            next.ServeHTTP(ww, r)
            
            // Registrar métricas
            duration := time.Since(start).Seconds()
            status := strconv.Itoa(ww.Status())
            
            requestCounter.WithLabelValues(r.Method, r.URL.Path, status).Inc()
            requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
        })
    }
    
    // Configurar router
    r := mux.NewRouter()
    r.Use(metricsMiddleware)
    
    // Endpoint para métricas Prometheus
    r.Handle("/metrics", promhttp.Handler())
    
    // ... resto de la configuración del servidor
}
```

## Ejercicios Prácticos

1. **Implementar una API RESTful completa**:
   - Crear una API para un blog con usuarios, posts y comentarios.
   - Implementar autenticación JWT.
   - Usar arquitectura limpia o hexagonal.
   - Añadir validación de datos de entrada.
   - Implementar paginación y filtrado.

2. **Crear un sistema de caché distribuido**:
   - Implementar un middleware de caché que use Redis.
   - Añadir invalidación de caché cuando los datos cambian.
   - Implementar diferentes estrategias de caché (time-based, LRU).
   - Medir y comparar el rendimiento con y sin caché.

3. **Desarrollar un sistema de websockets**:
   - Crear una aplicación de chat en tiempo real.
   - Implementar salas de chat.
   - Añadir autenticación para los usuarios.
   - Implementar reconexión automática.
   - Añadir notificaciones de "usuario escribiendo".

4. **Implementar un sistema de autorización basado en roles**:
   - Crear un sistema RBAC (Role-Based Access Control).
   - Implementar diferentes roles (admin, editor, usuario).
   - Añadir permisos granulares.
   - Crear middleware para verificar permisos.
   - Implementar auditoría de acciones.

5. **Crear un proxy inverso y balanceador de carga**:
   - Implementar un proxy que distribuya solicitudes entre múltiples servidores.
   - Añadir health checks para los servidores backend.
   - Implementar diferentes estrategias de balanceo (round-robin, least connections).
   - Añadir circuit breaker para manejar fallos en los backends.
   - Implementar rate limiting por cliente.

## Conclusiones

El desarrollo web avanzado en Go combina la simplicidad y eficiencia del lenguaje con patrones arquitectónicos robustos y técnicas de optimización. Al aprovechar las características de Go como la concurrencia y el rendimiento, podemos crear aplicaciones web escalables, seguras y mantenibles.

Los conceptos clave a recordar incluyen:

1. **Arquitectura**: Separar las preocupaciones usando arquitectura limpia o hexagonal mejora la mantenibilidad y testabilidad.

2. **Middleware**: Utilizar middleware para funcionalidades transversales como logging, autenticación y métricas.

3. **Rendimiento**: Implementar técnicas como caché, compresión y pools para optimizar el rendimiento.

4. **Seguridad**: Aplicar prácticas de seguridad como protección contra CSRF, XSS y cabeceras de seguridad.

5. **Operaciones**: Configurar correctamente el despliegue con graceful shutdown y monitoreo.

Al dominar estas técnicas avanzadas, podrás desarrollar aplicaciones web en Go que sean robustas, escalables y de alto rendimiento.

## Referencias

1. Mat Ryer. (2019). Go Programming Blueprints. Packt Publishing.
2. Sau Sheong Chang. (2019). Go Web Programming. Manning Publications.
3. Nic Jackson. (2018). Building Microservices with Go. Packt Publishing.
4. Jon Bodner. (2021). Learning Go: An Idiomatic Approach to Real-World Go Programming. O'Reilly Media.
5. Sam Newman. (2021). Building Microservices. O'Reilly Media.
6. Martin Fowler. (2018). Patterns of Enterprise Application Architecture. Addison-Wesley Professional.
7. OWASP. (2021). OWASP Top Ten. https://owasp.org/www-project-top-ten/
8. Go Documentation. (2023). net/http package. https://golang.org/pkg/net/http/