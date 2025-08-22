#!/usr/bin/env python3
"""
Analizador de Logs - Nivel 1

Esta aplicación de consola analiza archivos de log, filtra entradas por nivel de severidad
y genera estadísticas básicas.
"""

import argparse
import re
import os
import sys
from datetime import datetime
from collections import Counter, defaultdict
from typing import Dict, List, Tuple, Optional

# Definición de colores para la salida en consola
class Colors:
    RESET = "\033[0m"
    RED = "\033[91m"
    GREEN = "\033[92m"
    YELLOW = "\033[93m"
    BLUE = "\033[94m"
    MAGENTA = "\033[95m"
    CYAN = "\033[96m"

# Definición de niveles de log y sus colores asociados
LOG_LEVELS = {
    "DEBUG": Colors.BLUE,
    "INFO": Colors.GREEN,
    "WARNING": Colors.YELLOW,
    "ERROR": Colors.RED,
    "CRITICAL": Colors.MAGENTA
}

class LogEntry:
    """Representa una entrada individual de log."""
    
    def __init__(self, timestamp: datetime, level: str, component: str, message: str):
        self.timestamp = timestamp
        self.level = level
        self.component = component
        self.message = message
    
    def __str__(self) -> str:
        """Devuelve una representación en string de la entrada de log."""
        color = LOG_LEVELS.get(self.level, Colors.RESET)
        timestamp_str = self.timestamp.strftime("%Y-%m-%d %H:%M:%S")
        return f"{timestamp_str} {color}{self.level}{Colors.RESET} [{self.component}] {self.message}"

class LogAnalyzer:
    """Analizador de archivos de log."""
    
    def __init__(self, log_file_path: str):
        """Inicializa el analizador con la ruta al archivo de log."""
        self.log_file_path = log_file_path
        self.entries: List[LogEntry] = []
        self.log_pattern = re.compile(
            r"(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) (\w+) \[([^\]]+)\] (.+)"
        )
    
    def parse_log_file(self) -> None:
        """Lee y parsea el archivo de log."""
        try:
            with open(self.log_file_path, 'r') as file:
                for line in file:
                    entry = self._parse_log_line(line.strip())
                    if entry:
                        self.entries.append(entry)
        except FileNotFoundError:
            print(f"Error: El archivo {self.log_file_path} no existe.")
            sys.exit(1)
        except Exception as e:
            print(f"Error al leer el archivo de log: {e}")
            sys.exit(1)
    
    def _parse_log_line(self, line: str) -> Optional[LogEntry]:
        """Parsea una línea de log y devuelve un objeto LogEntry."""
        match = self.log_pattern.match(line)
        if not match:
            return None
        
        timestamp_str, level, component, message = match.groups()
        try:
            timestamp = datetime.strptime(timestamp_str, "%Y-%m-%d %H:%M:%S")
            return LogEntry(timestamp, level, component, message)
        except ValueError:
            return None
    
    def filter_by_level(self, level: str) -> List[LogEntry]:
        """Filtra las entradas de log por nivel de severidad."""
        return [entry for entry in self.entries if entry.level == level]
    
    def filter_by_component(self, component: str) -> List[LogEntry]:
        """Filtra las entradas de log por componente."""
        return [entry for entry in self.entries if component.lower() in entry.component.lower()]
    
    def filter_by_keyword(self, keyword: str) -> List[LogEntry]:
        """Filtra las entradas de log por palabra clave en el mensaje."""
        return [entry for entry in self.entries if keyword.lower() in entry.message.lower()]
    
    def get_level_statistics(self) -> Dict[str, int]:
        """Genera estadísticas de cantidad de entradas por nivel."""
        counter = Counter(entry.level for entry in self.entries)
        return dict(counter)
    
    def get_component_statistics(self) -> Dict[str, int]:
        """Genera estadísticas de cantidad de entradas por componente."""
        counter = Counter(entry.component for entry in self.entries)
        return dict(counter)
    
    def get_hourly_distribution(self) -> Dict[int, int]:
        """Genera estadísticas de distribución de entradas por hora."""
        distribution = defaultdict(int)
        for entry in self.entries:
            distribution[entry.timestamp.hour] += 1
        return dict(distribution)
    
    def get_error_summary(self) -> Dict[str, List[str]]:
        """Genera un resumen de errores agrupados por componente."""
        error_summary = defaultdict(list)
        for entry in self.entries:
            if entry.level in ["ERROR", "CRITICAL"]:
                error_summary[entry.component].append(entry.message)
        return dict(error_summary)

class LogAnalyzerCLI:
    """Interfaz de línea de comandos para el analizador de logs."""
    
    def __init__(self):
        self.parser = self._create_parser()
        self.args = self.parser.parse_args()
    
    def _create_parser(self) -> argparse.ArgumentParser:
        """Crea y configura el parser de argumentos."""
        parser = argparse.ArgumentParser(
            description="Analizador de archivos de log",
            formatter_class=argparse.RawDescriptionHelpFormatter
        )
        
        parser.add_argument(
            "log_file",
            help="Ruta al archivo de log a analizar"
        )
        
        parser.add_argument(
            "-l", "--level",
            help="Filtrar por nivel de log (DEBUG, INFO, WARNING, ERROR, CRITICAL)"
        )
        
        parser.add_argument(
            "-c", "--component",
            help="Filtrar por componente"
        )
        
        parser.add_argument(
            "-k", "--keyword",
            help="Filtrar por palabra clave en el mensaje"
        )
        
        parser.add_argument(
            "-s", "--stats",
            action="store_true",
            help="Mostrar estadísticas del archivo de log"
        )
        
        parser.add_argument(
            "-e", "--errors",
            action="store_true",
            help="Mostrar resumen de errores"
        )
        
        return parser
    
    def run(self) -> None:
        """Ejecuta el analizador de logs según los argumentos proporcionados."""
        analyzer = LogAnalyzer(self.args.log_file)
        analyzer.parse_log_file()
        
        # Aplicar filtros si se especificaron
        filtered_entries = analyzer.entries
        
        if self.args.level:
            filtered_entries = [entry for entry in filtered_entries 
                               if entry.level == self.args.level.upper()]
        
        if self.args.component:
            filtered_entries = [entry for entry in filtered_entries 
                               if self.args.component.lower() in entry.component.lower()]
        
        if self.args.keyword:
            filtered_entries = [entry for entry in filtered_entries 
                               if self.args.keyword.lower() in entry.message.lower()]
        
        # Mostrar estadísticas si se solicitaron
        if self.args.stats:
            self._show_statistics(analyzer)
        
        # Mostrar resumen de errores si se solicitó
        if self.args.errors:
            self._show_error_summary(analyzer)
        
        # Si no se solicitaron estadísticas ni resumen de errores, mostrar entradas filtradas
        if not (self.args.stats or self.args.errors):
            self._show_entries(filtered_entries)
    
    def _show_entries(self, entries: List[LogEntry]) -> None:
        """Muestra las entradas de log en la consola."""
        if not entries:
            print("No se encontraron entradas que coincidan con los criterios de filtrado.")
            return
        
        print(f"Mostrando {len(entries)} entradas:")
        print("-" * 80)
        for entry in entries:
            print(entry)
    
    def _show_statistics(self, analyzer: LogAnalyzer) -> None:
        """Muestra estadísticas del archivo de log."""
        level_stats = analyzer.get_level_statistics()
        component_stats = analyzer.get_component_statistics()
        hourly_distribution = analyzer.get_hourly_distribution()
        
        print("\nEstadísticas del archivo de log:")
        print("-" * 40)
        
        print("\nDistribución por nivel:")
        for level, count in sorted(level_stats.items()):
            color = LOG_LEVELS.get(level, Colors.RESET)
            print(f"  {color}{level}{Colors.RESET}: {count}")
        
        print("\nDistribución por componente:")
        for component, count in sorted(component_stats.items(), key=lambda x: x[1], reverse=True):
            print(f"  {component}: {count}")
        
        print("\nDistribución por hora:")
        for hour in sorted(hourly_distribution.keys()):
            count = hourly_distribution[hour]
            bar = "█" * (count // 2 + 1)
            print(f"  {hour:02d}:00 - {hour:02d}:59: {count} {bar}")
    
    def _show_error_summary(self, analyzer: LogAnalyzer) -> None:
        """Muestra un resumen de los errores encontrados."""
        error_summary = analyzer.get_error_summary()
        
        if not error_summary:
            print("No se encontraron errores en el archivo de log.")
            return
        
        print("\nResumen de errores:")
        print("-" * 40)
        
        for component, errors in sorted(error_summary.items()):
            print(f"\n{Colors.CYAN}{component}{Colors.RESET}:")
            for i, error in enumerate(errors, 1):
                print(f"  {i}. {Colors.RED}{error}{Colors.RESET}")

def main():
    """Función principal del programa."""
    cli = LogAnalyzerCLI()
    cli.run()

if __name__ == "__main__":
    main()