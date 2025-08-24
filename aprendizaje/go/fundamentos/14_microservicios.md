# Microservicios en Go

## Introducción a los Microservicios

Los microservicios son un enfoque arquitectónico para el desarrollo de software donde una aplicación se construye como un conjunto de servicios pequeños e independientes. Cada servicio se ejecuta en su propio proceso y se comunica con otros servicios a través de mecanismos ligeros, generalmente API HTTP.

En contraste con las aplicaciones monolíticas tradicionales, los microservicios ofrecen:

- **Despliegue independiente**: Cada servicio puede ser desplegado sin afectar a otros.
- **Escalabilidad selectiva**: Los servicios con mayor demanda pueden escalarse independientemente.
- **Diversidad tecnológica**: Diferentes servicios pueden utilizar diferentes tecnologías.
- **Resiliencia mejorada**: El fallo de un servicio no necesariamente afecta a toda la aplicación.
- **Equipos autónomos**: Diferentes equipos pueden trabajar en diferentes servicios.

Go, con su simplicidad, eficiencia y soporte nativo para concurrencia, es una excelente opción para implementar microservicios.

## Frameworks para Microservicios en Go

### Go estándar (net/http)

Go incluye un potente paquete HTTP en su biblioteca estándar que es suficiente para construir microservicios simples.

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 999.99},
	{ID: 2, Name: "Smartphone", Price: 699.99},
	{ID: 3, Name: "Tablet", Price: 399.99},
}

func main() {
	// Rutas
	http.HandleFunc("/products", handleProducts)
	http.HandleFunc("/products/", handleProduct)

	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Devolver todos los productos
		json.NewEncoder(w).Encode(products)
	case "POST":
		// Añadir un nuevo producto
		var product Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Asignar un nuevo ID
		product.ID = len(products) + 1
		products = append(products, product)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handleProduct(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del producto de la URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Buscar producto por ID
	var product *Product
	for i, p := range products {
		if p.ID == id {
			product = &products[i]
			break
		}
	}

	if product == nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		// Devolver el producto
		json.NewEncoder(w).Encode(product)
	case "PUT":
		// Actualizar el producto
		var updatedProduct Product
		err := json.NewDecoder(r.Body).Decode(&updatedProduct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Mantener el mismo ID
		updatedProduct.ID = id
		*product = updatedProduct

		json.NewEncoder(w).Encode(product)
	case "DELETE":
		// Eliminar el producto
		for i, p := range products {
			if p.ID == id {
				products = append(products[:i], products[i+1:]...)
				break
			}
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
```

### Gin

Gin es un framework web HTTP de alto rendimiento que facilita la creación de API RESTful.

```go
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 999.99},
	{ID: 2, Name: "Smartphone", Price: 699.99},
	{ID: 3, Name: "Tablet", Price: 399.99},
}

func main() {
	r := gin.Default()

	// Rutas
	r.GET("/products", getProducts)
	r.GET("/products/:id", getProduct)
	r.POST("/products", createProduct)
	r.PUT("/products/:id", updateProduct)
	r.DELETE("/products/:id", deleteProduct)

	// Iniciar servidor
	r.Run(":8080")
}

func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

func getProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for _, product := range products {
		if product.ID == id {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}

func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Asignar un nuevo ID
	product.ID = len(products) + 1
	products = append(products, product)

	c.JSON(http.StatusCreated, product)
}

func updateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var productIndex = -1
	for i, product := range products {
		if product.ID == id {
			productIndex = i
			break
		}
	}

	if productIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	var updatedProduct Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mantener el mismo ID
	updatedProduct.ID = id
	products[productIndex] = updatedProduct

	c.JSON(http.StatusOK, updatedProduct)
}

func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}
```

### Echo

Echo es un framework web de alto rendimiento y minimalista para Go.

```go
package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 999.99},
	{ID: 2, Name: "Smartphone", Price: 699.99},
	{ID: 3, Name: "Tablet", Price: 399.99},
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Rutas
	e.GET("/products", getProducts)
	e.GET("/products/:id", getProduct)
	e.POST("/products", createProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)

	// Iniciar servidor
	e.Logger.Fatal(e.Start(":8080"))
}

func getProducts(c echo.Context) error {
	return c.JSON(http.StatusOK, products)
}

func getProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	for _, product := range products {
		if product.ID == id {
			return c.JSON(http.StatusOK, product)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Producto no encontrado"})
}

func createProduct(c echo.Context) error {
	var product Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Asignar un nuevo ID
	product.ID = len(products) + 1
	products = append(products, product)

	return c.JSON(http.StatusCreated, product)
}

func updateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	var productIndex = -1
	for i, product := range products {
		if product.ID == id {
			productIndex = i
			break
		}
	}

	if productIndex == -1 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Producto no encontrado"})
	}

	var updatedProduct Product
	if err := c.Bind(&updatedProduct); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Mantener el mismo ID
	updatedProduct.ID = id
	products[productIndex] = updatedProduct

	return c.JSON(http.StatusOK, updatedProduct)
}

func deleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Producto no encontrado"})
}
```

### Go Kit

Go Kit es un conjunto de paquetes que ayudan a estructurar y construir microservicios robustos.

```go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// Modelo
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Repositorio
type ProductRepository interface {
	GetProducts() []Product
	GetProduct(id int) (Product, error)
	CreateProduct(product Product) Product
	UpdateProduct(id int, product Product) (Product, error)
	DeleteProduct(id int) error
}

type inMemoryProductRepository struct {
	products []Product
}

func NewInMemoryProductRepository() ProductRepository {
	return &inMemoryProductRepository{
		products: []Product{
			{ID: 1, Name: "Laptop", Price: 999.99},
			{ID: 2, Name: "Smartphone", Price: 699.99},
			{ID: 3, Name: "Tablet", Price: 399.99},
		},
	}
}

func (r *inMemoryProductRepository) GetProducts() []Product {
	return r.products
}

func (r *inMemoryProductRepository) GetProduct(id int) (Product, error) {
	for _, product := range r.products {
		if product.ID == id {
			return product, nil
		}
	}
	return Product{}, errors.New("producto no encontrado")
}

func (r *inMemoryProductRepository) CreateProduct(product Product) Product {
	// Asignar un nuevo ID
	product.ID = len(r.products) + 1
	r.products = append(r.products, product)
	return product
}

func (r *inMemoryProductRepository) UpdateProduct(id int, product Product) (Product, error) {
	for i, p := range r.products {
		if p.ID == id {
			// Mantener el mismo ID
			product.ID = id
			r.products[i] = product
			return product, nil
		}
	}
	return Product{}, errors.New("producto no encontrado")
}

func (r *inMemoryProductRepository) DeleteProduct(id int) error {
	for i, product := range r.products {
		if product.ID == id {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}
	return errors.New("producto no encontrado")
}

// Servicio
type ProductService interface {
	GetProducts(ctx context.Context) ([]Product, error)
	GetProduct(ctx context.Context, id int) (Product, error)
	CreateProduct(ctx context.Context, product Product) (Product, error)
	UpdateProduct(ctx context.Context, id int, product Product) (Product, error)
	DeleteProduct(ctx context.Context, id int) error
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetProducts(ctx context.Context) ([]Product, error) {
	return s.repo.GetProducts(), nil
}

func (s *productService) GetProduct(ctx context.Context, id int) (Product, error) {
	return s.repo.GetProduct(id)
}

func (s *productService) CreateProduct(ctx context.Context, product Product) (Product, error) {
	return s.repo.CreateProduct(product), nil
}

func (s *productService) UpdateProduct(ctx context.Context, id int, product Product) (Product, error) {
	return s.repo.UpdateProduct(id, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.DeleteProduct(id)
}

// Endpoints
type getProductsRequest struct{}

type getProductsResponse struct {
	Products []Product `json:"products"`
	Err      string    `json:"err,omitempty"`
}

func makeGetProductsEndpoint(s ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		products, err := s.GetProducts(ctx)
		if err != nil {
			return getProductsResponse{Err: err.Error()}, nil
		}
		return getProductsResponse{Products: products}, nil
	}
}

type getProductRequest struct {
	ID int
}

type getProductResponse struct {
	Product Product `json:"product,omitempty"`
	Err     string  `json:"err,omitempty"`
}

func makeGetProductEndpoint(s ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getProductRequest)
		product, err := s.GetProduct(ctx, req.ID)
		if err != nil {
			return getProductResponse{Err: err.Error()}, nil
		}
		return getProductResponse{Product: product}, nil
	}
}

type createProductRequest struct {
	Product Product
}

type createProductResponse struct {
	Product Product `json:"product"`
	Err     string  `json:"err,omitempty"`
}

func makeCreateProductEndpoint(s ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createProductRequest)
		product, err := s.CreateProduct(ctx, req.Product)
		if err != nil {
			return createProductResponse{Err: err.Error()}, nil
		}
		return createProductResponse{Product: product}, nil
	}
}

// Transporte HTTP
func decodeGetProductsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return getProductsRequest{}, nil
}

func decodeGetProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}
	return getProductRequest{ID: id}, nil
}

func decodeCreateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, err
	}
	return createProductRequest{Product: product}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	// Crear repositorio y servicio
	repo := NewInMemoryProductRepository()
	svc := NewProductService(repo)

	// Crear endpoints
	getProductsEndpoint := makeGetProductsEndpoint(svc)
	getProductEndpoint := makeGetProductEndpoint(svc)
	createProductEndpoint := makeCreateProductEndpoint(svc)

	// Configurar transporte HTTP
	r := mux.NewRouter()

	r.Methods("GET").Path("/products").Handler(http.NewServer(
		getProductsEndpoint,
		decodeGetProductsRequest,
		encodeResponse,
		[]http.ServerOption{
			http.ServerErrorEncoder(encodeError),
		}...,
	))

	r.Methods("GET").Path("/products/{id}").Handler(http.NewServer(
		getProductEndpoint,
		decodeGetProductRequest,
		encodeResponse,
		[]http.ServerOption{
			http.ServerErrorEncoder(encodeError),
		}...,
	))

	r.Methods("POST").Path("/products").Handler(http.NewServer(
		createProductEndpoint,
		decodeCreateProductRequest,
		encodeResponse,
		[]http.ServerOption{
			http.ServerErrorEncoder(encodeError),
		}...,
	))

	// Iniciar servidor
	logger.Log("transport", "HTTP", "addr", ":8080")
	logger.Log("msg", "HTTP server started")
	logger.Log("err", http.ListenAndServe(":8080", r))
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
```

## Comunicación entre Microservicios

### API REST

La comunicación a través de API REST es uno de los métodos más comunes para la interacción entre microservicios.

```go
// order_service.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Order struct {
	ID        int     `json:"id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var orders = []Order{}

func main() {
	r := gin.Default()

	r.POST("/orders", createOrder)
	r.GET("/orders", getOrders)
	r.GET("/orders/:id", getOrder)

	r.Run(":8081")
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener información del producto desde el servicio de productos
	productURL := fmt.Sprintf("http://localhost:8080/products/%d", order.ProductID)
	resp, err := http.Get(productURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al comunicarse con el servicio de productos"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Producto no encontrado"})
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer la respuesta del servicio de productos"})
		return
	}

	var product Product
	if err := json.Unmarshal(body, &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la respuesta del servicio de productos"})
		return
	}

	// Calcular el total
	order.Total = product.Price * float64(order.Quantity)

	// Asignar un nuevo ID
	order.ID = len(orders) + 1
	orders = append(orders, order)

	c.JSON(http.StatusCreated, order)
}

func getOrders(c *gin.Context) {
	c.JSON(http.StatusOK, orders)
}

func getOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for _, order := range orders {
		if order.ID == id {
			c.JSON(http.StatusOK, order)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
}
```

### gRPC

gRPC es un framework de RPC (Remote Procedure Call) de alto rendimiento que puede conectar servicios en y entre centros de datos.

```protobuf
// product.proto
syntax = "proto3";

package product;

service ProductService {
    rpc GetProduct (GetProductRequest) returns (GetProductResponse) {}
    rpc CreateProduct (CreateProductRequest) returns (CreateProductResponse) {}
}

message GetProductRequest {
    int32 id = 1;
}

message GetProductResponse {
    Product product = 1;
}

message CreateProductRequest {
    Product product = 1;
}

message CreateProductResponse {
    Product product = 1;
}

message Product {
    int32 id = 1;
    string name = 2;
    double price = 3;
}
```

```go
// product_server.go
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "path/to/product"
)

type server struct {
	pb.UnimplementedProductServiceServer
	products []*pb.Product
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	for _, product := range s.products {
		if product.Id == req.Id {
			return &pb.GetProductResponse{Product: product}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Producto con ID %d no encontrado", req.Id)
}

func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := req.GetProduct()
	
	// Asignar un nuevo ID
	product.Id = int32(len(s.products) + 1)
	
	// Añadir el producto
	s.products = append(s.products, product)
	
	return &pb.CreateProductResponse{Product: product}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{
		products: []*pb.Product{
			{Id: 1, Name: "Laptop", Price: 999.99},
			{Id: 2, Name: "Smartphone", Price: 699.99},
			{Id: 3, Name: "Tablet", Price: 399.99},
		},
	})
	
	log.Println("Server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
```

```go
// product_client.go
package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "path/to/product"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	
	client := pb.NewProductServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	// Obtener un producto
	productResp, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: 1})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %v", productResp.GetProduct())
	
	// Crear un producto
	newProduct := &pb.Product{
		Name:  "Headphones",
		Price: 149.99,
	}
	
	createResp, err := client.CreateProduct(ctx, &pb.CreateProductRequest{Product: newProduct})
	if err != nil {
		log.Fatalf("Could not create product: %v", err)
	}
	log.Printf("Created product: %v", createResp.GetProduct())
}
```

### Mensajería Asíncrona

La mensajería asíncrona permite la comunicación entre servicios sin bloquear el servicio que envía el mensaje.

```go
// publisher.go (usando RabbitMQ)
package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Order struct {
	OrderID    string  `json:"order_id"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	CustomerID string  `json:"customer_id"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declarar una cola
	q, err := ch.QueueDeclare(
		"order_queue", // nombre
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Crear una orden
	order := Order{
		OrderID:    "12345",
		ProductID:  1,
		Quantity:   2,
		CustomerID: "cust789",
	}

	// Convertir a JSON
	body, err := json.Marshal(order)
	failOnError(err, "Failed to encode order to JSON")

	// Publicar mensaje
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf("Sent order: %s", order.OrderID)
}
```

```go
// consumer.go (usando RabbitMQ)
package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Order struct {
	OrderID    string  `json:"order_id"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	CustomerID string  `json:"customer_id"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declarar la misma cola
	q, err := ch.QueueDeclare(
		"order_queue", // nombre
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Configurar la calidad de servicio
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	// Consumir mensajes
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var order Order
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("Error decoding order: %s", err)
				d.Nack(false, false) // rechazar el mensaje
				continue
			}

			log.Printf("Received order: %s for product %d", order.OrderID, order.ProductID)

			// Procesar la orden
			// ...

			d.Ack(false) // confirmar el mensaje
		}
	}()

	log.Printf("Waiting for orders. To exit press CTRL+C")
	<-forever
}
```

## Patrones de Diseño para Microservicios

### API Gateway

El patrón API Gateway proporciona un punto de entrada único para los clientes, enrutando las solicitudes a los servicios apropiados.

```go
// api_gateway.go
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Configuración de servicios
var services = map[string]string{
	"products": "http://localhost:8080",
	"orders":   "http://localhost:8081",
	"payments": "http://localhost:8082",
	"users":    "http://localhost:8083",
}

func main() {
	r := gin.Default()

	// Middleware para logging
	r.Use(func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		log.Printf("[API Gateway] %s %s %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			latency,
		)
	})

	// Health check
	r.GET("/health", healthCheck)

	// Proxy para servicios
	r.Any("/:service/*path", routeToService)

	// Rutas específicas para mejorar la experiencia del cliente
	r.GET("/products", getProducts)
	r.GET("/orders/:id/details", getOrderDetails)

	r.Run(":8000")
}

func healthCheck(c *gin.Context) {
	// Verificar la salud de todos los servicios
	health := make(map[string]string)

	for service, url := range services {
		status := "DOWN"

		client := http.Client{
			Timeout: 1 * time.Second,
		}

		resp, err := client.Get(url + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			status = "UP"
		}

		health[service] = status
	}

	// Determinar el estado general
	overall := "UP"
	for _, status := range health {
		if status != "UP" {
			overall = "DEGRADED"
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   overall,
		"services": health,
	})
}

func routeToService(c *gin.Context) {
	service := c.Param("service")
	path := c.Param("path")

	// Verificar si el servicio existe
	serviceURL, exists := services[service]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	// Construir la URL de destino
	target, err := url.Parse(serviceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid service URL"})
		return
	}

	// Crear proxy reverso
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Actualizar la solicitud
	c.Request.URL.Path = path
	c.Request.URL.Host = target.Host
	c.Request.URL.Scheme = target.Scheme
	c.Request.Host = target.Host

	// Servir la solicitud
	proxy.ServeHTTP(c.Writer, c.Request)
}

func getProducts(c *gin.Context) {
	// Obtener productos desde el servicio de productos
	resp, err := http.Get(services["products"] + "/products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error communicating with product service"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response from product service"})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

func getOrderDetails(c *gin.Context) {
	orderID := c.Param("id")

	// Obtener la orden desde el servicio de órdenes
	orderResp, err := http.Get(services["orders"] + "/orders/" + orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error communicating with order service"})
		return
	}
	defer orderResp.Body.Close()

	if orderResp.StatusCode != http.StatusOK {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	orderBody, err := ioutil.ReadAll(orderResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response from order service"})
		return
	}

	// Obtener información de pago desde el servicio de pagos
	paymentResp, err := http.Get(services["payments"] + "/payments?order_id=" + orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error communicating with payment service"})
		return
	}
	defer paymentResp.Body.Close()

	paymentBody, err := ioutil.ReadAll(paymentResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response from payment service"})
		return
	}

	// Combinar la información
	result := gin.H{
		"order":    json.RawMessage(orderBody),
		"payments": json.RawMessage(paymentBody),
	}

	c.JSON(http.StatusOK, result)
}
```

### Circuit Breaker

El patrón Circuit Breaker evita que una aplicación intente repetidamente una operación que probablemente fallará.

```go
// circuit_breaker.go
package main

import (
	"errors"
	"sync"
	"time"
)

type CircuitBreaker struct {
	threshold   int
	timeout     time.Duration
	failures    int
	state       string
	lastFailure time.Time
	mutex       sync.Mutex
}

const (
	StateClosed   = "CLOSED"
	StateOpen     = "OPEN"
	StateHalfOpen = "HALF-OPEN"
)

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		timeout:   timeout,
		state:     StateClosed,
	}
}

func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	cb.mutex.Lock()

	if cb.state == StateOpen {
		// Comprobar si ha pasado el tiempo de timeout
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = StateHalfOpen
			cb.mutex.Unlock()
		} else {
			cb.mutex.Unlock()
			return nil, errors.New("circuit breaker is open")
		}
	} else {
		cb.mutex.Unlock()
	}

	result, err := fn()

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		// Registrar el fallo
		cb.failures++
		cb.lastFailure = time.Now()

		// Si alcanzamos el umbral de fallos, abrir el circuito
		if cb.state == StateClosed && cb.failures >= cb.threshold {
			cb.state = StateOpen
		}

		// Si el circuito estaba semi-abierto y falló, volver a abrirlo
		if cb.state == StateHalfOpen {
			cb.state = StateOpen
		}

		return nil, err
	}

	// Si el circuito estaba semi-abierto y la llamada tuvo éxito, cerrarlo
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		cb.failures = 0
	}

	return result, nil
}

// Ejemplo de uso
func main() {
	cb := NewCircuitBreaker(3, 10*time.Second)

	// Función que puede fallar
	callExternalService := func() (interface{}, error) {
		// Simular una llamada a un servicio externo
		// ...
		return nil, errors.New("service unavailable")
	}

	// Usar el circuit breaker
	for i := 0; i < 10; i++ {
		result, err := cb.Execute(callExternalService)
		if err != nil {
			log.Printf("Error: %s", err)
		} else {
			log.Printf("Result: %v", result)
		}

		time.Sleep(1 * time.Second)
	}
}
```

### Service Discovery

El patrón Service Discovery permite a los servicios encontrarse entre sí dinámicamente.

```go
// service_registry.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceInfo struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Address       string            `json:"address"`
	Port          int               `json:"port"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	LastHeartbeat time.Time         `json:"last_heartbeat"`
}

type ServiceRegistry struct {
	services map[string]ServiceInfo
	mutex    sync.RWMutex
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]ServiceInfo),
	}
}

func (sr *ServiceRegistry) Register(service ServiceInfo) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	service.LastHeartbeat = time.Now()
	sr.services[service.ID] = service
}

func (sr *ServiceRegistry) Unregister(serviceID string) bool {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	_, exists := sr.services[serviceID]
	if !exists {
		return false
	}

	delete(sr.services, serviceID)
	return true
}

func (sr *ServiceRegistry) Heartbeat(serviceID string) bool {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	service, exists := sr.services[serviceID]
	if !exists {
		return false
	}

	service.LastHeartbeat = time.Now()
	sr.services[serviceID] = service
	return true
}

func (sr *ServiceRegistry) GetServices(serviceName string) []ServiceInfo {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	var result []ServiceInfo

	for _, service := range sr.services {
		if serviceName == "" || service.Name == serviceName {
			result = append(result, service)
		}
	}

	return result
}

func (sr *ServiceRegistry) GetService(serviceID string) (ServiceInfo, bool) {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	service, exists := sr.services[serviceID]
	return service, exists
}

func (sr *ServiceRegistry) CleanupInactiveServices(timeout time.Duration) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	now := time.Now()
	for id, service := range sr.services {
		if now.Sub(service.LastHeartbeat) > timeout {
			log.Printf("Removing inactive service: %s (%s)", service.Name, id)
			delete(sr.services, id)
		}
	}
}

func main() {
	registry := NewServiceRegistry()

	// Iniciar el hilo de limpieza
	go func() {
		for {
			time.Sleep(10 * time.Second)
			registry.CleanupInactiveServices(30 * time.Second)
		}
	}()

	r := gin.Default()

	// Registrar un servicio
	r.POST("/register", func(c *gin.Context) {
		var service ServiceInfo
		if err := c.ShouldBindJSON(&service); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service data"})
			return
		}

		if service.ID == "" {
			service.ID = uuid.New().String()
		}

		registry.Register(service)
		c.JSON(http.StatusOK, gin.H{"id": service.ID})
	})

	// Enviar heartbeat
	r.PUT("/heartbeat/:id", func(c *gin.Context) {
		id := c.Param("id")
		if !registry.Heartbeat(id) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		c.Status(http.StatusOK)
	})

	// Dar de baja un servicio
	r.DELETE("/unregister/:id", func(c *gin.Context) {
		id := c.Param("id")
		if !registry.Unregister(id) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		c.Status(http.StatusNoContent)
	})

	// Obtener todos los servicios o filtrar por nombre
	r.GET("/services", func(c *gin.Context) {
		name := c.Query("name")
		services := registry.GetServices(name)
		c.JSON(http.StatusOK, services)
	})

	// Obtener un servicio específico
	r.GET("/services/:id", func(c *gin.Context) {
		id := c.Param("id")
		service, exists := registry.GetService(id)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		c.JSON(http.StatusOK, service)
	})

	r.Run(":8500")
}
```

```go
// service_client.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ServiceInfo struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type ServiceClient struct {
	registryURL string
	serviceInfo ServiceInfo
	client      http.Client
	stopCh      chan struct{}
}

func NewServiceClient(registryURL string, serviceInfo ServiceInfo) *ServiceClient {
	return &ServiceClient{
		registryURL: registryURL,
		serviceInfo: serviceInfo,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		stopCh: make(chan struct{}),
	}
}

func (sc *ServiceClient) Register() (string, error) {
	body, err := json.Marshal(sc.serviceInfo)
	if err != nil {
		return "", err
	}

	resp, err := sc.client.Post(sc.registryURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	sc.serviceInfo.ID = result["id"]
	return sc.serviceInfo.ID, nil
}

func (sc *ServiceClient) Unregister() error {
	req, err := http.NewRequest("DELETE", sc.registryURL+"/unregister/"+sc.serviceInfo.ID, nil)
	if err != nil {
		return err
	}

	resp, err := sc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (sc *ServiceClient) StartHeartbeat(interval time.Duration) {
	sc.stopCh = make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := sc.sendHeartbeat(); err != nil {
					log.Printf("Error sending heartbeat: %v", err)
				}
			case <-sc.stopCh:
				return
			}
		}
	}()
}

func (sc *ServiceClient) StopHeartbeat() {
	close(sc.stopCh)
}

func (sc *ServiceClient) sendHeartbeat() error {
	req, err := http.NewRequest("PUT", sc.registryURL+"/heartbeat/"+sc.serviceInfo.ID, nil)
	if err != nil {
		return err
	}

	resp, err := sc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (sc *ServiceClient) DiscoverService(serviceName string) ([]ServiceInfo, error) {
	resp, err := sc.client.Get(sc.registryURL + "/services?name=" + serviceName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var services []ServiceInfo
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return nil, err
	}

	return services, nil
}
```

## Despliegue de Microservicios

### Docker

Docker es una plataforma que permite empaquetar, distribuir y ejecutar aplicaciones en contenedores.

```dockerfile
# Dockerfile para un microservicio en Go
FROM golang:1.17-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Imagen final
FROM alpine:latest

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/main .

# Exponer el puerto
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./main"]
```

```yaml
# docker-compose.yml
version: '3'

services:
  product-service:
    build: ./product-service
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=products
    depends_on:
      - postgres

  order-service:
    build: ./order-service
    ports:
      - "8081:8080"
    environment:
      - PRODUCT_SERVICE_URL=http://product-service:8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=orders
    depends_on:
      - postgres
      - product-service

  api-gateway:
    build: ./api-gateway
    ports:
      - "8000:8000"
    environment:
      - PRODUCT_SERVICE_URL=http://product-service:8080
      - ORDER_SERVICE_URL=http://order-service:8080
    depends_on:
      - product-service
      - order-service

  postgres:
    image: postgres:13
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_MULTIPLE_DATABASES=products,orders
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./init-db.sh:/docker-entrypoint-initdb.d/init-db.sh

volumes:
  postgres-data:
```

### Kubernetes

Kubernetes es una plataforma de orquestación de contenedores que facilita la automatización del despliegue, escalado y gestión de aplicaciones en contenedores.

```yaml
# product-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: product-service
  labels:
    app: product-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: product-service
  template:
    metadata:
      labels:
        app: product-service
    spec:
      containers:
      - name: product-service
        image: product-service:latest
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: postgres
        - name: DB_PORT
          value: "5432"
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: password
        - name: DB_NAME
          value: products
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: product-service
spec:
  selector:
    app: product-service
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: microservices-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$1
spec:
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /products(/|$)(.*)
        pathType: Prefix
        backend:
          service:
            name: product-service
            port:
              number: 80
      - path: /orders(/|$)(.*)
        pathType: Prefix
        backend:
          service:
            name: order-service
            port:
              number: 80
```

## Monitoreo y Observabilidad

### Logging

El logging es esencial para entender el comportamiento de los microservicios en producción.

```go
// logger.go
package logger

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type contextKey string

const requestIDKey contextKey = "requestID"

// Logger es una extensión de logrus.Logger
type Logger struct {
	*logrus.Logger
}

// New crea una nueva instancia de Logger
func New() *Logger {
	logger := logrus.New()
	logger.SetFormatter(&jsonFormatter{})
	logger.SetOutput(os.Stdout)

	return &Logger{Logger: logger}
}

// WithRequestID añade un ID de solicitud al contexto
func WithRequestID(ctx context.Context) context.Context {
	requestID := uuid.New().String()
	return context.WithValue(ctx, requestIDKey, requestID)
}

// GetRequestID obtiene el ID de solicitud del contexto
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// WithContext crea una entrada de log con información del contexto
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	return l.WithField("request_id", GetRequestID(ctx))
}

// jsonFormatter es un formateador personalizado para logrus
type jsonFormatter struct {
}

func (f *jsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+4)

	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	data["timestamp"] = entry.Time.Format(time.RFC3339)
	data["level"] = entry.Level.String()
	data["message"] = entry.Message

	// Añadir información sobre el archivo y la línea
	if _, file, line, ok := runtime.Caller(6); ok {
		data["file"] = file
		data["line"] = line
	}

	var b []byte
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	b = append(b, '\n')
	return b, nil
}
```

```go
// middleware.go (para Gin)
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"path/to/logger"
)

// LoggingMiddleware es un middleware para Gin que registra información sobre las solicitudes
func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Crear un contexto con ID de solicitud
		ctx := logger.WithRequestID(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		// Tiempo de inicio
		start := time.Now()

		// Procesar la solicitud
		c.Next()

		// Tiempo de finalización
		end := time.Now()
		latency := end.Sub(start)

		// Registrar la solicitud
		log.WithContext(ctx).WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user_agent": c.Request.UserAgent(),
		}).Info("Request processed")
	}
}
```

### Métricas

Las métricas permiten medir el rendimiento y la salud de los microservicios.

```go
// metrics.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics contiene todas las métricas de la aplicación
type Metrics struct {
	RequestsTotal      *prometheus.CounterVec
	RequestDuration    *prometheus.HistogramVec
	DatabaseQueriesTotal *prometheus.CounterVec
	DatabaseQueryDuration *prometheus.HistogramVec
	ActiveConnections  prometheus.Gauge
	ErrorsTotal        *prometheus.CounterVec
}

// New crea una nueva instancia de Metrics
func New(namespace string) *Metrics {
	return &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "requests_total",
				Help:      "Total number of requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "request_duration_seconds",
				Help:      "Request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		DatabaseQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "database_queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "table"},
		),
		DatabaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "database_query_duration_seconds",
				Help:      "Database query duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		),
		ActiveConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "active_connections",
				Help:      "Number of active connections",
			},
		),
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "errors_total",
				Help:      "Total number of errors",
			},
			[]string{"type"},
		),
	}
}
```

```go
// middleware.go (para Gin)
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"path/to/metrics"
)

// MetricsMiddleware es un middleware para Gin que registra métricas sobre las solicitudes
func MetricsMiddleware(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Incrementar el contador de conexiones activas
		m.ActiveConnections.Inc()
		defer m.ActiveConnections.Dec()

		// Tiempo de inicio
		start := time.Now()

		// Procesar la solicitud
		c.Next()

		// Tiempo de finalización
		duration := time.Since(start).Seconds()

		// Registrar métricas
		m.RequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, c.Writer.Status()).Inc()
		m.RequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(duration)

		// Registrar errores
		if c.Writer.Status() >= 500 {
			m.ErrorsTotal.WithLabelValues("server").Inc()
		} else if c.Writer.Status() >= 400 {
			m.ErrorsTotal.WithLabelValues("client").Inc()
		}
	}
}
```

### Tracing

El tracing permite seguir el flujo de una solicitud a través de múltiples microservicios.

```go
// tracing.go
package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.7.0"
)

// InitTracer inicializa el tracer de OpenTelemetry
func InitTracer(serviceName, jaegerEndpoint string) (func(), error) {
	// Crear exportador Jaeger
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)),
	)
	if err != nil {
		return nil, err
	}

	// Crear proveedor de trazas
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Establecer el proveedor de trazas global
	otel.SetTracerProvider(tp)

	// Establecer el propagador global
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Devolver una función para limpiar
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}, nil
}
```

```go
// middleware.go (para Gin)
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// TracingMiddleware es un middleware para Gin que añade tracing a las solicitudes
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}
```

## Ejemplo Práctico: Sistema de Comercio Electrónico

A continuación, se presenta un ejemplo práctico de un sistema de comercio electrónico implementado con microservicios en Go.

### Estructura del Proyecto

```
ecommerce/
├── api-gateway/
│   ├── Dockerfile
│   ├── main.go
│   └── ...
├── product-service/
│   ├── Dockerfile
│   ├── main.go
│   ├── models/
│   ├── handlers/
│   ├── repository/
│   └── ...
├── order-service/
│   ├── Dockerfile
│   ├── main.go
│   ├── models/
│   ├── handlers/
│   ├── repository/
│   └── ...
├── payment-service/
│   ├── Dockerfile
│   ├── main.go
│   ├── models/
│   ├── handlers/
│   ├── repository/
│   └── ...
└── docker-compose.yml
```

### Servicio de Productos

```go
// product-service/models/product.go
package models

type Product struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	Stock       int     `json:"stock" db:"stock"`
	Category    string  `json:"category" db:"category"`
	ImageURL    string  `json:"image_url" db:"image_url"`
}
```

```go
// product-service/repository/product_repository.go
package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"path/to/models"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id string) (models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.Product) error
	DeleteProduct(ctx context.Context, id string) error
}

type postgresProductRepository struct {
	db *sqlx.DB
}

func NewPostgresProductRepository(db *sqlx.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) GetProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product

	query := `SELECT id, name, description, price, stock, category, image_url FROM products`
	err := r.db.SelectContext(ctx, &products, query)
	if err != nil {
		return nil, errors.Wrap(err, "error getting products")
	}

	return products, nil
}

func (r *postgresProductRepository) GetProductByID(ctx context.Context, id string) (models.Product, error) {
	var product models.Product

	query := `SELECT id, name, description, price, stock, category, image_url FROM products WHERE id = $1`
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Product{}, errors.New("product not found")
		}
		return models.Product{}, errors.Wrap(err, "error getting product")
	}

	return product, nil
}

func (r *postgresProductRepository) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	query := `INSERT INTO products (id, name, description, price, stock, category, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, product.ID, product.Name, product.Description, product.Price, product.Stock, product.Category, product.ImageURL).Scan(&product.ID)
	if err != nil {
		return models.Product{}, errors.Wrap(err, "error creating product")
	}

	return product, nil
}

func (r *postgresProductRepository) UpdateProduct(ctx context.Context, product models.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, stock = $4, category = $5, image_url = $6 WHERE id = $7`

	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Stock, product.Category, product.ImageURL, product.ID)
	if err != nil {
		return errors.Wrap(err, "error updating product")
	}

	return nil
}

func (r *postgresProductRepository) DeleteProduct(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "error deleting product")
	}

	return nil
}
```

```go
// product-service/handlers/product_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"path/to/models"
	"path/to/repository"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.repo.GetProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.repo.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generar un ID único
	product.ID = uuid.New().String()

	createdProduct, err := h.repo.CreateProduct(c.Request.Context(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// Verificar si el producto existe
	_, err := h.repo.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer el ID del producto
	product.ID = id

	if err := h.repo.UpdateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Verificar si el producto existe
	_, err := h.repo.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
```

```go
// product-service/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"path/to/handlers"
	"path/to/repository"
	"path/to/tracing"
)

func main() {
	// Configuración de la base de datos
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Configuración de tracing
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = "http://jaeger:14268/api/traces"
	}

	cleanup, err := tracing.InitTracer("product-service", jaegerEndpoint)
	if err != nil {
		log.Fatalf("Could not initialize tracer: %v", err)
	}
	defer cleanup()

	// Configuración del repositorio y handler
	productRepo := repository.NewPostgresProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepo)

	// Configuración del router
	r := gin.Default()

	// Middleware
	r.Use(middleware.TracingMiddleware("product-service"))

	// Rutas
	r.GET("/products", productHandler.GetProducts)
	r.GET("/products/:id", productHandler.GetProductByID)
	r.POST("/products", productHandler.CreateProduct)
	r.PUT("/products/:id", productHandler.UpdateProduct)
	r.DELETE("/products/:id", productHandler.DeleteProduct)

	// Ruta para métricas de Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Ruta para health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Iniciar el servidor
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Iniciar el servidor en una goroutine
	go func() {
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
```

## Mejores Prácticas para Microservicios en Go

### Diseño

1. **Principio de Responsabilidad Única**: Cada microservicio debe tener una única responsabilidad y razón para cambiar.
2. **Independencia de Datos**: Cada microservicio debe tener su propia base de datos o esquema.
3. **API Bien Definidas**: Utiliza contratos claros y versionados para las API.
4. **Tamaño Adecuado**: No hagas microservicios demasiado pequeños ni demasiado grandes.
5. **Diseño Orientado al Dominio (DDD)**: Utiliza DDD para identificar los límites de los microservicios.

### Desarrollo

1. **Código Limpio**: Mantén el código limpio, bien estructurado y con pruebas.
2. **Manejo de Errores**: Implementa un manejo de errores consistente y robusto.
3. **Configuración Externalizada**: Utiliza variables de entorno o archivos de configuración.
4. **Versionado de API**: Versiona tus API para permitir cambios sin romper la compatibilidad.
5. **Pruebas Automatizadas**: Implementa pruebas unitarias, de integración y de contrato.

### Operaciones

1. **Observabilidad**: Implementa logging, métricas y tracing.
2. **Despliegue Automatizado**: Utiliza CI/CD para automatizar el despliegue.
3. **Contenedores**: Utiliza Docker para empaquetar tus microservicios.
4. **Orquestación**: Utiliza Kubernetes para orquestar tus contenedores.
5. **Monitoreo y Alertas**: Implementa monitoreo y alertas para detectar problemas.

## Conclusión

Los microservicios en Go ofrecen una forma poderosa y eficiente de construir aplicaciones escalables y mantenibles. Go, con su simplicidad, eficiencia y soporte nativo para concurrencia, es una excelente opción para implementar microservicios.

En este documento, hemos explorado los conceptos fundamentales de los microservicios en Go, desde los frameworks disponibles hasta los patrones de diseño, pasando por la comunicación entre servicios, el despliegue y el monitoreo. También hemos visto un ejemplo práctico de un sistema de comercio electrónico implementado con microservicios en Go.

## Recursos Adicionales

- [Go Kit](https://gokit.io/): Un conjunto de paquetes para construir microservicios en Go.
- [Gin Web Framework](https://github.com/gin-gonic/gin): Un framework web HTTP de alto rendimiento para Go.
- [gRPC](https://grpc.io/): Un framework de RPC de alto rendimiento.
- [Docker](https://www.docker.com/): Una plataforma para empaquetar, distribuir y ejecutar aplicaciones en contenedores.
- [Kubernetes](https://kubernetes.io/): Una plataforma de orquestación de contenedores.
- [Prometheus](https://prometheus.io/): Un sistema de monitoreo y alerta.
- [Jaeger](https://www.jaegertracing.io/): Un sistema de tracing distribuido.
- [Building Microservices with Go](https://www.amazon.com/Building-Microservices-Go-Sam-Newman/dp/1491950358): Un libro sobre microservicios en Go.
- [Microservices Patterns](https://microservices.io/patterns/index.html): Patrones de diseño para microservicios.