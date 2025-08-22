package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Colores para la salida en consola
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[91m"
	ColorGreen  = "\033[92m"
	ColorYellow = "\033[93m"
	ColorBlue   = "\033[94m"
	ColorPurple = "\033[95m"
	ColorCyan   = "\033[96m"
)

// Niveles de log y sus colores asociados
var logLevelColors = map[string]string{
	"DEBUG":    ColorBlue,
	"INFO":     ColorGreen,
	"WARNING":  ColorYellow,
	"ERROR":    ColorRed,
	"CRITICAL": ColorPurple,
}

// LogEntry representa una entrada individual de log
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Component string
	Message   string
}

// String devuelve una representación en string de la entrada de log
func (e LogEntry) String() string {
	color := logLevelColors[e.Level]
	timestampStr := e.Timestamp.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s %s%s%s [%s] %s", 
		timestampStr, color, e.Level, ColorReset, e.Component, e.Message)
}

// LogAnalyzer analiza archivos de log
type LogAnalyzer struct {
	LogFilePath string
	Entries     []LogEntry
	LogPattern  *regexp.Regexp
}

// NewLogAnalyzer crea un nuevo analizador de logs
func NewLogAnalyzer(logFilePath string) *LogAnalyzer {
	return &LogAnalyzer{
		LogFilePath: logFilePath,
		Entries:     []LogEntry{},
		LogPattern:  regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) (\w+) \[([^\]]+)\] (.+)`),
	}
}

// ParseLogFile lee y parsea el archivo de log
func (a *LogAnalyzer) ParseLogFile() error {
	file, err := os.Open(a.LogFilePath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := a.parseLogLine(line)
		if err == nil {
			a.Entries = append(a.Entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error al leer el archivo: %w", err)
	}

	return nil
}

// parseLogLine parsea una línea de log y devuelve un objeto LogEntry
func (a *LogAnalyzer) parseLogLine(line string) (LogEntry, error) {
	matches := a.LogPattern.FindStringSubmatch(line)
	if matches == nil || len(matches) != 5 {
		return LogEntry{}, fmt.Errorf("formato de línea inválido")
	}

	timestampStr, level, component, message := matches[1], matches[2], matches[3], matches[4]
	timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		return LogEntry{}, fmt.Errorf("formato de timestamp inválido: %w", err)
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Component: component,
		Message:   message,
	}, nil
}

// FilterByLevel filtra las entradas de log por nivel de severidad
func (a *LogAnalyzer) FilterByLevel(level string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range a.Entries {
		if entry.Level == level {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// FilterByComponent filtra las entradas de log por componente
func (a *LogAnalyzer) FilterByComponent(component string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range a.Entries {
		if strings.Contains(strings.ToLower(entry.Component), strings.ToLower(component)) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// FilterByKeyword filtra las entradas de log por palabra clave en el mensaje
func (a *LogAnalyzer) FilterByKeyword(keyword string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range a.Entries {
		if strings.Contains(strings.ToLower(entry.Message), strings.ToLower(keyword)) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// GetLevelStatistics genera estadísticas de cantidad de entradas por nivel
func (a *LogAnalyzer) GetLevelStatistics() map[string]int {
	stats := make(map[string]int)
	for _, entry := range a.Entries {
		stats[entry.Level]++
	}
	return stats
}

// GetComponentStatistics genera estadísticas de cantidad de entradas por componente
func (a *LogAnalyzer) GetComponentStatistics() map[string]int {
	stats := make(map[string]int)
	for _, entry := range a.Entries {
		stats[entry.Component]++
	}
	return stats
}

// GetHourlyDistribution genera estadísticas de distribución de entradas por hora
func (a *LogAnalyzer) GetHourlyDistribution() map[int]int {
	distribution := make(map[int]int)
	for _, entry := range a.Entries {
		distribution[entry.Timestamp.Hour()]++
	}
	return distribution
}

// GetErrorSummary genera un resumen de errores agrupados por componente
func (a *LogAnalyzer) GetErrorSummary() map[string][]string {
	errorSummary := make(map[string][]string)
	for _, entry := range a.Entries {
		if entry.Level == "ERROR" || entry.Level == "CRITICAL" {
			errorSummary[entry.Component] = append(errorSummary[entry.Component], entry.Message)
		}
	}
	return errorSummary
}

func main() {
	// Definir flags de línea de comandos
	levelFlag := flag.String("level", "", "Filtrar por nivel de log (DEBUG, INFO, WARNING, ERROR, CRITICAL)")
	componentFlag := flag.String("component", "", "Filtrar por componente")
	keywordFlag := flag.String("keyword", "", "Filtrar por palabra clave en el mensaje")
	statsFlag := flag.Bool("stats", false, "Mostrar estadísticas del archivo de log")
	errorsFlag := flag.Bool("errors", false, "Mostrar resumen de errores")

	// Parsear flags
	flag.Parse()

	// Verificar que se proporcionó un archivo de log
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: Debe proporcionar la ruta al archivo de log")
		fmt.Println("Uso: go run main.go [opciones] archivo_log")
		flag.PrintDefaults()
		os.Exit(1)
	}

	logFilePath := args[0]

	// Crear y configurar el analizador de logs
	analyzer := NewLogAnalyzer(logFilePath)

	// Parsear el archivo de log
	err := analyzer.ParseLogFile()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// Aplicar filtros si se especificaron
	filteredEntries := analyzer.Entries

	if *levelFlag != "" {
		filteredEntries = analyzer.FilterByLevel(strings.ToUpper(*levelFlag))
	}

	if *componentFlag != "" {
		var componentFiltered []LogEntry
		for _, entry := range filteredEntries {
			if strings.Contains(strings.ToLower(entry.Component), strings.ToLower(*componentFlag)) {
				componentFiltered = append(componentFiltered, entry)
			}
		}
		filteredEntries = componentFiltered
	}

	if *keywordFlag != "" {
		var keywordFiltered []LogEntry
		for _, entry := range filteredEntries {
			if strings.Contains(strings.ToLower(entry.Message), strings.ToLower(*keywordFlag)) {
				keywordFiltered = append(keywordFiltered, entry)
			}
		}
		filteredEntries = keywordFiltered
	}

	// Mostrar estadísticas si se solicitaron
	if *statsFlag {
		showStatistics(analyzer)
	}

	// Mostrar resumen de errores si se solicitó
	if *errorsFlag {
		showErrorSummary(analyzer)
	}

	// Si no se solicitaron estadísticas ni resumen de errores, mostrar entradas filtradas
	if !*statsFlag && !*errorsFlag {
		showEntries(filteredEntries)
	}
}

// showEntries muestra las entradas de log en la consola
func showEntries(entries []LogEntry) {
	if len(entries) == 0 {
		fmt.Println("No se encontraron entradas que coincidan con los criterios de filtrado.")
		return
	}

	fmt.Printf("Mostrando %d entradas:\n", len(entries))
	fmt.Println(strings.Repeat("-", 80))
	for _, entry := range entries {
		fmt.Println(entry.String())
	}
}

// showStatistics muestra estadísticas del archivo de log
func showStatistics(analyzer *LogAnalyzer) {
	levelStats := analyzer.GetLevelStatistics()
	componentStats := analyzer.GetComponentStatistics()
	hourlyDistribution := analyzer.GetHourlyDistribution()

	fmt.Println("\nEstadísticas del archivo de log:")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("\nDistribución por nivel:")
	// Ordenar niveles para mostrarlos en un orden consistente
	levels := make([]string, 0, len(levelStats))
	for level := range levelStats {
		levels = append(levels, level)
	}
	sort.Strings(levels)

	for _, level := range levels {
		count := levelStats[level]
		color := logLevelColors[level]
		fmt.Printf("  %s%s%s: %d\n", color, level, ColorReset, count)
	}

	fmt.Println("\nDistribución por componente:")
	// Ordenar componentes por cantidad (de mayor a menor)
	type componentCount struct {
		component string
		count     int
	}
	componentCounts := make([]componentCount, 0, len(componentStats))
	for component, count := range componentStats {
		componentCounts = append(componentCounts, componentCount{component, count})
	}
	sort.Slice(componentCounts, func(i, j int) bool {
		return componentCounts[i].count > componentCounts[j].count
	})

	for _, cc := range componentCounts {
		fmt.Printf("  %s: %d\n", cc.component, cc.count)
	}

	fmt.Println("\nDistribución por hora:")
	// Ordenar horas
	hours := make([]int, 0, len(hourlyDistribution))
	for hour := range hourlyDistribution {
		hours = append(hours, hour)
	}
	sort.Ints(hours)

	for _, hour := range hours {
		count := hourlyDistribution[hour]
		bar := strings.Repeat("█", count/2+1)
		fmt.Printf("  %02d:00 - %02d:59: %d %s\n", hour, hour, count, bar)
	}
}

// showErrorSummary muestra un resumen de los errores encontrados
func showErrorSummary(analyzer *LogAnalyzer) {
	errorSummary := analyzer.GetErrorSummary()

	if len(errorSummary) == 0 {
		fmt.Println("No se encontraron errores en el archivo de log.")
		return
	}

	fmt.Println("\nResumen de errores:")
	fmt.Println(strings.Repeat("-", 40))

	// Ordenar componentes alfabéticamente
	components := make([]string, 0, len(errorSummary))
	for component := range errorSummary {
		components = append(components, component)
	}
	sort.Strings(components)

	for _, component := range components {
		errors := errorSummary[component]
		fmt.Printf("\n%s%s%s:\n", ColorCyan, component, ColorReset)
		for i, errorMsg := range errors {
			fmt.Printf("  %d. %s%s%s\n", i+1, ColorRed, errorMsg, ColorReset)
		}
	}
}