# Desarrollo de Aplicaciones de Escritorio con Go

## Introducción

Go es conocido principalmente como un lenguaje para desarrollo backend y herramientas de línea de comandos, pero también ofrece capacidades para crear aplicaciones de escritorio multiplataforma. En este documento, exploraremos las diferentes opciones disponibles para desarrollar interfaces gráficas de usuario (GUI) con Go, desde bibliotecas nativas hasta frameworks que utilizan tecnologías web.

Aprenderemos sobre las ventajas y desventajas de cada enfoque, patrones de diseño para aplicaciones de escritorio, y cómo implementar funcionalidades comunes como persistencia de datos, comunicación con servicios externos y distribución de aplicaciones.

## Bibliotecas y Frameworks GUI para Go

### Opciones Nativas

#### Fyne

[Fyne](https://fyne.io/) es un toolkit GUI nativo para Go que utiliza OpenGL para renderizar interfaces de usuario. Es completamente nativo de Go y no depende de CGo, lo que facilita la compilación cruzada.

**Ejemplo básico con Fyne:**

```go
package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Crear una nueva aplicación
	a := app.New()
	
	// Crear una ventana
	w := a.NewWindow("¡Hola, Fyne!")
	
	// Establecer el contenido de la ventana
	w.SetContent(container.NewVBox(
		widget.NewLabel("¡Bienvenido a Fyne!"),
		widget.NewButton("Salir", func() {
			a.Quit()
		}),
	))
	
	// Establecer el tamaño de la ventana
	w.Resize(fyne.NewSize(400, 200))
	
	// Mostrar la ventana y ejecutar la aplicación
	w.ShowAndRun()
}
```

**Ventajas de Fyne:**
- Completamente nativo de Go
- Fácil compilación cruzada
- Apariencia consistente en todas las plataformas
- Soporte para temas claro y oscuro
- Diseño responsivo

**Desventajas:**
- Conjunto limitado de widgets en comparación con frameworks más maduros
- Personalización visual limitada
- Rendimiento puede ser un problema en aplicaciones complejas

#### Gio

[Gio](https://gioui.org/) es un toolkit GUI inmediato para Go que utiliza GPU para renderizar. Está diseñado para ser portable y eficiente.

**Ejemplo básico con Gio:**

```go
package main

import (
	"image/color"
	"log"
	"os"

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
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops
	var button widget.Clickable

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return material.H3(th, "¡Hola, Gio!").Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return material.Button(th, &button, "Salir").Layout(gtx)
					},
				),
			)

			if button.Clicked() {
				return nil
			}

			e.Frame(gtx.Ops)
		}
	}
}
```

**Ventajas de Gio:**
- Modelo de programación inmediato (similar a React)
- Alto rendimiento con aceleración GPU
- Soporte para múltiples plataformas, incluyendo móviles
- No depende de CGo

**Desventajas:**
- Curva de aprendizaje pronunciada
- Documentación limitada
- Ecosistema menos maduro

#### ui

[ui](https://github.com/andlabs/ui) es un binding de Go para la biblioteca [libui](https://github.com/andlabs/libui), que proporciona controles nativos en cada plataforma.

**Ejemplo básico con ui:**

```go
package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

func main() {
	err := ui.Main(func() {
		// Crear una ventana
		window := ui.NewWindow("¡Hola, ui!", 300, 200, true)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		ui.OnShouldQuit(func() bool {
			window.Destroy()
			return true
		})

		// Crear controles
		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel("¡Bienvenido a ui!"), false)
		button := ui.NewButton("Salir")
		button.OnClicked(func(*ui.Button) {
			ui.Quit()
		})
		box.Append(button, false)

		// Establecer el contenido de la ventana
		window.SetChild(box)
		window.SetMargined(true)

		// Mostrar la ventana
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
```

**Ventajas de ui:**
- Utiliza controles nativos de cada plataforma
- Apariencia nativa en cada sistema operativo
- API relativamente simple

**Desventajas:**
- Requiere CGo, lo que complica la compilación cruzada
- Desarrollo menos activo
- Conjunto limitado de widgets

### Enfoques Híbridos

#### Wails

[Wails](https://wails.io/) permite crear aplicaciones de escritorio utilizando Go para el backend y tecnologías web (HTML/CSS/JavaScript) para la interfaz de usuario.

**Estructura básica de un proyecto Wails:**

```go
// main.go
package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("¡Hola %s, es un placer conocerte!", name)
}

// ShowDialog muestra un diálogo nativo
func (a *App) ShowDialog() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Selecciona un archivo",
		Filters: []runtime.FileFilter{
			{DisplayName: "Imágenes", Pattern: "*.png;*.jpg;*.jpeg"},
			{DisplayName: "Documentos", Pattern: "*.pdf;*.doc;*.docx"},
		},
	})
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Archivo seleccionado: " + selection
}
```

En el frontend (JavaScript):

```javascript
// frontend/src/App.jsx
import { useState } from 'react';
import { Greet, ShowDialog } from '../wailsjs/go/main/App';

function App() {
  const [name, setName] = useState('');
  const [result, setResult] = useState('');
  const [file, setFile] = useState('');

  function greet() {
    Greet(name).then((result) => {
      setResult(result);
    });
  }

  function openDialog() {
    ShowDialog().then((result) => {
      setFile(result);
    });
  }

  return (
    <div>
      <h1>¡Hola Wails!</h1>
      <div>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Ingresa tu nombre"
        />
        <button onClick={greet}>Saludar</button>
      </div>
      <div>{result}</div>
      <div>
        <button onClick={openDialog}>Seleccionar archivo</button>
      </div>
      <div>{file}</div>
    </div>
  );
}

export default App;
```

**Ventajas de Wails:**
- Combina la potencia de Go con la flexibilidad de las tecnologías web
- Permite utilizar frameworks web populares como React, Vue o Svelte
- Acceso a APIs nativas a través de Go
- Distribución sencilla de aplicaciones

**Desventajas:**
- Mayor tamaño de la aplicación final
- Requiere conocimientos de desarrollo web
- Rendimiento potencialmente menor que soluciones nativas

#### Webview

[Webview](https://github.com/webview/webview) es una biblioteca minimalista que proporciona un componente webview para mostrar contenido web en una aplicación nativa.

**Ejemplo básico con Webview:**

```go
package main

import (
	"github.com/webview/webview"
)

func main() {
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Ejemplo de Webview")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://golang.org/")
	w.Run()
}
```

**Ejemplo con contenido HTML local y comunicación bidireccional:**

```go
package main

import (
	"fmt"
	"github.com/webview/webview"
)

func main() {
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Aplicación Webview")
	w.SetSize(800, 600, webview.HintNone)

	// Definir una función que puede ser llamada desde JavaScript
	w.Bind("saludar", func(nombre string) string {
		return fmt.Sprintf("¡Hola, %s!", nombre)
	})

	// Cargar HTML directamente
	w.Navigate("data:text/html," + `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Aplicación Webview</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 20px; }
			button { padding: 8px 16px; }
		</style>
	</head>
	<body>
		<h1>Ejemplo de Webview con Go</h1>
		<input type="text" id="nombre" placeholder="Ingresa tu nombre">
		<button onclick="saludarDesdeJS()">Saludar</button>
		<div id="resultado"></div>

		<script>
			async function saludarDesdeJS() {
				const nombre = document.getElementById('nombre').value || 'Invitado';
				const resultado = await window.saludar(nombre);
				document.getElementById('resultado').textContent = resultado;
			}
		</script>
	</body>
	</html>
	`)

	w.Run()
}
```

**Ventajas de Webview:**
- Extremadamente ligero
- Fácil de integrar
- Utiliza el motor de renderizado web del sistema

**Desventajas:**
- Funcionalidad limitada
- Diferencias entre plataformas
- Requiere CGo

#### Lorca

[Lorca](https://github.com/zserge/lorca) utiliza Chrome instalado en el sistema para renderizar la interfaz de usuario.

**Ejemplo básico con Lorca:**

```go
package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/zserge/lorca"
)

func main() {
	// Verificar si Chrome está instalado
	ui, err := lorca.New("", "", 800, 600)
	if err != nil {
		fmt.Println("No se pudo iniciar Chrome:", err)
		return
	}
	defer ui.Close()

	// Cargar HTML
	ui.Load("data:text/html," + url.PathEscape(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Lorca Example</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 20px; }
			button { padding: 8px 16px; }
		</style>
	</head>
	<body>
		<h1>Ejemplo de Lorca con Go</h1>
		<input type="text" id="nombre" placeholder="Ingresa tu nombre">
		<button onclick="saludar()">Saludar</button>
		<div id="resultado"></div>

		<script>
			function saludar() {
				const nombre = document.getElementById('nombre').value || 'Invitado';
				// Llamar a una función de Go
				document.getElementById('resultado').textContent = window.saludarDesdeGo(nombre);
			}
		</script>
	</body>
	</html>
	`))

	// Exponer funciones de Go a JavaScript
	ui.Bind("saludarDesdeGo", func(nombre string) string {
		return fmt.Sprintf("¡Hola, %s! Saludos desde Go.", nombre)
	})

	// Esperar a que se cierre la ventana
	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigc:
	case <-ui.Done():
	}
}
```

**Ventajas de Lorca:**
- No requiere empaquetado de Chromium
- Comunicación bidireccional entre Go y JavaScript
- Soporte para múltiples plataformas

**Desventajas:**
- Requiere Chrome instalado en el sistema
- Tamaño de aplicación más pequeño pero dependencia externa
- No es adecuado para aplicaciones que requieren distribución sin dependencias

## Patrones de Diseño para Aplicaciones de Escritorio

### Arquitectura MVC/MVVM

La separación de responsabilidades es crucial en aplicaciones de escritorio. El patrón Modelo-Vista-Controlador (MVC) o Modelo-Vista-VistaModelo (MVVM) ayuda a organizar el código.

**Ejemplo de MVC con Fyne:**

```go
package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Modelo: Representa los datos y la lógica de negocio
type TaskModel struct {
	tasks []string
}

func NewTaskModel() *TaskModel {
	return &TaskModel{tasks: []string{}}
}

func (m *TaskModel) AddTask(task string) {
	m.tasks = append(m.tasks, task)
}

func (m *TaskModel) GetTasks() []string {
	return m.tasks
}

// Controlador: Maneja la interacción entre el modelo y la vista
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
	if task != "" {
		c.model.AddTask(task)
		c.view.UpdateTaskList()
	}
}

// Vista: Maneja la interfaz de usuario
type TaskView struct {
	controller *TaskController
	model      *TaskModel
	taskList   *widget.List
	taskEntry  *widget.Entry
}

func NewTaskView(controller *TaskController, model *TaskModel) *TaskView {
	return &TaskView{
		controller: controller,
		model:      model,
	}
}

func (v *TaskView) CreateUI(window fyne.Window) {
	// Crear componentes de la UI
	v.taskEntry = widget.NewEntry()
	v.taskEntry.SetPlaceHolder("Nueva tarea...")

	addButton := widget.NewButton("Añadir", func() {
		v.controller.AddTask(v.taskEntry.Text)
		v.taskEntry.SetText("")
	})

	v.taskList = widget.NewList(
		func() int {
			return len(v.model.GetTasks())
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Elemento")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(v.model.GetTasks()[id])
		},
	)

	// Organizar la UI
	content := container.NewBorder(
		container.NewBorder(nil, nil, nil, addButton, v.taskEntry),
		nil, nil, nil,
		v.taskList,
	)

	window.SetContent(content)
}

func (v *TaskView) UpdateTaskList() {
	v.taskList.Refresh()
}

func main() {
	a := app.New()
	w := a.NewWindow("Gestor de Tareas MVC")
	w.Resize(fyne.NewSize(400, 300))

	// Crear componentes MVC
	model := NewTaskModel()
	controller := NewTaskController(model)
	view := NewTaskView(controller, model)

	// Conectar componentes
	controller.SetView(view)

	// Crear la UI
	view.CreateUI(w)

	w.ShowAndRun()
}
```

### Gestión de Estado

La gestión de estado es fundamental en aplicaciones de escritorio. Aquí hay un ejemplo utilizando un enfoque similar a Redux con Wails:

```go
// backend/store/store.go
package store

import (
	"sync"
)

// Estado de la aplicación
type AppState struct {
	Counter int    `json:"counter"`
	Message string `json:"message"`
}

// Acciones
const (
	INCREMENT_COUNTER = "INCREMENT_COUNTER"
	DECREMENT_COUNTER = "DECREMENT_COUNTER"
	SET_MESSAGE      = "SET_MESSAGE"
)

type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Reducer
func Reducer(state AppState, action Action) AppState {
	switch action.Type {
	case INCREMENT_COUNTER:
		state.Counter++
	case DECREMENT_COUNTER:
		state.Counter--
	case SET_MESSAGE:
		if payload, ok := action.Payload.(string); ok {
			state.Message = payload
		}
	}
	return state
}

// Store
type Store struct {
	state     AppState
	listeners []func(AppState)
	mutex     sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		state: AppState{
			Counter: 0,
			Message: "Bienvenido",
		},
		listeners: []func(AppState){},
	}
}

func (s *Store) GetState() AppState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state
}

func (s *Store) Dispatch(action Action) {
	s.mutex.Lock()
	s.state = Reducer(s.state, action)
	s.mutex.Unlock()

	// Notificar a los listeners
	for _, listener := range s.listeners {
		listener(s.GetState())
	}
}

func (s *Store) Subscribe(listener func(AppState)) func() {
	s.mutex.Lock()
	s.listeners = append(s.listeners, listener)
	index := len(s.listeners) - 1
	s.mutex.Unlock()

	// Devolver función para cancelar la suscripción
	return func() {
		s.mutex.Lock()
		s.listeners = append(s.listeners[:index], s.listeners[index+1:]...)
		s.mutex.Unlock()
	}
}
```

En el archivo principal:

```go
// main.go
package main

import (
	"context"
	"encoding/json"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"myapp/store"
)

// App struct
type App struct {
	ctx   context.Context
	store *store.Store
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		store: store.NewStore(),
	}
}

// startup is called when the app starts. The context is saved
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Suscribirse a cambios de estado y notificar al frontend
	a.store.Subscribe(func(state store.AppState) {
		jsonState, _ := json.Marshal(state)
		runtime.EventsEmit(a.ctx, "state-changed", string(jsonState))
	})
}

// GetState devuelve el estado actual
func (a *App) GetState() store.AppState {
	return a.store.GetState()
}

// Dispatch envía una acción al store
func (a *App) Dispatch(action store.Action) {
	a.store.Dispatch(action)
}

// Métodos de conveniencia
func (a *App) IncrementCounter() {
	a.store.Dispatch(store.Action{Type: store.INCREMENT_COUNTER})
}

func (a *App) DecrementCounter() {
	a.store.Dispatch(store.Action{Type: store.DECREMENT_COUNTER})
}

func (a *App) SetMessage(message string) {
	a.store.Dispatch(store.Action{
		Type:    store.SET_MESSAGE,
		Payload: message,
	})
}
```

En el frontend (React):

```jsx
// frontend/src/App.jsx
import { useState, useEffect } from 'react';
import { GetState, IncrementCounter, DecrementCounter, SetMessage } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime';

function App() {
  const [state, setState] = useState({ counter: 0, message: '' });

  useEffect(() => {
    // Obtener estado inicial
    GetState().then(setState);

    // Suscribirse a cambios de estado
    EventsOn('state-changed', (jsonState) => {
      setState(JSON.parse(jsonState));
    });
  }, []);

  return (
    <div>
      <h1>Gestor de Estado</h1>
      <div>
        <h2>Contador: {state.counter}</h2>
        <button onClick={IncrementCounter}>Incrementar</button>
        <button onClick={DecrementCounter}>Decrementar</button>
      </div>
      <div>
        <h2>Mensaje: {state.message}</h2>
        <input
          type="text"
          value={state.message}
          onChange={(e) => SetMessage(e.target.value)}
          placeholder="Escribe un mensaje"
        />
      </div>
    </div>
  );
}

export default App;
```

## Funcionalidades Comunes en Aplicaciones de Escritorio

### Persistencia de Datos

#### SQLite con Go

SQLite es una excelente opción para aplicaciones de escritorio que necesitan almacenar datos localmente.

```go
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type TaskDB struct {
	db *sql.DB
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func NewTaskDB(appDataDir string) (*TaskDB, error) {
	// Asegurar que el directorio existe
	os.MkdirAll(appDataDir, 0755)

	// Abrir la base de datos
	dbPath := filepath.Join(appDataDir, "tasks.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	// Crear la tabla si no existe
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT 0
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("error al crear la tabla: %w", err)
	}

	return &TaskDB{db: db}, nil
}

func (t *TaskDB) Close() error {
	return t.db.Close()
}

func (t *TaskDB) AddTask(task Task) (int64, error) {
	result, err := t.db.Exec(
		"INSERT INTO tasks (title, description, completed) VALUES (?, ?, ?)",
		task.Title, task.Description, task.Completed,
	)
	if err != nil {
		return 0, fmt.Errorf("error al añadir tarea: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener ID: %w", err)
	}

	return id, nil
}

func (t *TaskDB) GetTasks() ([]Task, error) {
	rows, err := t.db.Query("SELECT id, title, description, completed FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("error al consultar tareas: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed)
		if err != nil {
			return nil, fmt.Errorf("error al escanear tarea: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *TaskDB) UpdateTask(task Task) error {
	_, err := t.db.Exec(
		"UPDATE tasks SET title = ?, description = ?, completed = ? WHERE id = ?",
		task.Title, task.Description, task.Completed, task.ID,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar tarea: %w", err)
	}

	return nil
}

func (t *TaskDB) DeleteTask(id int) error {
	_, err := t.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error al eliminar tarea: %w", err)
	}

	return nil
}
```

#### Almacenamiento de Configuración

Para almacenar configuraciones de aplicación, podemos usar formatos como JSON o YAML:

```go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type AppConfig struct {
	Theme       string `json:"theme"`
	Language    string `json:"language"`
	FontSize    int    `json:"fontSize"`
	AutoSave    bool   `json:"autoSave"`
	RecentFiles []string `json:"recentFiles"`
}

func DefaultConfig() AppConfig {
	return AppConfig{
		Theme:       "light",
		Language:    "es",
		FontSize:    12,
		AutoSave:    true,
		RecentFiles: []string{},
	}
}

func LoadConfig(appDataDir string) (AppConfig, error) {
	configPath := filepath.Join(appDataDir, "config.json")

	// Verificar si el archivo existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Crear configuración por defecto
		config := DefaultConfig()
		err := SaveConfig(appDataDir, config)
		return config, err
	}

	// Leer el archivo
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("error al leer configuración: %w", err)
	}

	// Deserializar JSON
	var config AppConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("error al parsear configuración: %w", err)
	}

	return config, nil
}

func SaveConfig(appDataDir string, config AppConfig) error {
	// Asegurar que el directorio existe
	os.MkdirAll(appDataDir, 0755)

	configPath := filepath.Join(appDataDir, "config.json")

	// Serializar a JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar configuración: %w", err)
	}

	// Escribir al archivo
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar configuración: %w", err)
	}

	return nil
}
```

### Comunicación con Servicios Externos

#### Consumo de APIs REST

```go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *APIClient) SetToken(token string) {
	c.Token = token
}

func (c *APIClient) Get(endpoint string, result interface{}) error {
	return c.Request(http.MethodGet, endpoint, nil, result)
}

func (c *APIClient) Post(endpoint string, body interface{}, result interface{}) error {
	return c.Request(http.MethodPost, endpoint, body, result)
}

func (c *APIClient) Put(endpoint string, body interface{}, result interface{}) error {
	return c.Request(http.MethodPut, endpoint, body, result)
}

func (c *APIClient) Delete(endpoint string) error {
	return c.Request(http.MethodDelete, endpoint, nil, nil)
}

func (c *APIClient) Request(method, endpoint string, body, result interface{}) error {
	// Preparar URL
	url := c.BaseURL + endpoint

	// Preparar cuerpo de la solicitud
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error al serializar cuerpo: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Crear solicitud
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("error al crear solicitud: %w", err)
	}

	// Establecer headers
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	// Enviar solicitud
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error al leer respuesta: %w", err)
	}

	// Verificar código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("error de API: código %d, respuesta: %s", resp.StatusCode, string(respBody))
	}

	// Deserializar respuesta si se proporcionó un destino
	if result != nil && len(respBody) > 0 {
		err = json.Unmarshal(respBody, result)
		if err != nil {
			return fmt.Errorf("error al deserializar respuesta: %w", err)
		}
	}

	return nil
}
```

Ejemplo de uso:

```go
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	client := api.NewAPIClient("https://api.example.com/v1")

	// Autenticación
	type AuthResponse struct {
		Token string `json:"token"`
	}

	var authResp AuthResponse
	err := client.Post("/auth/login", map[string]string{
		"username": "usuario",
		"password": "contraseña",
	}, &authResp)

	if err != nil {
		fmt.Printf("Error de autenticación: %v\n", err)
		return
	}

	client.SetToken(authResp.Token)

	// Obtener usuarios
	var users []User
	err = client.Get("/users", &users)
	if err != nil {
		fmt.Printf("Error al obtener usuarios: %v\n", err)
		return
	}

	fmt.Printf("Usuarios obtenidos: %+v\n", users)
}
```

#### Comunicación en Tiempo Real con WebSockets

```go
package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type WSClient struct {
	URL             string
	Conn            *websocket.Conn
	SendCh          chan Message
	RecvCh          chan Message
	Handlers        map[string]func(interface{})
	handlersMutex   sync.RWMutex
	ReconnectPeriod time.Duration
	Ctx             context.Context
	Cancel          context.CancelFunc
	Connected       bool
	connectedMutex  sync.RWMutex
}

func NewWSClient(url string) *WSClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &WSClient{
		URL:             url,
		SendCh:          make(chan Message, 10),
		RecvCh:          make(chan Message, 10),
		Handlers:        make(map[string]func(interface{})),
		ReconnectPeriod: 5 * time.Second,
		Ctx:             ctx,
		Cancel:          cancel,
	}
}

func (c *WSClient) Connect() error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(c.URL, nil)
	if err != nil {
		return fmt.Errorf("error al conectar: %w", err)
	}

	c.Conn = conn
	c.setConnected(true)

	// Iniciar goroutines para enviar y recibir
	go c.readPump()
	go c.writePump()

	return nil
}

func (c *WSClient) Close() {
	c.Cancel()
	if c.Conn != nil {
		c.Conn.Close()
	}
}

func (c *WSClient) RegisterHandler(messageType string, handler func(interface{})) {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()
	c.Handlers[messageType] = handler
}

func (c *WSClient) Send(messageType string, payload interface{}) {
	c.SendCh <- Message{Type: messageType, Payload: payload}
}

func (c *WSClient) isConnected() bool {
	c.connectedMutex.RLock()
	defer c.connectedMutex.RUnlock()
	return c.Connected
}

func (c *WSClient) setConnected(connected bool) {
	c.connectedMutex.Lock()
	defer c.connectedMutex.Unlock()
	c.Connected = connected
}

func (c *WSClient) readPump() {
	defer func() {
		c.setConnected(false)
		c.reconnect()
	}()

	for {
		select {
		case <-c.Ctx.Done():
			return
		default:
			_, data, err := c.Conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error al leer mensaje: %v\n", err)
				return
			}

			var message Message
			err = json.Unmarshal(data, &message)
			if err != nil {
				fmt.Printf("Error al deserializar mensaje: %v\n", err)
				continue
			}

			// Enviar a canal de recepción
			c.RecvCh <- message

			// Procesar con handler registrado
			c.handlersMutex.RLock()
			handler, exists := c.Handlers[message.Type]
			c.handlersMutex.RUnlock()

			if exists {
				handler(message.Payload)
			}
		}
	}
}

func (c *WSClient) writePump() {
	for {
		select {
		case <-c.Ctx.Done():
			return
		case msg := <-c.SendCh:
			if !c.isConnected() {
				continue
			}

			data, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error al serializar mensaje: %v\n", err)
				continue
			}

			err = c.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Printf("Error al enviar mensaje: %v\n", err)
				c.setConnected(false)
				c.reconnect()
				return
			}
		}
	}
}

func (c *WSClient) reconnect() {
	if c.ReconnectPeriod <= 0 {
		return
	}

	go func() {
		for {
			select {
			case <-c.Ctx.Done():
				return
			case <-time.After(c.ReconnectPeriod):
				if c.isConnected() {
					return
				}

				fmt.Println("Intentando reconectar...")
				err := c.Connect()
				if err != nil {
					fmt.Printf("Error al reconectar: %v\n", err)
					continue
				}

				fmt.Println("Reconexión exitosa")
				return
			}
		}
	}()
}
```

Ejemplo de uso:

```go
func main() {
	client := websocket.NewWSClient("ws://example.com/ws")

	// Registrar handlers
	client.RegisterHandler("chat_message", func(payload interface{}) {
		msg, ok := payload.(map[string]interface{})
		if !ok {
			return
		}
		fmt.Printf("Mensaje recibido de %s: %s\n", msg["user"], msg["text"])
	})

	client.RegisterHandler("user_joined", func(payload interface{}) {
		user, ok := payload.(string)
		if !ok {
			return
		}
		fmt.Printf("Usuario conectado: %s\n", user)
	})

	// Conectar
	err := client.Connect()
	if err != nil {
		fmt.Printf("Error al conectar: %v\n", err)
		return
	}

	// Enviar mensaje
	client.Send("chat_message", map[string]string{
		"text": "¡Hola a todos!",
	})

	// Mantener la aplicación en ejecución
	select {}
}
```

### Manejo de Archivos y Diálogos

#### Diálogos Nativos con Fyne

```go
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Manejo de Archivos")
	w.Resize(fyne.NewSize(600, 400))

	// Área de texto para mostrar contenido
	textArea := widget.NewMultiLineEntry()
	textArea.Wrapping = fyne.TextWrapWord

	// Botones para acciones de archivo
	openButton := widget.NewButton("Abrir", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return // Usuario canceló
			}
			defer reader.Close()

			// Leer contenido del archivo
			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// Mostrar contenido en el área de texto
			textArea.SetText(string(data))
			w.SetTitle(fmt.Sprintf("Manejo de Archivos - %s", filepath.Base(reader.URI().String())))
		}, w)

		// Filtrar por tipos de archivo
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".md", ".go"}))
		fd.Show()
	})

	saveButton := widget.NewButton("Guardar", func() {
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil {
				return // Usuario canceló
			}
			defer writer.Close()

			// Escribir contenido al archivo
			_, err = writer.Write([]byte(textArea.Text))
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			dialog.ShowInformation("Éxito", "Archivo guardado correctamente", w)
			w.SetTitle(fmt.Sprintf("Manejo de Archivos - %s", filepath.Base(writer.URI().String())))
		}, w)

		// Filtrar por tipos de archivo
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".md", ".go"}))
		fd.Show()
	})

	// Botón para mostrar diálogo de selección de directorio
	dirButton := widget.NewButton("Seleccionar Directorio", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if list == nil {
				return // Usuario canceló
			}

			// Listar archivos en el directorio
			items, err := list.List()
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			var fileList string
			for _, item := range items {
				fileList += filepath.Base(item.String()) + "\n"
			}

			textArea.SetText(fmt.Sprintf("Directorio: %s\n\nArchivos:\n%s", list.String(), fileList))
		}, w)
	})

	// Botón para mostrar diálogo de confirmación
	confirmButton := widget.NewButton("Confirmar Acción", func() {
		dialog.ShowConfirm("Confirmar", "¿Estás seguro de que deseas realizar esta acción?", func(ok bool) {
			if ok {
				textArea.SetText(textArea.Text + "\n\nAcción confirmada")
			} else {
				textArea.SetText(textArea.Text + "\n\nAcción cancelada")
			}
		}, w)
	})

	// Organizar la interfaz
	buttonContainer := container.NewHBox(openButton, saveButton, dirButton, confirmButton)
	content := container.NewBorder(buttonContainer, nil, nil, nil, textArea)
	w.SetContent(content)

	w.ShowAndRun()
}
```

#### Diálogos Nativos con Wails

```go
// main.go
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// OpenFile abre un diálogo para seleccionar un archivo y devuelve su contenido
func (a *App) OpenFile() (string, error) {
	// Mostrar diálogo de selección de archivo
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Seleccionar archivo",
		Filters: []runtime.FileFilter{
			{DisplayName: "Documentos de texto", Pattern: "*.txt;*.md"},
			{DisplayName: "Todos los archivos", Pattern: "*.*"},
		},
	})

	if err != nil {
		return "", fmt.Errorf("error al abrir diálogo: %w", err)
	}

	if filePath == "" {
		return "", nil // Usuario canceló
	}

	// Leer contenido del archivo
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error al leer archivo: %w", err)
	}

	return string(data), nil
}

// SaveFile muestra un diálogo para guardar un archivo
func (a *App) SaveFile(content string) (string, error) {
	// Mostrar diálogo de guardar archivo
	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Guardar archivo",
		Filters: []runtime.FileFilter{
			{DisplayName: "Documentos de texto", Pattern: "*.txt"},
			{DisplayName: "Markdown", Pattern: "*.md"},
		},
		DefaultFilename: "documento.txt",
	})

	if err != nil {
		return "", fmt.Errorf("error al abrir diálogo: %w", err)
	}

	if filePath == "" {
		return "", nil // Usuario canceló
	}

	// Escribir contenido al archivo
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("error al escribir archivo: %w", err)
	}

	return filePath, nil
}

// SelectDirectory muestra un diálogo para seleccionar un directorio
func (a *App) SelectDirectory() (string, error) {
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Seleccionar directorio",
	})

	if err != nil {
		return "", fmt.Errorf("error al abrir diálogo: %w", err)
	}

	if dirPath == "" {
		return "", nil // Usuario canceló
	}

	return dirPath, nil
}

// ListDirectory lista los archivos en un directorio
func (a *App) ListDirectory(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer directorio: %w", err)
	}

	var fileList []string
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}

	return fileList, nil
}

// ShowMessage muestra un mensaje al usuario
func (a *App) ShowMessage(title, message string) error {
	return runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})
}

// ShowConfirmation muestra un diálogo de confirmación
func (a *App) ShowConfirmation(title, message string) (bool, error) {
	return runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:      runtime.QuestionDialog,
		Title:     title,
		Message:   message,
		Buttons:   []string{"Sí", "No"},
		DefaultButton: "No",
		CancelButton:  "No",
	})
}
```

En el frontend (React):

```jsx
// frontend/src/App.jsx
import { useState, useEffect } from 'react';
import {
  OpenFile,
  SaveFile,
  SelectDirectory,
  ListDirectory,
  ShowMessage,
  ShowConfirmation
} from '../wailsjs/go/main/App';

function App() {
  const [fileContent, setFileContent] = useState('');
  const [filePath, setFilePath] = useState('');
  const [directoryPath, setDirectoryPath] = useState('');
  const [files, setFiles] = useState([]);

  const handleOpenFile = async () => {
    try {
      const content = await OpenFile();
      if (content) {
        setFileContent(content);
      }
    } catch (error) {
      console.error('Error al abrir archivo:', error);
    }
  };

  const handleSaveFile = async () => {
    try {
      const path = await SaveFile(fileContent);
      if (path) {
        setFilePath(path);
        ShowMessage('Éxito', `Archivo guardado en: ${path}`);
      }
    } catch (error) {
      console.error('Error al guardar archivo:', error);
    }
  };

  const handleSelectDirectory = async () => {
    try {
      const path = await SelectDirectory();
      if (path) {
        setDirectoryPath(path);
        const fileList = await ListDirectory(path);
        setFiles(fileList);
      }
    } catch (error) {
      console.error('Error al seleccionar directorio:', error);
    }
  };

  const handleConfirmDelete = async () => {
    try {
      const confirmed = await ShowConfirmation(
        'Confirmar eliminación',
        '¿Estás seguro de que deseas eliminar este archivo?'
      );
      
      if (confirmed) {
        ShowMessage('Información', 'Eliminación confirmada');
      }
    } catch (error) {
      console.error('Error en confirmación:', error);
    }
  };

  return (
    <div className="container">
      <h1>Editor de Texto</h1>
      
      <div className="button-group">
        <button onClick={handleOpenFile}>Abrir Archivo</button>
        <button onClick={handleSaveFile}>Guardar Archivo</button>
        <button onClick={handleSelectDirectory}>Seleccionar Directorio</button>
        <button onClick={handleConfirmDelete}>Eliminar</button>
      </div>
      
      {filePath && <p>Archivo actual: {filePath}</p>}
      
      <textarea
        value={fileContent}
        onChange={(e) => setFileContent(e.target.value)}
        rows={15}
        cols={80}
      />
      
      {directoryPath && (
        <div className="file-explorer">
          <h3>Contenido de: {directoryPath}</h3>
          <ul>
            {files.map((file, index) => (
              <li key={index}>{file}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

export default App;
```

## Distribución de Aplicaciones de Escritorio

### Empaquetado con Fyne

Fyne proporciona herramientas para empaquetar aplicaciones para diferentes plataformas:

```bash
# Instalar la herramienta fyne
go install fyne.io/fyne/v2/cmd/fyne@latest

# Empaquetar para la plataforma actual
fyne package -icon icon.png

# Empaquetar para Windows
fyne package -os windows -icon icon.png

# Empaquetar para macOS
fyne package -os darwin -icon icon.png

# Empaquetar para Linux
fyne package -os linux -icon icon.png

# Crear un instalador para Windows
fyne package -os windows -icon icon.png -install
```

### Empaquetado con Wails

Wails también proporciona herramientas para empaquetar aplicaciones:

```bash
# Construir para desarrollo
wails dev

# Construir para producción (plataforma actual)
wails build

# Construir para Windows
wails build -platform windows/amd64

# Construir para macOS
wails build -platform darwin/universal

# Construir para Linux
wails build -platform linux/amd64

# Crear un instalador para Windows
wails build -platform windows/amd64 -nsis
```

## Optimización y Mejores Prácticas

### Rendimiento

1. **Evitar bloqueos en la UI**: Ejecutar operaciones pesadas en goroutines separadas.

```go
func (a *App) ProcessLargeFile(filePath string) {
	// Mostrar indicador de progreso
	runtime.EventsEmit(a.ctx, "processing-started")
	
	// Ejecutar en una goroutine separada
	go func() {
		// Simular procesamiento largo
		result, err := processFile(filePath)
		
		// Actualizar UI cuando termine
		if err != nil {
			runtime.EventsEmit(a.ctx, "processing-error", err.Error())
		} else {
			runtime.EventsEmit(a.ctx, "processing-complete", result)
		}
	}()
}
```

2. **Reutilización de recursos**: Evitar crear y destruir recursos frecuentemente.

```go
// Mal ejemplo: Crear una nueva conexión para cada consulta
func (a *App) BadQueryDatabase() []Record {
	db, _ := sql.Open("sqlite3", "app.db")
	defer db.Close()
	
	// Ejecutar consulta
	// ...
}

// Buen ejemplo: Reutilizar conexión
type App struct {
	ctx context.Context
	db  *sql.DB
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.db, _ = sql.Open("sqlite3", "app.db")
}

func (a *App) shutdown(ctx context.Context) {
	a.db.Close()
}

func (a *App) GoodQueryDatabase() []Record {
	// Usar la conexión existente
	// ...
}
```

### Seguridad

1. **Validación de entrada**: Validar todas las entradas del usuario.

```go
func (a *App) SaveUserData(name, email string, age int) error {
	// Validar nombre
	if name == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	
	// Validar email con expresión regular
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("email inválido")
	}
	
	// Validar edad
	if age < 0 || age > 120 {
		return errors.New("edad inválida")
	}
	
	// Proceder con el guardado
	// ...
	return nil
}
```

2. **Manejo seguro de archivos**: Validar rutas y permisos.

```go
func (a *App) ReadUserFile(filePath string) (string, error) {
	// Validar que la ruta esté dentro del directorio permitido
	userDataDir := filepath.Join(a.GetUserDataDir(), "files")
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}
	
	// Verificar que la ruta absoluta comienza con el directorio permitido
	if !strings.HasPrefix(absPath, userDataDir) {
		return "", errors.New("acceso denegado: ruta fuera del directorio permitido")
	}
	
	// Leer archivo
	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}
```

### Experiencia de Usuario

1. **Retroalimentación visual**: Proporcionar indicadores de progreso para operaciones largas.

```jsx
// En React (Wails)
import { useState } from 'react';
import { ProcessLargeFile } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime';

function ProcessingComponent() {
  const [status, setStatus] = useState('idle');
  const [progress, setProgress] = useState(0);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);
  
  useEffect(() => {
    // Suscribirse a eventos
    EventsOn('processing-started', () => {
      setStatus('processing');
      setProgress(0);
      setError(null);
    });
    
    EventsOn('processing-progress', (percent) => {
      setProgress(percent);
    });
    
    EventsOn('processing-complete', (result) => {
      setStatus('complete');
      setResult(result);
    });
    
    EventsOn('processing-error', (errorMsg) => {
      setStatus('error');
      setError(errorMsg);
    });
  }, []);
  
  const handleProcess = () => {
    ProcessLargeFile('/path/to/file.dat');
  };
  
  return (
    <div>
      <button onClick={handleProcess} disabled={status === 'processing'}>
        Procesar Archivo
      </button>
      
      {status === 'processing' && (
        <div className="progress-bar">
          <div className="progress" style={{ width: `${progress}%` }}></div>
          <span>{progress}%</span>
        </div>
      )}
      
      {status === 'complete' && (
        <div className="result">
          <h3>Procesamiento completado</h3>
          <pre>{JSON.stringify(result, null, 2)}</pre>
        </div>
      )}
      
      {status === 'error' && (
        <div className="error">
          <h3>Error</h3>
          <p>{error}</p>
        </div>
      )}
    </div>
  );
}
```

2. **Accesibilidad**: Asegurar que la aplicación sea accesible para todos los usuarios.

```go
// En Fyne
package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.New()
	
	// Establecer tema de alto contraste para mejor accesibilidad
	a.Settings().SetTheme(theme.HighContrastTheme())
	
	w := a.NewWindow("Aplicación Accesible")
	
	// Usar etiquetas descriptivas
	nameLabel := widget.NewLabel("Nombre:")
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Ingresa tu nombre")
	
	// Establecer descripciones para lectores de pantalla
	nameEntry.Importance = widget.HighImportance
	
	// Botón con texto claro
	submitButton := widget.NewButton("Enviar Formulario", func() {
		// Acción
	})
	
	// Usar tamaño de texto más grande
	heading := canvas.NewText("Formulario de Registro", theme.PrimaryColor())
	heading.TextSize = 24
	heading.TextStyle.Bold = true
	
	content := container.NewVBox(
		heading,
		nameLabel,
		nameEntry,
		submitButton,
	)
	
	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}
```

## Ejercicios Prácticos

### Ejercicio 1: Editor de Texto Simple

Crea un editor de texto simple con las siguientes funcionalidades:

- Abrir y guardar archivos
- Edición básica de texto
- Búsqueda y reemplazo
- Configuración de fuente y tema

### Ejercicio 2: Gestor de Tareas

Desarrolla una aplicación de gestión de tareas con:

- Creación, edición y eliminación de tareas
- Categorización y priorización
- Recordatorios y notificaciones
- Persistencia de datos con SQLite

### Ejercicio 3: Visor de Imágenes

Implementa un visor de imágenes con:

- Navegación por directorios
- Zoom y rotación de imágenes
- Presentación de diapositivas
- Edición básica (recorte, ajuste de brillo/contraste)

## Conclusiones

Go ofrece diversas opciones para desarrollar aplicaciones de escritorio, desde bibliotecas nativas como Fyne y Gio hasta enfoques híbridos como Wails y Webview. Cada enfoque tiene sus ventajas y desventajas, y la elección dependerá de los requisitos específicos del proyecto.

Para aplicaciones que requieren una apariencia nativa y un rendimiento óptimo, las bibliotecas como Fyne son una excelente opción. Para proyectos que necesitan interfaces de usuario más ricas o donde el equipo tiene experiencia en desarrollo web, los enfoques híbridos como Wails pueden ser más adecuados.

Independientemente del enfoque elegido, es importante seguir buenas prácticas de diseño, rendimiento y seguridad para crear aplicaciones de escritorio robustas y amigables para el usuario.

## Referencias

1. Fyne - Toolkit GUI nativo para Go: https://fyne.io/
2. Gio - Toolkit GUI inmediato para Go: https://gioui.org/
3. ui - Binding de Go para libui: https://github.com/andlabs/ui
4. Wails - Framework para aplicaciones de escritorio con Go y Web: https://wails.io/
5. Webview - Componente webview minimalista: https://github.com/webview/webview
6. Lorca - Biblioteca para crear aplicaciones de escritorio con Go y Chrome: https://github.com/zserge/lorca
7. "Building Cross-Platform GUI Applications with Fyne" - Andrew Williams
8. "Desktop Apps with Go" - Mark Bates
9. "Building Desktop Applications with Go" - Packt Publishing