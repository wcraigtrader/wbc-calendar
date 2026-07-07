package main

import (
	"log"
	"os"
	"path/filepath"
	"slices"

	"wbc-calendar/internal/calendar"
	"wbc-calendar/internal/config"
	"wbc-calendar/internal/excel"
	"wbc-calendar/internal/website"
)

func main() {
	config := config.ParseCommandLine()

	// Create Excel reader
	reader, err := excel.NewExcelReader(config.ExcelFilePath)
	if err != nil {
		log.Fatalf("Error creating Excel reader: %v", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing Excel reader: %v", err)
		}
	}()

	// Get available sheets
	sheets := reader.GetSheetNames()

	if config.Verbose || config.ListOnly {
		log.Printf("Available sheets: %v\n", sheets)
	}

	// If list-only mode, exit after showing sheets
	if config.ListOnly {
		return
	}

	data, err := reader.ReadSheet(config.SheetName, config.Zone, config.Year)
	if err != nil {
		log.Fatalf("Error reading sheet '%s': %v", config.SheetName, err)
	}

	log.Printf("Read %d events from sheet '%s'\n", len(data), config.SheetName)

	if err := CreateOutputDirectory(config); err != nil {
		log.Fatalf("Error creating output directory '%s': %v", config.OutputDirectory, err)
	}

	schedule := calendar.NewSchedule()
	for _, event := range data {
		if len(config.Include) > 0 {
			if !event.IsTournament() {
				continue
			}
			if !slices.Contains(config.Include, event.EventCode) {
				continue
			}
		}
		if len(config.Exclude) > 0 {
			if slices.Contains(config.Exclude, event.EventCode) {
				continue
			}
		}
		schedule.AddEvent(event)
	}

	schedule.Cleanup()
	schedule.WriteAllWebCalendars(config.OutputDirectory)
	schedule.WriteOtherSchedules(config.OutputDirectory)

	website.CreateWebsite(schedule, config)

	log.Printf("Schedule created with %d tournaments and %d calendars\n", len(schedule.Tournaments), len(schedule.Calendars))
}

func CreateOutputDirectory(config *config.Config) error {
	if err := os.MkdirAll(config.OutputDirectory, 0o755); err != nil {
		return err
	}

	if config.Clean {
		entries, err := os.ReadDir(config.OutputDirectory)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			for _, entry := range entries {
				target := filepath.Join(config.OutputDirectory, entry.Name())
				if err := os.RemoveAll(target); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
