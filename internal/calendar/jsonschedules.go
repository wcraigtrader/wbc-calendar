package calendar

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type JsonEvent struct {
	Code       string  `json:"Code"`
	Prize      *int    `json:"Prize"`
	Class      *string `json:"Class"`
	Format     *string `json:"Format"`
	Continuous string  `json:"Continuous"`
	Event      string  `json:"Event"`
	GM         string  `json:"GM"`
	Location   string  `json:"Location"`
	Date       string  `json:"Date"`
	Time       float64 `json:"Time"`
	Duration   float64 `json:"Duration"`
}

func (s *Schedule) CreateJsonSchedule(outputDir string) {
	// [
	//   {"Code": "7WS", "Prize": null, "Class": null, "Format": null, "Continuous": "", "Event": "7 Wonders Demo 1/1", "GM": "Shea, Elizabeth", "Location": "EH Annex #2", "Date": "07/26/2025", "Time": 14, "Duration": 1.0},
	//   {"Code": "7WS", "Prize": 3.0, "Class": "B", "Format": "HMW-P", "Continuous": "", "Event": "7 Wonders Heat 1/3", "GM": "Shea, Elizabeth", "Location": "Grand Ballroom", "Date": "07/26/2025", "Time": 15, "Duration": 2.0},
	//   {"Code": "7WS", "Prize": 3.0, "Class": "B", "Format": "HMW-P", "Continuous": "", "Event": "7 Wonders Heat 2/3", "GM": "Shea, Elizabeth", "Location": "Grand Ballroom", "Date": "07/26/2025", "Time": 19, "Duration": 2.0},
	//   {"Code": "7WS", "Prize": 3.0, "Class": "B", "Format": "HMW-P", "Continuous": "", "Event": "7 Wonders Heat 3/3", "GM": "Shea, Elizabeth", "Location": "Grand Ballroom", "Date": "07/27/2025", "Time": 9, "Duration": 2.0},
	//   {"Code": "7WS", "Prize": 3.0, "Class": "B", "Format": "HMW-P", "Continuous": "", "Event": "7 Wonders Quarterfinal", "GM": "Shea, Elizabeth", "Location": "Grand Ballroom", "Date": "07/27/2025", "Time": 14, "Duration": 2.0},
	//   {"Code": "7WS", "Prize": 3.0, "Class": "B", "Format": "HMW-P", "Continuous": "", "Event": "7 Wonders Semifinal", "GM": "Shea, Elizabeth", "Location": "Grand Ballroom", "Date": "07/27/2025", "Time": 16, "Duration": 2.0},
	//   {"Code": "7WS", "Prize": 3.0, "Class": "B", "Format": "HMW-P", "Continuous": "", "Event": "7 Wonders Final", "GM": "Shea, Elizabeth", "Location": "Laurel", "Date": "07/28/2025", "Time": 9, "Duration": 2.0},
	// ]

	filename := "details.json"

	log.Printf("Creating JSON details file in %s", outputDir)

	file, err := os.Create(fmt.Sprintf("%s/%s", outputDir, filename))
	if err != nil {
		log.Printf("Error creating JSON file: %v", err)
		return
	}
	defer file.Close()

	file.WriteString("[\n")
	for i, e := range s.Everything {
		eventname := e.EventName
		if e.Type == "Tournament" && e.Session != nil {
			eventname += " " + e.Session.Name
		}
		continuous := ""
		if e.Style == "Continuous" {
			continuous = "Y"
		}

		record := JsonEvent{
			Code:       e.EventCode,
			Prize:      &e.Prizes,
			Class:      &e.Class,
			Format:     &e.Format,
			Continuous: continuous,
			Event:      eventname,
			GM:         e.GM,
			Location:   e.Location,
			Date:       e.Date.Format("01/02/2006"),
			Time:       float64(e.Start.Hour()) + float64(e.Start.Minute())/60.0,
			Duration:   e.Duration.Hours(),
		}

		jsonData, err := json.Marshal(record)
		if err != nil {
			log.Printf("Error marshaling JSON for event %s: %v", e, err)
			continue
		}

		file.WriteString("  ")
		file.Write(jsonData)
		if i < len(s.Everything)-1 {
			file.WriteString(",\n")
		} else {
			file.WriteString("\n")
		}
	}
	file.WriteString("]\n")

	log.Printf("JSON schedule created: %s/%s", outputDir, filename)

	s.AddOtherSchedule(filename, "Detailed schedule [JSON]")

}
