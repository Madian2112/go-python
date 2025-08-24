# Introducción a Go (Golang)

## ¿Qué es Go?

Go (o Golang) es un lenguaje de programación de código abierto desarrollado por Google en 2007 y lanzado públicamente en 2009. Fue creado por Robert Griesemer, Rob Pike y Ken Thompson con el objetivo de ser un lenguaje eficiente, confiable y fácil de usar para construir software a gran escala.

## Características principales

- **Simplicidad**: Sintaxis clara y concisa con pocas palabras clave.
- **Compilado**: Se compila directamente a código máquina para diferentes plataformas.
- **Tipado estático**: Los tipos se verifican en tiempo de compilación.
- **Recolección de basura**: Gestión automática de memoria.
- **Concurrencia integrada**: Goroutines y canales para programación concurrente.
- **Biblioteca estándar robusta**: Incluye paquetes para HTTP, JSON, criptografía, etc.
- **Compilación rápida**: Tiempos de compilación extremadamente cortos.
- **Multiplataforma**: Funciona en Windows, macOS, Linux y otros sistemas.

## ¿Por qué usar Go?

- **Rendimiento**: Ofrece velocidad cercana a C/C++ con mayor facilidad de desarrollo.
- **Simplicidad**: Fácil de aprender y mantener.
- **Concurrencia**: Modelo de concurrencia simple y potente.
- **Compilación rápida**: Mejora la productividad del desarrollador.
- **Despliegue sencillo**: Genera binarios estáticos sin dependencias externas.
- **Ecosistema creciente**: Amplia adopción en la industria, especialmente en infraestructura y DevOps.

## Instalación

### Windows

1. Visita [golang.org/dl](https://golang.org/dl/)
2. Descarga el instalador MSI para Windows
3. Ejecuta el instalador y sigue las instrucciones
4. Verifica la instalación abriendo una terminal y escribiendo:
   ```
   go version
   ```

### macOS

1. **Usando Homebrew**:
   ```
   brew install go
   ```

2. **Instalación manual**:
   - Descarga el paquete desde [golang.org/dl](https://golang.org/dl/)
   - Abre el archivo .pkg y sigue las instrucciones
   - Verifica con `go version`

### Linux

**Ubuntu/Debian**:
```
sudo apt update
sudo apt install golang-go
```

**Fedora**:
```
sudo dnf install golang
```

**Arch Linux**:
```
sudo pacman -S go
```

## Configuración del entorno

Go utiliza una variable de entorno llamada `GOPATH` que define el espacio de trabajo. A partir de Go 1.8, si no se establece, se usa un valor predeterminado (`$HOME/go` en Unix/Linux/macOS o `%USERPROFILE%\go` en Windows).

### Estructura de directorios recomendada

```
$GOPATH/
  ├── bin/       # Ejecutables compilados
  ├── pkg/       # Paquetes compilados
  └── src/       # Código fuente
      └── github.com/
          └── tuusuario/
              └── tuproyecto/
                  ├── main.go
                  └── ...
```

A partir de Go 1.11, puedes usar módulos de Go para gestionar dependencias fuera del `GOPATH`.

## Tu primer programa en Go

Crea un archivo llamado `hola.go` con el siguiente contenido:

```go
package main

import "fmt"

func main() {
    fmt.Println("¡Hola, Mundo!")
}
```

Compila y ejecuta el programa:

```
go run hola.go
```

Deberías ver el mensaje "¡Hola, Mundo!" en la pantalla.

También puedes compilar el programa para generar un ejecutable:

```
go build hola.go
```

Esto creará un archivo ejecutable (`hola` en Unix/Linux/macOS o `hola.exe` en Windows) que puedes ejecutar directamente.

## Explicación del código

```go
package main  // Declara que este archivo pertenece al paquete main

import "fmt"  // Importa el paquete fmt para entrada/salida formateada

// La función main es el punto de entrada del programa
func main() {
    fmt.Println("¡Hola, Mundo!")  // Imprime un mensaje en la consola
}
```

- **package main**: Todo programa en Go debe pertenecer a un paquete. El paquete `main` es especial porque indica que el programa es ejecutable.
- **import "fmt"**: Importa el paquete `fmt` de la biblioteca estándar, que proporciona funciones para entrada/salida formateada.
- **func main()**: La función `main` es el punto de entrada de un programa ejecutable en Go.
- **fmt.Println()**: Función que imprime texto en la consola y añade un salto de línea.

## Herramientas de Go

Go viene con varias herramientas útiles:

- **go run**: Compila y ejecuta un programa Go
- **go build**: Compila paquetes y dependencias
- **go install**: Compila e instala paquetes y dependencias
- **go test**: Ejecuta pruebas
- **go fmt**: Formatea código Go según el estilo estándar
- **go get**: Descarga e instala paquetes y dependencias
- **go mod**: Gestiona módulos y dependencias
- **go doc**: Muestra documentación de paquetes

## Entornos de Desarrollo Integrados (IDEs)

Algunos IDEs populares para Go:

- **Visual Studio Code** con la extensión Go
- **GoLand** de JetBrains
- **Vim/Neovim** con plugins para Go
- **Sublime Text** con plugins para Go
- **Atom** con plugins para Go

## Filosofía de Go

Go fue diseñado con ciertos principios en mente:

- **Simplicidad**: Menos es más. Go tiene pocas características pero bien pensadas.
- **Legibilidad**: El código debe ser fácil de leer y entender.
- **Eficiencia**: Debe ser rápido tanto en tiempo de compilación como de ejecución.
- **Concurrencia**: Debe facilitar la programación concurrente.

## Recursos adicionales

- [Sitio oficial de Go](https://golang.org/)
- [Tour de Go](https://tour.golang.org/) - Tutorial interactivo oficial
- [Effective Go](https://golang.org/doc/effective_go) - Guía de mejores prácticas
- [Go by Example](https://gobyexample.com/) - Ejemplos prácticos
- [Documentación oficial](https://golang.org/doc/)
- [Playground de Go](https://play.golang.org/) - Prueba código Go en línea

---

En la siguiente sección, aprenderemos sobre variables, tipos de datos y operadores en Go.