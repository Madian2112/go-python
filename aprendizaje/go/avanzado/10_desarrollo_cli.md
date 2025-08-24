# Desarrollo de Aplicaciones CLI en Go

## Introducción

Las aplicaciones de línea de comandos (CLI) son herramientas fundamentales para desarrolladores, administradores de sistemas y usuarios avanzados. Go es un lenguaje excelente para crear aplicaciones CLI debido a su compilación estática, rendimiento, facilidad de distribución y rica biblioteca estándar.

En este documento, exploraremos cómo diseñar y desarrollar aplicaciones CLI robustas, intuitivas y potentes utilizando Go, desde conceptos básicos hasta patrones avanzados y mejores prácticas.

## Fundamentos de Aplicaciones CLI

### Estructura Básica

Una aplicación CLI simple en Go puede comenzar así:

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    // Definir flags
    verbose := flag.Bool("verbose", false, "Habilitar salida detallada")
    name := flag.String("name", "mundo", "Nombre para saludar")
    
    // Parsear flags
    flag.Parse()
    
    // Usar los valores
    if *verbose {
        fmt.Printf("Ejecutando con verbose=%v y name=%s\n", *verbose, *name)
    }
    
    fmt.Printf("¡Hola, %s!\n", *name)
}
```

Para ejecutar esta aplicación:

```bash
# Compilar
go build -o saludar

# Ejecutar con diferentes opciones
./saludar
./saludar -name "Gopher"
./saludar -verbose -name "Desarrollador"
```

### Manejo de Argumentos y Flags

#### Usando el paquete flag

El paquete `flag` de la biblioteca estándar proporciona funcionalidad básica para procesar flags de línea de comandos:

```go
package main

import (
    "flag"
    "fmt"
    "strings"
)

func main() {
    // Flags básicos
    strFlag := flag.String("string", "valor por defecto", "Descripción del flag")
    intFlag := flag.Int("int", 42, "Un número entero")
    boolFlag := flag.Bool("bool", false, "Un valor booleano")
    
    // Flag con nombre corto y largo
    var verbose bool
    flag.BoolVar(&verbose, "verbose", false, "Modo detallado")
    flag.BoolVar(&verbose, "v", false, "Modo detallado (shorthand)")
    
    // Flag personalizado
    var languages []string
    flag.Func("lang", "Lenguajes de programación (separados por comas)", func(s string) error {
        languages = strings.Split(s, ",")
        return nil
    })
    
    // Parsear flags
    flag.Parse()
    
    // Acceder a argumentos posicionales (después de los flags)
    args := flag.Args()
    
    // Mostrar valores
    fmt.Printf("String: %s\n", *strFlag)
    fmt.Printf("Int: %d\n", *intFlag)
    fmt.Printf("Bool: %v\n", *boolFlag)
    fmt.Printf("Verbose: %v\n", verbose)
    fmt.Printf("Languages: %v\n", languages)
    fmt.Printf("Args: %v\n", args)
}
```

#### Usando bibliotecas de terceros

Para aplicaciones CLI más complejas, bibliotecas como `github.com/spf13/cobra` ofrecen funcionalidades avanzadas:

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
)

func main() {
    var verbose bool
    var source string
    
    // Comando raíz
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "Una aplicación CLI de ejemplo",
        Long:  "Una aplicación CLI de ejemplo creada con Cobra en Go",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("Ejecutando con verbose=%v y source=%s\n", verbose, source)
            if len(args) > 0 {
                fmt.Printf("Argumentos adicionales: %v\n", args)
            }
        },
    }
    
    // Flags globales
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Salida detallada")
    rootCmd.PersistentFlags().StringVar(&source, "source", "local", "Fuente de datos")
    
    // Subcomando
    var recursive bool
    var destination string
    
    copyCmd := &cobra.Command{
        Use:   "copy [flags] source destination",
        Short: "Copia archivos",
        Args:  cobra.ExactArgs(2),
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("Copiando de %s a %s (recursive=%v)\n", args[0], args[1], recursive)
        },
    }
    
    // Flags específicos del subcomando
    copyCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Copiar recursivamente")
    
    // Añadir subcomando al comando raíz
    rootCmd.AddCommand(copyCmd)
    
    // Ejecutar
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

### Entrada y Salida

#### Lectura de Entrada Estándar

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    // Leer una línea
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Ingresa tu nombre: ")
    nombre, _ := reader.ReadString('\n')
    nombre = strings.TrimSpace(nombre)
    
    fmt.Printf("Hola, %s!\n", nombre)
    
    // Leer múltiples líneas hasta EOF (Ctrl+D en Unix, Ctrl+Z en Windows)
    fmt.Println("Ingresa texto (Ctrl+D para finalizar):")
    scanner := bufio.NewScanner(os.Stdin)
    var lines []string
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    
    fmt.Printf("Ingresaste %d líneas:\n", len(lines))
    for i, line := range lines {
        fmt.Printf("%d: %s\n", i+1, line)
    }
}
```

#### Manejo de Salida Formateada

```go
package main

import (
    "fmt"
    "os"
    "text/tabwriter"
)

func main() {
    // Salida básica
    fmt.Println("Salida estándar")
    fmt.Fprintln(os.Stderr, "Salida de error")
    
    // Salida formateada con tabwriter
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "Nombre\tEdad\tCiudad")
    fmt.Fprintln(w, "----\t---\t------")
    fmt.Fprintln(w, "Alice\t28\tNueva York")
    fmt.Fprintln(w, "Bob\t35\tSan Francisco")
    fmt.Fprintln(w, "Charlie\t45\tLondres")
    w.Flush()
}
```

## Diseño de Interfaces CLI Efectivas

### Principios de Diseño

1. **Consistencia**: Mantén una estructura de comandos y flags coherente.
2. **Simplicidad**: Haz que las tareas comunes sean simples y las complejas sean posibles.
3. **Progresividad**: Permite a los usuarios descubrir gradualmente la funcionalidad.
4. **Retroalimentación**: Proporciona información clara sobre lo que está sucediendo.
5. **Tolerancia**: Maneja errores de manera elegante y proporciona sugerencias útiles.

### Ayuda y Documentación

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "myapp [command]",
        Short: "Una aplicación CLI de ejemplo",
        Long: `Una aplicación CLI de ejemplo que demuestra cómo crear
una interfaz de línea de comandos efectiva con documentación
detallada y ayuda integrada.`,
        Example: `  myapp serve --port 8080
  myapp config --set key=value
  myapp help [command]`,
    }
    
    // Comando con ejemplos y ayuda detallada
    serveCmd := &cobra.Command{
        Use:   "serve [flags]",
        Short: "Inicia el servidor",
        Long: `Inicia el servidor HTTP en el puerto especificado.
Por defecto, el servidor escucha en todas las interfaces
y utiliza el puerto 8080.`,
        Example: `  myapp serve
  myapp serve --port 3000
  myapp serve --address 127.0.0.1 --port 8000`,
        Run: func(cmd *cobra.Command, args []string) {
            // Implementación
        },
    }
    
    var port int
    var address string
    
    serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Puerto para escuchar")
    serveCmd.Flags().StringVarP(&address, "address", "a", "0.0.0.0", "Dirección para escuchar")
    
    rootCmd.AddCommand(serveCmd)
    
    // Añadir comando de ayuda personalizado
    rootCmd.SetHelpCommand(&cobra.Command{
        Use:   "ayuda [comando]",
        Short: "Ayuda sobre cualquier comando",
        Long: `Ayuda proporciona información detallada sobre cualquier comando
en la aplicación. Simplemente ejecuta myapp ayuda [comando] para
obtener detalles completos.`,
        Run: func(cmd *cobra.Command, args []string) {
            cmd.Root().Help()
        },
    })
    
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

### Feedback Visual

```go
package main

import (
    "fmt"
    "os"
    "time"
    
    "github.com/briandowns/spinner"
    "github.com/fatih/color"
)

func main() {
    // Colores para diferentes tipos de mensajes
    success := color.New(color.FgGreen).PrintlnFunc()
    warning := color.New(color.FgYellow).PrintlnFunc()
    error := color.New(color.FgRed, color.Bold).PrintlnFunc()
    info := color.New(color.FgCyan).PrintfFunc()
    
    // Mensajes con colores
    success("✓ Operación completada con éxito")
    warning("⚠ Advertencia: espacio en disco bajo")
    error("✗ Error: no se pudo conectar al servidor")
    info("ℹ Procesando archivo: %s\n", "datos.json")
    
    // Spinner para operaciones largas
    s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
    s.Prefix = "Descargando archivos "
    s.Suffix = " por favor espere..."
    s.Start()
    
    // Simular operación larga
    time.Sleep(3 * time.Second)
    
    s.Stop()
    success("✓ Descarga completada")
    
    // Barra de progreso simple
    fmt.Println("Procesando:")
    for i := 0; i <= 100; i += 10 {
        drawProgressBar(i)
        time.Sleep(200 * time.Millisecond)
    }
    fmt.Println("\nProcesamiento completado")
}

func drawProgressBar(percent int) {
    width := 40
    completed := width * percent / 100
    remaining := width - completed
    
    fmt.Printf("\r[%s%s] %d%%", 
        strings.Repeat("=", completed),
        strings.Repeat(" ", remaining),
        percent)
    os.Stdout.Sync()
}
```

## Patrones Avanzados para CLI

### Estructura de Proyectos CLI

Una estructura recomendada para proyectos CLI más grandes:

```
myapp/
├── cmd/
│   ├── myapp/
│   │   └── main.go       # Punto de entrada principal
│   └── myappctl/
│       └── main.go       # Herramienta de control (opcional)
├── internal/
│   ├── cli/
│   │   ├── root.go       # Comando raíz
│   │   ├── serve.go      # Comando serve
│   │   └── version.go    # Comando version
│   ├── config/
│   │   └── config.go     # Manejo de configuración
│   └── handler/
│       └── handler.go    # Lógica de negocio
├── pkg/
│   ├── api/              # Paquetes públicos reutilizables
│   └── utils/
├── go.mod
└── go.sum
```

Ejemplo de implementación:

```go
// cmd/myapp/main.go
package main

import (
    "os"
    
    "myapp/internal/cli"
)

func main() {
    if err := cli.Execute(); err != nil {
        os.Exit(1)
    }
}

// internal/cli/root.go
package cli

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "MyApp - Una aplicación CLI",
}

func Execute() error {
    return rootCmd.Execute()
}

// internal/cli/serve.go
package cli

import (
    "github.com/spf13/cobra"
    
    "myapp/internal/handler"
)

func init() {
    serveCmd := &cobra.Command{
        Use:   "serve",
        Short: "Inicia el servidor",
        RunE: func(cmd *cobra.Command, args []string) error {
            return handler.Serve()
        },
    }
    
    rootCmd.AddCommand(serveCmd)
}
```

### Configuración Avanzada

#### Usando archivos de configuración

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

type Config struct {
    Server struct {
        Port    int
        Host    string
        LogFile string
    }
    Database struct {
        Host     string
        Port     int
        User     string
        Password string
        Name     string
    }
    Features struct {
        EnableCache bool
        MaxWorkers  int
    }
}

var cfgFile string
var config Config

func main() {
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "Aplicación con configuración avanzada",
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            // Inicializar configuración
            return initConfig()
        },
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("Configuración cargada: %+v\n", config)
        },
    }
    
    // Flag para archivo de configuración
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "archivo de configuración (por defecto es $HOME/.myapp.yaml)")
    
    // Flags que pueden sobrescribir la configuración
    rootCmd.PersistentFlags().Int("port", 8080, "puerto del servidor")
    rootCmd.PersistentFlags().String("host", "0.0.0.0", "host del servidor")
    
    // Vincular flags con viper
    viper.BindPFlag("server.port", rootCmd.PersistentFlags().Lookup("port"))
    viper.BindPFlag("server.host", rootCmd.PersistentFlags().Lookup("host"))
    
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func initConfig() error {
    if cfgFile != "" {
        // Usar archivo de configuración especificado
        viper.SetConfigFile(cfgFile)
    } else {
        // Buscar en ubicaciones predeterminadas
        home, err := os.UserHomeDir()
        if err != nil {
            return err
        }
        
        viper.AddConfigPath(home)
        viper.AddConfigPath(".")
        viper.SetConfigName(".myapp")
    }
    
    // Leer variables de entorno
    viper.SetEnvPrefix("MYAPP")
    viper.AutomaticEnv()
    
    // Leer configuración
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return err
        }
        // Archivo no encontrado, usar valores por defecto
    } else {
        fmt.Println("Usando archivo de configuración:", viper.ConfigFileUsed())
    }
    
    // Establecer valores por defecto
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.host", "0.0.0.0")
    viper.SetDefault("server.logfile", "server.log")
    viper.SetDefault("features.enablecache", true)
    viper.SetDefault("features.maxworkers", 5)
    
    // Cargar configuración en la estructura
    if err := viper.Unmarshal(&config); err != nil {
        return err
    }
    
    return nil
}
```

Ejemplo de archivo de configuración (`.myapp.yaml`):

```yaml
server:
  port: 9000
  host: 127.0.0.1
  logfile: /var/log/myapp.log

database:
  host: localhost
  port: 5432
  user: postgres
  password: secret
  name: myapp

features:
  enablecache: true
  maxworkers: 10
```

### Interactividad

#### Prompts y Menús Interactivos

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/AlecAivazis/survey/v2"
)

func main() {
    // Preguntas simples
    name := ""
    prompt := &survey.Input{
        Message: "¿Cuál es tu nombre?",
    }
    survey.AskOne(prompt, &name)
    
    // Selección múltiple
    languages := []string{}
    multiSelect := &survey.MultiSelect{
        Message: "¿Qué lenguajes de programación conoces?",
        Options: []string{"Go", "Python", "JavaScript", "Rust", "Java", "C++"},
    }
    survey.AskOne(multiSelect, &languages)
    
    // Confirmación
    confirm := false
    confirmPrompt := &survey.Confirm{
        Message: "¿Deseas continuar?",
    }
    survey.AskOne(confirmPrompt, &confirm)
    
    if !confirm {
        fmt.Println("Operación cancelada")
        os.Exit(0)
    }
    
    // Contraseña
    password := ""
    passwordPrompt := &survey.Password{
        Message: "Ingresa tu contraseña:",
    }
    survey.AskOne(passwordPrompt, &password)
    
    // Mostrar resultados
    fmt.Printf("Nombre: %s\n", name)
    fmt.Printf("Lenguajes: %v\n", languages)
    fmt.Printf("Contraseña: %s\n", password)
    
    // Encuesta completa
    answers := struct {
        Name     string
        Role     string
        Database string
        Features []string
        Comments string
    }{}
    
    questions := []*survey.Question{
        {
            Name:     "name",
            Prompt:   &survey.Input{Message: "Nombre del proyecto:"},
            Validate: survey.Required,
        },
        {
            Name: "role",
            Prompt: &survey.Select{
                Message: "Rol principal:",
                Options: []string{"Backend", "Frontend", "Full Stack", "DevOps"},
                Default: "Backend",
            },
        },
        {
            Name: "database",
            Prompt: &survey.Select{
                Message: "Base de datos:",
                Options: []string{"PostgreSQL", "MySQL", "MongoDB", "SQLite"},
            },
        },
        {
            Name: "features",
            Prompt: &survey.MultiSelect{
                Message: "Características:",
                Options: []string{"API REST", "GraphQL", "WebSockets", "Autenticación", "Docker"},
            },
        },
        {
            Name: "comments",
            Prompt: &survey.Editor{
                Message: "Comentarios adicionales:",
            },
        },
    }
    
    if err := survey.Ask(questions, &answers); err != nil {
        fmt.Println(err.Error())
        return
    }
    
    fmt.Printf("\nConfigurando proyecto %s...\n", answers.Name)
    // Usar las respuestas para configurar el proyecto
}
```

### Manejo de Señales del Sistema

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Crear un contexto cancelable
    ctx, cancel := context.WithCancel(context.Background())
    
    // Canal para señales
    signalCh := make(chan os.Signal, 1)
    signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
    
    // Goroutine para manejar señales
    go func() {
        sig := <-signalCh
        fmt.Printf("\nRecibida señal %s, iniciando apagado graceful...\n", sig)
        cancel() // Cancelar el contexto
        
        // Si recibimos una segunda señal, salir inmediatamente
        sig = <-signalCh
        fmt.Printf("\nRecibida segunda señal %s, forzando salida\n", sig)
        os.Exit(1)
    }()
    
    // Simular una tarea larga
    fmt.Println("Iniciando tarea larga (presiona Ctrl+C para cancelar)...")
    if err := runLongTask(ctx); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Tarea completada con éxito")
}

func runLongTask(ctx context.Context) error {
    // Simular una tarea que toma tiempo
    for i := 1; i <= 10; i++ {
        // Verificar si el contexto ha sido cancelado
        select {
        case <-ctx.Done():
            fmt.Println("Tarea cancelada, limpiando recursos...")
            time.Sleep(500 * time.Millisecond) // Simular limpieza
            return ctx.Err()
        default:
            // Continuar con la tarea
        }
        
        fmt.Printf("Procesando paso %d/10...\n", i)
        time.Sleep(1 * time.Second)
    }
    
    return nil
}
```

## Testing de Aplicaciones CLI

### Testing de Comandos

```go
// cli/root_test.go
package cli

import (
    "bytes"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
    // Capturar salida
    output := &bytes.Buffer{}
    rootCmd.SetOut(output)
    rootCmd.SetErr(output)
    
    // Ejecutar comando sin argumentos
    rootCmd.SetArgs([]string{})
    err := rootCmd.Execute()
    
    // Verificar que no hay error
    assert.NoError(t, err)
    
    // Verificar que la salida contiene el mensaje esperado
    assert.Contains(t, output.String(), "MyApp - Una aplicación CLI")
}

func TestVersionCommand(t *testing.T) {
    // Capturar salida
    output := &bytes.Buffer{}
    rootCmd.SetOut(output)
    
    // Ejecutar comando version
    rootCmd.SetArgs([]string{"version"})
    err := rootCmd.Execute()
    
    // Verificar que no hay error
    assert.NoError(t, err)
    
    // Verificar que la salida contiene la versión
    assert.Contains(t, output.String(), "v0.1.0")
}
```

### Testing de Interacción con el Usuario

```go
package cli

import (
    "bytes"
    "io"
    "strings"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

// Función a probar
func askForConfirmation(r io.Reader, w io.Writer) bool {
    var response string
    
    fmt.Fprint(w, "¿Estás seguro? [s/N]: ")
    _, err := fmt.Fscanln(r, &response)
    if err != nil {
        return false
    }
    
    response = strings.ToLower(strings.TrimSpace(response))
    return response == "s" || response == "si" || response == "y" || response == "yes"
}

func TestAskForConfirmation(t *testing.T) {
    testCases := []struct {
        input    string
        expected bool
    }{
        {"s\n", true},
        {"si\n", true},
        {"y\n", true},
        {"yes\n", true},
        {"n\n", false},
        {"no\n", false},
        {"\n", false},      // Enter sin texto
        {"invalid\n", false}, // Respuesta inválida
    }
    
    for _, tc := range testCases {
        // Preparar entrada y salida
        input := strings.NewReader(tc.input)
        output := &bytes.Buffer{}
        
        // Ejecutar función
        result := askForConfirmation(input, output)
        
        // Verificar resultado
        assert.Equal(t, tc.expected, result)
        assert.Contains(t, output.String(), "¿Estás seguro?")
    }
}
```

## Distribución y Empaquetado

### Cross-Compilación

Go facilita la compilación para diferentes plataformas:

```bash
# Compilar para Windows desde Linux/Mac
GOOS=windows GOARCH=amd64 go build -o myapp.exe

# Compilar para macOS desde Linux/Windows
GOOS=darwin GOARCH=amd64 go build -o myapp-mac

# Compilar para Linux desde Windows/Mac
GOOS=linux GOARCH=amd64 go build -o myapp-linux

# Compilar para ARM (ej. Raspberry Pi)
GOOS=linux GOARCH=arm go build -o myapp-arm
```

### Empaquetado con GoReleaser

GoReleaser automatiza la creación de releases para múltiples plataformas.

Ejemplo de configuración (`.goreleaser.yml`):

```yaml
before:
  hooks:
    - go mod tidy

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}

archives:
  - replacements:
      darwin: macOS
      linux: Linux
      windows: Windows
      amd64: x86_64
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{ .Tag }}-next"

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^ci:'
      - Merge pull request
      - Merge branch
```

Para crear un release:

```bash
# Instalar GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Crear un release de prueba (sin publicar)
goreleaser release --snapshot --rm-dist

# Crear y publicar un release
export GITHUB_TOKEN=your_github_token
goreleaser release
```

## Mejores Prácticas

### Diseño de Comandos

1. **Estructura jerárquica**: Organiza los comandos en una estructura lógica.
2. **Nombres claros**: Usa nombres de comandos y flags descriptivos y consistentes.
3. **Documentación**: Proporciona ayuda detallada para cada comando y flag.
4. **Valores por defecto**: Establece valores por defecto sensatos para minimizar la configuración.
5. **Validación**: Valida la entrada del usuario y proporciona mensajes de error claros.

### Manejo de Errores

```go
package main

import (
    "errors"
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
)

// Errores personalizados
var (
    ErrInvalidInput = errors.New("entrada inválida")
    ErrNotFound     = errors.New("recurso no encontrado")
    ErrPermission   = errors.New("permiso denegado")
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "Ejemplo de manejo de errores",
    }
    
    // Comando con manejo de errores
    processCmd := &cobra.Command{
        Use:   "process [file]",
        Short: "Procesa un archivo",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            filename := args[0]
            return processFile(filename)
        },
    }
    
    rootCmd.AddCommand(processCmd)
    
    // Personalizar manejo de errores
    rootCmd.SilenceErrors = true
    rootCmd.SilenceUsage = true
    
    if err := rootCmd.Execute(); err != nil {
        handleError(err)
        os.Exit(1)
    }
}

func processFile(filename string) error {
    // Verificar si el archivo existe
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return fmt.Errorf("%w: %s", ErrNotFound, filename)
    }
    
    // Verificar permisos
    file, err := os.Open(filename)
    if err != nil {
        if os.IsPermission(err) {
            return fmt.Errorf("%w: %s", ErrPermission, filename)
        }
        return err
    }
    defer file.Close()
    
    // Procesar archivo
    fmt.Printf("Procesando %s...\n", filename)
    
    return nil
}

func handleError(err error) {
    // Determinar tipo de error y mostrar mensaje apropiado
    switch {
    case errors.Is(err, ErrInvalidInput):
        fmt.Fprintln(os.Stderr, "Error: Entrada inválida. Por favor verifica los parámetros.")
    case errors.Is(err, ErrNotFound):
        fmt.Fprintln(os.Stderr, "Error: Recurso no encontrado. Verifica que el archivo o recurso exista.")
    case errors.Is(err, ErrPermission):
        fmt.Fprintln(os.Stderr, "Error: Permiso denegado. Verifica los permisos del archivo o recurso.")
    default:
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    }
}
```

### Rendimiento y Eficiencia

1. **Carga perezosa**: Carga recursos solo cuando son necesarios.
2. **Paralelismo**: Usa goroutines para operaciones concurrentes.
3. **Buffers**: Usa buffers para operaciones de E/S eficientes.
4. **Reutilización**: Reutiliza objetos para reducir la presión del GC.
5. **Perfilado**: Usa herramientas de perfilado para identificar cuellos de botella.

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "runtime/pprof"
    "sync"
)

func main() {
    // Perfilado de CPU (opcional)
    if os.Getenv("PROFILE") == "true" {
        f, _ := os.Create("cpu.prof")
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    
    // Procesar archivos en paralelo
    files := []string{"file1.txt", "file2.txt", "file3.txt"}
    results := processFilesParallel(files)
    
    // Mostrar resultados
    for file, count := range results {
        fmt.Printf("%s: %d líneas\n", file, count)
    }
}

func processFilesParallel(files []string) map[string]int {
    var wg sync.WaitGroup
    var mu sync.Mutex
    results := make(map[string]int)
    
    for _, file := range files {
        wg.Add(1)
        go func(filename string) {
            defer wg.Done()
            
            count, err := countLines(filename)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error procesando %s: %v\n", filename, err)
                return
            }
            
            mu.Lock()
            results[filename] = count
            mu.Unlock()
        }(file)
    }
    
    wg.Wait()
    return results
}

func countLines(filename string) (int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return 0, err
    }
    defer file.Close()
    
    // Usar buffer para lectura eficiente
    reader := bufio.NewReader(file)
    count := 0
    
    for {
        _, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            return count, err
        }
        count++
    }
    
    return count, nil
}
```

## Ejercicios Prácticos

### Ejercicio 1: Gestor de Tareas CLI

Crea una aplicación CLI para gestionar tareas con las siguientes funcionalidades:

- Añadir tareas con título, descripción y fecha límite
- Listar tareas con filtros (completadas, pendientes, por fecha)
- Marcar tareas como completadas
- Eliminar tareas
- Exportar/importar tareas en formato JSON

### Ejercicio 2: Herramienta de Monitoreo de Sistema

Desarrolla una herramienta CLI que muestre información del sistema:

- Uso de CPU y memoria
- Espacio en disco
- Procesos en ejecución
- Conexiones de red activas
- Historial de uso (con gráficos en terminal)

### Ejercicio 3: Cliente API REST

Implementa un cliente CLI para interactuar con una API REST:

- Autenticación con tokens
- CRUD de recursos
- Filtrado y paginación
- Formato de salida configurable (JSON, tabla, CSV)
- Caché local de respuestas

## Conclusiones

Go es un lenguaje excelente para desarrollar aplicaciones CLI debido a su compilación estática, rendimiento y facilidad de distribución. Al seguir los principios y patrones descritos en este documento, puedes crear herramientas CLI potentes, intuitivas y mantenibles.

Recuerda estos puntos clave:

1. **Diseño centrado en el usuario**: Crea interfaces intuitivas y consistentes.
2. **Documentación clara**: Proporciona ayuda detallada y ejemplos.
3. **Manejo robusto de errores**: Ofrece mensajes claros y sugerencias útiles.
4. **Configuración flexible**: Permite múltiples fuentes de configuración.
5. **Pruebas exhaustivas**: Asegura que tu aplicación funcione correctamente en diferentes escenarios.

Con estas prácticas, tus aplicaciones CLI en Go serán herramientas valiosas para desarrolladores y usuarios finales.

## Referencias

1. Documentación oficial de Go: https://golang.org/doc/
2. Cobra - Biblioteca para aplicaciones CLI: https://github.com/spf13/cobra
3. Viper - Gestión de configuración: https://github.com/spf13/viper
4. Survey - Prompts interactivos: https://github.com/AlecAivazis/survey
5. GoReleaser - Automatización de releases: https://goreleaser.com/
6. "Command-Line Tools in Go" - Fatih Arslan
7. "Building Command-Line Applications in Go" - Alex Ellis
8. "12 Factor CLI Apps" - Jeff Dickey
9. "Terminal User Interfaces with Go" - Charm.sh