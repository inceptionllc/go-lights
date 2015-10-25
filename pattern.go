package lights

import (
	"errors"
	"image/color"
	"strconv"
	"strings"
	"time"
)

// Pattern captures all data needed for a light pattern.
type Pattern struct {
	ID    string
	Loops int
	Slots []*Slot
}

// NewPattern creates a pattern from a pattern specification string.
func NewPattern(pattern string) (*Pattern, error) {
	p := &Pattern{}
	var err error
	parts := strings.Split(pattern, "|")
	if len(parts) == 0 {
		return nil, errors.New("Empty pattern " + pattern)
	}

	// Determine the number of loops will be of the form `:ID[:loops]`
	headers := strings.Split(parts[0], ":")
	p.Loops = -1
	switch len(headers) {
	case 3:
		loop := strings.TrimSpace(headers[2])
		if len(loop) > 0 {
			p.Loops, err = strconv.Atoi(loop)
			if err != nil {
				return nil, errors.New("Loop count was not an integer")
			}
		}
		fallthrough
	case 2:
		p.ID = strings.TrimSpace(headers[1])
	}

	if len(parts) == 1 {
		p.Slots = []*Slot{}
		return p, nil // No slots - we don't need to do anything.
	}

	for _, slot := range parts[1:] {
		// Each slot is a color, fade, hold, transition separated by ","
		s, err := NewSlot(slot)
		if err != nil {
			return nil, err
		}
		p.Slots = append(p.Slots, s)
	}
	return p, nil
}

// Slot captures the information about a single slot in a pattern.
type Slot struct {
	Color      color.Color
	Fade       time.Duration
	Hold       time.Duration
	Transition string
}

// NewSlot creates a slot from a slot specification.
func NewSlot(slot string) (s *Slot, err error) {
	items := strings.Split(slot, ",")
	s = &Slot{Transition: "ease"}
	switch len(items) {
	case 0:
		// Empty slot found - ignore
		return nil, errors.New("No slot information found")
	case 4:
		// color, fade, hold, and transition
		value := strings.TrimSpace(items[3])
		if len(value) > 0 {
			s.Transition = value
		}
		fallthrough
	case 3:
		value := strings.TrimSpace(items[2])
		if len(value) > 0 {
			s.Hold, err = time.ParseDuration(value)
			if err != nil {
				return nil, err
			}
		}
		fallthrough
	case 2:
		value := strings.TrimSpace(items[1])
		if len(value) > 0 {
			s.Fade, err = time.ParseDuration(value)
			if err != nil {
				return nil, err
			}
		}
		fallthrough
	case 1:
		value := strings.TrimSpace(items[0])
		if len(value) > 0 {
			s.Color, err = ParseColorCode(value)
			if err != nil {
				return nil, err
			}
		}
	}
	return
}
