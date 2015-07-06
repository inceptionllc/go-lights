package schedule

import "io"

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
