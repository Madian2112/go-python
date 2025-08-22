# Aplicación de Gestión de Tareas (TODO) en Go

Esta es una aplicación de línea de comandos para gestionar tareas (TODO), implementada en Go siguiendo las mejores prácticas de desarrollo.

## Características

- Añadir nuevas tareas
- Listar tareas (todas o solo pendientes)
- Marcar tareas como completadas
- Eliminar tareas
- Persistencia de datos en formato JSON

## Requisitos

- Go 1.13 o superior

## Estructura del Proyecto

```
todo_app/
├── main.go         # Código principal de la aplicación
├── README.md       # Este archivo
└── data/           # Directorio para almacenar datos (creado automáticamente)
    └── tasks.json  # Archivo de datos (creado automáticamente)
```

## Compilación y Uso

### Compilar la aplicación

```bash
go build -o todo
```

### Listar tareas

```bash
./todo list         # Listar tareas pendientes
./todo list --all   # Listar todas las tareas
```

### Añadir una tarea

```bash
./todo add "Completar informe mensual"
```

### Marcar una tarea como completada

```bash
./todo complete 1   # Completar la tarea con ID 1
```

### Eliminar una tarea

```bash
./todo delete 1     # Eliminar la tarea con ID 1
```

## Buenas Prácticas Implementadas

### 1. Estructura Modular

El código está organizado en tipos y funciones con responsabilidades bien definidas:

- `Task`: Representa una tarea individual
- `TaskRepository`: Gestiona la persistencia de las tareas (patrón Repository)
- `TodoApp`: Implementa la lógica de la aplicación

**Referencia**: [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

### 2. Manejo de Errores Explícito

Siguiendo la filosofía de Go, los errores se manejan de forma explícita y se propagan adecuadamente.

**Referencia**: [Error Handling in Go](https://go.dev/blog/error-handling-and-go)

### 3. Uso de Structs y Métodos

Se utilizan structs para modelar los datos y métodos asociados para implementar comportamientos.

**Referencia**: [Effective Go - Structs](https://go.dev/doc/effective_go#structs)

### 4. Uso de Interfaces Implícitas

Aunque no se definen interfaces explícitas en este ejemplo simple, el diseño permite la implementación de interfaces para facilitar pruebas y extensibilidad.

**Referencia**: [Effective Go - Interfaces](https://go.dev/doc/effective_go#interfaces)

### 5. Manejo de Flags para Argumentos de Línea de Comandos

Se utiliza el paquete `flag` de la biblioteca estándar para procesar argumentos de línea de comandos.

**Referencia**: [Command-line flags in Go](https://gobyexample.com/command-line-flags)

### 6. Serialización JSON

Se utiliza el paquete `encoding/json` para la serialización y deserialización de datos.

**Referencia**: [JSON and Go](https://go.dev/blog/json)

### 7. Gestión de Recursos

Se implementa una gestión adecuada de recursos, como el manejo de archivos y directorios.

**Referencia**: [Working with Files in Go](https://www.devdungeon.com/content/working-files-go)

## Principios de Diseño Aplicados

1. **Simplicidad**: Go favorece la simplicidad y la claridad sobre la complejidad y la abstracción excesiva.

2. **Composición sobre Herencia**: Se utiliza la composición de tipos en lugar de la herencia.

3. **Manejo Explícito de Errores**: Los errores se manejan de forma explícita y se propagan adecuadamente.

4. **Concurrencia Segura**: Aunque este ejemplo no utiliza concurrencia, el diseño permite añadirla de forma segura si fuera necesario.

**Referencia**: [Go Proverbs](https://go-proverbs.github.io/)

## Mejoras Potenciales

- Implementar pruebas unitarias
- Añadir funcionalidad para editar tareas existentes
- Implementar categorías o etiquetas para las tareas
- Añadir fechas límite para las tareas
- Utilizar una base de datos SQL o NoSQL para la persistencia
- Implementar una API REST para acceder a las tareas desde otros sistemas
- Añadir soporte para concurrencia en operaciones de lectura/escritura

## Comparación con la Versión en Python

### Ventajas de Go

- Rendimiento superior, especialmente en aplicaciones con alta concurrencia
- Compilación a un único binario ejecutable sin dependencias externas
- Tipado estático que ayuda a prevenir errores en tiempo de compilación
- Manejo de errores explícito que fomenta la robustez

### Ventajas de Python

- Sintaxis más concisa y expresiva
- Mayor velocidad de desarrollo para prototipos y aplicaciones simples
- Ecosistema más amplio de bibliotecas para diversos dominios
- Curva de aprendizaje más suave para principiantes

**Referencia**: [Go vs Python: Choosing the Right Language](https://www.bacancytechnology.com/blog/go-vs-python)