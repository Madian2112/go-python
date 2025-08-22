# Aplicación de Gestión de Tareas (TODO) en Python

Esta es una aplicación de línea de comandos para gestionar tareas (TODO), implementada en Python siguiendo las mejores prácticas de desarrollo.

## Características

- Añadir nuevas tareas
- Listar tareas (todas o solo pendientes)
- Marcar tareas como completadas
- Eliminar tareas
- Persistencia de datos en formato JSON

## Requisitos

- Python 3.6 o superior

## Estructura del Proyecto

```
todo_app/
├── todo.py         # Script principal de la aplicación
├── README.md       # Este archivo
└── data/           # Directorio para almacenar datos (creado automáticamente)
    └── tasks.json  # Archivo de datos (creado automáticamente)
```

## Uso

### Listar tareas

```bash
python todo.py list         # Listar tareas pendientes
python todo.py list --all   # Listar todas las tareas
```

### Añadir una tarea

```bash
python todo.py add "Completar informe mensual"
```

### Marcar una tarea como completada

```bash
python todo.py complete 1   # Completar la tarea con ID 1
```

### Eliminar una tarea

```bash
python todo.py delete 1     # Eliminar la tarea con ID 1
```

## Buenas Prácticas Implementadas

### 1. Estructura Modular y Orientada a Objetos

El código está organizado en clases con responsabilidades bien definidas:

- `Task`: Representa una tarea individual
- `TaskRepository`: Gestiona la persistencia de las tareas (patrón Repository)
- `TodoApp`: Implementa la lógica de la aplicación

**Referencia**: [Clean Architecture in Python](https://www.thedigitalcatonline.com/blog/2016/11/14/clean-architectures-in-python-a-step-by-step-example/)

### 2. Tipado Estático con Type Hints

Se utilizan anotaciones de tipo para mejorar la legibilidad y permitir la verificación estática de tipos.

**Referencia**: [Python Type Checking Guide](https://realpython.com/python-type-checking/)

### 3. Documentación Completa

Se utilizan docstrings para documentar clases, métodos y funciones, siguiendo el estilo de Google.

**Referencia**: [Google Python Style Guide - Docstrings](https://google.github.io/styleguide/pyguide.html#38-comments-and-docstrings)

### 4. Manejo de Errores

Se implementa un manejo adecuado de excepciones para proporcionar mensajes de error claros y evitar fallos inesperados.

**Referencia**: [Python Exception Handling Techniques](https://realpython.com/python-exceptions/)

### 5. Separación de Preocupaciones

El código separa claramente:
- Lógica de negocio (clases Task y TodoApp)
- Persistencia de datos (TaskRepository)
- Interfaz de usuario (función main y parse_arguments)

**Referencia**: [The Importance of Separation of Concerns](https://dev.to/tamerlang/separation-of-concerns-the-simple-way-4jp2)

### 6. Uso de Argumentos de Línea de Comandos

Se utiliza el módulo `argparse` para procesar argumentos de línea de comandos de forma robusta.

**Referencia**: [Command-Line Parsing in Python](https://realpython.com/command-line-interfaces-python-argparse/)

### 7. Gestión de Recursos

Se utilizan bloques `with` para el manejo de archivos, asegurando que se cierren correctamente incluso en caso de error.

**Referencia**: [Context Managers and the 'with' Statement](https://realpython.com/python-with-statement/)

## Principios de Diseño Aplicados

1. **DRY (Don't Repeat Yourself)**: El código evita la duplicación mediante la abstracción adecuada.

2. **SOLID**:
   - **S** (Responsabilidad Única): Cada clase tiene una única responsabilidad.
   - **O** (Abierto/Cerrado): El diseño permite extender la funcionalidad sin modificar el código existente.
   - **L** (Sustitución de Liskov): Las subclases podrían sustituir a sus clases base.
   - **I** (Segregación de Interfaces): No se fuerzan dependencias innecesarias.
   - **D** (Inversión de Dependencias): Las dependencias son inyectadas (por ejemplo, el repositorio).

3. **Fail Fast**: Los errores se detectan y reportan lo antes posible.

**Referencia**: [SOLID Principles in Python](https://www.digitalocean.com/community/conceptual-articles/s-o-l-i-d-the-first-five-principles-of-object-oriented-design)

## Mejoras Potenciales

- Implementar pruebas unitarias con pytest
- Añadir funcionalidad para editar tareas existentes
- Implementar categorías o etiquetas para las tareas
- Añadir fechas límite para las tareas
- Crear una interfaz de usuario más interactiva con una biblioteca como prompt_toolkit