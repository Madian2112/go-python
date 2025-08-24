# Desarrollo de Aplicaciones Móviles con Go

## Introducción

Go es un lenguaje que tradicionalmente se ha utilizado para desarrollo backend, pero existen diversas opciones para utilizarlo en el desarrollo de aplicaciones móviles. En este documento, exploraremos las diferentes alternativas para crear aplicaciones móviles con Go, desde enfoques nativos hasta soluciones híbridas.

## Opciones para Desarrollo Móvil con Go

### 1. Gomobile

Gomobile es una herramienta oficial del proyecto Go que permite utilizar código Go en aplicaciones móviles de dos maneras principales:

#### Bind: Creación de Bibliotecas Nativas

Permite crear bibliotecas nativas que pueden ser utilizadas desde código Java/Kotlin (Android) o Objective-C/Swift (iOS).

```bash
# Instalación de gomobile
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest

# Inicialización
gomobile init

# Generar bindings para Android e iOS
gomobile bind -target=android -o=mylib.aar github.com/miusuario/mipaquete
gomobile bind -target=ios -o=MyLib.xcframework github.com/miusuario/mipaquete
```

Ejemplo de paquete Go que puede ser utilizado en aplicaciones móviles:

```go
// github.com/miusuario/mipaquete/calc/calc.go
package calc

// Suma dos números y devuelve el resultado
func Suma(a, b int) int {
	return a + b
}

// Resta dos números y devuelve el resultado
func Resta(a, b int) int {
	return a - b
}

// Multiplica dos números y devuelve el resultado
func Multiplica(a, b int) int {
	return a * b
}

// Divide dos números y devuelve el resultado
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("división por cero")
	}
	return a / b, nil
}
```

Uso en Android (Kotlin):

```kotlin
// Importar la biblioteca generada
import calc.Calc

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        
        // Usar funciones de Go
        val resultado = Calc.suma(5, 3)
        Log.d("GoLib", "Resultado: $resultado")
        
        try {
            val division = Calc.divide(10, 2)
            Log.d("GoLib", "División: $division")
        } catch (e: Exception) {
            Log.e("GoLib", "Error: ${e.message}")
        }
    }
}
```

Uso en iOS (Swift):

```swift
// Importar la biblioteca generada
import MyLib

class ViewController: UIViewController {
    override func viewDidLoad() {
        super.viewDidLoad()
        
        // Usar funciones de Go
        let resultado = CalcSuma(5, 3)
        print("Resultado: \(resultado)")
        
        do {
            let division = try CalcDivide(10, 2)
            print("División: \(division)")
        } catch {
            print("Error: \(error)")
        }
    }
}
```

#### Build: Aplicaciones Completas

Permite crear aplicaciones móviles completas utilizando Go y OpenGL para la interfaz gráfica.

```bash
# Crear una aplicación para Android
gomobile build -target=android github.com/miusuario/miapp

# Crear una aplicación para iOS
gomobile build -target=ios github.com/miusuario/miapp
```

Ejemplo de una aplicación básica con gomobile:

```go
// github.com/miusuario/miapp/main.go
package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
				case lifecycle.CrossOff:
					glctx = nil
				}
			case size.Event:
				sz = e
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}
				glctx.ClearColor(1, 0, 0, 1) // Fondo rojo
				glctx.Clear(gl.COLOR_BUFFER_BIT)
				a.Publish()
			case touch.Event:
				// Manejar eventos táctiles
			}
		}
	})
}
```

### 2. Fyne Mobile

Fyne, el toolkit de GUI para Go, también soporta aplicaciones móviles:

```go
// github.com/miusuario/miapp/main.go
package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Mi App Móvil")
	
	hello := widget.NewLabel("¡Hola desde Go!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Saludar", func() {
			hello.SetText("¡Hola, mundo!")
		}),
	))
	
	w.ShowAndRun()
}
```

Para compilar para móviles:

```bash
# Instalar herramientas
go install fyne.io/fyne/v2/cmd/fyne@latest

# Empaquetar para Android
fyne package -os android -appID com.miempresa.miapp

# Empaquetar para iOS (requiere macOS)
fyne package -os ios -appID com.miempresa.miapp
```

### 3. Gioui

Gio es un toolkit de UI inmediato para Go que soporta múltiples plataformas, incluyendo Android e iOS:

```go
// github.com/miusuario/miapp/main.go
package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := app.NewWindow()
		th := material.NewTheme()
		var ops op.Ops
		var button widget.Clickable
		
		for e := range w.Events() {
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							return material.H3(th, "¡Hola desde Gio!").Layout(gtx)
						},
					),
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							return material.Button(th, &button, "Presionar").Layout(gtx)
						},
					),
				)
				
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
```

Para compilar para móviles:

```bash
# Para Android
gio -target android -appid com.miempresa.miapp -o miapp.apk .

# Para iOS (requiere macOS)
gio -target ios -appid com.miempresa.miapp .
```

### 4. Enfoques Híbridos

#### Flutter con Go

Utilizar Go para la lógica de negocio y Flutter para la UI:

```go
// github.com/miusuario/miapp/backend/api.go
package backend

import "C"

//export GetData
func GetData() *C.char {
	// Lógica de negocio en Go
	data := map[string]interface{}{
		"nombre": "Juan",
		"edad":   30,
	}
	
	jsonData, _ := json.Marshal(data)
	return C.CString(string(jsonData))
}

//export FreeString
func FreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func main() {}
```

En Flutter (Dart):

```dart
import 'dart:ffi';
import 'dart:convert';

// Definir las funciones de FFI
typedef GetDataFunc = Pointer<Utf8> Function();
typedef FreeStringFunc = Void Function(Pointer<Utf8>);

class GoBackend {
  late final DynamicLibrary _lib;
  late final GetDataFunc _getData;
  late final FreeStringFunc _freeString;
  
  GoBackend() {
    // Cargar la biblioteca
    _lib = DynamicLibrary.open('libbackend.so'); // o .dylib en iOS
    
    // Obtener referencias a las funciones
    _getData = _lib.lookupFunction<GetDataFunc, GetDataFunc>('GetData');
    _freeString = _lib.lookupFunction<FreeStringFunc, FreeStringFunc>('FreeString');
  }
  
  Map<String, dynamic> getData() {
    final ptr = _getData();
    final jsonStr = ptr.toDartString();
    _freeString(ptr);
    return jsonDecode(jsonStr);
  }
}
```

#### React Native con Go

Utilizar Go para la lógica de negocio y React Native para la UI:

```go
// Similar al ejemplo de Flutter, pero con bindings para React Native
```

En React Native (JavaScript):

```javascript
import { NativeModules } from 'react-native';

const { GoBackend } = NativeModules;

// Usar las funciones de Go
GoBackend.getData()
  .then(data => console.log(data))
  .catch(error => console.error(error));
```

## Consideraciones para el Desarrollo Móvil con Go

### Ventajas

1. **Rendimiento**: Go ofrece un rendimiento cercano al nativo, especialmente para operaciones intensivas.
2. **Código Compartido**: Permite compartir código entre backend y aplicaciones móviles.
3. **Concurrencia**: El modelo de concurrencia de Go es útil para operaciones en segundo plano.
4. **Tamaño de Binarios**: Los binarios de Go son relativamente pequeños y eficientes.

### Desafíos

1. **Ecosistema UI**: El ecosistema de UI para móviles en Go es menos maduro que alternativas como React Native o Flutter.
2. **Curva de Aprendizaje**: Combinar Go con frameworks móviles puede requerir conocimiento de múltiples tecnologías.
3. **Integración con APIs Nativas**: Puede ser más complejo acceder a ciertas APIs específicas de la plataforma.
4. **Soporte de la Comunidad**: Menos recursos y ejemplos disponibles comparado con soluciones más populares.

## Patrones de Arquitectura para Aplicaciones Móviles con Go

### Patrón de Biblioteca Nativa

Utilizar Go para implementar funcionalidades específicas como bibliotecas nativas:

```
+-------------------+
|  UI Nativa        |
| (Swift/Kotlin)    |
+--------+----------+
         |
         v
+-------------------+
|  Biblioteca Go    |
| (Lógica compleja) |
+-------------------+
```

Ejemplo de arquitectura:

```go
// github.com/miusuario/milibreria/api.go
package api

// Estructura exportada
type Usuario struct {
	ID       int
	Nombre   string
	Apellido string
	Email    string
}

// Servicio de usuarios
type UsuarioService struct {
	// campos privados
	db *sql.DB
}

// Constructor exportado
func NewUsuarioService() *UsuarioService {
	// Inicialización
	return &UsuarioService{
		// inicializar campos
	}
}

// Método exportado
func (s *UsuarioService) ObtenerUsuario(id int) (Usuario, error) {
	// Implementación
	return Usuario{ID: id, Nombre: "Juan", Apellido: "Pérez", Email: "juan@ejemplo.com"}, nil
}

// Método exportado
func (s *UsuarioService) GuardarUsuario(u Usuario) error {
	// Implementación
	return nil
}
```

### Patrón de Aplicación Completa

Implementar toda la aplicación en Go, utilizando bibliotecas de UI como Fyne o Gio:

```
+-------------------+
|  UI Go (Fyne/Gio) |
+--------+----------+
         |
         v
+-------------------+
|  Lógica Go        |
+--------+----------+
         |
         v
+-------------------+
|  Persistencia Go  |
+-------------------+
```

Ejemplo de arquitectura con Fyne:

```go
// github.com/miusuario/miapp/main.go
package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Modelo
type TaskModel struct {
	tasks []string
}

func (m *TaskModel) AddTask(task string) {
	m.tasks = append(m.tasks, task)
}

func (m *TaskModel) GetTasks() []string {
	return m.tasks
}

// Controlador
type TaskController struct {
	model *TaskModel
	view  *TaskView
}

func NewTaskController(model *TaskModel) *TaskController {
	return &TaskController{model: model}
}

func (c *TaskController) SetView(view *TaskView) {
	c.view = view
}

func (c *TaskController) AddTask(task string) {
	c.model.AddTask(task)
	c.view.UpdateTaskList()
}

// Vista
type TaskView struct {
	controller *TaskController
	model      *TaskModel
	taskList   *widget.List
	input      *widget.Entry
}

func NewTaskView(controller *TaskController, model *TaskModel) *TaskView {
	return &TaskView{
		controller: controller,
		model:      model,
	}
}

func (v *TaskView) BuildUI(window fyne.Window) {
	v.input = widget.NewEntry()
	v.input.SetPlaceHolder("Nueva tarea...")
	
	addButton := widget.NewButton("Añadir", func() {
		if v.input.Text != "" {
			v.controller.AddTask(v.input.Text)
			v.input.SetText("")
		}
	})
	
	v.taskList = widget.NewList(
		func() int { return len(v.model.GetTasks()) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(v.model.GetTasks()[id])
		},
	)
	
	content := container.NewBorder(
		container.NewBorder(nil, nil, nil, addButton, v.input),
		nil, nil, nil,
		v.taskList,
	)
	
	window.SetContent(content)
}

func (v *TaskView) UpdateTaskList() {
	v.taskList.Refresh()
}

func main() {
	// Inicializar MVC
	model := &TaskModel{}
	controller := NewTaskController(model)
	view := NewTaskView(controller, model)
	controller.SetView(view)
	
	// Inicializar aplicación
	a := app.New()
	w := a.NewWindow("Lista de Tareas")
	view.BuildUI(w)
	w.Resize(fyne.NewSize(300, 400))
	w.ShowAndRun()
}
```

### Patrón de Backend Go con UI Nativa

Utilizar Go como un servicio backend local que se comunica con la UI nativa:

```
+-------------------+
|  UI Nativa        |
| (Swift/Kotlin)    |
+--------+----------+
         |
         v
+-------------------+
|  API Local        |
| (HTTP/gRPC)       |
+--------+----------+
         |
         v
+-------------------+
|  Servicio Go      |
+-------------------+
```

Ejemplo de servicio local en Go:

```go
// github.com/miusuario/miapp/backend/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

type TaskService struct {
	tasks  []Task
	nextID int
	mutex  sync.Mutex
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}

func (s *TaskService) GetTasks() []Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.tasks
}

func (s *TaskService) AddTask(title string) Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	task := Task{
		ID:     s.nextID,
		Title:  title,
		Status: false,
	}
	
	s.tasks = append(s.tasks, task)
	s.nextID++
	
	return task
}

func (s *TaskService) UpdateTask(id int, status bool) (Task, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Status = status
			return s.tasks[i], true
		}
	}
	
	return Task{}, false
}

func main() {
	service := NewTaskService()
	
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tasks := service.GetTasks()
			json.NewEncoder(w).Encode(tasks)
			
		case http.MethodPost:
			var req struct {
				Title string `json:"title"`
			}
			
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			
			task := service.AddTask(req.Title)
			json.NewEncoder(w).Encode(task)
			
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		
		// Extraer ID de la URL
		id := 0 // Extraer de la URL
		
		var req struct {
			Status bool `json:"status"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		task, found := service.UpdateTask(id, req.Status)
		if !found {
			http.Error(w, "Tarea no encontrada", http.StatusNotFound)
			return
		}
		
		json.NewEncoder(w).Encode(task)
	})
	
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
```

## Ejercicios Prácticos

### Ejercicio 1: Biblioteca de Cálculos Matemáticos

Crea una biblioteca en Go que implemente funciones matemáticas avanzadas y utilízala en una aplicación Android o iOS.

1. Implementa funciones para cálculos estadísticos (media, mediana, desviación estándar).
2. Añade funciones para operaciones matriciales.
3. Utiliza gomobile para generar los bindings.
4. Crea una aplicación móvil simple que utilice estas funciones.

### Ejercicio 2: Aplicación de Lista de Tareas

Desarrolla una aplicación de lista de tareas completa utilizando Fyne:

1. Implementa la creación, edición y eliminación de tareas.
2. Añade categorías y prioridades a las tareas.
3. Implementa persistencia local con SQLite.
4. Empaqueta la aplicación para Android.

### Ejercicio 3: Cliente de API REST

Crea una aplicación móvil que consuma una API REST:

1. Implementa la lógica de comunicación con la API en Go.
2. Crea una interfaz de usuario con Gio.
3. Implementa caché local para funcionamiento offline.
4. Añade autenticación y manejo de sesiones.

## Conclusiones

Go ofrece diversas opciones para el desarrollo de aplicaciones móviles, desde la creación de bibliotecas nativas hasta aplicaciones completas. Aunque el ecosistema de UI móvil en Go no es tan maduro como otras alternativas, proporciona ventajas significativas en términos de rendimiento y capacidad de compartir código entre plataformas.

La elección del enfoque dependerá de los requisitos específicos del proyecto, el conocimiento del equipo y las restricciones de tiempo y recursos. Para funcionalidades específicas que requieren alto rendimiento, utilizar Go como una biblioteca nativa puede ser la mejor opción, mientras que para aplicaciones completas, los enfoques híbridos o las bibliotecas como Fyne y Gio ofrecen soluciones viables.

## Referencias

1. Gomobile - Herramientas para móviles en Go: https://pkg.go.dev/golang.org/x/mobile
2. Fyne - Toolkit GUI multiplataforma: https://fyne.io/
3. Gio - Toolkit UI inmediato: https://gioui.org/
4. "Mobile Programming with Go" - Packt Publishing
5. "Cross-Platform GUI Programming with Fyne" - Andrew Williams
6. "Building Mobile Applications with Go" - O'Reilly Media
7. Documentación oficial de Gomobile: https://github.com/golang/go/wiki/Mobile
8. Ejemplos de Gomobile: https://github.com/golang/mobile/tree/master/example
9. Fyne Mobile: https://developer.fyne.io/started/mobile
10. Gio Mobile: https://gioui.org/doc/mobile