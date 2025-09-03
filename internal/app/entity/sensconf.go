package entity

import "strings"

const (
	SensconfCtxKey     = "sensorconf"     // key for sensorconf value in context
	SensconfItemCtxKey = "sensorconfItem" // key for sensorconf item value in context

	_inactivePrefix = "#" // prefix for config line of inactive sensor
)

// Sensconf is a slice of sensorconf items.
// It's the result of parsing the sensors' config file.
type Sensconf []*SensconfItem

// SensconfItem is a model for station sensor config line.
// The line looks like `@include "/etc/ssc-station.d/GPS.conf"` for active sensor and
// `# @include "/etc/ssc-station.d/GPS.conf"` for inactive sensor.
type SensconfItem struct {
	// index of sensor
	Idx int
	// config line (e.g. @include "/etc/ssc-station.d/GPS.conf")
	Line string
	// sensor name (filename of sensor's config file without extension, e.g. GPS)
	Name string
	// true if sensor is active
	Active bool
}

// ChangeActive sets active to true if Active is false and vice versa.
// It updates sensor's config line.
func (s *SensconfItem) ChangeActive() {
	if s.Active {
		s.Line = _inactivePrefix + s.Line
	} else {
		s.Line = strings.TrimPrefix(s.Line, _inactivePrefix)
	}
	s.Active = !s.Active
}
