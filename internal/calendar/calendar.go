package calendar

import (
	"log"
	"sort"

	"wbc-calendar/internal/event"
)

// ----- Data Types -----------------------------------------------------------

type Calendar struct {
	Events []*event.Event
	Name   string
}

type Tournament struct {
	Calendar *Calendar
	Code     string
}

type OtherSchedule struct {
	Filename    string
	Description string
}

type Schedule struct {
	Everything  []*event.Event
	Tournaments map[string]*Tournament
	Calendars   map[string]*Calendar
	Others      []*OtherSchedule

	WBCSite string
	Title   string
}

// ----- Calendar Methods -----------------------------------------------------

func NewCalendar(name string) *Calendar {
	return &Calendar{
		Name:   name,
		Events: make([]*event.Event, 0),
	}
}

func (c *Calendar) AddEvent(e *event.Event) {
	for _, o := range c.Events {
		if o.Matches(e) {
			log.Printf("Duplicate event found in %s: %s and %s", c.Name, o, e)
			return
		}
	}

	c.Events = append(c.Events, e)
}

// ----- Tournament Methods ---------------------------------------------------

func NewTournament(e *event.Event) *Tournament {
	return &Tournament{
		Code:     e.EventCode,
		Calendar: NewCalendar(e.EventName),
	}
}

func (t *Tournament) AddEvent(e *event.Event) {
	t.Calendar.AddEvent(e)
}

// ----- Schedule Methods -----------------------------------------------------

func NewSchedule() *Schedule {

	return &Schedule{
		Everything:  make([]*event.Event, 0),
		Tournaments: make(map[string]*Tournament),
		Calendars:   make(map[string]*Calendar),
		Others:      make([]*OtherSchedule, 0),
	}
}

func (s *Schedule) AddEvent(e *event.Event) {
	s.Everything = append(s.Everything, e)

	if e.Type == "Tournament" {
		if t, ok := s.Tournaments[e.EventCode]; ok {
			t.AddEvent(e)
		} else {
			t := NewTournament(e)
			t.AddEvent(e)
			s.Tournaments[e.EventCode] = t
		}
	} else {
		if c, ok := s.Calendars[e.Category]; ok {
			c.AddEvent(e)
		} else {
			c := NewCalendar(e.Category)
			c.AddEvent(e)
			s.Calendars[e.Category] = c
		}
	}
}

func (s *Schedule) AddOtherSchedule(Filename, Description string) {
	flat := OtherSchedule{
		Filename:    Filename,
		Description: Description,
	}
	s.Others = append(s.Others, &flat)
}

func (s *Schedule) Cleanup() {
	// Placeholder for cleaning up the schedules, removing duplicates, etc.
}

func (s *Schedule) WriteOtherSchedules(outputDir string) {
	s.CreateCSVDetails(outputDir)
	s.CreateCSVSchedule(outputDir)
	s.CreateJsonSchedule(outputDir)
}

func (s *Schedule) SortedTournamentList() []*Tournament {
	list := make([]*Tournament, 0, len(s.Tournaments))
	for _, t := range s.Tournaments {
		list = append(list, t)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Calendar.Name < list[j].Calendar.Name
	})

	// sort.Slice(list, func(i, j int) bool {
	// 	var left, right string
	// 	if list[i] != nil && list[i].Calendar != nil {
	// 		left = list[i].Calendar.Name
	// 	}
	// 	if list[j] != nil && list[j].Calendar != nil {
	// 		right = list[j].Calendar.Name
	// 	}

	// 	// Primary sort: calendar name (case-insensitive)
	// 	li := strings.ToLower(left)
	// 	rj := strings.ToLower(right)
	// 	if li != rj {
	// 		return li < rj
	// 	}

	// 	// Tie-breaker: tournament code (keeps ordering deterministic)
	// 	var lc, rc string
	// 	if list[i] != nil {
	// 		lc = strings.ToLower(list[i].Code)
	// 	}
	// 	if list[j] != nil {
	// 		rc = strings.ToLower(list[j].Code)
	// 	}
	// 	return lc < rc
	// })

	return list
}
