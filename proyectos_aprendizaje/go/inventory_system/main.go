package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Enumeraciones
type ProductCategory string
type TransactionType string

const (
	// Categorías de productos
	Electronics ProductCategory = "ELECTRONICS"
	Clothing    ProductCategory = "CLOTHING"
	Food        ProductCategory = "FOOD"
	Books       ProductCategory = "BOOKS"
	Other       ProductCategory = "OTHER"

	// Tipos de transacciones
	Purchase TransactionType = "PURCHASE"
	Sale     TransactionType = "SALE"
	Adjust   TransactionType = "ADJUST"

	// Archivos de datos
	ProductsFile     = "products.json"
	TransactionsFile = "transactions.json"
)

// Modelos
type Product struct {
	SKU         string          `json:"sku"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Category    ProductCategory `json:"category"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Transaction struct {
	ID          string          `json:"id"`
	ProductSKU  string          `json:"product_sku"`
	Type        TransactionType `json:"type"`
	Quantity    int             `json:"quantity"`
	Notes       string          `json:"notes"`
	Timestamp   time.Time       `json:"timestamp"`
}

// Interfaces DAO
type ProductDAO interface {
	Save(product Product) error
	FindAll() ([]Product, error)
	FindBySKU(sku string) (Product, error)
	FindByCategory(category ProductCategory) ([]Product, error)
	Delete(sku string) error
}

type TransactionDAO interface {
	Save(transaction Transaction) error
	FindAll() ([]Transaction, error)
	FindByProductSKU(sku string) ([]Transaction, error)
	FindByType(transactionType TransactionType) ([]Transaction, error)
}

// Implementaciones DAO
type JSONProductDAO struct {
	FilePath string
}

func NewJSONProductDAO() *JSONProductDAO {
	return &JSONProductDAO{FilePath: ProductsFile}
}

func (dao *JSONProductDAO) Save(product Product) error {
	products, err := dao.FindAll()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Actualizar o añadir producto
	found := false
	for i, p := range products {
		if p.SKU == product.SKU {
			products[i] = product
			found = true
			break
		}
	}

	if !found {
		products = append(products, product)
	}

	// Guardar en archivo
	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}

	// Asegurar que el directorio existe
	dir := filepath.Dir(dao.FilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return ioutil.WriteFile(dao.FilePath, data, 0644)
}

func (dao *JSONProductDAO) FindAll() ([]Product, error) {
	data, err := ioutil.ReadFile(dao.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Product{}, nil
		}
		return nil, err
	}

	var products []Product
	err = json.Unmarshal(data, &products)
	return products, err
}

func (dao *JSONProductDAO) FindBySKU(sku string) (Product, error) {
	products, err := dao.FindAll()
	if err != nil {
		return Product{}, err
	}

	for _, p := range products {
		if p.SKU == sku {
			return p, nil
		}
	}

	return Product{}, errors.New("product not found")
}

func (dao *JSONProductDAO) FindByCategory(category ProductCategory) ([]Product, error) {
	products, err := dao.FindAll()
	if err != nil {
		return nil, err
	}

	var result []Product
	for _, p := range products {
		if p.Category == category {
			result = append(result, p)
		}
	}

	return result, nil
}

func (dao *JSONProductDAO) Delete(sku string) error {
	products, err := dao.FindAll()
	if err != nil {
		return err
	}

	var newProducts []Product
	found := false
	for _, p := range products {
		if p.SKU != sku {
			newProducts = append(newProducts, p)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("product not found")
	}

	data, err := json.MarshalIndent(newProducts, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dao.FilePath, data, 0644)
}

type JSONTransactionDAO struct {
	FilePath string
}

func NewJSONTransactionDAO() *JSONTransactionDAO {
	return &JSONTransactionDAO{FilePath: TransactionsFile}
}

func (dao *JSONTransactionDAO) Save(transaction Transaction) error {
	transactions, err := dao.FindAll()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Añadir transacción
	transactions = append(transactions, transaction)

	// Guardar en archivo
	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return err
	}

	// Asegurar que el directorio existe
	dir := filepath.Dir(dao.FilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return ioutil.WriteFile(dao.FilePath, data, 0644)
}

func (dao *JSONTransactionDAO) FindAll() ([]Transaction, error) {
	data, err := ioutil.ReadFile(dao.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Transaction{}, nil
		}
		return nil, err
	}

	var transactions []Transaction
	err = json.Unmarshal(data, &transactions)
	return transactions, err
}

func (dao *JSONTransactionDAO) FindByProductSKU(sku string) ([]Transaction, error) {
	transactions, err := dao.FindAll()
	if err != nil {
		return nil, err
	}

	var result []Transaction
	for _, t := range transactions {
		if t.ProductSKU == sku {
			result = append(result, t)
		}
	}

	return result, nil
}

func (dao *JSONTransactionDAO) FindByType(transactionType TransactionType) ([]Transaction, error) {
	transactions, err := dao.FindAll()
	if err != nil {
		return nil, err
	}

	var result []Transaction
	for _, t := range transactions {
		if t.Type == transactionType {
			result = append(result, t)
		}
	}

	return result, nil
}

// Servicio de Inventario (Facade)
type InventoryService struct {
	productDAO     ProductDAO
	transactionDAO TransactionDAO
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		productDAO:     NewJSONProductDAO(),
		transactionDAO: NewJSONTransactionDAO(),
	}
}

// Métodos de gestión de productos
func (s *InventoryService) AddProduct(name, description string, category ProductCategory, price float64, quantity int) (Product, error) {
	// Validar datos
	if name == "" {
		return Product{}, errors.New("name is required")
	}
	if price < 0 {
		return Product{}, errors.New("price must be positive")
	}
	if quantity < 0 {
		return Product{}, errors.New("quantity must be positive")
	}

	// Generar SKU único
	sku := generateSKU(name, category)

	// Crear producto
	product := Product{
		SKU:         sku,
		Name:        name,
		Description: description,
		Category:    category,
		Price:       price,
		Quantity:    quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Guardar producto
	err := s.productDAO.Save(product)
	if err != nil {
		return Product{}, err
	}

	// Registrar transacción si hay cantidad inicial
	if quantity > 0 {
		_, err = s.RecordTransaction(sku, Purchase, quantity, "Initial stock")
		if err != nil {
			return product, fmt.Errorf("product created but failed to record transaction: %v", err)
		}
	}

	return product, nil
}

func (s *InventoryService) UpdateProduct(sku, name, description string, category ProductCategory, price float64) (Product, error) {
	// Obtener producto existente
	product, err := s.productDAO.FindBySKU(sku)
	if err != nil {
		return Product{}, err
	}

	// Actualizar campos
	if name != "" {
		product.Name = name
	}
	if description != "" {
		product.Description = description
	}
	if category != "" {
		product.Category = category
	}
	if price >= 0 {
		product.Price = price
	}
	product.UpdatedAt = time.Now()

	// Guardar producto actualizado
	err = s.productDAO.Save(product)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (s *InventoryService) GetProduct(sku string) (Product, error) {
	return s.productDAO.FindBySKU(sku)
}

func (s *InventoryService) DeleteProduct(sku string) error {
	return s.productDAO.Delete(sku)
}

func (s *InventoryService) GetAllProducts() ([]Product, error) {
	return s.productDAO.FindAll()
}

func (s *InventoryService) GetProductsByCategory(category ProductCategory) ([]Product, error) {
	return s.productDAO.FindByCategory(category)
}

// Métodos de gestión de stock
func (s *InventoryService) AddStock(sku string, quantity int, notes string) (Product, error) {
	// Validar datos
	if quantity <= 0 {
		return Product{}, errors.New("quantity must be positive")
	}

	// Obtener producto
	product, err := s.productDAO.FindBySKU(sku)
	if err != nil {
		return Product{}, err
	}

	// Actualizar cantidad
	product.Quantity += quantity
	product.UpdatedAt = time.Now()

	// Guardar producto
	err = s.productDAO.Save(product)
	if err != nil {
		return Product{}, err
	}

	// Registrar transacción
	_, err = s.RecordTransaction(sku, Purchase, quantity, notes)
	if err != nil {
		return product, fmt.Errorf("stock added but failed to record transaction: %v", err)
	}

	return product, nil
}

func (s *InventoryService) RemoveStock(sku string, quantity int, notes string) (Product, error) {
	// Validar datos
	if quantity <= 0 {
		return Product{}, errors.New("quantity must be positive")
	}

	// Obtener producto
	product, err := s.productDAO.FindBySKU(sku)
	if err != nil {
		return Product{}, err
	}

	// Verificar stock suficiente
	if product.Quantity < quantity {
		return Product{}, errors.New("insufficient stock")
	}

	// Actualizar cantidad
	product.Quantity -= quantity
	product.UpdatedAt = time.Now()

	// Guardar producto
	err = s.productDAO.Save(product)
	if err != nil {
		return Product{}, err
	}

	// Registrar transacción
	_, err = s.RecordTransaction(sku, Sale, quantity, notes)
	if err != nil {
		return product, fmt.Errorf("stock removed but failed to record transaction: %v", err)
	}

	return product, nil
}

func (s *InventoryService) AdjustStock(sku string, newQuantity int, notes string) (Product, error) {
	// Obtener producto
	product, err := s.productDAO.FindBySKU(sku)
	if err != nil {
		return Product{}, err
	}

	// Calcular diferencia
	diff := newQuantity - product.Quantity

	// Actualizar cantidad
	oldQuantity := product.Quantity
	product.Quantity = newQuantity
	product.UpdatedAt = time.Now()

	// Guardar producto
	err = s.productDAO.Save(product)
	if err != nil {
		return Product{}, err
	}

	// Registrar transacción
	_, err = s.RecordTransaction(sku, Adjust, diff, fmt.Sprintf("%s (Adjusted from %d to %d)", notes, oldQuantity, newQuantity))
	if err != nil {
		return product, fmt.Errorf("stock adjusted but failed to record transaction: %v", err)
	}

	return product, nil
}

// Métodos de gestión de transacciones
func (s *InventoryService) RecordTransaction(productSKU string, transactionType TransactionType, quantity int, notes string) (Transaction, error) {
	// Validar que el producto existe
	_, err := s.productDAO.FindBySKU(productSKU)
	if err != nil {
		return Transaction{}, err
	}

	// Crear transacción
	transaction := Transaction{
		ID:         uuid.New().String(),
		ProductSKU: productSKU,
		Type:       transactionType,
		Quantity:   quantity,
		Notes:      notes,
		Timestamp:  time.Now(),
	}

	// Guardar transacción
	err = s.transactionDAO.Save(transaction)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func (s *InventoryService) GetAllTransactions() ([]Transaction, error) {
	return s.transactionDAO.FindAll()
}

func (s *InventoryService) GetTransactionsByProduct(sku string) ([]Transaction, error) {
	return s.transactionDAO.FindByProductSKU(sku)
}

func (s *InventoryService) GetTransactionsByType(transactionType TransactionType) ([]Transaction, error) {
	return s.transactionDAO.FindByType(transactionType)
}

// Métodos de reportes
func (s *InventoryService) GetLowStockProducts(threshold int) ([]Product, error) {
	products, err := s.productDAO.FindAll()
	if err != nil {
		return nil, err
	}

	var lowStockProducts []Product
	for _, p := range products {
		if p.Quantity <= threshold && p.Quantity > 0 {
			lowStockProducts = append(lowStockProducts, p)
		}
	}

	return lowStockProducts, nil
}

func (s *InventoryService) GetOutOfStockProducts() ([]Product, error) {
	products, err := s.productDAO.FindAll()
	if err != nil {
		return nil, err
	}

	var outOfStockProducts []Product
	for _, p := range products {
		if p.Quantity == 0 {
			outOfStockProducts = append(outOfStockProducts, p)
		}
	}

	return outOfStockProducts, nil
}

func (s *InventoryService) GetInventoryValue() (float64, error) {
	products, err := s.productDAO.FindAll()
	if err != nil {
		return 0, err
	}

	var totalValue float64
	for _, p := range products {
		totalValue += p.Price * float64(p.Quantity)
	}

	return totalValue, nil
}

func (s *InventoryService) GetTransactionSummary() (map[TransactionType]int, error) {
	transactions, err := s.transactionDAO.FindAll()
	if err != nil {
		return nil, err
	}

	summary := make(map[TransactionType]int)
	summary[Purchase] = 0
	summary[Sale] = 0
	summary[Adjust] = 0

	for _, t := range transactions {
		summary[t.Type] += t.Quantity
	}

	return summary, nil
}

// Funciones auxiliares
func generateSKU(name string, category ProductCategory) string {
	// Generar un SKU basado en la categoría y el nombre
	prefix := string(category)[:2]
	namePrefix := ""
	if len(name) >= 3 {
		namePrefix = strings.ToUpper(strings.ReplaceAll(name[:3], " ", ""))
	} else {
		namePrefix = strings.ToUpper(strings.ReplaceAll(name, " ", ""))
	}

	// Añadir un timestamp para garantizar unicidad
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 36)

	return fmt.Sprintf("%s%s%s", prefix, namePrefix, timestamp)
}

// CLI
type CLI struct {
	service *InventoryService
}

func NewCLI() *CLI {
	return &CLI{
		service: NewInventoryService(),
	}
}

func (cli *CLI) Run() {
	// Definir comandos principales
	addProductCmd := flag.NewFlagSet("add-product", flag.ExitOnError)
	updateProductCmd := flag.NewFlagSet("update-product", flag.ExitOnError)
	showProductCmd := flag.NewFlagSet("show-product", flag.ExitOnError)
	deleteProductCmd := flag.NewFlagSet("delete-product", flag.ExitOnError)
	listProductsCmd := flag.NewFlagSet("list-products", flag.ExitOnError)
	addStockCmd := flag.NewFlagSet("add-stock", flag.ExitOnError)
	removeStockCmd := flag.NewFlagSet("remove-stock", flag.ExitOnError)
	adjustStockCmd := flag.NewFlagSet("adjust-stock", flag.ExitOnError)
	listTransactionsCmd := flag.NewFlagSet("list-transactions", flag.ExitOnError)
	reportCmd := flag.NewFlagSet("report", flag.ExitOnError)

	// Definir flags para add-product
	addProductName := addProductCmd.String("name", "", "Product name")
	addProductDesc := addProductCmd.String("description", "", "Product description")
	addProductCategory := addProductCmd.String("category", "", "Product category (ELECTRONICS, CLOTHING, FOOD, BOOKS, OTHER)")
	addProductPrice := addProductCmd.Float64("price", 0, "Product price")
	addProductQuantity := addProductCmd.Int("quantity", 0, "Initial product quantity")

	// Definir flags para update-product
	updateProductSKU := updateProductCmd.String("sku", "", "Product SKU")
	updateProductName := updateProductCmd.String("name", "", "Product name")
	updateProductDesc := updateProductCmd.String("description", "", "Product description")
	updateProductCategory := updateProductCmd.String("category", "", "Product category (ELECTRONICS, CLOTHING, FOOD, BOOKS, OTHER)")
	updateProductPrice := updateProductCmd.Float64("price", -1, "Product price")

	// Definir flags para show-product
	showProductSKU := showProductCmd.String("sku", "", "Product SKU")
	showProductTransactions := showProductCmd.Bool("transactions", false, "Show product transactions")

	// Definir flags para delete-product
	deleteProductSKU := deleteProductCmd.String("sku", "", "Product SKU")

	// Definir flags para list-products
	listProductsCategory := listProductsCmd.String("category", "", "Filter by category")

	// Definir flags para add-stock
	addStockSKU := addStockCmd.String("sku", "", "Product SKU")
	addStockQuantity := addStockCmd.Int("quantity", 0, "Quantity to add")
	addStockNotes := addStockCmd.String("notes", "", "Transaction notes")

	// Definir flags para remove-stock
	removeStockSKU := removeStockCmd.String("sku", "", "Product SKU")
	removeStockQuantity := removeStockCmd.Int("quantity", 0, "Quantity to remove")
	removeStockNotes := removeStockCmd.String("notes", "", "Transaction notes")

	// Definir flags para adjust-stock
	adjustStockSKU := adjustStockCmd.String("sku", "", "Product SKU")
	adjustStockQuantity := adjustStockCmd.Int("quantity", 0, "New quantity")
	adjustStockNotes := adjustStockCmd.String("notes", "", "Transaction notes")

	// Definir flags para list-transactions
	listTransactionsProduct := listTransactionsCmd.String("product", "", "Filter by product SKU")
	listTransactionsType := listTransactionsCmd.String("type", "", "Filter by transaction type (PURCHASE, SALE, ADJUST)")

	// Definir flags para report
	reportType := reportCmd.String("type", "", "Report type (low-stock, out-of-stock, inventory-value, transaction-summary)")
	reportThreshold := reportCmd.Int("threshold", 5, "Threshold for low-stock report")

	// Verificar argumentos
	if len(os.Args) < 2 {
		fmt.Println("Expected subcommand")
		printUsage()
		os.Exit(1)
	}

	// Parsear subcomando
	switch os.Args[1] {
	case "add-product":
		addProductCmd.Parse(os.Args[2:])
		if addProductCmd.Parsed() {
			if *addProductName == "" {
				fmt.Println("--name is required")
				addProductCmd.PrintDefaults()
				os.Exit(1)
			}
			if *addProductCategory == "" {
				fmt.Println("--category is required")
				addProductCmd.PrintDefaults()
				os.Exit(1)
			}
			if *addProductPrice <= 0 {
				fmt.Println("--price must be positive")
				addProductCmd.PrintDefaults()
				os.Exit(1)
			}

			category := ProductCategory(*addProductCategory)
			product, err := cli.service.AddProduct(*addProductName, *addProductDesc, category, *addProductPrice, *addProductQuantity)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Product added successfully with SKU: %s\n", product.SKU)
		}

	case "update-product":
		updateProductCmd.Parse(os.Args[2:])
		if updateProductCmd.Parsed() {
			if *updateProductSKU == "" {
				fmt.Println("--sku is required")
				updateProductCmd.PrintDefaults()
				os.Exit(1)
			}

			category := ProductCategory(*updateProductCategory)
			product, err := cli.service.UpdateProduct(*updateProductSKU, *updateProductName, *updateProductDesc, category, *updateProductPrice)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Product updated successfully: %s - %s\n", product.SKU, product.Name)
		}

	case "show-product":
		showProductCmd.Parse(os.Args[2:])
		if showProductCmd.Parsed() {
			if *showProductSKU == "" {
				fmt.Println("--sku is required")
				showProductCmd.PrintDefaults()
				os.Exit(1)
			}

			product, err := cli.service.GetProduct(*showProductSKU)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			// Mostrar detalles del producto
			fmt.Println("Product Details:")
			fmt.Printf("  SKU: %s\n", product.SKU)
			fmt.Printf("  Name: %s\n", product.Name)
			fmt.Printf("  Description: %s\n", product.Description)
			fmt.Printf("  Category: %s\n", product.Category)
			fmt.Printf("  Price: $%.2f\n", product.Price)
			fmt.Printf("  Quantity: %d\n", product.Quantity)
			fmt.Printf("  Created: %s\n", product.CreatedAt.Format(time.RFC3339))
			fmt.Printf("  Updated: %s\n", product.UpdatedAt.Format(time.RFC3339))

			// Mostrar transacciones si se solicita
			if *showProductTransactions {
				transactions, err := cli.service.GetTransactionsByProduct(*showProductSKU)
				if err != nil {
					fmt.Printf("Error getting transactions: %v\n", err)
				} else {
					fmt.Println("\nTransactions:")
					if len(transactions) == 0 {
						fmt.Println("  No transactions found")
					} else {
						for _, t := range transactions {
							fmt.Printf("  %s | %s | Qty: %d | %s | %s\n",
								t.Timestamp.Format("2006-01-02 15:04:05"),
								t.Type,
								t.Quantity,
								t.ID[:8],
								t.Notes)
						}
					}
				}
			}
		}

	case "delete-product":
		deleteProductCmd.Parse(os.Args[2:])
		if deleteProductCmd.Parsed() {
			if *deleteProductSKU == "" {
				fmt.Println("--sku is required")
				deleteProductCmd.PrintDefaults()
				os.Exit(1)
			}

			err := cli.service.DeleteProduct(*deleteProductSKU)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Product deleted successfully")
		}

	case "list-products":
		listProductsCmd.Parse(os.Args[2:])
		if listProductsCmd.Parsed() {
			var products []Product
			var err error

			if *listProductsCategory != "" {
				category := ProductCategory(*listProductsCategory)
				products, err = cli.service.GetProductsByCategory(category)
			} else {
				products, err = cli.service.GetAllProducts()
			}

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			if len(products) == 0 {
				fmt.Println("No products found")
			} else {
				fmt.Println("Products:")
				fmt.Printf("%-12s | %-30s | %-15s | %-8s | %-10s\n", "SKU", "Name", "Category", "Price", "Quantity")
				fmt.Println(strings.Repeat("-", 85))

				for _, p := range products {
					fmt.Printf("%-12s | %-30s | %-15s | $%-7.2f | %-10d\n",
						p.SKU,
						truncateString(p.Name, 30),
						p.Category,
						p.Price,
						p.Quantity)
				}
			}
		}

	case "add-stock":
		addStockCmd.Parse(os.Args[2:])
		if addStockCmd.Parsed() {
			if *addStockSKU == "" {
				fmt.Println("--sku is required")
				addStockCmd.PrintDefaults()
				os.Exit(1)
			}
			if *addStockQuantity <= 0 {
				fmt.Println("--quantity must be positive")
				addStockCmd.PrintDefaults()
				os.Exit(1)
			}

			product, err := cli.service.AddStock(*addStockSKU, *addStockQuantity, *addStockNotes)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Stock added successfully. New quantity for %s: %d\n", product.Name, product.Quantity)
		}

	case "remove-stock":
		removeStockCmd.Parse(os.Args[2:])
		if removeStockCmd.Parsed() {
			if *removeStockSKU == "" {
				fmt.Println("--sku is required")
				removeStockCmd.PrintDefaults()
				os.Exit(1)
			}
			if *removeStockQuantity <= 0 {
				fmt.Println("--quantity must be positive")
				removeStockCmd.PrintDefaults()
				os.Exit(1)
			}

			product, err := cli.service.RemoveStock(*removeStockSKU, *removeStockQuantity, *removeStockNotes)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Stock removed successfully. New quantity for %s: %d\n", product.Name, product.Quantity)
		}

	case "adjust-stock":
		adjustStockCmd.Parse(os.Args[2:])
		if adjustStockCmd.Parsed() {
			if *adjustStockSKU == "" {
				fmt.Println("--sku is required")
				adjustStockCmd.PrintDefaults()
				os.Exit(1)
			}
			if *adjustStockQuantity < 0 {
				fmt.Println("--quantity must be non-negative")
				adjustStockCmd.PrintDefaults()
				os.Exit(1)
			}

			product, err := cli.service.AdjustStock(*adjustStockSKU, *adjustStockQuantity, *adjustStockNotes)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Stock adjusted successfully. New quantity for %s: %d\n", product.Name, product.Quantity)
		}

	case "list-transactions":
		listTransactionsCmd.Parse(os.Args[2:])
		if listTransactionsCmd.Parsed() {
			var transactions []Transaction
			var err error

			if *listTransactionsProduct != "" {
				transactions, err = cli.service.GetTransactionsByProduct(*listTransactionsProduct)
			} else if *listTransactionsType != "" {
				transactions, err = cli.service.GetTransactionsByType(TransactionType(*listTransactionsType))
			} else {
				transactions, err = cli.service.GetAllTransactions()
			}

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			if len(transactions) == 0 {
				fmt.Println("No transactions found")
			} else {
				fmt.Println("Transactions:")
				fmt.Printf("%-20s | %-12s | %-10s | %-8s | %-30s\n", "Timestamp", "Product", "Type", "Quantity", "Notes")
				fmt.Println(strings.Repeat("-", 90))

				for _, t := range transactions {
					fmt.Printf("%-20s | %-12s | %-10s | %-8d | %-30s\n",
						t.Timestamp.Format("2006-01-02 15:04:05"),
						t.ProductSKU,
						t.Type,
						t.Quantity,
						truncateString(t.Notes, 30))
				}
			}
		}

	case "report":
		reportCmd.Parse(os.Args[2:])
		if reportCmd.Parsed() {
			if *reportType == "" {
				fmt.Println("--type is required")
				reportCmd.PrintDefaults()
				os.Exit(1)
			}

			switch *reportType {
			case "low-stock":
				products, err := cli.service.GetLowStockProducts(*reportThreshold)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Low Stock Products (Threshold: %d):\n", *reportThreshold)
				if len(products) == 0 {
					fmt.Println("No products with low stock")
				} else {
					fmt.Printf("%-12s | %-30s | %-15s | %-8s | %-10s\n", "SKU", "Name", "Category", "Price", "Quantity")
					fmt.Println(strings.Repeat("-", 85))

					for _, p := range products {
						fmt.Printf("%-12s | %-30s | %-15s | $%-7.2f | %-10d\n",
							p.SKU,
							truncateString(p.Name, 30),
							p.Category,
							p.Price,
							p.Quantity)
					}
				}

			case "out-of-stock":
				products, err := cli.service.GetOutOfStockProducts()
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				fmt.Println("Out of Stock Products:")
				if len(products) == 0 {
					fmt.Println("No products out of stock")
				} else {
					fmt.Printf("%-12s | %-30s | %-15s | %-8s\n", "SKU", "Name", "Category", "Price")
					fmt.Println(strings.Repeat("-", 75))

					for _, p := range products {
						fmt.Printf("%-12s | %-30s | %-15s | $%-7.2f\n",
							p.SKU,
							truncateString(p.Name, 30),
							p.Category,
							p.Price)
					}
				}

			case "inventory-value":
				value, err := cli.service.GetInventoryValue()
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Total Inventory Value: $%.2f\n", value)

			case "transaction-summary":
				summary, err := cli.service.GetTransactionSummary()
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				fmt.Println("Transaction Summary:")
				fmt.Printf("  Purchases: %d items\n", summary[Purchase])
				fmt.Printf("  Sales: %d items\n", summary[Sale])
				fmt.Printf("  Adjustments: %d items\n", summary[Adjust])

			default:
				fmt.Printf("Unknown report type: %s\n", *reportType)
				reportCmd.PrintDefaults()
				os.Exit(1)
			}
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  inventory_system <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  add-product     Add a new product")
	fmt.Println("  update-product  Update an existing product")
	fmt.Println("  show-product    Show product details")
	fmt.Println("  delete-product  Delete a product")
	fmt.Println("  list-products   List all products")
	fmt.Println("  add-stock       Add stock to a product")
	fmt.Println("  remove-stock    Remove stock from a product")
	fmt.Println("  adjust-stock    Adjust stock to a specific quantity")
	fmt.Println("  list-transactions List transactions")
	fmt.Println("  report          Generate reports")
	fmt.Println("\nRun 'inventory_system <command> --help' for more information on a command.")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func main() {
	cli := NewCLI()
	cli.Run()
}