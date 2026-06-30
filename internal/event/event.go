package event

import (
	"fmt"
	// "log"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	Name   string
	Type   string
	Number int
	Total  int
}

type Event struct {
	hds *Headings
	row []string

	Line   int
	Errors []error

	EventCode string
	EventName string

	Date     time.Time
	Time     time.Time
	Start    time.Time
	Duration time.Duration
	Day      string
	DayCode  string
	Type     string
	Session  *Session
	Prizes   int
	Class    string
	Style    string
	Format   string
	Location string
	GM       string
	Category string
}

var sessionPattern = regexp.MustCompile(`(Demo|Heat|Round) +[(]([0-9]+)/([0-9]+)[)]`)
var sessionDividers = regexp.MustCompile(`[ /()]`)

func NewSession(value string) (*Session, error) {
	number, total := 0, 0

	if slices.Contains(validUniqueSessionTypes, value) {
		return &Session{
			Name:   value,
			Type:   value,
			Number: number,
			Total:  total,
		}, nil
	}

	parts := sessionDividers.Split(value, -1)

	if slices.Contains(validMultipleSessionTypes, parts[0]) {
		if len(parts) == 3 {
			number, _ = strconv.Atoi(parts[1])
			total, _ = strconv.Atoi(parts[2])
			return &Session{
				Name:   value,
				Type:   parts[0],
				Number: number,
				Total:  total,
			}, nil
		}
	}

	return &Session{
		Name:   value,
		Type:   parts[0],
		Number: number,
		Total:  total,
	}, fmt.Errorf("Invalid session format '%s'", value)
}

func (s Session) HasMultiples() bool {
	return slices.Contains(validMultipleSessionTypes, s.Type)
}

func (s Session) String() string {
	return s.Name
}

func NewEvent(headings *Headings, line int, row []string, zone *time.Location, year int) (*Event, error) {
	errors := make([]error, 0)

	event := Event{
		hds:    headings,
		row:    row,
		Errors: errors,
		Line:   line,
	}

	// log.Printf("Processing line %d: %v\n", line, row)

	event.EventCode = event.get("EventCode")
	event.EventName = event.get("EventName")

	event.Date = event.getDate("Date", zone)
	event.Day = event.getRequired("Day", validDays)
	event.DayCode = event.getRequired("DayCode", validDayCodes)
	event.Time = event.getTime("Time", zone)
	event.Duration = event.getDuration("Duration")

	event.Type = event.getOptional("Type", validEventTypes)
	event.Session = event.getSession("Session")
	event.Prizes = event.getInt("Prizes")
	event.Class = event.getOptional("Class", validClasses)
	event.Style = event.getOptional("Style", validStyles)
	event.Format = event.getOptional("Format", validFormats)
	event.Location = event.get("Location")
	event.GM = event.get("GM")
	event.Category = event.getRequired("Category", validCategories)

	event.setStartTime(zone)
	event.Validate()

	if year != 0 && !event.Date.IsZero() && event.Date.Year() != year {
		return nil, fmt.Errorf("event date '%s' does not match specified year '%d'", event.Date.Format("2006-01-02"), year)
	}

	return &event, nil
}

func (e *Event) setStartTime(zone *time.Location) {
	for e.Time.Day() > 1 {
		e.Time = e.Time.Add(-24 * time.Hour)
		e.Date = e.Date.Add(24 * time.Hour)
	}

	e.Start = time.Date(
		e.Date.Year(), e.Date.Month(), e.Date.Day(),
		e.Time.Hour(), e.Time.Minute(), 0, 0, zone)
}

func (e *Event) String() string {
	if e.Type == "Tournament" {
		return fmt.Sprintf("%s (%s) %s @ %s", e.EventName, e.EventCode, e.Session, e.Start.Format("2006-01-02 15:04"))
	}
	return fmt.Sprintf("%s (%s) @ %s", e.EventName, e.Category, e.Start.Format("2006-01-02 15:04"))
}

func (e *Event) addError(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

func (e *Event) get(column string) string {
	idx, ok := e.hds.Columns[column]
	if !ok || idx < 0 || idx >= len(e.row) {
		e.addError(fmt.Errorf("missing column '%s'", column))
		return ""
	}
	value := strings.TrimSpace(e.row[idx])
	if value == "--" || value == "---" {
		return ""
	}
	return value
}

func (e *Event) getInt(column string) int {
	value := e.get(column)
	if value == "" {
		return 0
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	e.addError(fmt.Errorf("Invalid integer '%s' in column '%s'", value, e.hds.orig(column)))
	return 0
}

func (e *Event) getRequired(column string, validValues []string) string {
	value := e.get(column)
	if slices.Contains(validValues, value) {
		return value
	}
	e.addError(fmt.Errorf("Invalid value '%s' in column '%s'", value, e.hds.orig(column)))
	return value
}

func (e *Event) getOptional(column string, validValues []string) string {
	value := e.get(column)
	if value == "" {
		return ""
	}
	return e.getRequired(column, validValues)
}

func (e *Event) getDate(column string, zone *time.Location) time.Time {
	value := e.get(column)
	if dt, err := time.ParseInLocation("1/2/06", value, zone); err == nil {
		return dt
	}
	e.addError(fmt.Errorf("Invalid date '%s' in column '%s'", value, e.hds.orig(column)))
	return time.Time{}
}

func (e *Event) getTime(column string, zone *time.Location) time.Time {
	value := e.get(column)
	if start, err := strconv.ParseFloat(value, 64); err == nil {
		hours := int(start)
		minutes := int((start - float64(hours)) * 60)
		return time.Time{}.Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute)
	}
	if start, err := time.ParseInLocation("3:04", value, zone); err == nil {
		return start
	}
	e.addError(fmt.Errorf("Invalid time '%s' in column '%s'", value, e.hds.orig(column)))
	return time.Time{}
}

func (e *Event) getDuration(column string) time.Duration {
	value := e.get(column)
	if value == "" {
		return time.Duration(0)
	}
	duration, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return time.Duration(duration * float64(time.Hour))
	}
	e.addError(fmt.Errorf("Invalid duration '%s' in column '%s': %w", value, e.hds.orig(column), err))
	return time.Duration(0)
}

func (e *Event) getSession(column string) *Session {
	value := e.get(column)
	if value == "" {
		return nil
	}
	session, err := NewSession(value)
	if err != nil {
		e.addError(fmt.Errorf("Invalid session '%s' in column '%s'", value, e.hds.orig(column)))
	}
	return session
}

func (e *Event) Validate() {
	if e.Date.IsZero() {
		e.addError(fmt.Errorf("missing required Date"))
	}
	if e.Location == "" {
		e.addError(fmt.Errorf("missing required Location"))
	}

	switch e.Type {
	case "Tournament":
		if e.Duration <= 0 {
			e.addError(fmt.Errorf("missing required Duration"))
		}
		if e.EventCode == "" {
			e.addError(fmt.Errorf("missing required EventCode"))
		}
		if e.EventName == "" {
			e.addError(fmt.Errorf("missing required EventName"))
		}
		if e.GM == "" {
			e.addError(fmt.Errorf("missing required GM"))
		}
		if e.Category == "" {
			e.addError(fmt.Errorf("missing required Category"))
		}

		if e.Session == nil {
			e.addError(fmt.Errorf("missing required R/H"))
		} else if e.Session.Type != "Demo" && e.Session.Type != "Draft" {
			if e.Prizes <= 0 {
				e.addError(fmt.Errorf("missing required Prizes"))
			}
			if e.Class == "" {
				e.addError(fmt.Errorf("missing required Class"))
			}
			if e.Style == "" {
				e.addError(fmt.Errorf("missing required Style"))
			}
			if e.Format == "" {
				e.addError(fmt.Errorf("missing required Format"))
			}
		}
	case "Juniors":
		if e.Duration <= 0 {
			e.addError(fmt.Errorf("missing required Duration"))
		}
		if e.EventName == "" {
			e.addError(fmt.Errorf("missing required EventName"))
		}
		if e.GM == "" {
			e.addError(fmt.Errorf("missing required GM"))
		}
		if e.Prizes <= 0 {
			e.addError(fmt.Errorf("missing required Prizes"))
		}
		if e.Class == "" {
			e.addError(fmt.Errorf("missing required Class"))
		}
		if e.Format == "" {
			e.addError(fmt.Errorf("missing required Format"))
		}
	}
}

func (e *Event) Matches(o *Event) bool {
	matches := false
	if e == nil || o == nil || e.Type != o.Type {
		return false
	}
	if e.Type == "Tournament" {
		matches = e.EventCode == o.EventCode && e.Session != nil && o.Session != nil && e.Session.Name == o.Session.Name
	} else {
		matches = e.Category == o.Category && e.EventName == o.EventName
		matches = matches && e.Start.Equal(o.Start)
	}
	return matches
}

func (e *Event) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *Event) ErrorString() string {
	var b strings.Builder
	first := true

	for _, err := range e.Errors {
		if err == nil {
			continue
		}
		if !first {
			b.WriteString("\n\t")
		}
		b.WriteString(err.Error())
		first = false
	}

	return b.String()
}

func (e *Event) Description() string {
	// WAW: A World at War Demo 1/1 Demo 1/1\nPreview: http://www.boardgamers.org/wbc25/previews/WAW.html

	var b strings.Builder

	b.WriteString(e.EventCode)
	b.WriteString(": ")
	b.WriteString(e.EventName)

	if e.Type == "Tournament" && e.Session != nil {
		b.WriteString(" ")
		b.WriteString(e.Session.Name)
	}

	url := e.PreviewURL()
	if url != "" {
		b.WriteString(" Preview: ")
		b.WriteString(url)
	}
	return b.String()
}

func (e *Event) PreviewURL() string {
	// https://www.boardgamers.org/wbc26/previews/waw.html

	if e.Type == "Tournament" {
		year := e.Date.Format("06")
		code := strings.ToLower(e.EventCode)
		return fmt.Sprintf("https://www.boardgamers.org/wbc%s/previews/%s.html", year, code)
	}
	return ""
}

func (e *Event) UID() string {
	// UID:WBC 2025: A World at War Demo 1/1
	// UID:WBC 2025: Shuttle from Pittsburgh Airport
	if e.Type == "Tournament" {
		return fmt.Sprintf("WBC/%s/%s/%s/%s", e.Date.Format("2006"), e.EventCode, e.EventName, e.Session.Name)
	}
	return fmt.Sprintf("WBC/%s/%s/%s/%s/%s", e.Date.Format("2006"), e.Category, e.EventName, e.DayCode, e.Start.Format("15:04"))
}
