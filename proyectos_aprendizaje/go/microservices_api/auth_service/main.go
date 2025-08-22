package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Configuración
const (
	PORT       = "8000"
	SECRET_KEY = "super-secret-auth-key" // En producción, usar variables de entorno
)

// Modelos
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type TokenRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Base de datos simulada
var users = []User{
	{
		ID:       "1",
		Username: "testuser",
		Password: "$2a$10$XgXLGk7Vx3zoT9qRk5PKIOMXUa5Eq8RCoZ0CJJdKGQTD.QNXcXALW", // password123
		Email:    "test@example.com",
		Role:     "user",
	},
	{
		ID:       "2",
		Username: "admin",
		Password: "$2a$10$XgXLGk7Vx3zoT9qRk5PKIOMXUa5Eq8RCoZ0CJJdKGQTD.QNXcXALW", // password123
		Email:    "admin@example.com",
		Role:     "admin",
	},
}

// Utilidades de seguridad
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"username": user.Username,
		"email": user.Email,
		"role":  user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	return tokenString, err
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})
}

// Middleware de autenticación
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

// Controladores
func getUserByUsername(username string) (*User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func getUserByID(id string) (*User, error) {
	for _, user := range users {
		if user.ID == id {
			userCopy := user
			userCopy.Password = "" // No devolver la contraseña
			return &userCopy, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func loginHandler(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := getUserByUsername(req.Username)
	if err != nil || !checkPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := generateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   86400, // 24 horas en segundos
	})
}

func getCurrentUserHandler(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	claims := userClaims.(jwt.MapClaims)
	userID, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in token"})
		return
	}

	user, err := getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "auth-service",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func main() {
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
	r.POST("/token", loginHandler)
	r.GET("/health", healthCheckHandler)

	// Rutas protegidas
	protected := r.Group("/")
	protected.Use(authMiddleware())
	protected.GET("/users/me", getCurrentUserHandler)

	// Iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	log.Printf("Auth service starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}