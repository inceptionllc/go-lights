package schedule

import (
	"image/color"
	"time"
)

// NewCrontab creates a new crontab ready to add specs
func NewCrontab() *Crontab {
	return &Crontab{}
}

// Crontab describes a complete set of cron entries
type Crontab struct {
	Specs []*Spec
}

// Add a new spec to the crontab
func (c *Crontab) Add(expr, action string) {
	c.Specs = append(c.Specs, &Spec{Expression: expr, Action: action})
}

// Spec describes a cron spec entry
type Spec struct {
	Expression string
	Action     string
}

// --------------------- Legacy code below ------------------------

// PatternSlot represents a light state in the overall pattern.
type PatternSlot struct {
	Color      color.Color   `json:"color"`
	Transition time.Duration `json:"transition"`
	Hold       time.Duration `json:"hold"`
}

// NewPatternSlot creates a new PatternSlot for the provided parameters.
func NewPatternSlot(col color.Color, hold, transition time.Duration) *PatternSlot {
	return &PatternSlot{Color: col, Hold: hold, Transition: transition}
}

// Pattern represents a light display pattern composed of one or more
// PatternSlots.
type Pattern struct {
	ID    string         `json:"_id"`
	Name  string         `json:"name"`
	Slots []*PatternSlot `json:"slots"`
}

// NewPattern creates a new empty.
func NewPattern() *Pattern {
	return &Pattern{Slots: make([]*PatternSlot, 0, 20)}
}

// TODO(stephen): convert the Mon, Tue, etc to a crontab pattern.
// There are several cron libraries for go that should make parsing cron patterns easy.

// Schedule determines the pattern of schedules to display.
type Schedule struct {
	ID      string   `json:"_id"`
	Name    string   `json:"name"`
	Pattern *Pattern `json:"pattern"`
	Mon     bool     `json:"mon"`
	Tue     bool     `json:"tue"`
	Wed     bool     `json:"wed"`
	Thu     bool     `json:"thu"`
	Fri     bool     `json:"fri"`
	Sat     bool     `json:"sat"`
	Sun     bool     `json:"sun"`
}

// Config represents a complete configuration of a light controller node
// including patterns available for display and schedules that are active
// on the node.
type Config struct {
	Patterns  map[string]*Pattern  `json:"patterns"`  // Patterns by ID
	Schedules map[string]*Schedule `json:"schedules"` // Schedules by ID
}

// NewConfig creates a new configuration with no patterns or schedules.
func NewConfig() *Config {
	return &Config{map[string]*Pattern{}, map[string]*Schedule{}}
}

// PatternByName locates a pattern by name and returns a 'ok' flag
// if a pattern with the given name can't be found.
func (c *Config) PatternByName(name string) (*Pattern, bool) {
	for _, p := range c.Patterns {
		if p.Name == name {
			return p, true
		}
	}
	return nil, false
}

// AddPattern adds a pattern to the configuration.
func (c *Config) AddPattern(p *Pattern) {
	c.Patterns[p.ID] = p
}
