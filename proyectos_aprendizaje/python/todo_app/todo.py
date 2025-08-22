#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Aplicación de gestión de tareas (TODO) en consola.

Este script implementa una aplicación de línea de comandos para gestionar tareas,
siguiendo las mejores prácticas de Python:
- Uso de docstrings para documentación
- Estructura modular con clases y funciones bien definidas
- Manejo de errores con excepciones
- Uso de tipos de datos apropiados
- Implementación de patrones de diseño (Repository pattern)

Autor: Ejemplo
Versión: 1.0
"""

import argparse
import json
import os
import sys
from datetime import datetime
from typing import Dict, List, Optional, Any, Union


class Task:
    """Clase que representa una tarea en la aplicación.
    
    Attributes:
        id (int): Identificador único de la tarea.
        title (str): Título descriptivo de la tarea.
        completed (bool): Estado de completitud de la tarea.
        created_at (str): Fecha y hora de creación en formato ISO.
        completed_at (Optional[str]): Fecha y hora de completitud en formato ISO, o None.
    """
    
    def __init__(self, id: int, title: str, completed: bool = False,
                 created_at: Optional[str] = None, completed_at: Optional[str] = None) -> None:
        """Inicializa una nueva tarea.
        
        Args:
            id: Identificador único de la tarea.
            title: Título descriptivo de la tarea.
            completed: Estado de completitud de la tarea (por defecto False).
            created_at: Fecha y hora de creación (por defecto ahora).
            completed_at: Fecha y hora de completitud (por defecto None).
        """
        self.id = id
        self.title = title
        self.completed = completed
        self.created_at = created_at or datetime.now().isoformat()
        self.completed_at = completed_at
    
    def complete(self) -> None:
        """Marca la tarea como completada y establece la fecha de completitud."""
        if not self.completed:
            self.completed = True
            self.completed_at = datetime.now().isoformat()
    
    def to_dict(self) -> Dict[str, Any]:
        """Convierte la tarea a un diccionario para serialización.
        
        Returns:
            Dict[str, Any]: Representación de la tarea como diccionario.
        """
        return {
            'id': self.id,
            'title': self.title,
            'completed': self.completed,
            'created_at': self.created_at,
            'completed_at': self.completed_at
        }
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'Task':
        """Crea una instancia de Task a partir de un diccionario.
        
        Args:
            data: Diccionario con los datos de la tarea.
            
        Returns:
            Task: Nueva instancia de Task.
        """
        return cls(
            id=data['id'],
            title=data['title'],
            completed=data['completed'],
            created_at=data['created_at'],
            completed_at=data['completed_at']
        )
    
    def __str__(self) -> str:
        """Devuelve una representación en cadena de la tarea.
        
        Returns:
            str: Representación legible de la tarea.
        """
        status = "[X]" if self.completed else "[ ]"
        return f"{self.id}. {status} {self.title}"


class TaskRepository:
    """Repositorio para gestionar la persistencia de tareas.
    
    Esta clase implementa el patrón Repository para abstraer
    la lógica de persistencia de las tareas.
    
    Attributes:
        file_path (str): Ruta al archivo de almacenamiento de tareas.
        tasks (List[Task]): Lista de tareas cargadas.
    """
    
    def __init__(self, file_path: str) -> None:
        """Inicializa el repositorio de tareas.
        
        Args:
            file_path: Ruta al archivo JSON para almacenar las tareas.
        """
        self.file_path = file_path
        self.tasks: List[Task] = []
        self._load()
    
    def _load(self) -> None:
        """Carga las tareas desde el archivo JSON."""
        try:
            if os.path.exists(self.file_path):
                with open(self.file_path, 'r', encoding='utf-8') as f:
                    tasks_data = json.load(f)
                    self.tasks = [Task.from_dict(task_data) for task_data in tasks_data]
        except (json.JSONDecodeError, KeyError) as e:
            print(f"Error al cargar tareas: {e}", file=sys.stderr)
            # Crear un archivo nuevo si hay problemas con el existente
            self.tasks = []
            self._save()
    
    def _save(self) -> None:
        """Guarda las tareas en el archivo JSON."""
        try:
            # Asegurar que el directorio existe
            os.makedirs(os.path.dirname(self.file_path), exist_ok=True)
            
            with open(self.file_path, 'w', encoding='utf-8') as f:
                tasks_data = [task.to_dict() for task in self.tasks]
                json.dump(tasks_data, f, indent=2, ensure_ascii=False)
        except (IOError, OSError) as e:
            print(f"Error al guardar tareas: {e}", file=sys.stderr)
    
    def get_all(self) -> List[Task]:
        """Obtiene todas las tareas.
        
        Returns:
            List[Task]: Lista de todas las tareas.
        """
        return self.tasks
    
    def get_by_id(self, task_id: int) -> Optional[Task]:
        """Obtiene una tarea por su ID.
        
        Args:
            task_id: ID de la tarea a buscar.
            
        Returns:
            Optional[Task]: La tarea encontrada o None si no existe.
        """
        for task in self.tasks:
            if task.id == task_id:
                return task
        return None
    
    def add(self, title: str) -> Task:
        """Añade una nueva tarea.
        
        Args:
            title: Título de la nueva tarea.
            
        Returns:
            Task: La tarea creada.
        """
        # Generar un nuevo ID (el máximo actual + 1, o 1 si no hay tareas)
        new_id = max([task.id for task in self.tasks], default=0) + 1
        task = Task(id=new_id, title=title)
        self.tasks.append(task)
        self._save()
        return task
    
    def update(self, task: Task) -> None:
        """Actualiza una tarea existente.
        
        Args:
            task: La tarea actualizada.
        """
        for i, t in enumerate(self.tasks):
            if t.id == task.id:
                self.tasks[i] = task
                self._save()
                return
    
    def delete(self, task_id: int) -> bool:
        """Elimina una tarea por su ID.
        
        Args:
            task_id: ID de la tarea a eliminar.
            
        Returns:
            bool: True si se eliminó correctamente, False si no se encontró.
        """
        for i, task in enumerate(self.tasks):
            if task.id == task_id:
                del self.tasks[i]
                self._save()
                return True
        return False


class TodoApp:
    """Aplicación principal para gestión de tareas.
    
    Esta clase implementa la lógica de la aplicación y
    coordina las interacciones entre la interfaz de usuario
    y el repositorio de tareas.
    
    Attributes:
        repository (TaskRepository): Repositorio para la persistencia de tareas.
    """
    
    def __init__(self, repository: TaskRepository) -> None:
        """Inicializa la aplicación.
        
        Args:
            repository: Repositorio para la persistencia de tareas.
        """
        self.repository = repository
    
    def list_tasks(self, show_completed: bool = True) -> None:
        """Muestra la lista de tareas.
        
        Args:
            show_completed: Si es True, muestra todas las tareas;
                           si es False, solo muestra las pendientes.
        """
        tasks = self.repository.get_all()
        
        if not tasks:
            print("No hay tareas registradas.")
            return
        
        print("\nLista de Tareas:")
        print("-" * 40)
        
        for task in tasks:
            if show_completed or not task.completed:
                print(task)
        
        print("-" * 40)
    
    def add_task(self, title: str) -> None:
        """Añade una nueva tarea.
        
        Args:
            title: Título de la nueva tarea.
        """
        if not title.strip():
            print("Error: El título de la tarea no puede estar vacío.", file=sys.stderr)
            return
        
        task = self.repository.add(title)
        print(f"Tarea añadida: {task}")
    
    def complete_task(self, task_id: int) -> None:
        """Marca una tarea como completada.
        
        Args:
            task_id: ID de la tarea a completar.
        """
        task = self.repository.get_by_id(task_id)
        
        if not task:
            print(f"Error: No se encontró la tarea con ID {task_id}", file=sys.stderr)
            return
        
        if task.completed:
            print(f"La tarea {task_id} ya está completada.")
            return
        
        task.complete()
        self.repository.update(task)
        print(f"Tarea completada: {task}")
    
    def delete_task(self, task_id: int) -> None:
        """Elimina una tarea.
        
        Args:
            task_id: ID de la tarea a eliminar.
        """
        if self.repository.delete(task_id):
            print(f"Tarea {task_id} eliminada.")
        else:
            print(f"Error: No se encontró la tarea con ID {task_id}", file=sys.stderr)


def parse_arguments() -> argparse.Namespace:
    """Analiza los argumentos de línea de comandos.
    
    Returns:
        argparse.Namespace: Objeto con los argumentos analizados.
    """
    parser = argparse.ArgumentParser(description="Aplicación de gestión de tareas (TODO)")
    subparsers = parser.add_subparsers(dest="command", help="Comandos disponibles")
    
    # Comando: list
    list_parser = subparsers.add_parser("list", help="Listar tareas")
    list_parser.add_argument(
        "--all", action="store_true", 
        help="Mostrar todas las tareas (por defecto solo muestra las pendientes)"
    )
    
    # Comando: add
    add_parser = subparsers.add_parser("add", help="Añadir una nueva tarea")
    add_parser.add_argument("title", help="Título de la tarea")
    
    # Comando: complete
    complete_parser = subparsers.add_parser("complete", help="Marcar una tarea como completada")
    complete_parser.add_argument("id", type=int, help="ID de la tarea")
    
    # Comando: delete
    delete_parser = subparsers.add_parser("delete", help="Eliminar una tarea")
    delete_parser.add_argument("id", type=int, help="ID de la tarea")
    
    return parser.parse_args()


def main() -> None:
    """Función principal que inicia la aplicación."""
    # Definir la ruta del archivo de datos
    data_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)), "data")
    file_path = os.path.join(data_dir, "tasks.json")
    
    # Crear el repositorio y la aplicación
    repository = TaskRepository(file_path)
    app = TodoApp(repository)
    
    # Analizar argumentos
    args = parse_arguments()
    
    # Ejecutar el comando correspondiente
    if args.command == "list":
        app.list_tasks(show_completed=args.all)
    elif args.command == "add":
        app.add_task(args.title)
    elif args.command == "complete":
        app.complete_task(args.id)
    elif args.command == "delete":
        app.delete_task(args.id)
    else:
        # Si no se especifica un comando, mostrar la lista de tareas pendientes
        app.list_tasks(show_completed=False)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nOperación cancelada por el usuario.")
        sys.exit(0)
    except Exception as e:
        print(f"Error inesperado: {e}", file=sys.stderr)
        sys.exit(1)