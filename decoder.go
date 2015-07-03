package schedule

import (
	"fmt"
	"image/color"
	"io"
	"strconv"
)

// ParseColorCode parses a color code in either #RGB or #RRGGBB formats
// following CSS color specifications. It returns the color found or
// an error if the color is not a valid color code.
func ParseColorCode(colorCode string) (color.Color, error) {
	switch len(colorCode) {
	case 4:
		r, err := strconv.ParseInt(string(colorCode[1])+string(colorCode[1]), 16, 0)
		if err != nil {
			return nil, err
		}
		g, err := strconv.ParseInt(string(colorCode[2])+string(colorCode[2]), 16, 0)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseInt(string(colorCode[3])+string(colorCode[3]), 16, 0)
		if err != nil {
			return nil, err
		}
		// Should be uint8s, just cast from int64
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(0)}, nil
	case 7:
		r, err := strconv.ParseInt(string(colorCode[1])+string(colorCode[2]), 16, 0)
		if err != nil {
			return nil, err
		}
		g, err := strconv.ParseInt(string(colorCode[3])+string(colorCode[4]), 16, 0)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseInt(string(colorCode[5])+string(colorCode[6]), 16, 0)
		if err != nil {
			return nil, err
		}
		// Should be uint8s, just cast from int64
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(0)}, nil
	default:
		return nil, fmt.Errorf("Color code must be 4 or 6 characters - found %d", len(colorCode))
	}
}

// ParseScheduleFile reads in a json encoded schedule and produces a Config.
func ParseScheduleFile(r io.ReadCloser) ([]*Schedule, error) {
	return nil, nil
	/*
		schedules := make([]*schedule.Schedule, 0, 1)
		patterns := make([]*schedule.Pattern, 0, 1)

		var jsonObj SObj

		decoder := json.NewDecoder(r)
		err := decoder.Decode(&jsonObj)
		if err != nil {
			return schedules, fmt.Errorf("Error parsing schedule obj JSON: %s", err)
		}

		for _, pObj := range jsonObj.Patterns {
			p := schedule.NewPattern()
			p.name = pObj.Name
			for _, pSlotObj := range pObj.Slots {
				col, err := ParseColorCode(pSlotObj.Color)
				if err != nil {
					return schedules, fmt.Errorf("Error parsing color code: %s", pSlotObj.Color)
				}
				p.AddPatternSlot(schedule.NewPatternSlot(col, pSlotObj.Hold, pSlotObj.Transition))
			}
			patterns = append(patterns, p)
		}

		for _, sObj := range jsonObj.Schedules {
			s := schedule.NewSchedule()
			p, err := schedule.findPattern(patterns, sObj.Pattern)
			if err != nil {
				return schedules, err
			}

			s.pattern = p
			s.mon = sObj.Mon
			s.tue = sObj.Tue
			s.wed = sObj.Wed
			s.thu = sObj.Thu
			s.fri = sObj.Fri
			s.sat = sObj.Sat
			s.sun = sObj.Sun

			schedules = append(schedules, s)
		}

		return schedules, nil
	*/
}
