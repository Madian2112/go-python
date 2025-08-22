# Analizador de Logs en Python

## Descripción
Este proyecto es un analizador de archivos de log de nivel básico (Nivel 1) que permite filtrar entradas por nivel de severidad, componente o palabra clave, y generar estadísticas básicas. Es una aplicación de consola que demuestra conceptos fundamentales de Python como el manejo de archivos, expresiones regulares, estructuras de datos y programación orientada a objetos.

## Características
- Parseo de archivos de log con formato específico
- Filtrado de entradas por nivel (DEBUG, INFO, WARNING, ERROR, CRITICAL)
- Filtrado de entradas por componente o servicio
- Búsqueda de entradas por palabra clave
- Generación de estadísticas de distribución por nivel, componente y hora
- Resumen de errores agrupados por componente
- Salida en consola con colores para mejor visualización

## Requisitos
- Python 3.6 o superior
- No requiere bibliotecas externas (solo usa la biblioteca estándar de Python)

## Estructura del Proyecto
```
log_analyzer/
├── log_analyzer.py     # Aplicación principal
├── sample.log          # Archivo de log de ejemplo
└── README.md           # Este archivo
```

## Uso

### Ejecución básica
```bash
python log_analyzer.py sample.log
```

### Filtrar por nivel de log
```bash
python log_analyzer.py sample.log --level ERROR
```

### Filtrar por componente
```bash
python log_analyzer.py sample.log --component DatabaseService
```

### Buscar por palabra clave
```bash
python log_analyzer.py sample.log --keyword "failed"
```

### Mostrar estadísticas
```bash
python log_analyzer.py sample.log --stats
```

### Mostrar resumen de errores
```bash
python log_analyzer.py sample.log --errors
```

### Combinar filtros
```bash
python log_analyzer.py sample.log --level WARNING --component SecurityModule --stats
```

## Conceptos Aplicados

### Conceptos Básicos de Python
- Uso de clases y objetos
- Manejo de archivos
- Expresiones regulares
- Estructuras de datos (listas, diccionarios, defaultdict, Counter)
- Tipado estático con anotaciones de tipo
- Manejo de excepciones
- Argumentos de línea de comandos

### Buenas Prácticas
- Documentación con docstrings
- Separación de responsabilidades (clases LogEntry, LogAnalyzer, LogAnalyzerCLI)
- Código modular y reutilizable
- Manejo adecuado de errores
- Nombres descriptivos para variables y funciones
- Uso de constantes para valores fijos (colores, niveles de log)

## Problemas del Mundo Real que Resuelve
- Análisis de logs de servidores y aplicaciones
- Detección de errores y problemas en sistemas
- Monitoreo de actividad de componentes específicos
- Identificación de patrones de uso o error

## Posibles Mejoras
- Soporte para diferentes formatos de log
- Exportación de resultados a CSV o JSON
- Visualización gráfica de estadísticas
- Análisis de tendencias a lo largo del tiempo
- Detección automática de anomalías

## Licencia
Este proyecto es de código abierto y está disponible bajo la licencia MIT.