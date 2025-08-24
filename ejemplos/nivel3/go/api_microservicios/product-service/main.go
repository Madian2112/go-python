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
	Port     string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	JWTSecret string
}

// Obtener configuración desde variables de entorno
func getConfig() Config {
	return Config{
		Port:     getEnv("PORT", "8081"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASS", "postgres"),
		DBName:   getEnv("DB_NAME", "products"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
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

// Modelo de producto
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Category    string  `json:"category"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

// Repositorio de productos
type ProductRepository struct {
	db *sql.DB
}

// Crear un nuevo repositorio de productos
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Obtener todos los productos
func (r *ProductRepository) GetProducts() ([]Product, error) {
	rows, err := r.db.Query("SELECT id, name, description, price, stock, category, created_at, updated_at FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Category, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// Obtener un producto por ID
func (r *ProductRepository) GetProductByID(id string) (Product, error) {
	var p Product
	err := r.db.QueryRow("SELECT id, name, description, price, stock, category, created_at, updated_at FROM products WHERE id = $1", id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Category, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return Product{}, err
	}

	return p, nil
}

// Crear un nuevo producto
func (r *ProductRepository) CreateProduct(p Product) (Product, error) {
	p.ID = uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	p.CreatedAt = now
	p.UpdatedAt = now

	_, err := r.db.Exec(
		"INSERT INTO products (id, name, description, price, stock, category, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		p.ID, p.Name, p.Description, p.Price, p.Stock, p.Category, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return Product{}, err
	}

	return p, nil
}

// Actualizar un producto
func (r *ProductRepository) UpdateProduct(p Product) error {
	p.UpdatedAt = time.Now().Format(time.RFC3339)

	_, err := r.db.Exec(
		"UPDATE products SET name = $1, description = $2, price = $3, stock = $4, category = $5, updated_at = $6 WHERE id = $7",
		p.Name, p.Description, p.Price, p.Stock, p.Category, p.UpdatedAt, p.ID,
	)
	return err
}

// Eliminar un producto
func (r *ProductRepository) DeleteProduct(id string) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
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

	// Crear tabla de productos si no existe
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL,
			stock INTEGER NOT NULL,
			category VARCHAR(50),
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
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
	repo := NewProductRepository(db)

	// Crear router
	r := gin.Default()

	// Rutas públicas
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Métricas de Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Rutas de productos
	r.GET("/products", func(c *gin.Context) {
		products, err := repo.GetProducts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		product, err := repo.GetProductByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	r.POST("/products", func(c *gin.Context) {
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdProduct, err := repo.CreateProduct(product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdProduct)
	})

	r.PUT("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Verificar si el producto existe
		_, err := repo.GetProductByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Actualizar producto
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.ID = id

		if err := repo.UpdateProduct(product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	})

	r.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Verificar si el producto existe
		_, err := repo.GetProductByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := repo.DeleteProduct(id); err != nil {
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
		log.Printf("Product service listening on port %s\n", config.Port)
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