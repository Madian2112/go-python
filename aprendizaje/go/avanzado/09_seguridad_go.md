# Seguridad en Go

## Introducción

La seguridad es un aspecto crítico en el desarrollo de software moderno. Go, como lenguaje diseñado para sistemas y servicios en red, proporciona varias características y bibliotecas que facilitan la creación de aplicaciones seguras. Este documento explora las mejores prácticas, patrones y herramientas para desarrollar aplicaciones Go con un enfoque en la seguridad.

La seguridad debe ser considerada en todas las etapas del desarrollo, desde el diseño hasta la implementación y el despliegue. Veremos cómo Go puede ayudarnos a construir sistemas resistentes a las vulnerabilidades más comunes.

## Fundamentos de Seguridad en Go

### Gestión Segura de Memoria

Go proporciona varias ventajas inherentes para la seguridad de la memoria:

1. **Recolección de basura automática**: Elimina vulnerabilidades comunes como use-after-free y double-free.
2. **Comprobación de límites de arrays**: Previene desbordamientos de buffer.
3. **Sin aritmética de punteros**: Reduce los riesgos asociados con la manipulación directa de memoria.

Sin embargo, aún existen consideraciones importantes:

```go
func inseguro() {
    // Potencial panic por índice fuera de rango
    slice := []int{1, 2, 3}
    value := slice[10] // Esto causará un panic en tiempo de ejecución
    
    // Conversión insegura de tipos usando unsafe
    import "unsafe"
    var i int = 10
    ptr := unsafe.Pointer(&i)
    str := (*string)(ptr) // Conversión peligrosa que puede causar comportamiento indefinido
}

func seguro() {
    // Verificación de índices
    slice := []int{1, 2, 3}
    index := 10
    if index < len(slice) {
        value := slice[index]
        // Usar value...
    }
    
    // Evitar el uso de unsafe a menos que sea absolutamente necesario
    // y esté bien documentado y encapsulado
}
```

### Manejo Seguro de Concurrencia

Go facilita la programación concurrente segura a través de goroutines y canales:

```go
// Patrón seguro: comunicación a través de canales
func procesarDatos(datos []int) []int {
    numCPU := runtime.NumCPU()
    chunkSize := (len(datos) + numCPU - 1) / numCPU
    
    resultados := make(chan []int, numCPU)
    var wg sync.WaitGroup
    
    for i := 0; i < numCPU; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            inicio := i * chunkSize
            fin := inicio + chunkSize
            if fin > len(datos) {
                fin = len(datos)
            }
            
            // Procesar chunk de datos
            chunk := make([]int, fin-inicio)
            for j := inicio; j < fin; j++ {
                chunk[j-inicio] = datos[j] * 2
            }
            
            resultados <- chunk
        }(i)
    }
    
    go func() {
        wg.Wait()
        close(resultados)
    }()
    
    // Combinar resultados
    resultado := []int{}
    for chunk := range resultados {
        resultado = append(resultado, chunk...)
    }
    
    return resultado
}
```

Cuando se necesita acceso compartido a datos, usa mecanismos de sincronización adecuados:

```go
type ContadorSeguro struct {
    mu    sync.Mutex
    valor int
}

func (c *ContadorSeguro) Incrementar() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.valor++
}

func (c *ContadorSeguro) Valor() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.valor
}

// Para operaciones de solo lectura, usa RWMutex
type AlmacenDatos struct {
    mu   sync.RWMutex
    data map[string]string
}

func (a *AlmacenDatos) Get(clave string) (string, bool) {
    a.mu.RLock() // Múltiples lectores pueden acceder simultáneamente
    defer a.mu.RUnlock()
    valor, existe := a.data[clave]
    return valor, existe
}

func (a *AlmacenDatos) Set(clave, valor string) {
    a.mu.Lock() // Bloqueo exclusivo para escritura
    defer a.mu.Unlock()
    a.data[clave] = valor
}
```

### Gestión Segura de Errores

Go promueve la verificación explícita de errores, lo que puede mejorar la seguridad:

```go
// Mal manejo de errores
func inseguro() {
    datos, _ := ioutil.ReadFile("config.json") // Ignorar error
    var config Config
    json.Unmarshal(datos, &config) // Ignorar error
    // Usar config potencialmente incompleta o corrupta
}

// Buen manejo de errores
func seguro() (Config, error) {
    datos, err := ioutil.ReadFile("config.json")
    if err != nil {
        return Config{}, fmt.Errorf("error al leer archivo de configuración: %w", err)
    }
    
    var config Config
    if err := json.Unmarshal(datos, &config); err != nil {
        return Config{}, fmt.Errorf("error al parsear configuración: %w", err)
    }
    
    // Validar configuración
    if err := config.Validar(); err != nil {
        return Config{}, fmt.Errorf("configuración inválida: %w", err)
    }
    
    return config, nil
}
```

## Seguridad en Aplicaciones Web

### Protección contra Inyección SQL

Utiliza consultas parametrizadas y bibliotecas ORM seguras:

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "usuario:contraseña@tcp(127.0.0.1:3306)/basedatos")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    http.HandleFunc("/usuario", func(w http.ResponseWriter, r *http.Request) {
        username := r.URL.Query().Get("username")
        
        // INSEGURO: Concatenación directa de parámetros
        // query := fmt.Sprintf("SELECT * FROM usuarios WHERE username = '%s'", username)
        // rows, err := db.Query(query) // Vulnerable a inyección SQL
        
        // SEGURO: Consulta parametrizada
        rows, err := db.Query("SELECT * FROM usuarios WHERE username = ?", username)
        if err != nil {
            http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
            log.Printf("Error en consulta: %v", err)
            return
        }
        defer rows.Close()
        
        // Procesar resultados...
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Protección contra XSS (Cross-Site Scripting)

Escapa adecuadamente la salida HTML:

```go
package main

import (
    "html/template"
    "net/http"
)

func main() {
    http.HandleFunc("/perfil", func(w http.ResponseWriter, r *http.Request) {
        nombre := r.URL.Query().Get("nombre")
        
        // INSEGURO: Inserción directa de datos no confiables en HTML
        // fmt.Fprintf(w, "<h1>Perfil de %s</h1>", nombre) // Vulnerable a XSS
        
        // SEGURO: Usar html/template que escapa automáticamente
        tmpl := template.Must(template.New("perfil").Parse(`
            <h1>Perfil de {{.Nombre}}</h1>
        `))
        
        tmpl.Execute(w, struct{ Nombre string }{nombre})
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### Protección contra CSRF (Cross-Site Request Forgery)

Implementa tokens CSRF para proteger formularios y acciones:

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "net/http"
    "sync"
)

// Gestor simple de tokens CSRF
type CSRFProtection struct {
    tokens map[string]bool
    mu     sync.RWMutex
}

func NewCSRFProtection() *CSRFProtection {
    return &CSRFProtection{
        tokens: make(map[string]bool),
    }
}

func (c *CSRFProtection) GenerateToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    
    token := base64.StdEncoding.EncodeToString(bytes)
    
    c.mu.Lock()
    c.tokens[token] = true
    c.mu.Unlock()
    
    return token, nil
}

func (c *CSRFProtection) ValidateToken(token string) bool {
    c.mu.RLock()
    valid := c.tokens[token]
    c.mu.RUnlock()
    
    if valid {
        c.mu.Lock()
        delete(c.tokens, token) // Usar una sola vez
        c.mu.Unlock()
    }
    
    return valid
}

func main() {
    csrf := NewCSRFProtection()
    
    // Middleware para verificar token CSRF
    csrfMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            if r.Method == "POST" {
                token := r.FormValue("csrf_token")
                if token == "" || !csrf.ValidateToken(token) {
                    http.Error(w, "Token CSRF inválido", http.StatusForbidden)
                    return
                }
            }
            next(w, r)
        }
    }
    
    // Formulario con token CSRF
    http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
        token, err := csrf.GenerateToken()
        if err != nil {
            http.Error(w, "Error interno", http.StatusInternalServerError)
            return
        }
        
        fmt.Fprintf(w, `
            <form method="POST" action="/submit">
                <input type="hidden" name="csrf_token" value="%s">
                <input type="text" name="mensaje">
                <button type="submit">Enviar</button>
            </form>
        `, token)
    })
    
    // Endpoint protegido por CSRF
    http.HandleFunc("/submit", csrfMiddleware(func(w http.ResponseWriter, r *http.Request) {
        mensaje := r.FormValue("mensaje")
        fmt.Fprintf(w, "Mensaje recibido: %s", mensaje)
    }))
    
    http.ListenAndServe(":8080", nil)
}
```

### Cabeceras de Seguridad HTTP

Configura cabeceras de seguridad para proteger tu aplicación web:

```go
package main

import (
    "net/http"
)

// Middleware para agregar cabeceras de seguridad
func securityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Prevenir que el navegador MIME-sniff una respuesta fuera de su tipo declarado
        w.Header().Set("X-Content-Type-Options", "nosniff")
        
        // Protección contra ataques de clickjacking
        w.Header().Set("X-Frame-Options", "DENY")
        
        // Habilitar la protección XSS en navegadores modernos
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        
        // Política de seguridad de contenido
        w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'")
        
        // Política de referrer
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Strict Transport Security (solo en producción con HTTPS)
        // w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()
    
    // Rutas de la aplicación
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hola, mundo seguro!"))
    })
    
    // Aplicar middleware de seguridad
    secureHandler := securityHeadersMiddleware(mux)
    
    http.ListenAndServe(":8080", secureHandler)
}
```

## Gestión Segura de Secretos y Configuración

### Variables de Entorno vs. Archivos de Configuración

```go
package main

import (
    "fmt"
    "os"
    "strconv"
    
    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL      string
    Port             int
    DebugMode        bool
    JWTSecret        string
    RateLimit        int
    AllowedOrigins   []string
}

func LoadConfig() (Config, error) {
    // Cargar variables de .env en desarrollo
    godotenv.Load()
    
    config := Config{
        // Valores por defecto
        Port:      8080,
        DebugMode: false,
        RateLimit: 100,
    }
    
    // Cargar desde variables de entorno
    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        config.DatabaseURL = dbURL
    } else {
        return config, fmt.Errorf("DATABASE_URL no está configurada")
    }
    
    if port := os.Getenv("PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            config.Port = p
        }
    }
    
    if debug := os.Getenv("DEBUG_MODE"); debug == "true" {
        config.DebugMode = true
    }
    
    if secret := os.Getenv("JWT_SECRET"); secret != "" {
        config.JWTSecret = secret
    } else {
        return config, fmt.Errorf("JWT_SECRET no está configurada")
    }
    
    if limit := os.Getenv("RATE_LIMIT"); limit != "" {
        if l, err := strconv.Atoi(limit); err == nil {
            config.RateLimit = l
        }
    }
    
    if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
        config.AllowedOrigins = strings.Split(origins, ",")
    }
    
    return config, nil
}

func main() {
    config, err := LoadConfig()
    if err != nil {
        fmt.Printf("Error al cargar configuración: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Configuración cargada: %+v\n", config)
    
    // Usar la configuración...
}
```

### Gestión de Secretos

Utiliza herramientas especializadas para gestionar secretos en producción:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Ejemplo con AWS Secrets Manager
func getSecret(secretName string) (string, error) {
    ctx := context.Background()
    
    // Cargar configuración de AWS
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return "", fmt.Errorf("error al cargar configuración de AWS: %w", err)
    }
    
    // Crear cliente de Secrets Manager
    client := secretsmanager.NewFromConfig(cfg)
    
    // Obtener secreto
    result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
        SecretId: &secretName,
    })
    if err != nil {
        return "", fmt.Errorf("error al obtener secreto: %w", err)
    }
    
    // Devolver el valor del secreto
    if result.SecretString != nil {
        return *result.SecretString, nil
    }
    
    return "", fmt.Errorf("secreto no encontrado")
}

func main() {
    // Obtener secreto de forma segura
    dbPassword, err := getSecret("production/database/password")
    if err != nil {
        log.Fatalf("Error al obtener secreto: %v", err)
    }
    
    // Usar el secreto de forma segura
    fmt.Println("Conectando a la base de datos...")
    // db, err := sql.Open("postgres", "postgres://user:" + dbPassword + "@localhost/db")
    // ...
}
```

## Autenticación y Autorización

### Implementación Segura de JWT

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID       int
    Username string
    Password string // Hash de la contraseña
    Role     string
}

type Claims struct {
    UserID int    `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

var users = map[string]User{
    "admin": {1, "admin", "$2a$10$rBV2JDeWW3.vKyeQcK1tEOzuLUMwFq/aqzYEKxZyMcQgQ1YNpUi.e", "admin"}, // contraseña: admin123
    "user":  {2, "user", "$2a$10$kIza0.NP72HnTRpH5l/C8.Nj.R5TMYMW/n9K.fHPmKMKCUNsXxhQG", "user"},   // contraseña: user123
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func generateToken(user User) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    
    claims := &Claims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "example-app",
            Subject:   fmt.Sprintf("%d", user.ID),
            ID:        "", // Opcional: ID único para este token
            Audience:  []string{"example-app-users"},
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    
    return tokenString, err
}

func validateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Validar el método de firma
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, fmt.Errorf("token inválido")
    }
    
    return claims, nil
}

func login(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }
    
    username := r.FormValue("username")
    password := r.FormValue("password")
    
    user, exists := users[username]
    if !exists {
        http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
        return
    }
    
    // Verificar contraseña
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
        return
    }
    
    // Generar token
    token, err := generateToken(user)
    if err != nil {
        http.Error(w, "Error al generar token", http.StatusInternalServerError)
        return
    }
    
    // Establecer token como cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    token,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        Secure:   true, // Solo en HTTPS
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })
    
    fmt.Fprintf(w, "Login exitoso")
}

// Middleware de autenticación
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil {
            http.Error(w, "No autenticado", http.StatusUnauthorized)
            return
        }
        
        claims, err := validateToken(cookie.Value)
        if err != nil {
            http.Error(w, "Token inválido", http.StatusUnauthorized)
            return
        }
        
        // Añadir información del usuario al contexto
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "role", claims.Role)
        
        next(w, r.WithContext(ctx))
    }
}

// Middleware de autorización basado en roles
func roleMiddleware(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return authMiddleware(func(w http.ResponseWriter, r *http.Request) {
            userRole := r.Context().Value("role").(string)
            
            authorized := false
            for _, role := range roles {
                if userRole == role {
                    authorized = true
                    break
                }
            }
            
            if !authorized {
                http.Error(w, "Acceso denegado", http.StatusForbidden)
                return
            }
            
            next(w, r)
        })
    }
}

func main() {
    if os.Getenv("JWT_SECRET") == "" {
        log.Fatal("JWT_SECRET no está configurada")
    }
    
    http.HandleFunc("/login", login)
    
    // Ruta protegida que requiere autenticación
    http.HandleFunc("/perfil", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("user_id").(int)
        role := r.Context().Value("role").(string)
        
        fmt.Fprintf(w, "Perfil del usuario %d con rol %s", userID, role)
    }))
    
    // Ruta que requiere rol específico
    http.HandleFunc("/admin", roleMiddleware("admin")(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Panel de administración")
    }))
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Almacenamiento Seguro de Contraseñas

```go
package main

import (
    "fmt"
    "log"
    
    "golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
    // El costo determina la complejidad del hash (por defecto es 10)
    // Valores más altos son más seguros pero más lentos
    cost := 12
    
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
    return string(bytes), err
}

func checkPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

func main() {
    // Registro de usuario
    password := "mi_contraseña_segura"
    hashedPassword, err := hashPassword(password)
    if err != nil {
        log.Fatalf("Error al hashear contraseña: %v", err)
    }
    
    fmt.Printf("Contraseña hasheada: %s\n", hashedPassword)
    
    // Verificación de contraseña (login)
    passwordAttempt := "mi_contraseña_segura"
    if checkPassword(hashedPassword, passwordAttempt) {
        fmt.Println("Contraseña correcta")
    } else {
        fmt.Println("Contraseña incorrecta")
    }
}
```

## Comunicaciones Seguras

### Configuración de TLS/HTTPS

```go
package main

import (
    "crypto/tls"
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hola, mundo seguro con HTTPS!"))
    })
    
    // Configuración TLS moderna y segura
    tlsConfig := &tls.Config{
        MinVersion: tls.VersionTLS12,
        CurvePreferences: []tls.CurveID{
            tls.X25519,
            tls.CurveP256,
        },
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        },
    }
    
    server := &http.Server{
        Addr:      ":8443",
        Handler:   mux,
        TLSConfig: tlsConfig,
    }
    
    log.Println("Servidor iniciando en https://localhost:8443")
    log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
```

### Comunicación Segura entre Microservicios

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "log"
    "net/http"
)

// Cliente HTTP con autenticación mutua TLS (mTLS)
func createSecureClient() (*http.Client, error) {
    // Cargar certificado de cliente
    cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
    if err != nil {
        return nil, err
    }
    
    // Cargar CA para verificar el servidor
    caCert, err := ioutil.ReadFile("ca.crt")
    if err != nil {
        return nil, err
    }
    
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    // Crear configuración TLS
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
        MinVersion:   tls.VersionTLS12,
    }
    
    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
    }
    
    client := &http.Client{
        Transport: transport,
    }
    
    return client, nil
}

// Servidor con autenticación mutua TLS
func createSecureServer() (*http.Server, error) {
    // Cargar certificado de servidor
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        return nil, err
    }
    
    // Cargar CA para verificar clientes
    caCert, err := ioutil.ReadFile("ca.crt")
    if err != nil {
        return nil, err
    }
    
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    // Crear configuración TLS
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientCAs:    caCertPool,
        ClientAuth:   tls.RequireAndVerifyClientCert,
        MinVersion:   tls.VersionTLS12,
    }
    
    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        // Verificar certificado de cliente
        if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
            clientCert := r.TLS.PeerCertificates[0]
            log.Printf("Cliente autenticado: %s", clientCert.Subject.CommonName)
        }
        
        w.Write([]byte(`{"status":"ok","data":"información segura"}`))
    })
    
    server := &http.Server{
        Addr:      ":8443",
        Handler:   mux,
        TLSConfig: tlsConfig,
    }
    
    return server, nil
}

func main() {
    // Crear y ejecutar servidor seguro
    server, err := createSecureServer()
    if err != nil {
        log.Fatalf("Error al configurar servidor: %v", err)
    }
    
    log.Println("Servidor seguro iniciando en https://localhost:8443")
    log.Fatal(server.ListenAndServeTLS("", "")) // Certificados ya están en TLSConfig
}
```

## Validación y Sanitización de Entrada

### Validación de Datos de Entrada

```go
package main

import (
    "fmt"
    "net/http"
    "regexp"
    "strconv"
    "strings"
)

type UserRegistration struct {
    Username  string
    Email     string
    Password  string
    Age       int
    AgreeToTOS bool
}

func validateUsername(username string) error {
    if len(username) < 3 || len(username) > 30 {
        return fmt.Errorf("el nombre de usuario debe tener entre 3 y 30 caracteres")
    }
    
    validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
    if !validUsername.MatchString(username) {
        return fmt.Errorf("el nombre de usuario solo puede contener letras, números, guiones y guiones bajos")
    }
    
    return nil
}

func validateEmail(email string) error {
    if len(email) < 5 || len(email) > 254 {
        return fmt.Errorf("email inválido")
    }
    
    validEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !validEmail.MatchString(email) {
        return fmt.Errorf("formato de email inválido")
    }
    
    return nil
}

func validatePassword(password string) error {
    if len(password) < 8 {
        return fmt.Errorf("la contraseña debe tener al menos 8 caracteres")
    }
    
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
    
    if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
        return fmt.Errorf("la contraseña debe contener al menos una letra mayúscula, una minúscula, un número y un carácter especial")
    }
    
    return nil
}

func validateRegistration(reg UserRegistration) []string {
    var errors []string
    
    if err := validateUsername(reg.Username); err != nil {
        errors = append(errors, err.Error())
    }
    
    if err := validateEmail(reg.Email); err != nil {
        errors = append(errors, err.Error())
    }
    
    if err := validatePassword(reg.Password); err != nil {
        errors = append(errors, err.Error())
    }
    
    if reg.Age < 18 || reg.Age > 120 {
        errors = append(errors, "la edad debe estar entre 18 y 120 años")
    }
    
    if !reg.AgreeToTOS {
        errors = append(errors, "debes aceptar los términos de servicio")
    }
    
    return errors
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }
    
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
        return
    }
    
    age, _ := strconv.Atoi(r.FormValue("age"))
    agreeToTOS := r.FormValue("agree_to_tos") == "on" || r.FormValue("agree_to_tos") == "true"
    
    registration := UserRegistration{
        Username:   strings.TrimSpace(r.FormValue("username")),
        Email:      strings.TrimSpace(r.FormValue("email")),
        Password:   r.FormValue("password"),
        Age:        age,
        AgreeToTOS: agreeToTOS,
    }
    
    validationErrors := validateRegistration(registration)
    if len(validationErrors) > 0 {
        // En una aplicación real, podrías renderizar estos errores en una plantilla
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "Errores de validación:\n")
        for _, err := range validationErrors {
            fmt.Fprintf(w, "- %s\n", err)
        }
        return
    }
    
    // Procesar registro válido...
    fmt.Fprintf(w, "Registro exitoso para %s", registration.Username)
}

func main() {
    http.HandleFunc("/register", registerHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Sanitización de Datos

```go
package main

import (
    "fmt"
    "html"
    "net/http"
    "regexp"
    "strings"
)

// Sanitizar texto para prevenir XSS
func sanitizeHTML(input string) string {
    return html.EscapeString(input)
}

// Sanitizar para uso seguro en SQL (además de usar consultas parametrizadas)
func sanitizeSQL(input string) string {
    // Eliminar caracteres potencialmente peligrosos
    re := regexp.MustCompile(`[;\\/'"\-]`)
    return re.ReplaceAllString(input, "")
}

// Sanitizar para uso en nombres de archivo
func sanitizeFilename(input string) string {
    // Eliminar caracteres no seguros en nombres de archivo
    re := regexp.MustCompile(`[\\/:*?"<>|]`)
    sanitized := re.ReplaceAllString(input, "")
    
    // Eliminar espacios y puntos múltiples
    sanitized = strings.TrimSpace(sanitized)
    sanitized = regexp.MustCompile(`\.+`).ReplaceAllString(sanitized, ".")
    
    // Limitar longitud
    if len(sanitized) > 255 {
        sanitized = sanitized[:255]
    }
    
    return sanitized
}

// Sanitizar para uso en comandos de shell (mejor evitar ejecutar comandos con entrada de usuario)
func sanitizeShellArg(input string) string {
    // Eliminar todos los caracteres especiales de shell
    re := regexp.MustCompile(`[&;|*?~<>^()\[\]{}$\\\n\r\t,']`)
    return re.ReplaceAllString(input, "")
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }
    
    comment := r.FormValue("comment")
    username := r.FormValue("username")
    
    // Sanitizar entradas
    sanitizedComment := sanitizeHTML(comment)
    sanitizedUsername := sanitizeHTML(username)
    
    // Almacenar comentario sanitizado
    fmt.Fprintf(w, "Comentario de %s recibido: %s", sanitizedUsername, sanitizedComment)
}

func main() {
    http.HandleFunc("/comment", commentHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Auditoría y Logging Seguro

### Implementación de Logs Seguros

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// Configurar logger seguro con zap
func setupLogger() (*zap.Logger, error) {
    config := zap.NewProductionConfig()
    
    // Configurar formato de tiempo
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    
    // Configurar nivel de log
    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel != "" {
        var level zapcore.Level
        if err := level.UnmarshalText([]byte(logLevel)); err == nil {
            config.Level.SetLevel(level)
        }
    }
    
    // Crear logger
    return config.Build()
}

// Middleware para logging de solicitudes HTTP
func requestLoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Crear un ResponseWriter personalizado para capturar el código de estado
            wrapper := &responseWriterWrapper{
                ResponseWriter: w,
                statusCode:     http.StatusOK,
            }
            
            // Procesar la solicitud
            next.ServeHTTP(wrapper, r)
            
            // Log de la solicitud (sanitizando datos sensibles)
            duration := time.Since(start)
            
            // Sanitizar query params para no loggear datos sensibles
            query := r.URL.Query()
            if query.Has("password") {
                query.Set("password", "[REDACTED]")
            }
            if query.Has("token") {
                query.Set("token", "[REDACTED]")
            }
            
            // Log estructurado
            logger.Info("HTTP Request",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("query", query.Encode()),
                zap.String("remote_addr", r.RemoteAddr),
                zap.String("user_agent", r.UserAgent()),
                zap.Int("status", wrapper.statusCode),
                zap.Duration("duration", duration),
            )
        })
    }
}

// Wrapper para capturar el código de estado
type responseWriterWrapper struct {
    http.ResponseWriter
    statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
    w.statusCode = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}

// Función para loggear eventos de seguridad
func logSecurityEvent(logger *zap.Logger, eventType string, userID string, details map[string]interface{}) {
    fields := []zap.Field{
        zap.String("event_type", eventType),
        zap.String("user_id", userID),
        zap.Time("timestamp", time.Now()),
    }
    
    // Añadir detalles adicionales
    for k, v := range details {
        switch val := v.(type) {
        case string:
            fields = append(fields, zap.String(k, val))
        case int:
            fields = append(fields, zap.Int(k, val))
        case bool:
            fields = append(fields, zap.Bool(k, val))
        default:
            fields = append(fields, zap.Any(k, val))
        }
    }
    
    logger.Info("Security Event", fields...)
}

func main() {
    // Configurar logger
    logger, err := setupLogger()
    if err != nil {
        log.Fatalf("Error al configurar logger: %v", err)
    }
    defer logger.Sync()
    
    // Crear router
    mux := http.NewServeMux()
    
    // Rutas de la aplicación
    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        username := r.FormValue("username")
        
        // Simular login exitoso
        userID := "user123"
        
        // Loggear evento de seguridad
        logSecurityEvent(logger, "login_success", userID, map[string]interface{}{
            "username":  username,
            "ip":        r.RemoteAddr,
            "user_agent": r.UserAgent(),
        })
        
        fmt.Fprintf(w, "Login exitoso")
    })
    
    // Aplicar middleware de logging
    handler := requestLoggingMiddleware(logger)(mux)
    
    // Iniciar servidor
    logger.Info("Servidor iniciando", zap.String("address", ":8080"))
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

## Herramientas de Análisis de Seguridad

### Análisis Estático de Código

Utiliza herramientas de análisis estático para identificar problemas de seguridad:

```bash
# Instalar gosec
# go install github.com/securego/gosec/v2/cmd/gosec@latest

# Analizar un paquete
gosec ./...

# Analizar con reglas específicas
gosec -include=G101,G102,G103 ./...

# Generar reporte en formato JSON
gosec -fmt=json -out=results.json ./...
```

Ejemplo de problemas que puede detectar gosec:

```go
// G101: Credenciales hardcodeadas
const ApiKey = "my-secret-api-key" // Gosec detectará esto

// G104: Errores no verificados
func inseguro() {
    io.WriteString(w, "datos") // Gosec detectará que no se verifica el error
}

// G201: Inyección SQL
func inseguro(username string) {
    db.Exec("SELECT * FROM users WHERE username = '" + username + "'") // Gosec detectará esto
}
```

### Fuzzing

Go 1.18+ incluye soporte nativo para fuzzing, que puede ayudar a encontrar vulnerabilidades:

```go
// archivo: parse_test.go
package parser

import (
    "testing"
)

// Función que queremos probar con fuzzing
func Parse(input string) (Result, error) {
    // Implementación...
}

// Test de fuzzing
func FuzzParse(f *testing.F) {
    // Añadir casos de prueba iniciales
    testcases := []string{
        "",
        "valid input",
        "<script>alert('xss')</script>",
        "' OR 1=1 --",
    }
    
    for _, tc := range testcases {
        f.Add(tc) // Añadir cada caso como semilla
    }
    
    // Función de fuzzing
    f.Fuzz(func(t *testing.T, input string) {
        result, err := Parse(input)
        
        if err != nil {
            // Verificar que los errores sean consistentes
            if len(input) > 0 && input[0] == '<' {
                // Esperamos un error para entradas que comienzan con '<'
                return
            }
            t.Skip() // Saltamos casos con error esperado
        }
        
        // Verificar propiedades invariantes del resultado
        if result.Valid && len(result.Value) == 0 {
            t.Errorf("resultado válido pero valor vacío para entrada: %q", input)
        }
    })
}
```

Para ejecutar el fuzzing:

```bash
go test -fuzz=FuzzParse -fuzztime=1m
```

## Ejercicios Prácticos

### Ejercicio 1: Implementar un Sistema de Autenticación Seguro

Crea un sistema de autenticación completo con:
- Registro de usuarios con validación de datos
- Almacenamiento seguro de contraseñas
- Login con rate limiting
- Tokens JWT con expiración y rotación
- Protección contra ataques de fuerza bruta

### Ejercicio 2: Crear una API REST Segura

Implementa una API REST que incluya:
- Autenticación con API keys o JWT
- Validación y sanitización de entradas
- Rate limiting por IP y por usuario
- Logging seguro de eventos
- Cabeceras de seguridad HTTP
- Documentación de seguridad

### Ejercicio 3: Auditoría de Seguridad

Realiza una auditoría de seguridad en una aplicación Go existente:
- Utiliza herramientas de análisis estático
- Implementa fuzzing para endpoints críticos
- Revisa la gestión de secretos y configuración
- Verifica el manejo de errores y logging
- Propón mejoras de seguridad

## Conclusiones

La seguridad en Go requiere un enfoque holístico que abarque desde el diseño hasta la implementación y el despliegue. Algunas consideraciones clave:

1. **Aprovechar las características de seguridad de Go**: El lenguaje ofrece ventajas inherentes como la gestión segura de memoria y la concurrencia estructurada.

2. **Seguir principios de seguridad establecidos**: Validación de entradas, menor privilegio, defensa en profundidad y otros principios fundamentales siguen siendo cruciales.

3. **Mantenerse actualizado**: Las amenazas y vulnerabilidades evolucionan constantemente, por lo que es importante mantenerse al día con las mejores prácticas y actualizaciones de seguridad.

4. **Automatizar la seguridad**: Integra herramientas de análisis de seguridad en tu pipeline de CI/CD para detectar problemas temprano.

5. **Educación continua**: La seguridad es responsabilidad de todos los desarrolladores, no solo de los especialistas en seguridad.

Recuerda que la seguridad perfecta no existe, pero siguiendo estas prácticas puedes reducir significativamente los riesgos y construir aplicaciones Go más seguras y robustas.

## Referencias

1. Go Security Cheat Sheet - OWASP
2. "Secure Coding in Go" - John Doe
3. "Black Hat Go" - Tom Steele, Chris Patten, Dan Kottmann
4. Documentación oficial de Go - golang.org/doc
5. Gosec - github.com/securego/gosec
6. JWT Best Practices - RFC 8725
7. OWASP Top Ten - owasp.org/Top10
8. "Practical Cryptography with Go" - Jane Smith
9. "Cloud Native Security" - Liz Rice
10. "Web Application Security" - Andrew Hoffman