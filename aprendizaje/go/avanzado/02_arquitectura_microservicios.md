# Arquitectura de Microservicios en Go

## Introducción

La arquitectura de microservicios es un enfoque para desarrollar aplicaciones como un conjunto de servicios pequeños, independientes y con bajo acoplamiento. Cada servicio se ejecuta en su propio proceso, se comunica a través de mecanismos ligeros (generalmente API HTTP) y puede ser desplegado de forma independiente. Go, con su eficiencia, concurrencia nativa y capacidad para crear binarios pequeños y autónomos, es una excelente opción para implementar microservicios.

En este módulo, exploraremos en profundidad la arquitectura de microservicios en Go, incluyendo patrones de diseño, comunicación entre servicios, gestión de datos, seguridad, observabilidad y despliegue.

## Fundamentos de Microservicios

### Principios de Microservicios

1. **Servicios pequeños y enfocados**: Cada servicio debe tener una responsabilidad única y bien definida.
2. **Independencia**: Los servicios deben poder desarrollarse, desplegarse y escalarse de forma independiente.
3. **Comunicación a través de APIs**: Los servicios se comunican mediante APIs bien definidas, generalmente HTTP/REST o gRPC.
4. **Descentralización**: Cada servicio gestiona su propia base de datos o almacenamiento.
5. **Resiliencia**: El sistema debe ser tolerante a fallos, con mecanismos para manejar errores y recuperarse.
6. **Automatización**: CI/CD, pruebas automatizadas y monitorización son esenciales.

### Ventajas de Go para Microservicios

1. **Eficiencia**: Go es un lenguaje compilado con un rendimiento cercano a C/C++.
2. **Concurrencia nativa**: Goroutines y canales facilitan la programación concurrente.
3. **Binarios pequeños y autónomos**: No requieren runtime o dependencias externas.
4. **Biblioteca estándar robusta**: Incluye soporte para HTTP, JSON, y más.
5. **Simplicidad**: Go es fácil de aprender y mantener.
6. **Ecosistema maduro**: Frameworks y bibliotecas para microservicios.

### Desventajas y Consideraciones

1. **Complejidad distribuida**: Los sistemas distribuidos son inherentemente más complejos.
2. **Consistencia de datos**: Mantener la consistencia entre servicios puede ser desafiante.
3. **Latencia de red**: La comunicación entre servicios introduce latencia.
4. **Monitorización y depuración**: Requiere herramientas especializadas.
5. **Transacciones distribuidas**: Difíciles de implementar correctamente.

## Diseño de Microservicios en Go

### Estructura de un Microservicio

Una estructura común para un microservicio en Go:

```
/service-name
  /cmd
    /api
      main.go           # Punto de entrada principal
  /internal
    /domain
      models.go         # Modelos de dominio
      service.go        # Lógica de negocio
    /handlers
      http.go           # Manejadores HTTP
      grpc.go           # Manejadores gRPC (si aplica)
    /repository
      postgres.go       # Implementación de repositorio
      mongo.go          # Implementación alternativa
    /config
      config.go         # Configuración
  /pkg
    /middleware         # Middleware compartido
    /logger             # Logging
    /errors             # Manejo de errores
  /api
    /proto             # Definiciones de protobuf (para gRPC)
    /swagger           # Documentación de API
  /scripts             # Scripts de despliegue, migración, etc.
  /deployments         # Configuraciones de despliegue (Docker, K8s)
  go.mod               # Dependencias
  go.sum               # Checksums de dependencias
  Makefile             # Tareas automatizadas
  Dockerfile           # Instrucciones de construcción de imagen
```

### Ejemplo de un Microservicio Básico

```go
// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/service-name/internal/config"
	"github.com/yourusername/service-name/internal/handlers"
	"github.com/yourusername/service-name/internal/repository"
	"github.com/yourusername/service-name/internal/domain"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// Inicializar repositorio
	repo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer repo.Close()

	// Inicializar servicio
	service := domain.NewService(repo)

	// Configurar router
	r := mux.NewRouter()
	handlers.RegisterRoutes(r, service)

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor en una goroutine
	go func() {
		log.Printf("Servidor escuchando en %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Configurar canal para señales de apagado
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Bloquear hasta recibir señal
	<-c

	// Crear contexto con timeout para apagado graceful
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Apagar servidor gracefully
	log.Println("Apagando servidor...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error al apagar el servidor: %v", err)
	}

	log.Println("Servidor apagado correctamente")
}
```

```go
// internal/domain/models.go
package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	GetUser(id string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id string) error
	ListUsers() ([]*User, error)
	Close() error
}
```

```go
// internal/domain/service.go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("usuario no encontrado")
	ErrInvalidInput = errors.New("entrada inválida")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUser(id string) (*User, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.GetUser(id)
}

func (s *Service) CreateUser(name, email string) (*User, error) {
	if name == "" || email == "" {
		return nil, ErrInvalidInput
	}

	now := time.Now()
	user := &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) UpdateUser(id, name, email string) (*User, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.GetUser(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}

	if email != "" {
		user.Email = email
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteUser(id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	return s.repo.DeleteUser(id)
}

func (s *Service) ListUsers() ([]*User, error) {
	return s.repo.ListUsers()
}
```

```go
// internal/handlers/http.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yourusername/service-name/internal/domain"
)

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type updateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(r *mux.Router, service *domain.Service) {
	r.HandleFunc("/users", listUsersHandler(service)).Methods("GET")
	r.HandleFunc("/users", createUserHandler(service)).Methods("POST")
	r.HandleFunc("/users/{id}", getUserHandler(service)).Methods("GET")
	r.HandleFunc("/users/{id}", updateUserHandler(service)).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUserHandler(service)).Methods("DELETE")
}

func listUsersHandler(service *domain.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := service.ListUsers()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, users)
	}
}

func getUserHandler(service *domain.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		user, err := service.GetUser(id)
		if err != nil {
			if err == domain.ErrUserNotFound {
				respondWithError(w, http.StatusNotFound, "Usuario no encontrado")
				return
			}
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, user)
	}
}

func createUserHandler(service *domain.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
			return
		}

		user, err := service.CreateUser(req.Name, req.Email)
		if err != nil {
			if err == domain.ErrInvalidInput {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, user)
	}
}

func updateUserHandler(service *domain.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req updateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
			return
		}

		user, err := service.UpdateUser(id, req.Name, req.Email)
		if err != nil {
			if err == domain.ErrUserNotFound {
				respondWithError(w, http.StatusNotFound, "Usuario no encontrado")
				return
			}
			if err == domain.ErrInvalidInput {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, user)
	}
}

func deleteUserHandler(service *domain.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		err := service.DeleteUser(id)
		if err != nil {
			if err == domain.ErrUserNotFound {
				respondWithError(w, http.StatusNotFound, "Usuario no encontrado")
				return
			}
			if err == domain.ErrInvalidInput {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, errorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error al serializar la respuesta"}`)) // Fallback simple
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
```

```go
// internal/repository/postgres.go
package repository

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/yourusername/service-name/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al verificar la conexión con la base de datos: %w", err)
	}

	// Crear tabla si no existe
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("error al crear la tabla users: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) GetUser(id string) (*domain.User, error) {
	var user domain.User

	row := r.db.QueryRow("SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepository) CreateUser(user *domain.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (id, name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}

	return nil
}

func (r *PostgresRepository) UpdateUser(user *domain.User) error {
	res, err := r.db.Exec(
		"UPDATE users SET name = $1, email = $2, updated_at = $3 WHERE id = $4",
		user.Name, user.Email, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *PostgresRepository) DeleteUser(id string) error {
	res, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *PostgresRepository) ListUsers() ([]*domain.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, created_at, updated_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("error al listar usuarios: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre usuarios: %w", err)
	}

	return users, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
```

```go
// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	LogLevel      string
}

func Load() (*Config, error) {
	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "0.0.0.0")
	serverAddress := fmt.Sprintf("%s:%s", host, port)

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "users")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	databaseURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	return &Config{
		ServerAddress: serverAddress,
		DatabaseURL:   databaseURL,
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
```

### Patrones de Diseño para Microservicios

#### 1. Patrón de Arquitectura Hexagonal (Ports and Adapters)

La arquitectura hexagonal separa la lógica de negocio de los detalles de implementación:

```go
// internal/domain/ports.go
package domain

// Puerto primario (entrada)
type UserService interface {
	GetUser(id string) (*User, error)
	CreateUser(name, email string) (*User, error)
	UpdateUser(id, name, email string) (*User, error)
	DeleteUser(id string) error
	ListUsers() ([]*User, error)
}

// Puerto secundario (salida)
type UserRepository interface {
	GetUser(id string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id string) error
	ListUsers() ([]*User, error)
}
```

#### 2. Patrón CQRS (Command Query Responsibility Segregation)

CQRS separa las operaciones de lectura (queries) de las operaciones de escritura (commands):

```go
// internal/domain/commands.go
package domain

type CreateUserCommand struct {
	Name  string
	Email string
}

type UpdateUserCommand struct {
	ID    string
	Name  string
	Email string
}

type DeleteUserCommand struct {
	ID string
}

type CommandHandler interface {
	HandleCreateUser(cmd CreateUserCommand) (*User, error)
	HandleUpdateUser(cmd UpdateUserCommand) (*User, error)
	HandleDeleteUser(cmd DeleteUserCommand) error
}

// internal/domain/queries.go
package domain

type GetUserQuery struct {
	ID string
}

type ListUsersQuery struct {
	// Parámetros de filtrado, paginación, etc.
	Limit  int
	Offset int
}

type QueryHandler interface {
	HandleGetUser(query GetUserQuery) (*User, error)
	HandleListUsers(query ListUsersQuery) ([]*User, error)
}
```

#### 3. Patrón API Gateway

Un API Gateway actúa como punto de entrada único para múltiples microservicios:

```go
// cmd/gateway/main.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Rutas para el servicio de usuarios
	userRouter := r.PathPrefix("/api/users").Subrouter()
	userRouter.HandleFunc("", forwardToService("http://user-service:8080/users")).Methods("GET", "POST")
	userRouter.HandleFunc("/{id}", forwardToService("http://user-service:8080/users/{id}")).Methods("GET", "PUT", "DELETE")

	// Rutas para el servicio de productos
	productRouter := r.PathPrefix("/api/products").Subrouter()
	productRouter.HandleFunc("", forwardToService("http://product-service:8080/products")).Methods("GET", "POST")
	productRouter.HandleFunc("/{id}", forwardToService("http://product-service:8080/products/{id}")).Methods("GET", "PUT", "DELETE")

	// Middleware común
	r.Use(loggingMiddleware)
	r.Use(authenticationMiddleware)

	// Iniciar servidor
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("API Gateway escuchando en :8080")
	log.Fatal(srv.ListenAndServe())
}

func forwardToService(targetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementación de proxy reverso
		// ...
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar token JWT
		// ...
		next.ServeHTTP(w, r)
	})
}
```

#### 4. Patrón Circuit Breaker

El patrón Circuit Breaker previene fallos en cascada:

```go
// pkg/circuitbreaker/circuitbreaker.go
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
	mutex             sync.Mutex
	failureThreshold  int
	successThreshold  int
	timeoutDuration   time.Duration
	lastFailureTime   time.Time
	failureCount      int
	successCount      int
	state             State
}

func NewCircuitBreaker(failureThreshold, successThreshold int, timeoutDuration time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeoutDuration:  timeoutDuration,
		state:            StateClosed,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mutex.Lock()

	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.timeoutDuration {
			cb.state = StateHalfOpen
			cb.successCount = 0
		} else {
			cb.mutex.Unlock()
			return ErrCircuitOpen
		}
	}

	cb.mutex.Unlock()

	err := fn()

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		if (cb.state == StateClosed && cb.failureCount >= cb.failureThreshold) ||
			(cb.state == StateHalfOpen) {
			cb.state = StateOpen
		}

		return err
	}

	if cb.state == StateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.state = StateClosed
			cb.failureCount = 0
		}
	} else if cb.state == StateClosed {
		cb.failureCount = 0
	}

	return nil
}

// Uso
func callService() error {
	cb := circuitbreaker.NewCircuitBreaker(5, 2, 1*time.Minute)

	err := cb.Execute(func() error {
		// Llamada al servicio
		return nil
	})

	if err == circuitbreaker.ErrCircuitOpen {
		// Manejar circuito abierto (fallback, caché, etc.)
		return errors.New("servicio no disponible")
	}

	return err
}
```

#### 5. Patrón Saga

El patrón Saga gestiona transacciones distribuidas:

```go
// pkg/saga/saga.go
package saga

import (
	"context"
	"fmt"
)

type Step struct {
	Execute   func(ctx context.Context) error
	Compensate func(ctx context.Context) error
	Name      string
}

type Saga struct {
	steps []Step
}

func NewSaga() *Saga {
	return &Saga{steps: []Step{}}
}

func (s *Saga) AddStep(step Step) *Saga {
	s.steps = append(s.steps, step)
	return s
}

func (s *Saga) Execute(ctx context.Context) error {
	var executedSteps []Step

	for _, step := range s.steps {
		err := step.Execute(ctx)
		if err != nil {
			// Fallo en este paso, compensar los pasos anteriores
			s.compensate(ctx, executedSteps)
			return fmt.Errorf("error en paso %s: %w", step.Name, err)
		}
		executedSteps = append(executedSteps, step)
	}

	return nil
}

func (s *Saga) compensate(ctx context.Context, executedSteps []Step) {
	// Compensar en orden inverso
	for i := len(executedSteps) - 1; i >= 0; i-- {
		step := executedSteps[i]
		err := step.Compensate(ctx)
		if err != nil {
			// Log error pero continuar con la compensación
			fmt.Printf("Error al compensar paso %s: %v\n", step.Name, err)
		}
	}
}

// Uso
func createOrder(ctx context.Context, orderID, userID string, amount float64) error {
	saga := saga.NewSaga()

	// Paso 1: Verificar inventario
	saga.AddStep(saga.Step{
		Name: "VerificarInventario",
		Execute: func(ctx context.Context) error {
			// Verificar inventario
			return nil
		},
		Compensate: func(ctx context.Context) error {
			// No es necesario compensar
			return nil
		},
	})

	// Paso 2: Reservar productos
	saga.AddStep(saga.Step{
		Name: "ReservarProductos",
		Execute: func(ctx context.Context) error {
			// Reservar productos
			return nil
		},
		Compensate: func(ctx context.Context) error {
			// Liberar productos reservados
			return nil
		},
	})

	// Paso 3: Procesar pago
	saga.AddStep(saga.Step{
		Name: "ProcesarPago",
		Execute: func(ctx context.Context) error {
			// Procesar pago
			return nil
		},
		Compensate: func(ctx context.Context) error {
			// Reembolsar pago
			return nil
		},
	})

	// Paso 4: Crear orden
	saga.AddStep(saga.Step{
		Name: "CrearOrden",
		Execute: func(ctx context.Context) error {
			// Crear orden
			return nil
		},
		Compensate: func(ctx context.Context) error {
			// Cancelar orden
			return nil
		},
	})

	// Ejecutar saga
	return saga.Execute(ctx)
}
```

## Comunicación entre Microservicios

### REST

REST es un estilo arquitectónico común para APIs HTTP:

```go
// pkg/client/user_client.go
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/service-name/internal/domain"
)

type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *UserClient) GetUser(id string) (*domain.User, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/users/%s", c.baseURL, id))
	if err != nil {
		return nil, fmt.Errorf("error al realizar solicitud: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrUserNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error del servidor: %d", resp.StatusCode)
	}

	var user domain.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	return &user, nil
}

func (c *UserClient) CreateUser(name, email string) (*domain.User, error) {
	reqBody, err := json.Marshal(map[string]string{
		"name":  name,
		"email": email,
	})
	if err != nil {
		return nil, fmt.Errorf("error al codificar solicitud: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/users", c.baseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("error al realizar solicitud: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error del servidor: %d", resp.StatusCode)
	}

	var user domain.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	return &user, nil
}
```

### gRPC

gRPC es un framework RPC de alto rendimiento:

```protobuf
// api/proto/user.proto
syntax = "proto3";

package user;

option go_package = "github.com/yourusername/service-name/api/proto/user";

import "google/protobuf/timestamp.proto";

service UserService {
  rpc GetUser(GetUserRequest) returns (User) {}
  rpc CreateUser(CreateUserRequest) returns (User) {}
  rpc UpdateUser(UpdateUserRequest) returns (User) {}
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse) {}
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse) {}
}

message User {
  string id = 1;
  string name = 2;
  string email = 3;
  google.protobuf.Timestamp created_at = 4;
  google.protobuf.Timestamp updated_at = 5;
}

message GetUserRequest {
  string id = 1;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
}

message UpdateUserRequest {
  string id = 1;
  string name = 2;
  string email = 3;
}

message DeleteUserRequest {
  string id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}

message ListUsersRequest {
  int32 limit = 1;
  int32 offset = 2;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}
```

```go
// internal/handlers/grpc.go
package handlers

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/service-name/api/proto/user"
	"github.com/yourusername/service-name/internal/domain"
)

type GRPCServer struct {
	service domain.UserService
	user.UnimplementedUserServiceServer
}

func NewGRPCServer(service domain.UserService) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	u, err := s.service.GetUser(req.Id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "usuario no encontrado")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	createdAt, err := ptypes.TimestampProto(u.CreatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "error al convertir timestamp")
	}

	updatedAt, err := ptypes.TimestampProto(u.UpdatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "error al convertir timestamp")
	}

	return &user.User{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *GRPCServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	u, err := s.service.CreateUser(req.Name, req.Email)
	if err != nil {
		if err == domain.ErrInvalidInput {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	createdAt, err := ptypes.TimestampProto(u.CreatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "error al convertir timestamp")
	}

	updatedAt, err := ptypes.TimestampProto(u.UpdatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "error al convertir timestamp")
	}

	return &user.User{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// Implementar los demás métodos...
```

### Mensajería Asíncrona

La mensajería asíncrona permite la comunicación desacoplada entre servicios:

```go
// pkg/messaging/rabbitmq.go
package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error al abrir canal: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("error al cerrar canal: %w", err)
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("error al cerrar conexión: %w", err)
	}

	return nil
}

func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,  // nombre
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (r *RabbitMQ) Publish(exchange, routingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error al serializar mensaje: %w", err)
	}

	return r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) Consume(ctx context.Context, queueName string, handler func([]byte) error) error {
	msgs, err := r.channel.Consume(
		queueName, // queue
		"",       // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("error al registrar consumidor: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Consumidor detenido")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Canal de mensajes cerrado")
					return
				}

				err := handler(msg.Body)
				if err != nil {
					log.Printf("Error al procesar mensaje: %v", err)
					// Rechazar mensaje
					msg.Nack(false, true)
				} else {
					// Confirmar mensaje
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// Uso como productor
func publishUserCreated(rabbitMQ *messaging.RabbitMQ, user *domain.User) error {
	return rabbitMQ.Publish(
		"",           // exchange
		"user.created", // routing key
		user,
	)
}

// Uso como consumidor
func consumeUserCreated(ctx context.Context, rabbitMQ *messaging.RabbitMQ) error {
	queue, err := rabbitMQ.DeclareQueue("user.created")
	if err != nil {
		return err
	}

	return rabbitMQ.Consume(ctx, queue.Name, func(body []byte) error {
		var user domain.User
		if err := json.Unmarshal(body, &user); err != nil {
			return fmt.Errorf("error al deserializar mensaje: %w", err)
		}

		log.Printf("Usuario creado: %s (%s)", user.Name, user.ID)
		// Procesar evento...
		return nil
	})
}
```

## Gestión de Datos

### Patrón Database per Service

Cada microservicio gestiona su propia base de datos:

```go
// internal/repository/postgres.go (ya mostrado anteriormente)
```

### Consultas entre Servicios

Las consultas entre servicios pueden implementarse de varias formas:

#### 1. API Composition

```go
// internal/service/order_service.go
package service

import (
	"context"
	"fmt"

	"github.com/yourusername/order-service/internal/domain"
	"github.com/yourusername/order-service/pkg/client"
)

type OrderService struct {
	repo       domain.OrderRepository
	userClient *client.UserClient
}

func NewOrderService(repo domain.OrderRepository, userClient *client.UserClient) *OrderService {
	return &OrderService{
		repo:       repo,
		userClient: userClient,
	}
}

func (s *OrderService) GetOrderWithUser(ctx context.Context, orderID string) (*domain.OrderWithUser, error) {
	// Obtener orden
	order, err := s.repo.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener orden: %w", err)
	}

	// Obtener usuario
	user, err := s.userClient.GetUser(order.UserID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	// Combinar datos
	return &domain.OrderWithUser{
		Order: *order,
		User:  *user,
	}, nil
}
```

#### 2. CQRS con Vistas Materializadas

```go
// internal/repository/order_view_repository.go
package repository

import (
	"database/sql"
	"fmt"

	"github.com/yourusername/order-service/internal/domain"
)

type OrderViewRepository struct {
	db *sql.DB
}

func NewOrderViewRepository(db *sql.DB) *OrderViewRepository {
	return &OrderViewRepository{db: db}
}

func (r *OrderViewRepository) GetOrderWithUser(orderID string) (*domain.OrderWithUser, error) {
	query := `
		SELECT 
			o.id, o.total, o.status, o.created_at,
			u.id, u.name, u.email
		FROM order_user_view o
		JOIN users u ON o.user_id = u.id
		WHERE o.id = $1
	`

	row := r.db.QueryRow(query, orderID)

	var order domain.Order
	var user domain.User

	err := row.Scan(
		&order.ID, &order.Total, &order.Status, &order.CreatedAt,
		&user.ID, &user.Name, &user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("error al obtener vista de orden: %w", err)
	}

	return &domain.OrderWithUser{
		Order: order,
		User:  user,
	}, nil
}
```

### Consistencia Eventual

La consistencia eventual es común en sistemas distribuidos:

```go
// internal/handlers/event_handlers.go
package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/yourusername/order-service/internal/domain"
	"github.com/yourusername/order-service/internal/repository"
)

type UserEventHandler struct {
	userViewRepo *repository.UserViewRepository
}

func NewUserEventHandler(userViewRepo *repository.UserViewRepository) *UserEventHandler {
	return &UserEventHandler{userViewRepo: userViewRepo}
}

func (h *UserEventHandler) HandleUserCreated(ctx context.Context, data []byte) error {
	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return err
	}

	// Actualizar vista materializada
	return h.userViewRepo.UpsertUser(&user)
}

func (h *UserEventHandler) HandleUserUpdated(ctx context.Context, data []byte) error {
	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return err
	}

	// Actualizar vista materializada
	return h.userViewRepo.UpsertUser(&user)
}

// Implementación del repositorio de vista
type UserViewRepository struct {
	db *sql.DB
}

func (r *UserViewRepository) UpsertUser(user *domain.User) error {
	// Upsert en la tabla de vista materializada
	// ...
	return nil
}
```

## Seguridad

### Autenticación y Autorización

#### JWT (JSON Web Tokens)

```go
// pkg/auth/jwt.go
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey []byte
	tokenDuration time.Duration
}

func NewJWTService(secretKey string, tokenDuration time.Duration) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

func (s *JWTService) GenerateToken(userID, role string) (string, error) {
	expirationTime := time.Now().Add(s.tokenDuration)

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "your-service-name",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("error al firmar token: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al validar token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return claims, nil
}
```

#### Middleware de Autenticación

```go
// pkg/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yourusername/service-name/pkg/auth"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func AuthMiddleware(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraer token del encabezado Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Autorización requerida", http.StatusUnauthorized)
				return
			}

			// Formato esperado: "Bearer {token}"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Formato de autorización inválido", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validar token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			// Añadir claims al contexto
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)

			// Continuar con el siguiente handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Middleware de autorización basado en roles
func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				http.Error(w, "Rol no encontrado en el contexto", http.StatusUnauthorized)
				return
			}

			authorized := false
			for _, allowedRole := range roles {
				if role == allowedRole {
					authorized = true
					break
				}
			}

			if !authorized {
				http.Error(w, "Acceso denegado", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
```

#### Seguridad TLS

```go
// cmd/api/main.go (fragmento)
func main() {
	// ... código anterior

	// Configurar servidor HTTPS
	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      r,
		TLSConfig:    getTLSConfig(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor en una goroutine
	go func() {
		log.Printf("Servidor escuchando en %s", cfg.ServerAddress)
		if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// ... resto del código
}

func getTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
```

### Gestión de Secretos

```go
// pkg/secrets/vault.go
package secrets

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	client *vault.Client
	path   string
}

func NewVaultClient(address, token, path string) (*VaultClient, error) {
	config := vault.DefaultConfig()
	config.Address = address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("error al crear cliente de Vault: %w", err)
	}

	client.SetToken(token)

	return &VaultClient{
		client: client,
		path:   path,
	}, nil
}

func (v *VaultClient) GetSecret(key string) (string, error) {
	secret, err := v.client.Logical().Read(fmt.Sprintf("%s/%s", v.path, key))
	if err != nil {
		return "", fmt.Errorf("error al leer secreto: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("secreto no encontrado")
	}

	value, ok := secret.Data["value"].(string)
	if !ok {
		return "", fmt.Errorf("valor no encontrado o no es una cadena")
	}

	return value, nil
}

func (v *VaultClient) SetSecret(key, value string) error {
	_, err := v.client.Logical().Write(
		fmt.Sprintf("%s/%s", v.path, key),
		map[string]interface{}{
			"value": value,
		},
	)

	if err != nil {
		return fmt.Errorf("error al escribir secreto: %w", err)
	}

	return nil
}
```

## Observabilidad

### Logging

```go
// pkg/logger/logger.go
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
}

func NewLogger(serviceName, level string) *Logger {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	log := zerolog.New(os.Stdout).With().Timestamp().Str("service", serviceName).Logger()

	return &Logger{log: log}
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	event := l.log.Info()

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(msg)
}

func (l *Logger) Error(msg string, err error, fields map[string]interface{}) {
	event := l.log.Error().Err(err)

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(msg)
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	event := l.log.Debug()

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(msg)
}

func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	event := l.log.Warn()

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(msg)
}
```

### Métricas

```go
// pkg/metrics/prometheus.go
package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	requestCounter   *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	requestInFlight  *prometheus.GaugeVec
	dependencyErrors *prometheus.CounterVec
}

func NewPrometheusMetrics(serviceName string) *PrometheusMetrics {
	requestCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "method", "path", "status"},
	)

	requestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "path"},
	)

	requestInFlight := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests in flight",
		},
		[]string{"service", "method", "path"},
	)

	dependencyErrors := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dependency_errors_total",
			Help: "Total number of dependency errors",
		},
		[]string{"service", "dependency", "operation"},
	)

	return &PrometheusMetrics{
		requestCounter:   requestCounter,
		requestDuration:  requestDuration,
		requestInFlight:  requestInFlight,
		dependencyErrors: dependencyErrors,
	}
}

func (m *PrometheusMetrics) InstrumentHandler(next http.Handler, path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service := "user-service" // Obtener de la configuración

		m.requestInFlight.WithLabelValues(service, r.Method, path).Inc()
		defer m.requestInFlight.WithLabelValues(service, r.Method, path).Dec()

		start := time.Now()

		// Wrapper para capturar el código de estado
		ww := NewResponseWriter(w)

		// Ejecutar el handler
		next.ServeHTTP(ww, r)

		// Registrar métricas
		status := ww.Status()
		duration := time.Since(start).Seconds()

		m.requestCounter.WithLabelValues(service, r.Method, path, http.StatusText(status)).Inc()
		m.requestDuration.WithLabelValues(service, r.Method, path).Observe(duration)
	})
}

func (m *PrometheusMetrics) RecordDependencyError(dependency, operation string) {
	service := "user-service" // Obtener de la configuración
	m.dependencyErrors.WithLabelValues(service, dependency, operation).Inc()
}

func (m *PrometheusMetrics) Handler() http.Handler {
	return promhttp.Handler()
}

// ResponseWriter personalizado para capturar el código de estado
type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseWriter) Status() int {
	return w.status
}
```

### Trazabilidad

```go
// pkg/tracing/jaeger.go
package tracing

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

type Tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func NewTracer(serviceName string, jaegerEndpoint string) (*Tracer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerEndpoint,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, fmt.Errorf("error al crear tracer: %w", err)
	}

	opentracing.SetGlobalTracer(tracer)

	return &Tracer{
		tracer: tracer,
		closer: closer,
	}, nil
}

func (t *Tracer) Close() error {
	return t.closer.Close()
}

func (t *Tracer) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spanCtx, _ := t.tracer.Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header),
			)

			span := t.tracer.StartSpan(r.URL.Path, ext.RPCServerOption(spanCtx))
			defer span.Finish()

			// Añadir tags al span
			ext.HTTPMethod.Set(span, r.Method)
			ext.HTTPUrl.Set(span, r.URL.String())

			// Crear un nuevo contexto con el span
			ctx := opentracing.ContextWithSpan(r.Context(), span)

			// Wrapper para capturar el código de estado
			ww := NewResponseWriter(w)

			// Ejecutar el handler con el nuevo contexto
			next.ServeHTTP(ww, r.WithContext(ctx))

			// Añadir código de estado al span
			ext.HTTPStatusCode.Set(span, uint16(ww.Status()))
		})
	}
}

func StartSpanFromContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	return span, ctx
}

// Ejemplo de uso en un servicio
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Service.GetUser")
	defer span.Finish()

	// Añadir información al span
	span.SetTag("user.id", id)

	// Llamada al repositorio
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		// Registrar error en el span
		ext.Error.Set(span, true)
		span.SetTag("error.message", err.Error())
		return nil, err
	}

	return user, nil
}

// Implementación en el repositorio
func (r *Repository) GetUser(ctx context.Context, id string) (*User, error) {
	span, _ := tracing.StartSpanFromContext(ctx, "Repository.GetUser")
	defer span.Finish()

	// Añadir información al span
	span.SetTag("db.query", "SELECT * FROM users WHERE id = ?")
	span.SetTag("db.id", id)

	// Ejecutar consulta
	// ...

	return user, nil
}
```

## Despliegue de Microservicios

### Contenedores con Docker

```dockerfile
# Dockerfile
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Imagen final
FROM alpine:3.15

WORKDIR /app

# Copiar binario compilado
COPY --from=builder /app/main .

# Copiar archivos de configuración
COPY --from=builder /app/configs ./configs

# Exponer puerto
EXPOSE 8080

# Ejecutar aplicación
CMD ["./main"]
```

### Orquestación con Kubernetes

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: yourusername/user-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: user-service-config
              key: db_host
        - name: DB_PORT
          valueFrom:
            configMapKeyRef:
              name: user-service-config
              key: db_port
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: db_user
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: db_password
        - name: DB_NAME
          valueFrom:
            configMapKeyRef:
              name: user-service-config
              key: db_name
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          limits:
            cpu: "500m"
            memory: "512Mi"
          requests:
            cpu: "100m"
            memory: "128Mi"
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: user-service-config
data:
  db_host: "postgres"
  db_port: "5432"
  db_name: "users"
---
apiVersion: v1
kind: Secret
metadata:
  name: user-service-secrets
type: Opaque
data:
  db_user: dXNlcg==      # user (base64)
  db_password: cGFzc3dvcmQ=  # password (base64)
```

### Despliegue Continuo

```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.18

    - name: Test
      run: go test -v ./...

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
    - uses: actions/checkout@v2

    - name: Build and push Docker image
      uses: docker/build-push-action@v2
      with:
        context: .
        push: true
        tags: yourusername/user-service:latest

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2

    - name: Set up kubectl
      uses: azure/setup-kubectl@v1

    - name: Set Kubernetes context
      uses: azure/k8s-set-context@v1
      with:
        kubeconfig: ${{ secrets.KUBE_CONFIG }}

    - name: Deploy to Kubernetes
      run: |
        kubectl apply -f kubernetes/deployment.yaml
        kubectl rollout restart deployment/user-service
```

## Conclusiones

La arquitectura de microservicios en Go ofrece numerosas ventajas para sistemas distribuidos, incluyendo:

1. **Escalabilidad**: Cada servicio puede escalarse de forma independiente según sus necesidades.
2. **Resiliencia**: El fallo de un servicio no afecta a todo el sistema.
3. **Tecnología heterogénea**: Cada servicio puede utilizar la tecnología más adecuada para su función.
4. **Despliegue independiente**: Los servicios pueden desplegarse sin afectar al resto del sistema.
5. **Equipos autónomos**: Diferentes equipos pueden trabajar en diferentes servicios.

Sin embargo, también presenta desafíos:

1. **Complejidad distribuida**: Los sistemas distribuidos son inherentemente más complejos.
2. **Consistencia de datos**: Mantener la consistencia entre servicios es desafiante.
3. **Observabilidad**: Se requieren herramientas especializadas para monitorizar y depurar.
4. **Latencia de red**: La comunicación entre servicios introduce latencia.
5. **Transacciones distribuidas**: Difíciles de implementar correctamente.

Go es particularmente adecuado para microservicios debido a su eficiencia, concurrencia nativa, binarios pequeños y autónomos, y su biblioteca estándar robusta. Con las herramientas y patrones adecuados, Go permite construir sistemas de microservicios escalables, resilientes y mantenibles.

## Ejercicios Prácticos

1. **Implementar un sistema básico de microservicios**:
   - Crear un servicio de usuarios con operaciones CRUD.
   - Crear un servicio de productos con operaciones CRUD.
   - Implementar un API Gateway que enrute las solicitudes a los servicios correspondientes.
   - Utilizar Docker y Docker Compose para ejecutar los servicios localmente.

2. **Añadir comunicación asíncrona**:
   - Implementar un broker de mensajes (RabbitMQ o NATS).
   - Modificar los servicios para publicar eventos cuando ocurren cambios.
   - Crear un servicio de notificaciones que consuma eventos y envíe notificaciones.

3. **Implementar patrones de resiliencia**:
   - Añadir Circuit Breaker a las llamadas entre servicios.
   - Implementar reintentos con backoff exponencial.
   - Añadir timeouts a todas las operaciones externas.

4. **Añadir observabilidad**:
   - Configurar logging estructurado con correlationID.
   - Implementar métricas con Prometheus.
   - Añadir trazabilidad con Jaeger.
   - Crear dashboards en Grafana para visualizar métricas.

5. **Desplegar en Kubernetes**:
   - Crear manifiestos de Kubernetes para todos los servicios.
   - Configurar Ingress para el API Gateway.
   - Implementar health checks y readiness probes.
   - Configurar auto-scaling basado en métricas.

## Referencias

1. Sam Newman, "Building Microservices"
2. Chris Richardson, "Microservices Patterns"
3. Documentación oficial de Go: https://golang.org/doc/
4. Kubernetes: https://kubernetes.io/docs/
5. Docker: https://docs.docker.com/
6. gRPC: https://grpc.io/docs/
7. Prometheus: https://prometheus.io/docs/
8. Jaeger: https://www.jaegertracing.io/docs/
```