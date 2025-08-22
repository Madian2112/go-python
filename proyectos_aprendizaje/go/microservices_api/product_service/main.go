package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Configuración
const (
	PORT           = "8001"
	PRODUCTS_FILE  = "products.json"
	AUTH_SECRET_KEY = "super-secret-auth-key" // Debe coincidir con el servicio de autenticación
)

// Modelos
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Base de datos simulada
var (
	products []Product
	mutex    sync.RWMutex
)

// Middleware de autenticación
func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AUTH_SECRET_KEY), nil
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

// Funciones de persistencia
func loadProducts() error {
	mutex.Lock()
	defer mutex.Unlock()

	// Verificar si el archivo existe
	if _, err := os.Stat(PRODUCTS_FILE); os.IsNotExist(err) {
		// Crear productos de ejemplo si el archivo no existe
		products = []Product{
			{
				ID:          "1",
				Name:        "Laptop",
				Description: "Potente laptop para desarrollo",
				Price:       999.99,
				Category:    "Electronics",
				Stock:       10,
				CreatedAt:   time.Now().Format(time.RFC3339),
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
			{
				ID:          "2",
				Name:        "Smartphone",
				Description: "Teléfono inteligente de última generación",
				Price:       699.99,
				Category:    "Electronics",
				Stock:       15,
				CreatedAt:   time.Now().Format(time.RFC3339),
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
		}
		return saveProducts()
	}

	// Leer el archivo
	data, err := ioutil.ReadFile(PRODUCTS_FILE)
	if err != nil {
		return err
	}

	// Deserializar los productos
	return json.Unmarshal(data, &products)
}

func saveProducts() error {
	mutex.RLock()
	defer mutex.RUnlock()

	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(PRODUCTS_FILE, data, 0644)
}

// Controladores
func getProductByID(id string) (*Product, int) {
	mutex.RLock()
	defer mutex.RUnlock()

	for i, product := range products {
		if product.ID == id {
			return &products[i], i
		}
	}

	return nil, -1
}

func getProductsHandler(c *gin.Context) {
	mutex.RLock()
	defer mutex.RUnlock()

	c.JSON(http.StatusOK, products)
}

func getProductHandler(c *gin.Context) {
	id := c.Param("id")
	product, _ := getProductByID(id)

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func createProductHandler(c *gin.Context) {
	// Verificar autenticación
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generar ID y timestamps
	product.ID = uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	product.CreatedAt = now
	product.UpdatedAt = now

	// Agregar el producto
	mutex.Lock()
	products = append(products, product)
	mutex.Unlock()

	// Guardar en el archivo
	if err := saveProducts(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func updateProductHandler(c *gin.Context) {
	// Verificar autenticación
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id := c.Param("id")
	product, index := getProductByID(id)

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	var updatedProduct Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar campos
	updatedProduct.ID = id
	updatedProduct.CreatedAt = product.CreatedAt
	updatedProduct.UpdatedAt = time.Now().Format(time.RFC3339)

	// Actualizar el producto
	mutex.Lock()
	products[index] = updatedProduct
	mutex.Unlock()

	// Guardar en el archivo
	if err := saveProducts(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save product"})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func deleteProductHandler(c *gin.Context) {
	// Verificar autenticación
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id := c.Param("id")
	_, index := getProductByID(id)

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Eliminar el producto
	mutex.Lock()
	products = append(products[:index], products[index+1:]...)
	mutex.Unlock()

	// Guardar en el archivo
	if err := saveProducts(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

func getProductsByCategoryHandler(c *gin.Context) {
	category := c.Param("category")

	mutex.RLock()
	defer mutex.RUnlock()

	var filteredProducts []Product
	for _, product := range products {
		if strings.EqualFold(product.Category, category) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	c.JSON(http.StatusOK, filteredProducts)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "product-service",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func main() {
	// Cargar productos
	if err := loadProducts(); err != nil {
		log.Fatalf("Failed to load products: %v\n", err)
	}

	// Configurar el router
	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas públicas
	r.GET("/health", healthCheckHandler)
	r.GET("/products", getProductsHandler)
	r.GET("/products/:id", getProductHandler)
	r.GET("/products/category/:category", getProductsByCategoryHandler)

	// Rutas protegidas
	protected := r.Group("/")
	protected.Use(authMiddleware())
	protected.POST("/products", createProductHandler)
	protected.PUT("/products/:id", updateProductHandler)
	protected.DELETE("/products/:id", deleteProductHandler)

	// Iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	log.Printf("Product service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}