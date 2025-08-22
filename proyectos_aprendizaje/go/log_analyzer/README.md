# Analizador de Logs en Go

## Descripción
Este proyecto es un analizador de archivos de log de nivel básico (Nivel 1) implementado en Go. Permite filtrar entradas por nivel de severidad, componente o palabra clave, y generar estadísticas básicas. Es una aplicación de consola que demuestra conceptos fundamentales de Go como el manejo de archivos, expresiones regulares, estructuras de datos y programación orientada a interfaces.

## Características
- Parseo de archivos de log con formato específico
- Filtrado de entradas por nivel (DEBUG, INFO, WARNING, ERROR, CRITICAL)
- Filtrado de entradas por componente o servicio
- Búsqueda de entradas por palabra clave
- Generación de estadísticas de distribución por nivel, componente y hora
- Resumen de errores agrupados por componente
- Salida en consola con colores para mejor visualización

## Requisitos
- Go 1.16 o superior
- No requiere bibliotecas externas (solo usa la biblioteca estándar de Go)

## Estructura del Proyecto
```
log_analyzer/
├── main.go            # Aplicación principal
├── sample.log         # Archivo de log de ejemplo
└── README.md          # Este archivo
```

## Compilación y Uso

### Compilación
```bash
go build -o log_analyzer
```

### Ejecución básica
```bash
# Si está compilado
./log_analyzer sample.log

# O directamente con go run
go run main.go sample.log
```

### Filtrar por nivel de log
```bash
go run main.go --level ERROR sample.log
```

### Filtrar por componente
```bash
go run main.go --component DatabaseService sample.log
```

### Buscar por palabra clave
```bash
go run main.go --keyword "failed" sample.log
```

### Mostrar estadísticas
```bash
go run main.go --stats sample.log
```

### Mostrar resumen de errores
```bash
go run main.go --errors sample.log
```

### Combinar filtros
```bash
go run main.go --level WARNING --component SecurityModule --stats sample.log
```

## Conceptos Aplicados

### Conceptos Básicos de Go
- Uso de estructuras (structs) y métodos
- Manejo de archivos y buffers
- Expresiones regulares
- Estructuras de datos (slices, maps)
- Manejo de errores explícito
- Formateo de strings
- Argumentos de línea de comandos con el paquete flag
- Ordenamiento de datos

### Buenas Prácticas
- Documentación con comentarios
- Separación de responsabilidades (tipos LogEntry, LogAnalyzer)
- Código modular y reutilizable
- Manejo adecuado de errores
- Nombres descriptivos para variables y funciones
- Uso de constantes para valores fijos (colores, patrones)
- Cierre adecuado de recursos (defer)

## Problemas del Mundo Real que Resuelve
- Análisis de logs de servidores y aplicaciones
- Detección de errores y problemas en sistemas
- Monitoreo de actividad de componentes específicos
- Identificación de patrones de uso o error

## Diferencias con la Implementación en Python

### Ventajas de Go
- Compilación a un único binario ejecutable
- Mejor rendimiento para procesamiento de archivos grandes
- Manejo de errores explícito
- Tipado estático que previene errores en tiempo de compilación
- Concurrencia nativa (aunque no se utiliza en este ejemplo básico)

### Ventajas de Python
- Código más conciso y legible
- Mayor abstracción con clases y herencia
- Biblioteca estándar más rica (Counter, defaultdict)
- Desarrollo más rápido y prototipado
- Mayor flexibilidad en la estructura del código

## Posibles Mejoras
- Implementación de concurrencia para procesar archivos grandes
- Soporte para diferentes formatos de log
- Exportación de resultados a CSV o JSON
- Visualización gráfica de estadísticas
- Análisis de tendencias a lo largo del tiempo
- Detección automática de anomalías

## Licencia
Este proyecto es de código abierto y está disponible bajo la licencia MIT.