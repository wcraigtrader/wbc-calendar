package main

import (
	"log"
	"os"

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

	if err := os.MkdirAll(config.CalendarOutput, 0o755); err != nil {
		log.Fatalf("Error creating output directory '%s': %v", config.CalendarOutput, err)
	}

	data, err := reader.ReadSheet(config.SheetName, config.Zone, config.Year)
	if err != nil {
		log.Fatalf("Error reading sheet '%s': %v", config.SheetName, err)
	}

	log.Printf("Read %d events from sheet '%s'\n", len(data), config.SheetName)

	schedule := calendar.NewSchedule()
	for _, event := range data {
		schedule.AddEvent(event)
	}

	schedule.Cleanup()
	schedule.WriteAllWebCalendars(config.CalendarOutput)
	schedule.WriteOtherSchedules(config.CalendarOutput)
	
	website.CreateWebsite(schedule, config.Year, config.CalendarOutput)

	log.Printf("Schedule created with %d tournaments and %d calendars\n", len(schedule.Tournaments), len(schedule.Calendars))
}
