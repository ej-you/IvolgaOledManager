// Package entity contains all app entities.
package entity

// StationResult is a station result data.
type SensorData struct {
	// measurement name to output
	Title string
	// sensor data value to output
	Data string

	// raw sensor data
	RawData float64
}

// TODO: implement Renderer
