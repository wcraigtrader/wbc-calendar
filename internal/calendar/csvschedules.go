package calendar

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// ----- Constants -------------------------------------------------------------

var csvHeaders = []string{
	"Date",
	"Time",
	"Event",
	"Prize",
	"Class",
	"Format",
	"Duration",
	"Continuous",
	"GM",
	"Location",
	"Code",
}

func (s *Schedule) CreateCSVDetails(outputDir string) {
	// Date,Time,Event,Prize,Class,Format,Duration,Continuous,GM,Location,Code
	// 2025-07-26,14.0,7 Wonders Demo 1/1,,,,1.0,,"Shea, Elizabeth",EH Annex #2,7WS
	// 2025-07-26,15.0,7 Wonders Heat 1/3,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-26,19.0,7 Wonders Heat 2/3,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-27,9.0,7 Wonders Heat 3/3,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-27,14.0,7 Wonders Quarterfinal,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-27,16.0,7 Wonders Semifinal,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-28,9.0,7 Wonders Final,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Laurel,7WS

	filename := "details.csv"

	log.Printf("Creating CSV details file in %s", outputDir)

	file, err := os.Create(fmt.Sprintf("%s/%s", outputDir, filename))
	if err != nil {
		log.Printf("Error creating CSV file: %v", err)
		return
	}
	defer file.Close()

	w := csv.NewWriter(file)

	if err := w.Write(csvHeaders); err != nil {
		log.Printf("Error writing CSV headers: %v", err)
		return
	}

	for _, e := range s.Everything {
		eventname := e.EventName
		if e.Type == "Tournament" && e.Session != nil {
			eventname += " " + e.Session.Name
		}
		continuous := ""
		if e.Style == "Continuous" {
			continuous = "Y"
		}

		record := []string{
			e.Date.Format("2006-01-02"),
			fmt.Sprintf("%d", e.Start.Hour()+e.Start.Minute()/60.0),
			eventname,
			fmt.Sprintf("%d", e.Prizes),
			e.Class,
			e.Format,
			fmt.Sprintf("%.1f", e.Duration.Hours()),
			continuous,
			e.GM,
			e.Location,
			e.EventCode,
		}
		if err := w.Write(record); err != nil {
			log.Printf("Error writing CSV record for event %s: %v", e, err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Printf("Error flushing CSV writer: %v", err)
		return
	}

	s.AddOtherSchedule(filename, "Detailed schedule [CSV]")
}

func (s *Schedule) CreateCSVSchedule(outputDir string) {
	// Date,Time,Event,Prize,Class,Format,Duration,Continuous,GM,Location,Code
	// 2025-07-26,14:00:00,7 Wonders Demo 1/1,,,,1.0,,"Shea, Elizabeth",EH Annex #2,7WS
	// 2025-07-26,15:00:00,7 Wonders Heat 1/3,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-26,19:00:00,7 Wonders Heat 2/3,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-27,09:00:00,7 Wonders Heat 3/3,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-27,14:00:00,7 Wonders Quarterfinal,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-27,16:00:00,7 Wonders Semifinal,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Grand Ballroom,7WS
	// 2025-07-28,09:00:00,7 Wonders Final,3.0,B,HMW-P,2.0,,"Shea, Elizabeth",Laurel,7WS

	filename := "schedule.csv"

	log.Printf("Creating CSV schedule file in %s", outputDir)

	file, err := os.Create(fmt.Sprintf("%s/%s", outputDir, filename))
	if err != nil {
		log.Printf("Error creating CSV file: %v", err)
		return
	}
	defer file.Close()

	w := csv.NewWriter(file)

	if err := w.Write(csvHeaders); err != nil {
		log.Printf("Error writing CSV headers: %v", err)
		return
	}

	for _, e := range s.Everything {
		eventname := e.EventName
		if e.Type == "Tournament" && e.Session != nil {
			eventname += " " + e.Session.Name
		}
		continuous := ""
		if e.Style == "Continuous" {
			continuous = "Y"
		}

		record := []string{
			e.Date.Format("2006-01-02"),
			e.Start.Format("15:04"),
			eventname,
			fmt.Sprintf("%d", e.Prizes),
			e.Class,
			e.Format,
			fmt.Sprintf("%.1f", e.Duration.Hours()),
			continuous,
			e.GM,
			e.Location,
			e.EventCode,
		}
		if err := w.Write(record); err != nil {
			log.Printf("Error writing CSV record for event %s: %v", e, err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Printf("Error flushing CSV writer: %v", err)
		return
	}

	s.AddOtherSchedule(filename, "Clean Schedule [CSV]")
}
