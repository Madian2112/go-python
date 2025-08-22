# Sistema de Gestión de Inventario en Python

## Descripción
Este proyecto es un sistema de gestión de inventario de nivel intermedio (Nivel 2) que permite administrar productos, realizar seguimiento de transacciones y generar reportes. Es una aplicación de consola que demuestra conceptos intermedios de Python como patrones de diseño, programación orientada a objetos avanzada, y manejo de datos persistentes.

## Características
- Gestión completa de productos (añadir, actualizar, eliminar, listar)
- Control de inventario (añadir, reducir, ajustar stock)
- Seguimiento de transacciones (compras, ventas, ajustes)
- Generación de reportes (stock bajo, sin stock, valor del inventario)
- Búsqueda y filtrado de productos por diferentes criterios
- Persistencia de datos en archivos JSON
- Interfaz de línea de comandos completa

## Requisitos
- Python 3.7 o superior
- No requiere bibliotecas externas (solo usa la biblioteca estándar de Python)

## Estructura del Proyecto
```
inventory_system/
├── models.py          # Definición de modelos y DAOs
├── services.py        # Servicios y lógica de negocio
├── inventory_cli.py   # Interfaz de línea de comandos
├── products.json      # Archivo de persistencia para productos (generado automáticamente)
├── transactions.json  # Archivo de persistencia para transacciones (generado automáticamente)
└── README.md          # Este archivo
```

## Patrones de Diseño Implementados

### 1. Data Access Object (DAO)
Implementado en `models.py`, este patrón separa la lógica de acceso a datos de la lógica de negocio, proporcionando una interfaz abstracta para la persistencia de datos.

### 2. Facade
Implementado en `services.py`, este patrón proporciona una interfaz simplificada al sistema complejo de gestión de inventario, ocultando la complejidad de las interacciones entre los diferentes componentes.

### 3. Command
Implementado en `inventory_cli.py`, este patrón encapsula una solicitud como un objeto, permitiendo parametrizar clientes con diferentes solicitudes, encolar o registrar solicitudes, y soportar operaciones reversibles.

### 4. Factory Method
Implementado en `inventory_cli.py` con el método `_create_command()`, este patrón define una interfaz para crear un objeto, pero deja que las subclases decidan qué clase instanciar.

## Uso

### Añadir un producto
```bash
python inventory_cli.py add-product --name "Laptop HP" --description "Laptop HP 15.6 pulgadas" --category ELECTRONICS --price 599.99 --quantity 10
```

### Listar productos
```bash
python inventory_cli.py list-products
```

### Filtrar productos por categoría
```bash
python inventory_cli.py list-products --category ELECTRONICS
```

### Mostrar detalles de un producto
```bash
python inventory_cli.py show-product --sku ABC123 --transactions
```

### Añadir stock a un producto
```bash
python inventory_cli.py add-stock --sku ABC123 --quantity 5 --notes "Reposición de inventario"
```

### Reducir stock (venta)
```bash
python inventory_cli.py remove-stock --sku ABC123 --quantity 2 --notes "Venta a cliente"
```

### Generar reporte de productos con stock bajo
```bash
python inventory_cli.py report --type low-stock --threshold 3
```

### Generar reporte de valor del inventario
```bash
python inventory_cli.py report --type inventory-value
```

## Conceptos Aplicados

### Conceptos Intermedios de Python
- Patrones de diseño (DAO, Facade, Command, Factory Method)
- Clases abstractas e interfaces (ABC, abstractmethod)
- Dataclasses para modelos de datos
- Enumeraciones tipadas (Enum)
- Anotaciones de tipo avanzadas (Union, Optional)
- Serialización y deserialización de objetos
- Manejo de fechas y horas
- Generación de identificadores únicos (UUID)

### Buenas Prácticas
- Separación de responsabilidades (modelos, servicios, interfaz)
- Código modular y reutilizable
- Manejo adecuado de errores y excepciones
- Validación de datos de entrada
- Documentación con docstrings
- Nombres descriptivos para variables, métodos y clases
- Uso de constantes y enumeraciones para valores fijos

## Problemas del Mundo Real que Resuelve
- Gestión de inventario para pequeñas y medianas empresas
- Seguimiento de stock y transacciones
- Valoración de inventario
- Identificación de productos con stock bajo o agotado
- Registro de movimientos de inventario

## Posibles Mejoras
- Implementación de una interfaz gráfica de usuario
- Soporte para múltiples almacenes o ubicaciones
- Integración con sistemas de punto de venta
- Generación de reportes más avanzados y gráficos
- Exportación de datos a diferentes formatos
- Implementación de un sistema de usuarios y permisos
- Integración con bases de datos relacionales o NoSQL

## Licencia
Este proyecto es de código abierto y está disponible bajo la licencia MIT.