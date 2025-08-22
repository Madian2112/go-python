package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Configuración
const (
	PORT = "8080"
)

// Servicios
var SERVICES = map[string]string{
	"auth":    "http://localhost:8000",
	"product": "http://localhost:8001",
	"order":   "http://localhost:8002",
}

// Modelos
type ServiceStatus struct {
	Service  string  `json:"service"`
	Status   string  `json:"status"`
	URL      string  `json:"url"`
	LatencyMs float64 `json:"latency_ms"`
}

// Middleware para logging
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		log.Printf("%s | %d | %s | %s | %v",
			c.Request.Method,
			c.Writer.Status(),
			c.Request.URL.Path,
			c.ClientIP(),
			latency,
		)
	}
}

// Funciones de utilidad
func forwardRequest(c *gin.Context, service, path string) {
	serviceURL, exists := SERVICES[service]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("service '%s' not found", service)})
		return
	}

	targetURL := fmt.Sprintf("%s%s", serviceURL, path)

	// Crear una nueva solicitud
	var req *http.Request
	var err error

	// Leer el cuerpo si es POST, PUT
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
			return
		}

		req, err = http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(c.Request.Method, targetURL, nil)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Copiar headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Copiar query params
	req.URL.RawQuery = c.Request.URL.RawQuery

	// Realizar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("error communicating with service: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response body"})
		return
	}

	// Copiar headers de respuesta
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Establecer el código de estado
	c.Status(resp.StatusCode)

	// Escribir el cuerpo de la respuesta
	c.Writer.Write(respBody)
}

// Controladores
func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "API Gateway for Microservices"})
}

func healthCheckHandler(c *gin.Context) {
	results := []ServiceStatus{}

	for serviceName, serviceURL := range SERVICES {
		startTime := time.Now()
		status := "healthy"

		// Verificar el estado del servicio
		resp, err := http.Get(fmt.Sprintf("%s/health", serviceURL))
		latency := time.Since(startTime).Seconds() * 1000

		if err != nil || resp.StatusCode != http.StatusOK {
			status = "unhealthy"
			if err != nil {
				status = "unavailable"
			}
		} else {
			defer resp.Body.Close()
		}

		results = append(results, ServiceStatus{
			Service:   serviceName,
			Status:    status,
			URL:       serviceURL,
			LatencyMs: latency,
		})
	}

	c.JSON(http.StatusOK, results)
}

func main() {
	// Configurar el router
	r := gin.Default()

	// Middleware
	r.Use(loggerMiddleware())

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas
	r.GET("/", rootHandler)
	r.GET("/health", healthCheckHandler)

	// Rutas para el servicio de autenticación
	r.POST("/api/auth/token", func(c *gin.Context) {
		forwardRequest(c, "auth", "/token")
	})

	r.GET("/api/auth/users/me", func(c *gin.Context) {
		forwardRequest(c, "auth", "/users/me")
	})

	// Rutas para el servicio de productos
	r.GET("/api/products", func(c *gin.Context) {
		forwardRequest(c, "product", "/products")
	})

	r.GET("/api/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, "product", fmt.Sprintf("/products/%s", id))
	})

	r.POST("/api/products", func(c *gin.Context) {
		forwardRequest(c, "product", "/products")
	})

	r.PUT("/api/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, "product", fmt.Sprintf("/products/%s", id))
	})

	r.DELETE("/api/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, "product", fmt.Sprintf("/products/%s", id))
	})

	r.GET("/api/products/category/:category", func(c *gin.Context) {
		category := c.Param("category")
		forwardRequest(c, "product", fmt.Sprintf("/products/category/%s", category))
	})

	// Rutas para el servicio de órdenes
	r.GET("/api/orders", func(c *gin.Context) {
		forwardRequest(c, "order", "/orders")
	})

	r.GET("/api/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, "order", fmt.Sprintf("/orders/%s", id))
	})

	r.POST("/api/orders", func(c *gin.Context) {
		forwardRequest(c, "order", "/orders")
	})

	r.PUT("/api/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, "order", fmt.Sprintf("/orders/%s", id))
	})

	r.DELETE("/api/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, "order", fmt.Sprintf("/orders/%s", id))
	})

	r.GET("/api/orders/customer/:customer_id", func(c *gin.Context) {
		customerID := c.Param("customer_id")
		forwardRequest(c, "order", fmt.Sprintf("/orders/customer/%s", customerID))
	})

	// Iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	log.Printf("API Gateway starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}