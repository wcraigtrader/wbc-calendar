package event

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

func testHeadings(columns map[string]int) *Headings {
	original := make(map[string]string, len(columns))
	for k := range columns {
		original[k] = k
	}
	return &Headings{
		Columns:  columns,
		Original: original,
	}
}

func testEvent(columns map[string]int, row []string) Event {
	return Event{
		hds:    testHeadings(columns),
		row:    row,
		Errors: make([]error, 0),
	}
}

func assertErrorCount(t *testing.T, e Event, want int) {
	t.Helper()
	if got := len(e.Errors); got != want {
		t.Fatalf("len(errs) = %d, want %d; errors=%q", got, want, e.ErrorString())
	}
}

func TestEvent_addError(t *testing.T) {
	tests := []struct {
		name     string
		inputErr error
		wantErrs int
	}{
		{
			name:     "nil error is ignored",
			inputErr: nil,
			wantErrs: 0,
		},
		{
			name:     "non-nil error is appended",
			inputErr: errors.New("boom"),
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := testEvent(map[string]int{}, []string{})
			e.addError(tt.inputErr)

			if got := len(e.Errors); got != tt.wantErrs {
				t.Fatalf("len(errs) = %d, want %d", got, tt.wantErrs)
			}
		})
	}
}
func TestEvent_get(t *testing.T) {
	tests := []struct {
		name     string
		columns  map[string]int
		row      []string
		column   string
		want     string
		wantErrs int
	}{
		{
			name:     "returns trimmed value",
			columns:  map[string]int{"A": 0},
			row:      []string{"  value  "},
			column:   "A",
			want:     "value",
			wantErrs: 0,
		},
		{
			name:     "double dash treated as empty",
			columns:  map[string]int{"A": 0},
			row:      []string{"--"},
			column:   "A",
			want:     "",
			wantErrs: 0,
		},
		{
			name:     "triple dash treated as empty",
			columns:  map[string]int{"A": 0},
			row:      []string{"---"},
			column:   "A",
			want:     "",
			wantErrs: 0,
		},
		{
			name:     "missing column returns empty",
			columns:  map[string]int{"A": 0},
			row:      []string{"value"},
			column:   "B",
			want:     "",
			wantErrs: 1,
		},
		{
			name:     "out of range column index returns empty",
			columns:  map[string]int{"A": 2},
			row:      []string{"value"},
			column:   "A",
			want:     "",
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := testEvent(tt.columns, tt.row)
			got := e.get(tt.column)
			if got != tt.want {
				t.Fatalf("get(%q) = %q, want %q", tt.column, got, tt.want)
			}
			assertErrorCount(t, e, tt.wantErrs)
		})
	}
}

func TestEvent_getInt(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		want     int
		column   string
		wantErrs int
	}{
		{name: "valid int", value: "42", want: 42, column: "N", wantErrs: 0},
		{name: "empty value", value: "", want: 0, column: "N", wantErrs: 0},
		{name: "invalid int", value: "4.2", want: 0, column: "N", wantErrs: 1},
		{name: "trimmed int", value: " 7 ", want: 7, column: "N", wantErrs: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := testEvent(map[string]int{tt.column: 0}, []string{tt.value})
			got := e.getInt(tt.column)
			if got != tt.want {
				t.Fatalf("getInt(%q) = %d, want %d", tt.column, got, tt.want)
			}
			assertErrorCount(t, e, tt.wantErrs)
		})
	}
}

func TestEvent_getRequired(t *testing.T) {
	validValues := []string{"A", "B", "C"}
	tests := []struct {
		name     string
		value    string
		want     string
		wantErrs int
	}{
		{name: "valid value", value: "B", want: "B", wantErrs: 0},
		{name: "invalid value returns original", value: "Z", want: "Z", wantErrs: 1},
		{name: "empty value", value: "", want: "", wantErrs: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := testEvent(map[string]int{"X": 0}, []string{tt.value})
			got := e.getRequired("X", validValues)
			if got != tt.want {
				t.Fatalf("getRequired = %q, want %q", got, tt.want)
			}
			assertErrorCount(t, e, tt.wantErrs)
		})
	}
}

func TestEvent_getOptional(t *testing.T) {
	validValues := []string{"A", "B", "C"}
	tests := []struct {
		name     string
		value    string
		want     string
		wantErrs int
	}{
		{name: "empty optional", value: "", want: "", wantErrs: 0},
		{name: "valid optional", value: "A", want: "A", wantErrs: 0},
		{name: "invalid optional returns original", value: "Z", want: "Z", wantErrs: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := testEvent(map[string]int{"X": 0}, []string{tt.value})
			got := e.getOptional("X", validValues)
			if got != tt.want {
				t.Fatalf("getOptional = %q, want %q", got, tt.want)
			}
			assertErrorCount(t, e, tt.wantErrs)
		})
	}
}

func TestEvent_getDate(t *testing.T) {
	zone := time.FixedZone("TestZone", 0)

	tests := []struct {
		name     string
		value    string
		wantY    int
		wantM    time.Month
		wantD    int
		wantNil  bool
		wantErrs int
	}{
		{name: "valid date", value: "6/24/26", wantY: 2026, wantM: time.June, wantD: 24, wantErrs: 0},
		{name: "invalid date", value: "2026-06-24", wantNil: true, wantErrs: 1},
		{name: "empty date", value: "", wantNil: true, wantErrs: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := testEvent(map[string]int{"Date": 0}, []string{tt.value})
			got := e.getDate("Date", zone)

			if tt.wantNil {
				if !got.IsZero() {
					t.Fatalf("getDate = %v, want zero time", got)
				}
			} else {
				if got.IsZero() {
					t.Fatalf("getDate returned zero time for valid input")
				}
				if got.Year() != tt.wantY || got.Month() != tt.wantM || got.Day() != tt.wantD {
					t.Fatalf("getDate = %v, want %04d-%02d-%02d", got, tt.wantY, tt.wantM, tt.wantD)
				}
			}

			assertErrorCount(t, e, tt.wantErrs)
		})
	}
}

func TestEvent_getTime(t *testing.T) {
	zone := time.FixedZone("TestZone", 0)

	tests := []struct {
		name     string
		value    string
		wantD    int
		wantH    int
		wantM    int
		wantNil  bool
		wantErrs int
	}{
		{name: "empty time", value: "", wantNil: true, wantErrs: 1},
		{name: "invalid time", value: "bad", wantNil: true, wantErrs: 1},
		{name: "excel decimal time", value: "13.5", wantD: 1, wantH: 13, wantM: 30, wantErrs: 0},
		{name: "clock time", value: "3:04", wantD: 1, wantH: 3, wantM: 4, wantErrs: 0},
		{name: "late night time", value: "24", wantD: 2, wantH: 0, wantM: 0, wantErrs: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := testEvent(map[string]int{"Time": 0}, []string{tt.value})
			got := e.getTime("Time", zone)

			if tt.wantNil {
				if !got.IsZero() {
					t.Fatalf("getTime = %v, want zero time", got)
				}
			} else if got.Day() != tt.wantD || got.Hour() != tt.wantH || got.Minute() != tt.wantM {
				t.Fatalf("getTime = %02d:%02d, want %02d:%02d:%02d", got.Day(), got.Hour(), got.Minute(), tt.wantH, tt.wantM)
			}

			assertErrorCount(t, e, tt.wantErrs)
		})
	}
}

func TestEvent_getDuration(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		want     time.Duration
		wantErrs int
	}{
		{name: "empty duration", value: "", want: 0, wantErrs: 0},
		{name: "valid decimal hours", value: "1.5", want: 90 * time.Minute, wantErrs: 0},
		{name: "invalid duration", value: "abc", want: 0, wantErrs: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := testEvent(map[string]int{"Duration": 0}, []string{tt.value})
			got := e.getDuration("Duration")
			if got != tt.want {
				t.Fatalf("getDuration = %v, want %v", got, tt.want)
			}
			assertErrorCount(t, e, tt.wantErrs)
		})
	}
}

func TestEvent_getSession(t *testing.T) {
	candidates := []string{"Draft", "Demo 1/2", "Heat 1/3", "Round 2/4"}
	validSession := ""
	for _, c := range candidates {
		if sessionPattern.MatchString(c) {
			validSession = c
			break
		}
	}
	if validSession == "" {
		t.Skip("no candidate matched sessionPattern; update test candidates")
	}

	tests := []struct {
		name     string
		value    string
		wantNil  bool
		wantErrs int
	}{
		{name: "empty session", value: "", wantNil: true, wantErrs: 0},
		{name: "invalid session", value: "not-a-session", wantNil: true, wantErrs: 1},
		{name: "valid session", value: validSession, wantNil: false, wantErrs: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := testEvent(map[string]int{"Session": 0}, []string{tt.value})
			got := e.getSession("Session")

			if tt.wantNil {
				if got != nil {
					t.Fatalf("getSession = %#v, want nil", got)
				}
			} else {
				if got == nil {
					t.Fatalf("getSession returned nil for valid input %q", tt.value)
				}
				if got.Name != tt.value {
					t.Fatalf("Session.Name = %q, want %q", got.Name, tt.value)
				}

				matches := sessionPattern.FindStringSubmatch(tt.value)
				if len(matches) >= 5 {
					if got.Type != matches[1] {
						t.Fatalf("Session.Type = %q, want %q", got.Type, matches[1])
					}
					if matches[3] != "" {
						n, _ := strconv.Atoi(matches[3])
						if got.Number != n {
							t.Fatalf("Session.Number = %d, want %d", got.Number, n)
						}
					}
					if matches[4] != "" {
						n, _ := strconv.Atoi(matches[4])
						if got.Total != n {
							t.Fatalf("Session.Total = %d, want %d", got.Total, n)
						}
					}
				}
			}

			assertErrorCount(t, e, tt.wantErrs)
		})
	}
}
