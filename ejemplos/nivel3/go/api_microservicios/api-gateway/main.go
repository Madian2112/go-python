package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Configuración del API Gateway
type Config struct {
	Port              string
	ProductServiceURL string
	OrderServiceURL   string
	UserServiceURL    string
	JWTSecret         string
}

// Obtener configuración desde variables de entorno
func getConfig() Config {
	return Config{
		Port:              getEnv("PORT", "8080"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://product-service:8081"),
		OrderServiceURL:   getEnv("ORDER_SERVICE_URL", "http://order-service:8082"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "http://user-service:8083"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
	}
}

// Helper para obtener variables de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Configuración
	config := getConfig()

	// Crear router
	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware de autenticación
	authMiddleware := AuthMiddleware(config.JWTSecret)

	// Rutas públicas
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	r.POST("/login", handleLogin(config))

	// Métricas de Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Rutas protegidas
	api := r.Group("/api")
	api.Use(authMiddleware)

	// Proxy para el servicio de productos
	products := api.Group("/products")
	products.Any("/*path", createProxy(config.ProductServiceURL))

	// Proxy para el servicio de pedidos
	orders := api.Group("/orders")
	orders.Any("/*path", createProxy(config.OrderServiceURL))

	// Proxy para el servicio de usuarios
	users := api.Group("/users")
	users.Any("/*path", createProxy(config.UserServiceURL))

	// Iniciar el servidor
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}

	// Iniciar el servidor en una goroutine
	go func() {
		log.Printf("API Gateway listening on port %s\n", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Esperar señal para apagar el servidor
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Contexto con timeout para apagar el servidor
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Apagar el servidor
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// Middleware de autenticación JWT
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Verificar formato del token (Bearer <token>)
		const prefix = "Bearer "
		if len(authorization) < len(prefix) || authorization[:len(prefix)] != prefix {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Extraer token
		token := authorization[len(prefix):]

		// Verificar token (implementación simplificada)
		// En un caso real, se debería usar un paquete como jwt-go para verificar el token
		if token == "invalid" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Continuar con la solicitud
		c.Next()
	}
}

// Handler para login
func handleLogin(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Estructura para los datos de login
		var loginData struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		// Validar datos de entrada
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verificar credenciales (implementación simplificada)
		// En un caso real, se debería verificar contra el servicio de usuarios
		if loginData.Username != "admin" || loginData.Password != "password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generar token (implementación simplificada)
		// En un caso real, se debería usar un paquete como jwt-go para generar el token
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

		c.JSON(http.StatusOK, gin.H{
			"token":   token,
			"type":    "Bearer",
			"expires": 3600,
		})
	}
}

// Crear un proxy para reenviar solicitudes a los microservicios
func createProxy(targetURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Crear cliente HTTP
		client := &http.Client{}

		// Construir URL de destino
		path := c.Param("path")
		url := targetURL + path

		// Crear solicitud
		req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
			return
		}

		// Copiar headers
		for name, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// Enviar solicitud
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error forwarding request"})
			return
		}
		defer resp.Body.Close()

		// Copiar headers de respuesta
		for name, values := range resp.Header {
			for _, value := range values {
				c.Header(name, value)
			}
		}

		// Copiar status code
		c.Status(resp.StatusCode)

		// Copiar body
		c.Writer.Write([]byte{})
		c.Request.Body = resp.Body
		c.Writer.WriteHeader(resp.StatusCode)
		c.Writer.Write([]byte{})
	}
}