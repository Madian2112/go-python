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
	PORT           = "8002"
	ORDERS_FILE    = "orders.json"
	AUTH_SECRET_KEY = "super-secret-auth-key" // Debe coincidir con el servicio de autenticación
	PRODUCT_SERVICE_URL = "http://localhost:8001"
)

// Modelos
type OrderItem struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required"`
	Price     float64 `json:"price"`
	Name      string  `json:"name,omitempty"`
}

type Order struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id" binding:"required"`
	Items      []OrderItem `json:"items" binding:"required"`
	Total      float64     `json:"total"`
	Status     string      `json:"status"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
}

// Base de datos simulada
var (
	orders []Order
	mutex  sync.RWMutex
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
func loadOrders() error {
	mutex.Lock()
	defer mutex.Unlock()

	// Verificar si el archivo existe
	if _, err := os.Stat(ORDERS_FILE); os.IsNotExist(err) {
		// Crear órdenes de ejemplo si el archivo no existe
		orders = []Order{
			{
				ID:         "1",
				CustomerID: "1",
				Items: []OrderItem{
					{
						ProductID: "1",
						Quantity:  2,
						Price:     999.99,
						Name:      "Laptop",
					},
				},
				Total:     1999.98,
				Status:    "completed",
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
		}
		return saveOrders()
	}

	// Leer el archivo
	data, err := ioutil.ReadFile(ORDERS_FILE)
	if err != nil {
		return err
	}

	// Deserializar las órdenes
	return json.Unmarshal(data, &orders)
}

func saveOrders() error {
	mutex.RLock()
	defer mutex.RUnlock()

	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(ORDERS_FILE, data, 0644)
}

// Funciones de servicio
func getProductByID(id string, token string) (*Product, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/products/%s", PRODUCT_SERVICE_URL, id), nil)
	if err != nil {
		return nil, err
	}

	// Agregar token si está disponible
	if token != "" {
		req.Header.Add("Authorization", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get product: status %d", resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

// Controladores
func getOrderByID(id string) (*Order, int) {
	mutex.RLock()
	defer mutex.RUnlock()

	for i, order := range orders {
		if order.ID == id {
			return &orders[i], i
		}
	}

	return nil, -1
}

func getOrdersHandler(c *gin.Context) {
	// Verificar autenticación
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	mutex.RLock()
	defer mutex.RUnlock()

	c.JSON(http.StatusOK, orders)
}

func getOrderHandler(c *gin.Context) {
	// Verificar autenticación
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id := c.Param("id")
	order, _ := getOrderByID(id)

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func createOrderHandler(c *gin.Context) {
	// Verificar autenticación
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener token para comunicación con el servicio de productos
	auth := c.GetHeader("Authorization")

	// Verificar productos y completar información
	var total float64
	for i, item := range order.Items {
		product, err := getProductByID(item.ProductID, auth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid product: %v", err)})
			return
		}

		// Actualizar información del item
		order.Items[i].Price = product.Price
		order.Items[i].Name = product.Name
		total += product.Price * float64(item.Quantity)
	}

	// Generar ID y timestamps
	order.ID = uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	order.CreatedAt = now
	order.UpdatedAt = now
	order.Status = "pending"
	order.Total = total

	// Agregar la orden
	mutex.Lock()
	orders = append(orders, order)
	mutex.Unlock()

	// Guardar en el archivo
	if err := saveOrders(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func updateOrderHandler(c *gin.Context) {
	// Verificar autenticación
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id := c.Param("id")
	order, index := getOrderByID(id)

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	var updatedOrder Order
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar campos
	updatedOrder.ID = id
	updatedOrder.CreatedAt = order.CreatedAt
	updatedOrder.UpdatedAt = time.Now().Format(time.RFC3339)

	// Actualizar la orden
	mutex.Lock()
	orders[index] = updatedOrder
	mutex.Unlock()

	// Guardar en el archivo
	if err := saveOrders(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save order"})
		return
	}

	c.JSON(http.StatusOK, updatedOrder)
}

func deleteOrderHandler(c *gin.Context) {
	// Verificar autenticación
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id := c.Param("id")
	_, index := getOrderByID(id)

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Eliminar la orden
	mutex.Lock()
	orders = append(orders[:index], orders[index+1:]...)
	mutex.Unlock()

	// Guardar en el archivo
	if err := saveOrders(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}

func getOrdersByCustomerHandler(c *gin.Context) {
	// Verificar autenticación
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	customerID := c.Param("customer_id")

	mutex.RLock()
	defer mutex.RUnlock()

	var customerOrders []Order
	for _, order := range orders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}

	c.JSON(http.StatusOK, customerOrders)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "order-service",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func main() {
	// Cargar órdenes
	if err := loadOrders(); err != nil {
		log.Fatalf("Failed to load orders: %v\n", err)
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

	// Rutas protegidas
	protected := r.Group("/")
	protected.Use(authMiddleware())
	protected.GET("/orders", getOrdersHandler)
	protected.GET("/orders/:id", getOrderHandler)
	protected.POST("/orders", createOrderHandler)
	protected.PUT("/orders/:id", updateOrderHandler)
	protected.DELETE("/orders/:id", deleteOrderHandler)
	protected.GET("/orders/customer/:customer_id", getOrdersByCustomerHandler)

	// Iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	log.Printf("Order service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}