package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Configuración del servicio
type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
}

// Obtener configuración desde variables de entorno
func getConfig() Config {
	return Config{
		Port:      getEnv("PORT", "8083"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "postgres"),
		DBName:    getEnv("DB_NAME", "users"),
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

// Modelo de usuario
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password,omitempty" binding:"required,min=6"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Modelo para login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Modelo para respuesta de login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Claims para JWT
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Repositorio de usuarios
type UserRepository struct {
	db *sql.DB
}

// Crear un nuevo repositorio de usuarios
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Obtener todos los usuarios
func (r *UserRepository) GetUsers() ([]User, error) {
	rows, err := r.db.Query("SELECT id, username, email, role, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// Obtener un usuario por ID
func (r *UserRepository) GetUserByID(id string) (User, error) {
	var u User
	err := r.db.QueryRow("SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1", id).Scan(
		&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

// Obtener un usuario por nombre de usuario
func (r *UserRepository) GetUserByUsername(username string) (User, error) {
	var u User
	var hashedPassword string

	err := r.db.QueryRow("SELECT id, username, email, password, role, created_at, updated_at FROM users WHERE username = $1", username).Scan(
		&u.ID, &u.Username, &u.Email, &hashedPassword, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	// Guardar la contraseña hasheada para verificación
	u.Password = hashedPassword

	return u, nil
}

// Crear un nuevo usuario
func (r *UserRepository) CreateUser(u User) (User, error) {
	// Verificar si el usuario ya existe
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2", u.Username, u.Email).Scan(&count)
	if err != nil {
		return User{}, err
	}

	if count > 0 {
		return User{}, fmt.Errorf("username or email already exists")
	}

	// Generar ID y timestamps
	u.ID = uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	u.CreatedAt = now
	u.UpdatedAt = now

	// Establecer rol por defecto si no se proporciona
	if u.Role == "" {
		u.Role = "user"
	}

	// Hashear contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	// Insertar usuario
	_, err = r.db.Exec(
		"INSERT INTO users (id, username, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		u.ID, u.Username, u.Email, string(hashedPassword), u.Role, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	// No devolver la contraseña
	u.Password = ""

	return u, nil
}

// Actualizar un usuario
func (r *UserRepository) UpdateUser(u User) (User, error) {
	// Verificar si el usuario existe
	_, err := r.GetUserByID(u.ID)
	if err != nil {
		return User{}, err
	}

	u.UpdatedAt = time.Now().Format(time.RFC3339)

	// Actualizar usuario sin cambiar la contraseña
	_, err = r.db.Exec(
		"UPDATE users SET email = $1, role = $2, updated_at = $3 WHERE id = $4",
		u.Email, u.Role, u.UpdatedAt, u.ID,
	)
	if err != nil {
		return User{}, err
	}

	// No devolver la contraseña
	u.Password = ""

	return u, nil
}

// Eliminar un usuario
func (r *UserRepository) DeleteUser(id string) error {
	// Verificar si el usuario existe
	_, err := r.GetUserByID(id)
	if err != nil {
		return err
	}

	// Eliminar usuario
	_, err = r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// Verificar credenciales de usuario
func (r *UserRepository) VerifyCredentials(username, password string) (User, error) {
	// Obtener usuario por nombre de usuario
	user, err := r.GetUserByUsername(username)
	if err != nil {
		return User{}, fmt.Errorf("invalid credentials")
	}

	// Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return User{}, fmt.Errorf("invalid credentials")
	}

	// No devolver la contraseña
	user.Password = ""

	return user, nil
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
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL,
			role VARCHAR(20) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Generar token JWT
func generateToken(user User, secret string) (string, error) {
	// Crear claims
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "user-service",
			Subject:   user.ID,
		},
	}

	// Crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
	repo := NewUserRepository(db)

	// Crear router
	r := gin.Default()

	// Rutas públicas
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Métricas de Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Ruta de registro
	r.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdUser, err := repo.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdUser)
	})

	// Ruta de login
	r.POST("/login", func(c *gin.Context) {
		var loginReq LoginRequest
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verificar credenciales
		user, err := repo.VerifyCredentials(loginReq.Username, loginReq.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Generar token
		token, err := generateToken(user, config.JWTSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		// Devolver respuesta
		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
			User:  user,
		})
	})

	// Middleware de autenticación
	authMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			// Obtener token del header
			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
				c.Abort()
				return
			}

			// Eliminar prefijo "Bearer " si existe
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}

			// Parsear token
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Verificar método de firma
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}

			// Guardar claims en el contexto
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)

			c.Next()
		}
	}

	// Middleware de autorización para administradores
	adminMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			role, exists := c.Get("role")
			if !exists || role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
				c.Abort()
				return
			}
			c.Next()
		}
	}

	// Grupo de rutas protegidas
	protected := r.Group("/")
	protected.Use(authMiddleware())

	// Rutas de usuarios
	protected.GET("/users", func(c *gin.Context) {
		users, err := repo.GetUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	protected.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		user, err := repo.GetUserByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	protected.GET("/me", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		user, err := repo.GetUserByID(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// Grupo de rutas para administradores
	admin := protected.Group("/admin")
	admin.Use(adminMiddleware())

	admin.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Verificar si el usuario existe
		_, err := repo.GetUserByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Actualizar usuario
		var updateData struct {
			Email string `json:"email" binding:"required,email"`
			Role  string `json:"role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validar rol
		validRoles := map[string]bool{"user": true, "admin": true}
		if !validRoles[updateData.Role] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}

		user := User{
			ID:    id,
			Email: updateData.Email,
			Role:  updateData.Role,
		}

		updatedUser, err := repo.UpdateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedUser)
	})

	admin.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Verificar que no se elimine a sí mismo
		currentUserID, _ := c.Get("userID")
		if currentUserID.(string) == id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete yourself"})
			return
		}

		// Eliminar usuario
		err := repo.DeleteUser(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
		log.Printf("User service listening on port %s\n", config.Port)
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