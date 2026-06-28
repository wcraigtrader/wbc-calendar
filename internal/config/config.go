package config

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	ExcelFilePath  string
	CalendarOutput string
	SheetName      string
	Location       string

	Clean    bool
	ListOnly bool
	Verbose  bool

	Zone *time.Location
	Year int

	Include []string
	Exclude []string
}

func ParseCommandLine() *Config {
	var config Config

	var include string
	var exclude string

	flag.StringVar(&config.ExcelFilePath, "file", "", "Path to Excel file (required)")
	flag.StringVar(&config.ExcelFilePath, "f", "", "Path to Excel file (short form)")
	flag.StringVar(&config.CalendarOutput, "output", "build", "Output directory for calendar files")
	flag.StringVar(&config.CalendarOutput, "o", "build", "Output directory for calendar files (short form)")
	flag.StringVar(&config.SheetName, "sheet", "App Version", "Sheet name to read")
	flag.StringVar(&config.SheetName, "s", "App Version", "Sheet name to read (short form)")
	flag.StringVar(&config.Location, "zone", "America/New_York", "Time zone location (default: America/New_York)")
	flag.StringVar(&config.Location, "z", "America/New_York", "Time zone location (short form)")
	flag.StringVar(&include, "include", "", "Comma-separated list of tournament codes to include (default: all)")
	flag.StringVar(&include, "i", "", "Comma-separated list of tournament codes to include (short form)")
	flag.StringVar(&exclude, "exclude", "", "Comma-separated list of tournament codes to exclude (default: none)")
	flag.StringVar(&exclude, "e", "", "Comma-separated list of tournament codes to exclude (short form)")
	flag.IntVar(&config.Year, "year", 0, "Calendar year")
	flag.IntVar(&config.Year, "y", 0, "Calendar year (short form)")
	flag.BoolVar(&config.Clean, "clean", false, "Clean output directory before writing")
	flag.BoolVar(&config.Clean, "c", false, "Clean output directory before writing (short form)")
	flag.BoolVar(&config.ListOnly, "list", false, "List available sheets and exit")
	flag.BoolVar(&config.ListOnly, "l", false, "List available sheets and exit (short form)")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&config.Verbose, "v", false, "Enable verbose output (short form)")

	// Custom usage message
	flag.Usage = func() {
		log.Printf("Usage: %s [OPTIONS]\n\n", os.Args[0])
		log.Printf("Excel spreadsheet reader for WBC Calendar project.\n\n")
		log.Printf("Options:\n")
		flag.PrintDefaults()
		log.Printf("\nExamples:\n")
		log.Printf("  %s -y 2026 -f data.xlsx\n", os.Args[0])
		log.Printf("  %s -y 2026 -f data.xlsx -l\n", os.Args[0])
		log.Printf("  %s -y 2026 -f data.xlsx -s \"Sheet1\"\n", os.Args[0])
	}

	flag.Parse()

	// Validate required arguments
	errors_detected := false

	if zone, err := time.LoadLocation(config.Location); err != nil {
		log.Printf("Error: Invalid time zone location '%s': %v", config.Location, err)
		errors_detected = true
	} else {
		config.Zone = zone
	}

	if config.Year == 0 {
		config.Year = time.Now().In(config.Zone).Year()
	}

	if config.Year < 2000 {
		config.Year += 2000
	}

	if config.CalendarOutput == "" {
		log.Printf("Error: Output directory for calendar files is required\n")
		errors_detected = true
	}

	if config.ExcelFilePath == "" {
		log.Printf("Error: Excel file path is required\n")
		errors_detected = true
	}

	// Check if file exists
	if _, err := os.Stat(config.ExcelFilePath); os.IsNotExist(err) {
		log.Printf("Error: File '%s' does not exist", config.ExcelFilePath)
		errors_detected = true
	}

	if include != "" {
		config.Include = strings.Split(include, ",")
	}

	if exclude != "" {
		config.Exclude = strings.Split(exclude, ",")
	}

	if errors_detected {
		flag.Usage()
		os.Exit(1)
	}

	if config.Verbose {
		log.Printf("Configuration:")
		log.Printf("  Output: %s", config.CalendarOutput)
		log.Printf("  File:   %s", config.ExcelFilePath)
		log.Printf("  Sheet:  %s", config.SheetName)
		log.Printf("  Zone:   %s", config.Zone)
		log.Printf("  Year:   %d", config.Year)
		log.Printf("  Clean:  %t", config.Clean)
		log.Printf("  List Only: %t", config.ListOnly)
		log.Printf("  Verbose: %t", config.Verbose)
	}

	return &config
}
