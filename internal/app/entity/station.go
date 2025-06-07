package entity

import "strings"

const _inactivePrefix = "#" // prefix for inactive sensor line

// StationSensor is a model for station sensor config line.
// The line looks like `@include "/etc/ssc-station.d/GPS.conf"` for active sensor or
// `# @include "/etc/ssc-station.d/GPS.conf"` for inactive sensor.
type StationSensor struct {
	// config line (e.g. @include "/etc/ssc-station.d/GPS.conf")
	Line string
	// sensor name (filename of sensor config file without extension, e.g. GPS)
	Name string
	// true if sensor is active
	Active bool
}

// ChangeActive sets active to true if Active is false and
// sets active to false if Active is true.
// It updates sensor line.
func (s *StationSensor) ChangeActive() {
	if s.Active {
		s.Line = _inactivePrefix + s.Line
	} else {
		s.Line = s.Line[1:]
	}
	s.Active = !s.Active
}

// StationSensors contains sensor instances. It's result of config file parsing.
type StationSensors []*StationSensor

// CollectAll collects all sensor lines to a slice of bytes and returns it.
func (s StationSensors) CollectAll(prefix, suffix string) []byte {
	var builder strings.Builder
	builder.WriteString(prefix)
	// collect sensor lines
	for _, sensor := range s {
		builder.WriteString("\n")
		builder.WriteString(sensor.Line)
	}
	builder.WriteString(suffix)
	return []byte(builder.String())
}
