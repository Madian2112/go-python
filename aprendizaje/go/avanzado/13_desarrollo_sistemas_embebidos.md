# Desarrollo de Sistemas Embebidos con Go

## Introducción

Go es un lenguaje que, gracias a su eficiencia, simplicidad y capacidad para compilación cruzada, resulta cada vez más atractivo para el desarrollo de sistemas embebidos y aplicaciones IoT (Internet de las Cosas). En este documento, exploraremos cómo utilizar Go en entornos embebidos, desde dispositivos de bajo nivel hasta plataformas IoT más potentes.

## Ventajas de Go para Sistemas Embebidos

1. **Compilación Cruzada**: Go permite compilar fácilmente para diferentes arquitecturas desde un único sistema de desarrollo.
2. **Binarios Estáticos**: Los programas Go pueden compilarse como binarios estáticos sin dependencias externas.
3. **Recolección de Basura**: Aunque puede ser un desafío en sistemas con recursos limitados, la recolección de basura simplifica la gestión de memoria.
4. **Concurrencia**: El modelo de goroutines y canales es ideal para manejar múltiples sensores y operaciones de E/S.
5. **Tipado Estático**: Ayuda a prevenir errores en tiempo de compilación, crucial para sistemas embebidos.

## Plataformas y Dispositivos Compatibles

Go puede ejecutarse en diversas plataformas embebidas, aunque con diferentes niveles de soporte:

### Dispositivos ARM

- Raspberry Pi (todas las versiones)
- BeagleBone
- NVIDIA Jetson
- Muchas placas basadas en ARM Cortex-A

### Dispositivos MIPS

- Routers y dispositivos de red con OpenWrt
- Algunas placas de desarrollo MIPS

### Dispositivos x86 de Bajo Consumo

- Intel Edison
- Intel Galileo
- PC Engines APU

### Microcontroladores

Para microcontroladores con recursos muy limitados, existen opciones como:

- **TinyGo**: Un subconjunto de Go diseñado para microcontroladores y WebAssembly
- **Emgo**: Un dialecto de Go para programación de microcontroladores

## Configuración del Entorno de Desarrollo

### Compilación Cruzada con Go

Go facilita la compilación cruzada mediante variables de entorno:

```bash
# Compilar para Raspberry Pi (ARM)
GOOS=linux GOARCH=arm GOARM=7 go build -o myapp main.go

# Compilar para BeagleBone (ARM)
GOOS=linux GOARCH=arm GOARM=7 go build -o myapp main.go

# Compilar para MIPS (router con OpenWrt)
GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o myapp main.go

# Compilar para MIPS little-endian
GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -o myapp main.go

# Compilar para x86 de 32 bits
GOOS=linux GOARCH=386 go build -o myapp main.go
```

### Configuración para TinyGo

TinyGo es una alternativa para dispositivos con recursos muy limitados:

```bash
# Instalar TinyGo
go install tinygo.org/x/tools/cmd/tinygo@latest

# Compilar para Arduino Uno
tinygo build -o firmware.hex -target arduino main.go

# Compilar para micro:bit
tinygo build -o firmware.hex -target microbit main.go

# Compilar para ESP32
tinygo build -o firmware.bin -target esp32 main.go

# Flashear directamente a un dispositivo
tinygo flash -target arduino main.go
```

## Acceso a Hardware con Go

### GPIO (Pines de Entrada/Salida de Propósito General)

Para interactuar con GPIO en dispositivos como Raspberry Pi:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

func main() {
	// Inicializar periph
	if _, err := host.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al inicializar periph: %v\n", err)
		os.Exit(1)
	}

	// Obtener un pin, por ejemplo GPIO17 en Raspberry Pi
	pin := gpioreg.ByName("GPIO17")
	if pin == nil {
		fmt.Fprintf(os.Stderr, "Error: pin no encontrado\n")
		os.Exit(1)
	}

	// Configurar como salida
	if err := pin.Out(gpio.Low); err != nil {
		fmt.Fprintf(os.Stderr, "Error al configurar pin: %v\n", err)
		os.Exit(1)
	}

	// Parpadear LED
	for i := 0; i < 10; i++ {
		// Encender
		pin.Out(gpio.High)
		time.Sleep(500 * time.Millisecond)

		// Apagar
		pin.Out(gpio.Low)
		time.Sleep(500 * time.Millisecond)
	}
}
```

### I2C (Comunicación entre Circuitos Integrados)

Para comunicarse con sensores y dispositivos I2C:

```go
package main

import (
	"fmt"
	"os"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func main() {
	// Inicializar periph
	if _, err := host.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al inicializar periph: %v\n", err)
		os.Exit(1)
	}

	// Abrir el bus I2C
	bus, err := i2creg.Open("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al abrir bus I2C: %v\n", err)
		os.Exit(1)
	}
	defer bus.Close()

	// Dirección del dispositivo I2C (ejemplo: sensor BME280)
	dev := i2c.Dev{Bus: bus, Addr: 0x76}

	// Leer el ID del chip (registro 0xD0 para BME280)
	buf := []byte{0xD0}
	if err := dev.Tx(buf, buf); err != nil {
		fmt.Fprintf(os.Stderr, "Error en transacción I2C: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ID del chip: 0x%X\n", buf[0])
}
```

### SPI (Interfaz de Periféricos Serie)

Para comunicarse con dispositivos SPI:

```go
package main

import (
	"fmt"
	"os"

	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

func main() {
	// Inicializar periph
	if _, err := host.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al inicializar periph: %v\n", err)
		os.Exit(1)
	}

	// Abrir puerto SPI
	port, err := spireg.Open("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al abrir puerto SPI: %v\n", err)
		os.Exit(1)
	}
	defer port.Close()

	// Configurar conexión SPI
	conn, err := port.Connect(1*1000*1000, spi.Mode0, 8)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al configurar SPI: %v\n", err)
		os.Exit(1)
	}

	// Ejemplo: Leer datos de un sensor
	tx := []byte{0x8F, 0x00} // Comando de lectura + byte dummy
	rx := make([]byte, len(tx))

	if err := conn.Tx(tx, rx); err != nil {
		fmt.Fprintf(os.Stderr, "Error en transacción SPI: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Datos recibidos: %v\n", rx)
}
```

### UART (Transmisor-Receptor Asíncrono Universal)

Para comunicación serie:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"go.bug.st/serial"
)

func main() {
	// Listar puertos disponibles
	puertos, err := serial.GetPortsList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al obtener puertos: %v\n", err)
		os.Exit(1)
	}

	if len(puertos) == 0 {
		fmt.Println("No se encontraron puertos serie")
		os.Exit(1)
	}

	fmt.Println("Puertos disponibles:")
	for _, puerto := range puertos {
		fmt.Println(puerto)
	}

	// Configurar puerto (ajustar nombre según sistema)
	puertoNombre := "/dev/ttyUSB0" // En Linux
	// puertoNombre := "COM3" // En Windows

	config := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	// Abrir puerto
	puerto, err := serial.Open(puertoNombre, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al abrir puerto: %v\n", err)
		os.Exit(1)
	}
	defer puerto.Close()

	// Escribir datos
	n, err := puerto.Write([]byte("AT\r\n"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al escribir: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Enviados %d bytes\n", n)

	// Leer respuesta
	buf := make([]byte, 100)
	time.Sleep(100 * time.Millisecond) // Esperar respuesta

	n, err = puerto.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al leer: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Recibidos %d bytes: %s\n", n, buf[:n])
}
```

## Programación de Microcontroladores con TinyGo

TinyGo permite utilizar un subconjunto de Go en microcontroladores:

### Ejemplo para Arduino

```go
package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		led.High()
		time.Sleep(time.Millisecond * 500)

		led.Low()
		time.Sleep(time.Millisecond * 500)
	}
}
```

### Ejemplo para ESP32

```go
package main

import (
	"machine"
	"time"
)

const led = machine.GPIO2

func main() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		led.High()
		time.Sleep(time.Millisecond * 500)

		led.Low()
		time.Sleep(time.Millisecond * 500)
	}
}
```

### Lectura de Sensores con TinyGo

```go
package main

import (
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/bme280"
)

func main() {
	// Configurar I2C
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		Frequency: 400000,
	})

	// Configurar sensor BME280
	sensor := bme280.New(i2c)
	sensor.Configure()

	for {
		temp, _ := sensor.ReadTemperature()
		pres, _ := sensor.ReadPressure()
		hum, _ := sensor.ReadHumidity()

		// Convertir a valores legibles
		tempC := float32(temp) / 1000.0
		presPa := float32(pres) / 100.0
		humRel := float32(hum) / 1024.0

		fmt.Printf("Temp: %.2f°C, Presión: %.2fhPa, Humedad: %.2f%%\r\n",
			tempC, presPa, humRel)

		time.Sleep(time.Second)
	}
}
```

## Desarrollo de Aplicaciones IoT con Go

### Arquitectura Típica de una Aplicación IoT

```
+----------------+      +----------------+      +----------------+
|                |      |                |      |                |
|  Dispositivos  |<---->|  Gateway/Edge  |<---->|    Cloud/      |
|  (Sensores)    |      |  (Procesamiento|      |    Backend     |
|                |      |   local)       |      |                |
+----------------+      +----------------+      +----------------+
```

### Ejemplo de Aplicación Edge en Go

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

// Estructura para datos del sensor
type SensorData struct {
	DeviceID  string  `json:"device_id"`
	Timestamp string  `json:"timestamp"`
	Temp      float64 `json:"temperature"`
	Humidity  float64 `json:"humidity"`
	Pressure  float64 `json:"pressure"`
}

// Simulación de lectura de sensor
func readSensor() (float64, float64, float64, error) {
	// En un caso real, aquí leeríamos datos del sensor
	// Para este ejemplo, devolvemos valores simulados
	return 22.5, 45.3, 1013.2, nil
}

func main() {
	// Inicializar periph (para hardware real)
	if _, err := host.Init(); err != nil {
		log.Fatalf("Error al inicializar periph: %v\n", err)
	}

	// Configurar cliente MQTT
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://broker.hivemq.com:1883")
	opts.SetClientID("go-iot-device-1")
	opts.SetUsername("usuario") // Si es necesario
	opts.SetPassword("clave")   // Si es necesario

	// Configurar reconexión
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Minute)

	// Crear cliente MQTT
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error al conectar con broker MQTT: %v\n", token.Error())
	}
	defer client.Disconnect(250)

	// Canal para señales de terminación
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// WaitGroup para esperar a que terminen las goroutines
	var wg sync.WaitGroup

	// Bandera para indicar terminación
	terminar := false

	// Goroutine para leer sensores y publicar datos
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Leer datos del sensor
				temp, humidity, pressure, err := readSensor()
				if err != nil {
					log.Printf("Error al leer sensor: %v\n", err)
					continue
				}

				// Crear estructura de datos
				data := SensorData{
					DeviceID:  "device-001",
					Timestamp: time.Now().Format(time.RFC3339),
					Temp:      temp,
					Humidity:  humidity,
					Pressure:  pressure,
				}

				// Convertir a JSON
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Printf("Error al serializar JSON: %v\n", err)
					continue
				}

				// Publicar en MQTT
				topic := "sensors/environmental"
				token := client.Publish(topic, 0, false, jsonData)
				if token.Wait() && token.Error() != nil {
					log.Printf("Error al publicar: %v\n", token.Error())
				} else {
					log.Printf("Datos publicados: %s\n", string(jsonData))
				}

			case <-sigChan:
				fmt.Println("Señal de terminación recibida")
				terminar = true
				return
			}

			if terminar {
				break
			}
		}
	}()

	// Esperar señal de terminación
	<-sigChan
	terminar = true

	// Esperar a que terminen las goroutines
	wg.Wait()
	fmt.Println("Aplicación terminada correctamente")
}
```

### Comunicación MQTT

MQTT es un protocolo común en aplicaciones IoT:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

// Manejador de mensajes recibidos
var mensajeHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Recibido mensaje en tema: %s\n", msg.Topic())
	fmt.Printf("Mensaje: %s\n", msg.Payload())
}

func main() {
	// Configurar opciones del cliente
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	opts.SetClientID("go-mqtt-client")
	opts.SetDefaultPublishHandler(mensajeHandler)

	// Crear cliente
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Suscribirse a un tema
	if token := client.Subscribe("go-mqtt/ejemplo", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Publicar mensaje
	token := client.Publish("go-mqtt/ejemplo", 0, false, "Hola desde Go")
	token.Wait()

	// Capturar señal para terminar
	signal := make(chan os.Signal, 1)
	signal.Notify(signal, syscall.SIGINT, syscall.SIGTERM)
	<-signal

	// Desconectar
	client.Disconnect(250)
}
```

## Optimización para Sistemas Embebidos

### Reducción del Tamaño del Binario

```bash
# Compilar con optimizaciones
go build -ldflags="-s -w" -o myapp main.go

# Comprimir el binario con UPX
upx --best myapp
```

### Optimización de Memoria

```go
// Reutilizar buffers para evitar asignaciones frecuentes
var buffer = make([]byte, 1024)

func procesarDatos() {
	// Usar buffer existente en lugar de crear uno nuevo
	n, err := dispositivo.Read(buffer)
	if err != nil {
		// Manejar error
		return
	}
	
	// Procesar buffer[:n]
}
```

### Reducción de Uso de Goroutines

```go
// En lugar de crear una goroutine por conexión
func manejarConexiones(listener net.Listener) {
	// Crear un pool de workers
	const numWorkers = 5
	trabajos := make(chan net.Conn, 100)
	
	// Iniciar workers
	for i := 0; i < numWorkers; i++ {
		go worker(trabajos)
	}
	
	// Aceptar conexiones y enviarlas al pool
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Manejar error
			continue
		}
		trabajos <- conn
	}
}

func worker(trabajos <-chan net.Conn) {
	for conn := range trabajos {
		// Manejar conexión
		// ...
		conn.Close()
	}
}
```

## Patrones de Diseño para Sistemas Embebidos

### Patrón Productor-Consumidor

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Estructura de datos del sensor
type SensorData struct {
	Timestamp time.Time
	Value     float64
}

func main() {
	// Canal para datos del sensor
	dataChannel := make(chan SensorData, 100)
	
	// Canal para señales de terminación
	done := make(chan struct{})
	
	// WaitGroup para esperar a que terminen las goroutines
	var wg sync.WaitGroup
	
	// Productor: lee datos del sensor
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				// Simular lectura del sensor
				data := SensorData{
					Timestamp: time.Now(),
					Value:     23.5, // Valor simulado
				}
				
				// Enviar datos al canal
				dataChannel <- data
				
				// Esperar antes de la siguiente lectura
				time.Sleep(1 * time.Second)
			}
		}
	}()
	
	// Consumidor: procesa los datos
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case data := <-dataChannel:
				// Procesar datos
				fmt.Printf("Datos recibidos: %v, Valor: %.2f\n", 
					data.Timestamp.Format(time.RFC3339), data.Value)
				
				// Aquí podríamos almacenar, filtrar o transmitir los datos
			}
		}
	}()
	
	// Capturar señal para terminar
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	
	// Señalizar terminación y esperar
	close(done)
	wg.Wait()
	fmt.Println("Programa terminado correctamente")
}
```

### Patrón de Máquina de Estados

```go
package main

import (
	"fmt"
	"time"
)

// Estados posibles
type Estado int

const (
	Inactivo Estado = iota
	Midiendo
	Transmitiendo
	AhorroBateria
	Error
)

// Eventos posibles
type Evento int

const (
	IniciarMedicion Evento = iota
	MedicionCompleta
	IniciarTransmision
	TransmisionCompleta
	BateriaBaja
	ErrorDetectado
	Reiniciar
)

// Estructura del dispositivo
type Dispositivo struct {
	estadoActual Estado
	bateria     int // Porcentaje de batería
}

// Método para manejar eventos
func (d *Dispositivo) manejarEvento(evento Evento) {
	switch d.estadoActual {
	case Inactivo:
		switch evento {
		case IniciarMedicion:
			fmt.Println("Iniciando medición...")
			d.estadoActual = Midiendo
		case BateriaBaja:
			fmt.Println("Batería baja, entrando en modo ahorro")
			d.estadoActual = AhorroBateria
		}
		
	case Midiendo:
		switch evento {
		case MedicionCompleta:
			fmt.Println("Medición completada")
			d.estadoActual = Inactivo
		case IniciarTransmision:
			fmt.Println("Iniciando transmisión de datos...")
			d.estadoActual = Transmitiendo
		case ErrorDetectado:
			fmt.Println("Error durante la medición")
			d.estadoActual = Error
		}
		
	case Transmitiendo:
		switch evento {
		case TransmisionCompleta:
			fmt.Println("Transmisión completada")
			d.estadoActual = Inactivo
		case ErrorDetectado:
			fmt.Println("Error durante la transmisión")
			d.estadoActual = Error
		case BateriaBaja:
			fmt.Println("Batería baja durante transmisión, completando primero")
			// Continuar en el mismo estado
		}
		
	case AhorroBateria:
		switch evento {
		case Reiniciar:
			fmt.Println("Saliendo del modo ahorro de batería")
			d.estadoActual = Inactivo
		}
		
	case Error:
		switch evento {
		case Reiniciar:
			fmt.Println("Reiniciando después de error")
			d.estadoActual = Inactivo
		}
	}
}

func main() {
	// Crear dispositivo
	dispositivo := &Dispositivo{
		estadoActual: Inactivo,
		bateria:     80,
	}
	
	// Simular secuencia de eventos
	fmt.Printf("Estado inicial: %v\n", dispositivo.estadoActual)
	
	dispositivo.manejarEvento(IniciarMedicion)
	time.Sleep(1 * time.Second)
	
	dispositivo.manejarEvento(MedicionCompleta)
	time.Sleep(1 * time.Second)
	
	dispositivo.manejarEvento(IniciarMedicion)
	time.Sleep(1 * time.Second)
	
	dispositivo.manejarEvento(IniciarTransmision)
	time.Sleep(1 * time.Second)
	
	dispositivo.manejarEvento(ErrorDetectado)
	time.Sleep(1 * time.Second)
	
	dispositivo.manejarEvento(Reiniciar)
	time.Sleep(1 * time.Second)
	
	dispositivo.manejarEvento(BateriaBaja)
	time.Sleep(1 * time.Second)
	
	dispositivo.manejarEvento(Reiniciar)
	
	fmt.Printf("Estado final: %v\n", dispositivo.estadoActual)
}
```

## Ejercicios Prácticos

### Ejercicio 1: Sistema de Monitoreo Ambiental

Desarrolla un sistema que lea datos de temperatura, humedad y presión atmosférica de un sensor (o simúlalos) y los envíe a un servidor MQTT.

1. Implementa la lectura de sensores (o simulación).
2. Configura la conexión MQTT.
3. Establece un formato JSON para los datos.
4. Implementa manejo de errores y reconexión.
5. Añade un modo de bajo consumo cuando no hay actividad.

### Ejercicio 2: Control de Dispositivos con GPIO

Crea un programa que controle dispositivos conectados a los pines GPIO de una Raspberry Pi o similar.

1. Configura varios pines como salidas para controlar LEDs o relés.
2. Configura pines como entradas para leer botones o sensores.
3. Implementa un servidor web simple para controlar los dispositivos remotamente.
4. Añade autenticación básica al servidor web.

### Ejercicio 3: Gateway IoT con TinyGo

Desarrolla un gateway IoT utilizando TinyGo en un microcontrolador compatible.

1. Configura la comunicación con sensores mediante I2C o SPI.
2. Implementa un protocolo de comunicación inalámbrica (WiFi, BLE, etc.).
3. Establece un formato eficiente para transmitir datos.
4. Implementa un sistema de gestión de energía para maximizar la duración de la batería.

## Conclusiones

Go ofrece capacidades significativas para el desarrollo de sistemas embebidos y aplicaciones IoT, desde dispositivos de bajo nivel hasta plataformas más potentes. Sus características de compilación cruzada, binarios estáticos y modelo de concurrencia lo hacen especialmente adecuado para estos entornos.

Para dispositivos con recursos muy limitados, TinyGo proporciona una alternativa viable que mantiene gran parte de la sintaxis y características de Go, adaptadas a las restricciones de los microcontroladores.

Al desarrollar aplicaciones embebidas con Go, es importante considerar las limitaciones de recursos y optimizar el código en consecuencia, prestando especial atención al uso de memoria, goroutines y el tamaño del binario resultante.

## Referencias

1. Periph.io - Biblioteca para acceso a hardware: https://periph.io/
2. TinyGo - Go para microcontroladores: https://tinygo.org/
3. Emgo - Dialecto de Go para microcontroladores: https://github.com/ziutek/emgo
4. Paho MQTT - Cliente MQTT para Go: https://github.com/eclipse/paho.mqtt.golang
5. Go Serial - Biblioteca para comunicación serie: https://github.com/bugst/go-serial
6. "Programming Embedded Systems with Go" - O'Reilly Media
7. "IoT Development with Go" - Packt Publishing
8. "TinyGo - Small Places" - Apress
9. Documentación oficial de Go sobre compilación cruzada: https://golang.org/doc/install/source#environment
10. Ejemplos de TinyGo: https://github.com/tinygo-org/drivers