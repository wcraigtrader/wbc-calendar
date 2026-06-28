package event

import (
	"fmt"
)

type Headings struct {
	Columns  map[string]int
	Original map[string]string
}

func CreateHeadings(row []string) (*Headings, error) {
	columns := make(map[string]int)
	original := make(map[string]string)

	var unrecognized []string
	for i, cell := range row {
		if canonical, ok := columnNames[cell]; ok {
			columns[canonical] = i
			original[canonical] = cell
		} else {
			unrecognized = append(unrecognized, cell)
		}
	}

	required := make(map[string]bool)
	for _, canonical := range columnNames {
		required[canonical] = true
	}

	var missing []string
	for canonical := range required {
		if _, ok := columns[canonical]; !ok {
			missing = append(missing, canonical)
		}
	}

	hm := Headings{Columns: columns, Original: original}
	if len(missing) > 0 || len(unrecognized) > 0 {
		return &hm, fmt.Errorf(
			"invalid header row, missing required columns: %v; unrecognized columns: %v",
			missing, unrecognized,
		)
	}

	return &hm, nil
}

func (h *Headings) orig(canonical string) string {
	if orig, ok := h.Original[canonical]; ok {
		return orig
	}
	return canonical
}
