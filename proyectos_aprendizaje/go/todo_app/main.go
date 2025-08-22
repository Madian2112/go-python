package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Task representa una tarea individual
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskRepository gestiona la persistencia de las tareas
type TaskRepository struct {
	filePath string
	tasks    []Task
}

// NewTaskRepository crea una nueva instancia de TaskRepository
func NewTaskRepository(dataDir string) (*TaskRepository, error) {
	// Asegurar que el directorio de datos exista
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("error al crear directorio de datos: %w", err)
	}

	filePath := filepath.Join(dataDir, "tasks.json")
	repo := &TaskRepository{
		filePath: filePath,
		tasks:    []Task{},
	}

	// Cargar tareas existentes si el archivo existe
	if _, err := os.Stat(filePath); err == nil {
		if err := repo.load(); err != nil {
			return nil, fmt.Errorf("error al cargar tareas: %w", err)
		}
	}

	return repo, nil
}

// load carga las tareas desde el archivo
func (r *TaskRepository) load() error {
	data, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		r.tasks = []Task{}
		return nil
	}

	return json.Unmarshal(data, &r.tasks)
}

// save guarda las tareas en el archivo
func (r *TaskRepository) save() error {
	data, err := json.MarshalIndent(r.tasks, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(r.filePath, data, 0644)
}

// GetAll devuelve todas las tareas
func (r *TaskRepository) GetAll() []Task {
	return r.tasks
}

// GetPending devuelve solo las tareas pendientes
func (r *TaskRepository) GetPending() []Task {
	var pendingTasks []Task
	for _, task := range r.tasks {
		if !task.Completed {
			pendingTasks = append(pendingTasks, task)
		}
	}
	return pendingTasks
}

// Add añade una nueva tarea
func (r *TaskRepository) Add(title string) (Task, error) {
	if title == "" {
		return Task{}, errors.New("el título de la tarea no puede estar vacío")
	}

	// Generar un nuevo ID (el máximo ID actual + 1)
	nextID := 1
	for _, task := range r.tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}

	newTask := Task{
		ID:        nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	r.tasks = append(r.tasks, newTask)

	if err := r.save(); err != nil {
		return Task{}, fmt.Errorf("error al guardar tarea: %w", err)
	}

	return newTask, nil
}

// Complete marca una tarea como completada
func (r *TaskRepository) Complete(id int) error {
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks[i].Completed = true
			return r.save()
		}
	}

	return fmt.Errorf("tarea con ID %d no encontrada", id)
}

// Delete elimina una tarea
func (r *TaskRepository) Delete(id int) error {
	for i, task := range r.tasks {
		if task.ID == id {
			// Eliminar la tarea del slice
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return r.save()
		}
	}

	return fmt.Errorf("tarea con ID %d no encontrada", id)
}

// TodoApp implementa la lógica de la aplicación
type TodoApp struct {
	repo *TaskRepository
}

// NewTodoApp crea una nueva instancia de TodoApp
func NewTodoApp(repo *TaskRepository) *TodoApp {
	return &TodoApp{repo: repo}
}

// ListTasks muestra las tareas (todas o solo pendientes)
func (app *TodoApp) ListTasks(showAll bool) {
	var tasks []Task
	if showAll {
		tasks = app.repo.GetAll()
		fmt.Println("Todas las tareas:")
	} else {
		tasks = app.repo.GetPending()
		fmt.Println("Tareas pendientes:")
	}

	if len(tasks) == 0 {
		fmt.Println("No hay tareas para mostrar.")
		return
	}

	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "✓"
		}
		fmt.Printf("[%s] %d: %s (Creada: %s)\n", status, task.ID, task.Title, task.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

// AddTask añade una nueva tarea
func (app *TodoApp) AddTask(title string) {
	task, err := app.repo.Add(title)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Tarea añadida con ID %d: %s\n", task.ID, task.Title)
}

// CompleteTask marca una tarea como completada
func (app *TodoApp) CompleteTask(id int) {
	if err := app.repo.Complete(id); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Tarea %d marcada como completada\n", id)
}

// DeleteTask elimina una tarea
func (app *TodoApp) DeleteTask(id int) {
	if err := app.repo.Delete(id); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Tarea %d eliminada\n", id)
}

func main() {
	// Definir comandos y flags
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listAll := listCmd.Bool("all", false, "Mostrar todas las tareas (incluyendo completadas)")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// Verificar que se proporcionó un subcomando
	if len(os.Args) < 2 {
		fmt.Println("Se requiere un subcomando: list, add, complete o delete")
		os.Exit(1)
	}

	// Crear directorio de datos
	dataDir := filepath.Join("data")
	repo, err := NewTaskRepository(dataDir)
	if err != nil {
		fmt.Printf("Error al inicializar el repositorio: %s\n", err)
		os.Exit(1)
	}

	app := NewTodoApp(repo)

	// Procesar subcomandos
	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		app.ListTasks(*listAll)

	case "add":
		addCmd.Parse(os.Args[2:])
		if addCmd.NArg() < 1 {
			fmt.Println("Se requiere un título para la tarea")
			os.Exit(1)
		}
		app.AddTask(addCmd.Arg(0))

	case "complete":
		completeCmd.Parse(os.Args[2:])
		if completeCmd.NArg() < 1 {
			fmt.Println("Se requiere un ID de tarea")
			os.Exit(1)
		}
		id, err := strconv.Atoi(completeCmd.Arg(0))
		if err != nil {
			fmt.Println("El ID debe ser un número")
			os.Exit(1)
		}
		app.CompleteTask(id)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if deleteCmd.NArg() < 1 {
			fmt.Println("Se requiere un ID de tarea")
			os.Exit(1)
		}
		id, err := strconv.Atoi(deleteCmd.Arg(0))
		if err != nil {
			fmt.Println("El ID debe ser un número")
			os.Exit(1)
		}
		app.DeleteTask(id)

	default:
		fmt.Printf("Subcomando desconocido: %s\n", os.Args[1])
		os.Exit(1)
	}
}