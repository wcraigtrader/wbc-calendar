package calendar

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"wbc-calendar/internal/event"

	ics "github.com/arran4/golang-ical"
)

func (s *Schedule) WriteAllWebCalendars(outputDir string) {
	count := 0

	log.Printf("Creating %d tournament calendars in %s", len(s.Tournaments), outputDir)
	for _, t := range s.Tournaments {
		cal := createVCalendar(t.Calendar)
		filename := fmt.Sprintf("%s/%s.ics", outputDir, t.Code)
		if err := writeCalendarFile(cal, filename); err != nil {
			log.Printf("Error writing calendar file %s: %v", filename, err)
		} else {
			count += 1
		}
	}

	log.Printf("Creating %d other calendars in %s", len(s.Calendars), outputDir)
	for _, c := range s.Calendars {
		cal := createVCalendar(c)
		filename := fmt.Sprintf("%s/%s.ics", outputDir, c.Name)
		if err := writeCalendarFile(cal, filename); err != nil {
			log.Printf("Error writing calendar file %s: %v", filename, err)
		} else {
			count += 1
		}
	}

	log.Printf("Created %d calendar files in %s", count, outputDir)
}

func writeCalendarFile(cal *ics.Calendar, filename string) error {
	var buf bytes.Buffer
	buf.WriteString(cal.Serialize(ics.WithNewLineWindows))
	return os.WriteFile(filename, buf.Bytes(), 0644)
}

func createVCalendar(c *Calendar) *ics.Calendar {
	// BEGIN:VCALENDAR
	// VERSION:2.0
	// PRODID:-//WBC 2025 WAW//ct7//
	// DESCRIPTION:WBC 2025 WAW: A World at War
	// SUMMARY:A World at War
	// URL:http://www.boardgamers.org/wbc25/previews/WAW.html
	// ... events go here ...
	// END:VCALENDAR

	if len(c.Events) == 0 {
		log.Printf("No events in calendar %s, skipping creation", c.Name)
		return nil
	}

	firstEvent := c.Events[0]
	year := firstEvent.Date.Format("2006")
	code := firstEvent.EventCode

	name := c.Name
	if firstEvent.Type == "Tournament" {
		name = code + " " + name
	}

	cal := ics.NewCalendar()
	cal.SetVersion("2.0")
	cal.SetProductId(fmt.Sprintf("-//WBC %s Calendar Generator//ct7//", year))
	cal.SetDescription(fmt.Sprintf("WBC %s %s", year, name))

	url := firstEvent.PreviewURL()
	if url != "" {
		cal.SetUrl(url)
	}

	for _, e := range c.Events {
		cal.AddVEvent(createVEvent(e))
	}

	return cal
}

func createVEvent(e *event.Event) *ics.VEvent {

	// BEGIN:VEVENT
	// SUMMARY:A World at War Demo 1/1
	// DTSTART:20250727T130000Z
	// DURATION:PT1H
	// DTSTAMP:20250712T161756Z
	// UID:WBC 2025: A World at War Demo 1/1
	// COMMENT:{'Code': 'WAW'\, 'Prize': None\, 'Class': None\, 'Format': None\,
	//  'Continuous': ''}
	// CONTACT:Lewis\, Peter
	// DESCRIPTION:WAW: A World at War Demo 1/1 Demo 1/1\nPreview: http://www.boa
	//  rdgamers.org/wbc25/previews/WAW.html
	// LAST-MODIFIED:20250712T161756Z
	// LOCATION:Winterberry
	// URL:http://www.boardgamers.org/wbc25/previews/WAW.html
	// END:VEVENT

	v := ics.NewEvent(e.UID())
	v.SetStartAt(e.Start)
	v.SetDuration(e.Duration)
	v.SetDtStampTime(time.Now())
	v.SetDescription(e.Description())
	v.SetLocation(e.Location)
	v.SetProperty(ics.ComponentPropertyOrganizer, e.GM)

	summary := e.EventName
	if e.Type == "Tournament" && e.Session != nil {
		summary += " " + e.Session.Name
	}
	v.SetSummary(summary)

	url := e.PreviewURL()
	if url != "" {
		v.SetURL(url)
	}

	if e.Type == "Tournament" {
		continuous := ""
		if e.Style == "Continuous" {
			continuous = "Y"
		}
		comment := fmt.Sprintf("{'Code': '%s', 'Prize': %d, 'Class': '%s', 'Format': '%s', 'Continuous': '%s'}", e.EventCode, e.Prizes, e.Class, e.Format, continuous)
		v.SetProperty(ics.ComponentPropertyComment, comment)
	}

	return v
}
