package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/lib/pq"
)

// Configuración del servicio
type Config struct {
	Port              string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	ProductServiceURL string
	JWTSecret         string
}

// Obtener configuración desde variables de entorno
func getConfig() Config {
	return Config{
		Port:              getEnv("PORT", "8082"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPass:            getEnv("DB_PASS", "postgres"),
		DBName:            getEnv("DB_NAME", "orders"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://product-service:8081"),
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

// Modelo de producto (simplificado para usar con el servicio de productos)
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Modelo de item de pedido
type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	Price     float64 `json:"price"`
	Product   Product `json:"product,omitempty"`
}

// Modelo de pedido
type Order struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id" binding:"required"`
	Status     string      `json:"status"`
	Items      []OrderItem `json:"items" binding:"required,dive"`
	TotalPrice float64     `json:"total_price"`
	CreatedAt  string      `json:"created_at,omitempty"`
	UpdatedAt  string      `json:"updated_at,omitempty"`
}

// Repositorio de pedidos
type OrderRepository struct {
	db *sql.DB
}

// Crear un nuevo repositorio de pedidos
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Obtener todos los pedidos
func (r *OrderRepository) GetOrders() ([]Order, error) {
	// Obtener pedidos
	rows, err := r.db.Query("SELECT id, user_id, status, total_price, created_at, updated_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalPrice, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Obtener items del pedido
		items, err := r.getOrderItems(o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items

		orders = append(orders, o)
	}

	return orders, nil
}

// Obtener un pedido por ID
func (r *OrderRepository) GetOrderByID(id string) (Order, error) {
	var o Order
	err := r.db.QueryRow("SELECT id, user_id, status, total_price, created_at, updated_at FROM orders WHERE id = $1", id).Scan(
		&o.ID, &o.UserID, &o.Status, &o.TotalPrice, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return Order{}, err
	}

	// Obtener items del pedido
	items, err := r.getOrderItems(o.ID)
	if err != nil {
		return Order{}, err
	}
	o.Items = items

	return o, nil
}

// Obtener items de un pedido
func (r *OrderRepository) getOrderItems(orderID string) ([]OrderItem, error) {
	rows, err := r.db.Query("SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Crear un nuevo pedido
func (r *OrderRepository) CreateOrder(o Order) (Order, error) {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback()

	// Generar ID y timestamps
	o.ID = uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	o.CreatedAt = now
	o.UpdatedAt = now
	o.Status = "pending"

	// Calcular precio total
	var totalPrice float64
	for i := range o.Items {
		totalPrice += o.Items[i].Price * float64(o.Items[i].Quantity)
	}
	o.TotalPrice = totalPrice

	// Insertar pedido
	_, err = tx.Exec(
		"INSERT INTO orders (id, user_id, status, total_price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		o.ID, o.UserID, o.Status, o.TotalPrice, o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return Order{}, err
	}

	// Insertar items del pedido
	for i := range o.Items {
		o.Items[i].ID = uuid.New().String()
		o.Items[i].OrderID = o.ID

		_, err = tx.Exec(
			"INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4, $5)",
			o.Items[i].ID, o.Items[i].OrderID, o.Items[i].ProductID, o.Items[i].Quantity, o.Items[i].Price,
		)
		if err != nil {
			return Order{}, err
		}
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		return Order{}, err
	}

	return o, nil
}

// Actualizar un pedido
func (r *OrderRepository) UpdateOrder(o Order) error {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	o.UpdatedAt = time.Now().Format(time.RFC3339)

	// Actualizar pedido
	_, err = tx.Exec(
		"UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3",
		o.Status, o.UpdatedAt, o.ID,
	)
	if err != nil {
		return err
	}

	// Confirmar transacción
	return tx.Commit()
}

// Cancelar un pedido
func (r *OrderRepository) CancelOrder(id string) error {
	// Obtener pedido
	order, err := r.GetOrderByID(id)
	if err != nil {
		return err
	}

	// Verificar si el pedido ya está cancelado
	if order.Status == "cancelled" {
		return fmt.Errorf("order already cancelled")
	}

	// Actualizar estado
	order.Status = "cancelled"
	return r.UpdateOrder(order)
}

// Inicializar la base de datos
func initDB(config Config) (*sql.DB, error) {
	// Construir cadena de conexión
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName,
	)

	// Conectar a la base de datos
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Crear tablas si no existen
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			status VARCHAR(20) NOT NULL,
			total_price DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS order_items (
			id VARCHAR(36) PRIMARY KEY,
			order_id VARCHAR(36) NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			product_id VARCHAR(36) NOT NULL,
			quantity INTEGER NOT NULL,
			price DECIMAL(10, 2) NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Cliente para el servicio de productos
type ProductClient struct {
	baseURL string
	client  *http.Client
}

// Crear un nuevo cliente para el servicio de productos
func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

// Obtener un producto por ID
func (c *ProductClient) GetProduct(id string) (Product, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/products/%s", c.baseURL, id))
	if err != nil {
		return Product{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Product{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return Product{}, err
	}

	return product, nil
}

func main() {
	// Configuración
	config := getConfig()

	// Inicializar base de datos
	db, err := initDB(config)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Crear repositorio
	repo := NewOrderRepository(db)

	// Crear cliente para el servicio de productos
	productClient := NewProductClient(config.ProductServiceURL)

	// Crear router
	r := gin.Default()

	// Rutas públicas
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Métricas de Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Rutas de pedidos
	r.GET("/orders", func(c *gin.Context) {
		orders, err := repo.GetOrders()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	})

	r.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		order, err := repo.GetOrderByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Obtener información de productos para cada item
		for i := range order.Items {
			product, err := productClient.GetProduct(order.Items[i].ProductID)
			if err != nil {
				// No fallar si no se puede obtener el producto
				log.Printf("Error getting product %s: %v", order.Items[i].ProductID, err)
				continue
			}
			order.Items[i].Product = product
		}

		c.JSON(http.StatusOK, order)
	})

	r.POST("/orders", func(c *gin.Context) {
		var order Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verificar y obtener información de productos para cada item
		for i := range order.Items {
			product, err := productClient.GetProduct(order.Items[i].ProductID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid product ID: %s", order.Items[i].ProductID)})
				return
			}
			order.Items[i].Price = product.Price
			order.Items[i].Product = product
		}

		createdOrder, err := repo.CreateOrder(order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdOrder)
	})

	r.PUT("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Verificar si el pedido existe
		order, err := repo.GetOrderByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Actualizar estado del pedido
		var updateData struct {
			Status string `json:"status" binding:"required"`
		}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validar estado
		validStatuses := map[string]bool{"pending": true, "processing": true, "completed": true, "cancelled": true}
		if !validStatuses[updateData.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		order.Status = updateData.Status

		if err := repo.UpdateOrder(order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	})

	r.DELETE("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Cancelar pedido
		err := repo.CancelOrder(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			} else if err.Error() == "order already cancelled" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	// Iniciar el servidor
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}

	// Iniciar el servidor en una goroutine
	go func() {
		log.Printf("Order service listening on port %s\n", config.Port)
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