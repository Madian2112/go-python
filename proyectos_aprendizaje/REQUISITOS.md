# Requisitos para Ejecutar los Proyectos de Consola

Este documento detalla los requisitos necesarios para ejecutar las aplicaciones de consola desarrolladas en Python y Go incluidas en este repositorio.

## Requisitos para Proyectos Python

### Requisitos de Sistema

- **Python**: Versión 3.6 o superior
  - Puedes verificar tu versión con `python --version` o `python3 --version`
  - Descarga: [Python.org](https://www.python.org/downloads/)

### Configuración del Entorno

1. **Instalación de Python**:
   - Windows: Descarga el instalador desde [Python.org](https://www.python.org/downloads/) y sigue las instrucciones. Asegúrate de marcar la opción "Add Python to PATH".
   - macOS: Puedes usar Homebrew con `brew install python` o descargar el instalador desde Python.org.
   - Linux: La mayoría de las distribuciones incluyen Python. Puedes instalarlo con el gestor de paquetes de tu distribución (por ejemplo, `apt install python3` en Ubuntu).

2. **Verificación de la instalación**:
   ```bash
   python --version  # o python3 --version en algunos sistemas
   ```

3. **Entorno Virtual (opcional pero recomendado)**:
   ```bash
   # Crear entorno virtual
   python -m venv venv
   
   # Activar entorno virtual
   # En Windows
   venv\Scripts\activate
   # En macOS/Linux
   source venv/bin/activate
   ```

### Ejecución de Proyectos Python

1. **Navegar al directorio del proyecto**:
   ```bash
   cd ejemplos_consola/python/[nombre_proyecto]
   ```

2. **Ejecutar la aplicación**:
   ```bash
   python [nombre_archivo].py [argumentos]
   ```

   Por ejemplo, para la aplicación de tareas (TODO):
   ```bash
   python todo.py list
   python todo.py add "Nueva tarea"
   ```

## Requisitos para Proyectos Go

### Requisitos de Sistema

- **Go**: Versión 1.13 o superior
  - Puedes verificar tu versión con `go version`
  - Descarga: [golang.org](https://golang.org/dl/)

### Configuración del Entorno

1. **Instalación de Go**:
   - Windows: Descarga el instalador MSI desde [golang.org](https://golang.org/dl/) y sigue las instrucciones.
   - macOS: Puedes usar Homebrew con `brew install go` o descargar el instalador desde golang.org.
   - Linux: Descarga el archivo tar.gz desde golang.org y extráelo en /usr/local.

2. **Verificación de la instalación**:
   ```bash
   go version
   ```

3. **Configuración del GOPATH (opcional)**:
   Go Modules (a partir de Go 1.11) ha reducido la necesidad de configurar manualmente el GOPATH, pero si trabajas con proyectos antiguos, puede ser necesario.

### Ejecución de Proyectos Go

1. **Navegar al directorio del proyecto**:
   ```bash
   cd ejemplos_consola/go/[nombre_proyecto]
   ```

2. **Compilar y ejecutar directamente**:
   ```bash
   go run main.go [argumentos]
   ```

   Por ejemplo, para la aplicación de tareas (TODO):
   ```bash
   go run main.go list
   go run main.go add "Nueva tarea"
   ```

3. **Compilar a un ejecutable** (opcional):
   ```bash
   go build -o [nombre_ejecutable]
   ```

   Y luego ejecutar:
   ```bash
   # En Windows
   [nombre_ejecutable].exe [argumentos]
   # En macOS/Linux
   ./[nombre_ejecutable] [argumentos]
   ```

## Solución de Problemas Comunes

### Python

1. **Comando Python no encontrado**:
   - Asegúrate de que Python esté en tu PATH
   - En algunos sistemas, puede ser necesario usar `python3` en lugar de `python`

2. **Módulos no encontrados**:
   - Verifica que estás en el directorio correcto
   - Si el proyecto tiene un archivo `requirements.txt`, instala las dependencias con `pip install -r requirements.txt`

### Go

1. **Comando Go no encontrado**:
   - Asegúrate de que Go esté en tu PATH
   - Reinicia la terminal después de instalar Go

2. **Errores de compilación**:
   - Verifica que estás usando una versión compatible de Go
   - Asegúrate de que todas las dependencias estén correctamente instaladas con `go mod tidy`

## Recursos Adicionales

- [Documentación oficial de Python](https://docs.python.org/)
- [Documentación oficial de Go](https://golang.org/doc/)
- [Tutorial de línea de comandos en Python](https://realpython.com/command-line-interfaces-python-argparse/)
- [Tutorial de línea de comandos en Go](https://gobyexample.com/command-line-flags)