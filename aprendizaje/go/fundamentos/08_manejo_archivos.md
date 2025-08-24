# Manejo de Archivos en Go

## Introducción

El manejo de archivos es una parte fundamental de la programación, ya que permite almacenar y recuperar datos de forma persistente. Go proporciona un conjunto de paquetes en la biblioteca estándar que facilitan la lectura, escritura y manipulación de archivos. En esta sección, aprenderemos cómo trabajar con archivos en Go, desde operaciones básicas hasta técnicas más avanzadas.

## Operaciones Básicas con Archivos

### Abrir y Cerrar Archivos

Para trabajar con un archivo en Go, primero debes abrirlo usando la función `os.Open()` o `os.OpenFile()`. Estas funciones devuelven un puntero a un objeto `os.File` y un error que debes verificar.

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Abrir un archivo en modo lectura
    archivo, err := os.Open("datos.txt")
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    
    // Cerrar el archivo cuando hayas terminado
    // Es importante usar defer para asegurarse de que el archivo se cierre
    // incluso si ocurre un error más adelante
    defer archivo.Close()
    
    // Realizar operaciones con el archivo
    // ...
    
    fmt.Println("Archivo procesado correctamente")
}
```

Es importante cerrar los archivos después de usarlos para liberar recursos del sistema. La declaración `defer` es especialmente útil para esto, ya que garantiza que el archivo se cerrará cuando la función termine, independientemente de cómo termine (normalmente o con un error).

### Modos de Apertura de Archivos

La función `os.OpenFile()` te permite especificar el modo de apertura y los permisos del archivo:

```go
// os.OpenFile(nombre, flag, perm)
archivo, err := os.OpenFile("datos.txt", os.O_RDWR|os.O_CREATE, 0644)
if err != nil {
    fmt.Println("Error al abrir el archivo:", err)
    return
}
defer archivo.Close()
```

Los flags más comunes son:

| Flag | Descripción |
|------|-------------|
| `os.O_RDONLY` | Abre el archivo en modo solo lectura |
| `os.O_WRONLY` | Abre el archivo en modo solo escritura |
| `os.O_RDWR` | Abre el archivo en modo lectura y escritura |
| `os.O_APPEND` | Añade datos al final del archivo |
| `os.O_CREATE` | Crea el archivo si no existe |
| `os.O_TRUNC` | Trunca el archivo al abrirlo |
| `os.O_EXCL` | Usado con O_CREATE, falla si el archivo ya existe |

Los flags se pueden combinar usando el operador OR bit a bit (`|`).

El tercer parámetro especifica los permisos del archivo (en sistemas tipo Unix). El valor `0644` es común y significa que el propietario puede leer y escribir, mientras que los demás solo pueden leer.

### Funciones de Conveniencia

Go proporciona funciones de conveniencia para casos de uso comunes:

```go
// Crear un archivo (equivalente a os.OpenFile con os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
archivo, err := os.Create("nuevo.txt")
if err != nil {
    fmt.Println("Error al crear el archivo:", err)
    return
}
defer archivo.Close()

// Abrir para añadir (equivalente a os.OpenFile con os.O_WRONLY|os.O_CREATE|os.O_APPEND)
archivo, err = os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
if err != nil {
    fmt.Println("Error al abrir el archivo para añadir:", err)
    return
}
defer archivo.Close()
```

## Lectura de Archivos

Go ofrece varias formas de leer datos de un archivo:

### Leer Todo el Contenido

```go
package main

import (
    "fmt"
    "os"
    "io/ioutil" // En Go 1.16+, usar "io"
)

func main() {
    // Leer todo el archivo de una vez
    // En Go 1.16+: contenido, err := os.ReadFile("datos.txt")
    contenido, err := ioutil.ReadFile("datos.txt")
    if err != nil {
        fmt.Println("Error al leer el archivo:", err)
        return
    }
    
    // contenido es un slice de bytes
    fmt.Println(string(contenido)) // Convertir a string para imprimir
}
```

### Leer por Partes

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    archivo, err := os.Open("datos.txt")
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer archivo.Close()
    
    // Crear un buffer para leer
    buffer := make([]byte, 100) // Leer hasta 100 bytes a la vez
    
    for {
        // Leer en el buffer
        n, err := archivo.Read(buffer)
        if err != nil {
            // io.EOF indica que hemos llegado al final del archivo
            if err.Error() == "EOF" {
                break
            }
            fmt.Println("Error al leer el archivo:", err)
            return
        }
        
        // Procesar los bytes leídos (n es el número de bytes leídos)
        fmt.Print(string(buffer[:n]))
        
        // Si leímos menos bytes de los que caben en el buffer, hemos terminado
        if n < len(buffer) {
            break
        }
    }
}
```

### Leer Línea por Línea

Go no tiene una función integrada para leer línea por línea, pero puedes usar `bufio.Scanner`:

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    archivo, err := os.Open("datos.txt")
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer archivo.Close()
    
    // Crear un scanner
    scanner := bufio.NewScanner(archivo)
    
    // Escanear línea por línea
    for scanner.Scan() {
        linea := scanner.Text() // Obtener la línea como string
        fmt.Println(linea)
    }
    
    // Verificar si hubo errores durante el escaneo
    if err := scanner.Err(); err != nil {
        fmt.Println("Error al leer el archivo:", err)
    }
}
```

## Escritura en Archivos

Go proporciona varias formas de escribir datos en archivos:

### Escribir Bytes

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Crear o truncar el archivo
    archivo, err := os.Create("salida.txt")
    if err != nil {
        fmt.Println("Error al crear el archivo:", err)
        return
    }
    defer archivo.Close()
    
    // Escribir bytes
    datos := []byte("Hola, mundo!\n")
    n, err := archivo.Write(datos)
    if err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }
    
    fmt.Printf("Se escribieron %d bytes\n", n)
    
    // También puedes escribir una cadena directamente
    n, err = archivo.WriteString("Esta es otra línea.\n")
    if err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }
    
    fmt.Printf("Se escribieron %d bytes más\n", n)
}
```

### Escribir Todo el Contenido de una Vez

```go
package main

import (
    "fmt"
    "io/ioutil" // En Go 1.16+, usar "os"
)

func main() {
    // Escribir todo el contenido de una vez
    datos := []byte("Este es el contenido completo del archivo.\nTiene múltiples líneas.\n")
    
    // En Go 1.16+: err := os.WriteFile("salida.txt", datos, 0644)
    err := ioutil.WriteFile("salida.txt", datos, 0644)
    if err != nil {
        fmt.Println("Error al escribir el archivo:", err)
        return
    }
    
    fmt.Println("Archivo escrito correctamente")
}
```

### Usar un Writer Bufferizado

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    archivo, err := os.Create("salida.txt")
    if err != nil {
        fmt.Println("Error al crear el archivo:", err)
        return
    }
    defer archivo.Close()
    
    // Crear un writer bufferizado
    writer := bufio.NewWriter(archivo)
    
    // Escribir en el buffer
    writer.WriteString("Línea 1\n")
    writer.WriteString("Línea 2\n")
    writer.WriteString("Línea 3\n")
    
    // Vaciar el buffer al archivo
    // Es importante llamar a Flush() para asegurarse de que todos los datos
    // se escriban en el archivo
    err = writer.Flush()
    if err != nil {
        fmt.Println("Error al vaciar el buffer:", err)
        return
    }
    
    fmt.Println("Archivo escrito correctamente")
}
```

### Anexar a un Archivo

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    // Abrir el archivo en modo append
    archivo, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer archivo.Close()
    
    // Añadir una línea
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    linea := fmt.Sprintf("[%s] Nuevo registro\n", timestamp)
    
    _, err = archivo.WriteString(linea)
    if err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }
    
    fmt.Println("Registro añadido correctamente")
}
```

## Posicionamiento en Archivos

Puedes controlar la posición actual dentro de un archivo usando los métodos `Seek()` y `ReadAt()`/`WriteAt()`:

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    archivo, err := os.OpenFile("datos.txt", os.O_RDWR, 0644)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer archivo.Close()
    
    // Obtener la posición actual
    posicion, err := archivo.Seek(0, os.SEEK_CUR)
    if err != nil {
        fmt.Println("Error al obtener la posición:", err)
        return
    }
    fmt.Printf("Posición inicial: %d\n", posicion) // Normalmente 0
    
    // Leer algunos datos
    buffer := make([]byte, 10)
    n, err := archivo.Read(buffer)
    if err != nil {
        fmt.Println("Error al leer el archivo:", err)
        return
    }
    fmt.Printf("Datos leídos: %s\n", buffer[:n])
    
    // Nueva posición
    posicion, err = archivo.Seek(0, os.SEEK_CUR)
    if err != nil {
        fmt.Println("Error al obtener la posición:", err)
        return
    }
    fmt.Printf("Nueva posición: %d\n", posicion) // Ahora es 10
    
    // Mover a una posición específica desde el inicio
    _, err = archivo.Seek(0, os.SEEK_SET) // Volver al inicio
    if err != nil {
        fmt.Println("Error al cambiar la posición:", err)
        return
    }
    fmt.Printf("Volvemos a la posición: %d\n", 0)
    
    // Leer de nuevo
    buffer = make([]byte, 5)
    n, err = archivo.Read(buffer)
    if err != nil {
        fmt.Println("Error al leer el archivo:", err)
        return
    }
    fmt.Printf("Datos leídos: %s\n", buffer[:n]) // Primeros 5 caracteres
    
    // Mover relativamente a la posición actual
    _, err = archivo.Seek(5, os.SEEK_CUR) // Avanzar 5 posiciones desde la posición actual
    if err != nil {
        fmt.Println("Error al cambiar la posición:", err)
        return
    }
    
    // Mover relativamente al final del archivo
    _, err = archivo.Seek(0, os.SEEK_END) // Ir al final del archivo
    if err != nil {
        fmt.Println("Error al cambiar la posición:", err)
        return
    }
    
    // Leer o escribir en una posición específica sin cambiar la posición actual
    buffer = make([]byte, 5)
    n, err = archivo.ReadAt(buffer, 15) // Leer 5 bytes a partir de la posición 15
    if err != nil && err.Error() != "EOF" {
        fmt.Println("Error al leer el archivo:", err)
        return
    }
    fmt.Printf("Datos leídos en la posición 15: %s\n", buffer[:n])
    
    // Escribir en una posición específica sin cambiar la posición actual
    _, err = archivo.WriteAt([]byte("NUEVO"), 20) // Escribir "NUEVO" en la posición 20
    if err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }
}
```

El método `Seek(offset, whence)` acepta dos parámetros:
- `offset`: El número de bytes a mover
- `whence`: Punto de referencia
  - `os.SEEK_SET` (0): Desde el inicio del archivo
  - `os.SEEK_CUR` (1): Desde la posición actual
  - `os.SEEK_END` (2): Desde el final del archivo

## Trabajar con Rutas de Archivos

Go proporciona el paquete `path/filepath` para trabajar con rutas de archivos de manera compatible con diferentes sistemas operativos:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    // Unir componentes de ruta
    ruta := filepath.Join("directorio", "subdirectorio", "archivo.txt")
    fmt.Println("Ruta:", ruta)
    
    // Obtener el directorio y el nombre de archivo
    directorio := filepath.Dir(ruta)
    nombre := filepath.Base(ruta)
    extension := filepath.Ext(ruta)
    fmt.Printf("Directorio: %s, Nombre: %s, Extensión: %s\n", directorio, nombre, extension)
    
    // Comprobar si una ruta existe
    _, err := os.Stat(ruta)
    existe := !os.IsNotExist(err)
    fmt.Printf("¿Existe? %t\n", existe)
    
    // Comprobar si es un archivo o un directorio
    info, err := os.Stat(ruta)
    if err == nil {
        esDirectorio := info.IsDir()
        fmt.Printf("¿Es directorio? %t\n", esDirectorio)
    }
    
    // Obtener la ruta absoluta
    rutaAbsoluta, err := filepath.Abs(ruta)
    if err != nil {
        fmt.Println("Error al obtener la ruta absoluta:", err)
    } else {
        fmt.Println("Ruta absoluta:", rutaAbsoluta)
    }
    
    // Limpiar una ruta (eliminar elementos como ".." y ".")
    rutaSucia := filepath.Join("dir", "..", "otro", ".", "archivo.txt")
    rutaLimpia := filepath.Clean(rutaSucia)
    fmt.Printf("Ruta sucia: %s, Ruta limpia: %s\n", rutaSucia, rutaLimpia)
    
    // Obtener la ruta relativa
    rutaRelativa, err := filepath.Rel("/home/usuario", "/home/usuario/documentos/archivo.txt")
    if err != nil {
        fmt.Println("Error al obtener la ruta relativa:", err)
    } else {
        fmt.Println("Ruta relativa:", rutaRelativa) // "documentos/archivo.txt"
    }
}
```

## Operaciones con Directorios

Go proporciona funciones para trabajar con directorios:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    // Obtener el directorio de trabajo actual
    dirActual, err := os.Getwd()
    if err != nil {
        fmt.Println("Error al obtener el directorio actual:", err)
        return
    }
    fmt.Println("Directorio actual:", dirActual)
    
    // Cambiar el directorio de trabajo
    err = os.Chdir("/ruta/a/otro/directorio")
    if err != nil {
        fmt.Println("Error al cambiar de directorio:", err)
        // Continuar con el directorio actual
    }
    
    // Crear un directorio
    err = os.Mkdir("nuevo_directorio", 0755)
    if err != nil {
        fmt.Println("Error al crear el directorio:", err)
    }
    
    // Crear directorios anidados
    err = os.MkdirAll("ruta/a/nuevos/directorios", 0755)
    if err != nil {
        fmt.Println("Error al crear los directorios:", err)
    }
    
    // Listar archivos y directorios
    entradas, err := os.ReadDir(".")
    if err != nil {
        fmt.Println("Error al listar el directorio:", err)
        return
    }
    
    fmt.Println("Contenido del directorio:")
    for _, entrada := range entradas {
        tipo := "Archivo"
        if entrada.IsDir() {
            tipo = "Directorio"
        }
        fmt.Printf("- %s (%s)\n", entrada.Name(), tipo)
    }
    
    // Recorrer un directorio recursivamente
    fmt.Println("\nRecorrido recursivo:")
    err = filepath.Walk(".", func(ruta string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("Error al acceder a %s: %v\n", ruta, err)
            return nil // Continuar con el recorrido
        }
        
        tipo := "Archivo"
        if info.IsDir() {
            tipo = "Directorio"
        }
        fmt.Printf("- %s (%s)\n", ruta, tipo)
        return nil
    })
    
    if err != nil {
        fmt.Println("Error al recorrer el directorio:", err)
    }
    
    // Eliminar un archivo
    err = os.Remove("archivo_a_eliminar.txt")
    if err != nil {
        fmt.Println("Error al eliminar el archivo:", err)
    }
    
    // Eliminar un directorio vacío
    err = os.Remove("directorio_vacio")
    if err != nil {
        fmt.Println("Error al eliminar el directorio:", err)
    }
    
    // Eliminar un directorio y su contenido
    err = os.RemoveAll("directorio_con_contenido")
    if err != nil {
        fmt.Println("Error al eliminar el directorio y su contenido:", err)
    }
    
    // Renombrar o mover un archivo o directorio
    err = os.Rename("viejo_nombre.txt", "nuevo_nombre.txt")
    if err != nil {
        fmt.Println("Error al renombrar el archivo:", err)
    }
}
```

## Formatos de Archivo Comunes

Go proporciona paquetes para trabajar con formatos de archivo comunes:

### CSV (Valores Separados por Comas)

```go
package main

import (
    "encoding/csv"
    "fmt"
    "os"
)

func main() {
    // Escribir un archivo CSV
    archivo, err := os.Create("datos.csv")
    if err != nil {
        fmt.Println("Error al crear el archivo:", err)
        return
    }
    defer archivo.Close()
    
    escritor := csv.NewWriter(archivo)
    defer escritor.Flush() // Asegurarse de que todos los datos se escriban
    
    // Escribir encabezados
    err = escritor.Write([]string{"Nombre", "Edad", "Ciudad"})
    if err != nil {
        fmt.Println("Error al escribir en el CSV:", err)
        return
    }
    
    // Escribir datos
    registros := [][]string{
        {"Ana", "25", "Madrid"},
        {"Carlos", "30", "Barcelona"},
        {"Elena", "28", "Valencia"},
    }
    
    for _, registro := range registros {
        err := escritor.Write(registro)
        if err != nil {
            fmt.Println("Error al escribir en el CSV:", err)
            return
        }
    }
    
    fmt.Println("Archivo CSV escrito correctamente")
    
    // Leer un archivo CSV
    archivo, err = os.Open("datos.csv")
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer archivo.Close()
    
    lector := csv.NewReader(archivo)
    
    // Leer encabezados
    encabezados, err := lector.Read()
    if err != nil {
        fmt.Println("Error al leer el CSV:", err)
        return
    }
    fmt.Println("Encabezados:", encabezados)
    
    // Leer todos los registros
    registros, err = lector.ReadAll()
    if err != nil {
        fmt.Println("Error al leer el CSV:", err)
        return
    }
    
    fmt.Println("Registros:")
    for i, registro := range registros {
        fmt.Printf("  %d: %v\n", i+1, registro)
    }
    
    // Alternativamente, leer registro por registro
    archivo, err = os.Open("datos.csv")
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer archivo.Close()
    
    lector = csv.NewReader(archivo)
    
    // Saltar encabezados
    _, err = lector.Read()
    if err != nil {
        fmt.Println("Error al leer el CSV:", err)
        return
    }
    
    fmt.Println("\nLeyendo registro por registro:")
    for {
        registro, err := lector.Read()
        if err != nil {
            if err.Error() == "EOF" {
                break // Fin del archivo
            }
            fmt.Println("Error al leer el CSV:", err)
            return
        }
        
        fmt.Printf("  Nombre: %s, Edad: %s, Ciudad: %s\n", registro[0], registro[1], registro[2])
    }
}
```

### JSON (JavaScript Object Notation)

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

// Definir una estructura que coincida con el formato JSON
type Persona struct {
    Nombre     string   `json:"nombre"`
    Edad       int      `json:"edad"`
    Ciudad     string   `json:"ciudad"`
    Intereses  []string `json:"intereses,omitempty"`
    Activo     bool     `json:"activo"`
    Altura     float64  `json:"altura,omitempty"`
}

func main() {
    // Crear datos
    personas := []Persona{
        {
            Nombre:    "Ana",
            Edad:      25,
            Ciudad:    "Madrid",
            Intereses: []string{"programación", "música", "viajes"},
            Activo:    true,
            Altura:    1.65,
        },
        {
            Nombre:    "Carlos",
            Edad:      30,
            Ciudad:    "Barcelona",
            Intereses: []string{"deportes", "cine"},
            Activo:    false,
        },
    }
    
    // Convertir a JSON con formato
    datosJSON, err := json.MarshalIndent(personas, "", "    ")
    if err != nil {
        fmt.Println("Error al convertir a JSON:", err)
        return
    }
    
    // Imprimir JSON
    fmt.Println(string(datosJSON))
    
    // Escribir JSON a un archivo
    err = os.WriteFile("personas.json", datosJSON, 0644)
    if err != nil {
        fmt.Println("Error al escribir el archivo JSON:", err)
        return
    }
    
    fmt.Println("Archivo JSON escrito correctamente")
    
    // Leer JSON desde un archivo
    datosLeidos, err := os.ReadFile("personas.json")
    if err != nil {
        fmt.Println("Error al leer el archivo JSON:", err)
        return
    }
    
    // Convertir JSON a estructura
    var personasLeidas []Persona
    err = json.Unmarshal(datosLeidos, &personasLeidas)
    if err != nil {
        fmt.Println("Error al convertir JSON a estructura:", err)
        return
    }
    
    fmt.Println("\nPersonas leídas del archivo JSON:")
    for i, persona := range personasLeidas {
        fmt.Printf("  Persona %d: %s, %d años, %s\n", i+1, persona.Nombre, persona.Edad, persona.Ciudad)
        fmt.Printf("    Intereses: %v\n", persona.Intereses)
        fmt.Printf("    Activo: %t, Altura: %.2f\n", persona.Activo, persona.Altura)
    }
    
    // Trabajar con JSON dinámico (sin estructura predefinida)
    jsonString := `{"nombre":"Elena","datos":{"profesion":"ingeniera","experiencia":5}}`
    var resultado map[string]interface{}
    
    err = json.Unmarshal([]byte(jsonString), &resultado)
    if err != nil {
        fmt.Println("Error al convertir JSON a mapa:", err)
        return
    }
    
    fmt.Println("\nJSON dinámico:")
    fmt.Printf("  Nombre: %s\n", resultado["nombre"])
    
    // Acceder a datos anidados
    if datos, ok := resultado["datos"].(map[string]interface{}); ok {
        fmt.Printf("  Profesión: %s\n", datos["profesion"])
        fmt.Printf("  Experiencia: %.0f años\n", datos["experiencia"])
    }
}
```

### Gob (Serialización de Objetos Go)

```go
package main

import (
    "encoding/gob"
    "fmt"
    "os"
)

// Definir una estructura para serializar
type Persona struct {
    Nombre    string
    Edad      int
    Ciudad    string
    Intereses []string
}

func main() {
    // Crear datos
    personas := []Persona{
        {
            Nombre:    "Ana",
            Edad:      25,
            Ciudad:    "Madrid",
            Intereses: []string{"programación", "música", "viajes"},
        },
        {
            Nombre:    "Carlos",
            Edad:      30,
            Ciudad:    "Barcelona",
            Intereses: []string{"deportes", "cine"},
        },
    }
    
    // Serializar (guardar) objetos
    archivo, err := os.Create("personas.gob")
    if err != nil {
        fmt.Println("Error al crear el archivo:", err)
        return
    }
    defer archivo.Close()
    
    codificador := gob.NewEncoder(archivo)
    err = codificador.Encode(personas)
    if err != nil {
        fmt.Println("Error al codificar:", err)
        return
    }
    
    fmt.Println("Datos serializados correctamente")
    
    // Deserializar (cargar) objetos
    archivo, err = os.Open("personas.gob")
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer archivo.Close()
    
    var personasLeidas []Persona
    decodificador := gob.NewDecoder(archivo)
    err = decodificador.Decode(&personasLeidas)
    if err != nil {
        fmt.Println("Error al decodificar:", err)
        return
    }
    
    fmt.Println("\nPersonas deserializadas:")
    for i, persona := range personasLeidas {
        fmt.Printf("  Persona %d: %s, %d años, %s\n", i+1, persona.Nombre, persona.Edad, persona.Ciudad)
        fmt.Printf("    Intereses: %v\n", persona.Intereses)
    }
}
```

## Ejemplo Práctico: Analizador de Logs

Vamos a crear un programa que analice un archivo de log y genere un informe:

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "sort"
    "strconv"
    "strings"
    "time"
)

type LogEntry struct {
    IP        string
    Timestamp time.Time
    Method    string
    URL       string
    Status    int
    Size      int
}

func main() {
    // Abrir el archivo de log
    rutaLog := "access.log"
    archivo, err := os.Open(rutaLog)
    if err != nil {
        fmt.Printf("Error al abrir el archivo %s: %v\n", rutaLog, err)
        return
    }
    defer archivo.Close()
    
    // Patrones para extraer información
    ipRegex := regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
    timestampRegex := regexp.MustCompile(`\[(\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4})\]`)
    requestRegex := regexp.MustCompile(`"(\w+) ([^\s]+) HTTP/[\d.]+"`)
    statusSizeRegex := regexp.MustCompile(`" (\d{3}) (\d+)`)
    
    // Contadores
    ips := make(map[string]int)
    metodos := make(map[string]int)
    urls := make(map[string]int)
    codigos := make(map[int]int)
    horas := make(map[int]int)
    
    var entradas []LogEntry
    
    // Leer el archivo línea por línea
    scanner := bufio.NewScanner(archivo)
    for scanner.Scan() {
        linea := scanner.Text()
        
        // Extraer información
        ipMatch := ipRegex.FindStringSubmatch(linea)
        timestampMatch := timestampRegex.FindStringSubmatch(linea)
        requestMatch := requestRegex.FindStringSubmatch(linea)
        statusSizeMatch := statusSizeRegex.FindStringSubmatch(linea)
        
        if ipMatch != nil && timestampMatch != nil && requestMatch != nil && statusSizeMatch != nil {
            ip := ipMatch[1]
            ips[ip]++
            
            // Parsear timestamp
            timestampStr := timestampMatch[1]
            timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", timestampStr)
            if err == nil {
                horas[timestamp.Hour()]++
            }
            
            metodo := requestMatch[1]
            metodos[metodo]++
            
            url := requestMatch[2]
            urls[url]++
            
            status, _ := strconv.Atoi(statusSizeMatch[1])
            codigos[status]++
            
            size, _ := strconv.Atoi(statusSizeMatch[2])
            
            // Crear entrada de log
            entrada := LogEntry{
                IP:        ip,
                Timestamp: timestamp,
                Method:    metodo,
                URL:       url,
                Status:    status,
                Size:      size,
            }
            
            entradas = append(entradas, entrada)
        }
    }
    
    if err := scanner.Err(); err != nil {
        fmt.Printf("Error al leer el archivo: %v\n", err)
        return
    }
    
    // Generar informe
    fmt.Println("=== INFORME DE ANÁLISIS DE LOG ===")
    fmt.Printf("\nTotal de líneas procesadas: %d\n", len(entradas))
    
    // IPs más frecuentes
    fmt.Println("\nIPs más frecuentes:")
    ipsSorted := sortMapByValue(ips)
    for i, kv := range ipsSorted {
        if i >= 5 {
            break
        }
        fmt.Printf("  %s: %d\n", kv.Key, kv.Value)
    }
    
    // Métodos HTTP
    fmt.Println("\nMétodos HTTP:")
    for metodo, count := range metodos {
        fmt.Printf("  %s: %d\n", metodo, count)
    }
    
    // Códigos de respuesta
    fmt.Println("\nCódigos de respuesta:")
    for codigo, count := range codigos {
        fmt.Printf("  %d: %d\n", codigo, count)
    }
    
    // URLs más solicitadas
    fmt.Println("\nURLs más solicitadas:")
    urlsSorted := sortMapByValue(urls)
    for i, kv := range urlsSorted {
        if i >= 5 {
            break
        }
        fmt.Printf("  %s: %d\n", kv.Key, kv.Value)
    }
    
    // Distribución por hora
    fmt.Println("\nDistribución por hora:")
    for hora := 0; hora < 24; hora++ {
        fmt.Printf("  %02d:00 - %02d:00: %d\n", hora, (hora+1)%24, horas[hora])
    }
    
    // Guardar informe en un archivo
    informeArchivo, err := os.Create("informe_log.txt")
    if err != nil {
        fmt.Printf("Error al crear el archivo de informe: %v\n", err)
        return
    }
    defer informeArchivo.Close()
    
    writer := bufio.NewWriter(informeArchivo)
    defer writer.Flush()
    
    writer.WriteString("=== INFORME DE ANÁLISIS DE LOG ===\n")
    writer.WriteString(fmt.Sprintf("\nTotal de líneas procesadas: %d\n", len(entradas)))
    
    // Escribir el resto del informe...
    
    fmt.Println("\nInforme guardado en 'informe_log.txt'")
}

// KeyValue para ordenar mapas
type KeyValue struct {
    Key   string
    Value int
}

// Ordenar mapa por valor (de mayor a menor)
func sortMapByValue(m map[string]int) []KeyValue {
    var ss []KeyValue
    for k, v := range m {
        ss = append(ss, KeyValue{k, v})
    }
    
    sort.Slice(ss, func(i, j int) bool {
        return ss[i].Value > ss[j].Value
    })
    
    return ss
}
```

## Buenas Prácticas

1. **Usar `defer` para cerrar archivos**: Siempre usa `defer archivo.Close()` inmediatamente después de abrir un archivo para asegurarte de que se cierre correctamente, incluso si ocurren errores.

2. **Verificar errores**: Go utiliza valores de retorno explícitos para los errores. Siempre verifica los errores devueltos por las funciones de manejo de archivos.

3. **Usar buffers para operaciones frecuentes**: Para operaciones frecuentes de lectura o escritura, usa `bufio.Reader`, `bufio.Writer` o `bufio.Scanner` para mejorar el rendimiento.

4. **Usar `filepath` en lugar de concatenar rutas manualmente**: El paquete `path/filepath` garantiza que las rutas sean compatibles con el sistema operativo actual.

5. **Manejar archivos grandes por partes**: Para archivos grandes, lee o escribe por partes en lugar de cargar todo el archivo en memoria.

6. **Usar los paquetes adecuados**: Utiliza paquetes especializados como `encoding/csv`, `encoding/json` o `encoding/gob` para formatos específicos.

7. **Verificar permisos y existencia**: Usa `os.Stat()` para verificar si un archivo existe y sus permisos antes de intentar operaciones que puedan fallar.

8. **Usar nombres de archivo seguros**: Evita caracteres especiales o espacios en nombres de archivo para garantizar la compatibilidad entre sistemas.

9. **Hacer copias de seguridad**: Antes de modificar archivos importantes, haz una copia de seguridad.

10. **Usar `ioutil` o funciones de conveniencia**: Para operaciones simples, considera usar funciones de conveniencia como `ioutil.ReadFile()` y `ioutil.WriteFile()` (o `os.ReadFile()` y `os.WriteFile()` en Go 1.16+).

## Recursos Adicionales

- [Documentación del paquete os](https://golang.org/pkg/os/)
- [Documentación del paquete io](https://golang.org/pkg/io/)
- [Documentación del paquete bufio](https://golang.org/pkg/bufio/)
- [Documentación del paquete path/filepath](https://golang.org/pkg/path/filepath/)
- [Documentación del paquete encoding/csv](https://golang.org/pkg/encoding/csv/)
- [Documentación del paquete encoding/json](https://golang.org/pkg/encoding/json/)
- [Documentación del paquete encoding/gob](https://golang.org/pkg/encoding/gob/)
- [Effective Go - Manejo de errores](https://golang.org/doc/effective_go#errors)
- [Go by Example - Reading Files](https://gobyexample.com/reading-files)
- [Go by Example - Writing Files](https://gobyexample.com/writing-files)
- [Go by Example - JSON](https://gobyexample.com/json)

---

En la siguiente sección, exploraremos el manejo de errores en Go, que es un aspecto fundamental del lenguaje y sigue un enfoque diferente al de muchos otros lenguajes de programación.