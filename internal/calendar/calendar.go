package calendar

import (
	"fmt"
	"log"
	"sort"

	"wbc-calendar/internal/event"
)

// ----- Data Types -----------------------------------------------------------

type Calendar struct {
	Events   []*event.Event
	Filename string
	Name     string
	Errors     []error
}

type Tournament struct {
	Calendar  *Calendar
	Code      string
	Sessions map[string][]int
	Totals   map[string]int
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

func (c * Calendar) addError(err error) {
	if err != nil {
		c.Errors = append(c.Errors, err)
	}
}

func (c *Calendar) Validate() {
	sort.SliceStable(c.Events, func(i, j int) bool {
		return c.Events[i].Start.Before(c.Events[j].Start)
	})

	for i := 1; i < len(c.Events); i++ {
		if c.Events[i].Start.Before(c.Events[i-1].Start) {
			c.addError(fmt.Errorf("%s session %s starts before %s", c.Name, c.Events[i].Session.Name, c.Events[i-1].Session.Name))
		}
	}
}

func (c *Calendar) HasErrors() bool {
	count := len(c.Errors)
	for _, e := range c.Events {
		count += len(e.Errors)
	}
	return count > 0
}

func (c *Calendar) ReportErrors() {
	if c.HasErrors() {
		log.Printf("Errors found in %s:", c.Name)
		for _, err := range c.Errors {
			log.Printf("\t%s", err)
		}
		for _, e := range c.Events {
			for _, err := range e.Errors {
				log.Printf("\tRow %d: %s %s", e.Line, e, err)
			}
		}
	}
}

// ----- Tournament Methods ---------------------------------------------------

func NewTournament(e *event.Event) *Tournament {
	return &Tournament{
		Code:     e.EventCode,
		Calendar: NewCalendar(e.EventName),
		Sessions: map[string][]int{
			"Demo":  {},
			"Heat":  {},
			"Round": {},
		},
		Totals:   map[string]int {"Demo": 0, "Heat": 0, "Round": 0},
	}
}

func (t *Tournament) AddEvent(e *event.Event) {
	t.Calendar.AddEvent(e)

	if e.Session != nil && e.Session.HasMultiples() {
		t.Sessions[e.Session.Type] = append(t.Sessions[e.Session.Type], e.Session.Number)
		if t.Totals[e.Session.Type] == 0 {
			t.Totals[e.Session.Type] = e.Session.Total
		} else if t.Totals[e.Session.Type] != e.Session.Total {
			t.Calendar.addError(fmt.Errorf("inconsistent total sessions for %s: %d vs %d", e.Session.Type, t.Totals[e.Session.Type], e.Session.Total))
		}
		if e.Session.Number > t.Totals[e.Session.Type] {
			t.Calendar.addError(fmt.Errorf("session number %d exceeds total %d for %s", e.Session.Number, t.Totals[e.Session.Type], e.Session.Type))
		}
	}

}

func (t *Tournament) Validate() {
	t.Calendar.Validate()

	for sessionType, numbers := range t.Sessions {
		if len(numbers) > 0 {
			expectedTotal := t.Totals[sessionType]
			if len(numbers) != expectedTotal {
				t.Calendar.addError(fmt.Errorf("expected %d %s sessions, but found %d", expectedTotal, sessionType, len(numbers)))
			}
		}
	}
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

	if e.IsTournament() {
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
	for _, t := range s.Tournaments {
		t.Validate()
		t.Calendar.ReportErrors()
	}
	for _, c := range s.Calendars {
		c.Validate()
		c.ReportErrors()
	}
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

	return list
}
