# 19. Herramientas de Go (Go Toolchain)

## Introducción

Go incluye un conjunto de herramientas de línea de comandos potentes que simplifican el ciclo de desarrollo, desde la creación y prueba hasta la gestión de dependencias.

## Comandos Esenciales:

- **`go run`**: Compila y ejecuta uno o más archivos fuente de Go.
- **`go build`**: Compila los paquetes y sus dependencias, pero no los ejecuta.
- **`go install`**: Compila e instala paquetes y comandos.
- **`go test`**: Ejecuta las pruebas para el paquete actual.
- **`go mod`**: Proporciona acceso a operaciones en los módulos de Go.
  - `go mod init`
  - `go mod tidy`
  - `go mod download`
- **`go fmt`**: Formatea el código fuente según las convenciones de Go.
- **`go vet`**: Examina el código fuente de Go y reporta posibles errores.
- **`go doc`**: Muestra la documentación de un paquete o símbolo.
