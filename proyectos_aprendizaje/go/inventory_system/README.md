# Sistema de Gestión de Inventario en Go

## Descripción
Este proyecto es un sistema de gestión de inventario de nivel intermedio (Nivel 2) implementado en Go. Permite administrar productos, realizar seguimiento de transacciones y generar reportes a través de una interfaz de línea de comandos. Demuestra conceptos intermedios de Go como interfaces, estructuras, manejo de errores y persistencia de datos.

## Características
- Gestión completa de productos (añadir, actualizar, eliminar, listar)
- Control de inventario (añadir, reducir, ajustar stock)
- Seguimiento de transacciones (compras, ventas, ajustes)
- Generación de reportes (stock bajo, sin stock, valor del inventario)
- Búsqueda y filtrado de productos por diferentes criterios
- Persistencia de datos en archivos JSON
- Interfaz de línea de comandos completa

## Requisitos
- Go 1.13 o superior
- Dependencia externa: github.com/google/uuid (para generación de identificadores únicos)

## Estructura del Proyecto
```
inventory_system/
├── main.go           # Archivo principal que contiene toda la implementación
├── products.json     # Archivo de persistencia para productos (generado automáticamente)
├── transactions.json # Archivo de persistencia para transacciones (generado automáticamente)
└── README.md         # Este archivo
```

## Patrones de Diseño Implementados

### 1. Data Access Object (DAO)
Implementado a través de las interfaces `ProductDAO` y `TransactionDAO`, este patrón separa la lógica de acceso a datos de la lógica de negocio, proporcionando una interfaz abstracta para la persistencia de datos.

### 2. Facade
Implementado en la estructura `InventoryService`, este patrón proporciona una interfaz simplificada al sistema complejo de gestión de inventario, ocultando la complejidad de las interacciones entre los diferentes componentes.

### 3. Command
Implementado en la estructura `CLI` y su método `Run()`, este patrón encapsula una solicitud como un objeto, permitiendo parametrizar clientes con diferentes solicitudes.

## Compilación
```bash
# Instalar dependencias
go get github.com/google/uuid

# Compilar el programa
go build -o inventory_system main.go
```

## Uso

### Añadir un producto
```bash
./inventory_system add-product --name "Laptop HP" --description "Laptop HP 15.6 pulgadas" --category ELECTRONICS --price 599.99 --quantity 10
```

### Listar productos
```bash
./inventory_system list-products
```

### Filtrar productos por categoría
```bash
./inventory_system list-products --category ELECTRONICS
```

### Mostrar detalles de un producto
```bash
./inventory_system show-product --sku ABC123 --transactions
```

### Añadir stock a un producto
```bash
./inventory_system add-stock --sku ABC123 --quantity 5 --notes "Reposición de inventario"
```

### Reducir stock (venta)
```bash
./inventory_system remove-stock --sku ABC123 --quantity 2 --notes "Venta a cliente"
```

### Generar reporte de productos con stock bajo
```bash
./inventory_system report --type low-stock --threshold 3
```

### Generar reporte de valor del inventario
```bash
./inventory_system report --type inventory-value
```

## Conceptos Aplicados

### Conceptos Intermedios de Go
- Interfaces y polimorfismo
- Estructuras y métodos
- Manejo de errores con el patrón de retorno de error
- Punteros y referencias
- Marshaling y unmarshaling de JSON
- Manejo de fechas y horas
- Flags y argumentos de línea de comandos
- Generación de identificadores únicos (UUID)
- Closures y funciones anónimas
- Manejo de archivos y directorios

### Buenas Prácticas
- Separación de responsabilidades (modelos, servicios, interfaz)
- Código modular y reutilizable
- Manejo explícito de errores
- Validación de datos de entrada
- Documentación con comentarios
- Nombres descriptivos para variables, métodos y estructuras
- Uso de constantes y enumeraciones para valores fijos

## Problemas del Mundo Real que Resuelve
- Gestión de inventario para pequeñas y medianas empresas
- Seguimiento de stock y transacciones
- Valoración de inventario
- Identificación de productos con stock bajo o agotado
- Registro de movimientos de inventario

## Diferencias con la Implementación en Python
- Uso de un único archivo en lugar de múltiples módulos
- Implementación basada en estructuras e interfaces en lugar de clases
- Manejo explícito de errores en lugar de excepciones
- Uso de punteros para modificar estructuras
- Tipado estático en lugar de tipado dinámico
- Compilación a un ejecutable binario en lugar de interpretación

## Posibles Mejoras
- Separación del código en múltiples archivos para mejor organización
- Implementación de pruebas unitarias
- Soporte para múltiples almacenes o ubicaciones
- Integración con bases de datos relacionales o NoSQL
- Implementación de una API REST para acceso remoto
- Soporte para múltiples usuarios y autenticación
- Generación de reportes más avanzados y exportación a diferentes formatos

## Licencia
Este proyecto es de código abierto y está disponible bajo la licencia MIT.