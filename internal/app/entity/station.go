// Package entity contains all app entities.
package entity

// StationResult is a station result data.
type StationResult struct {
	Title      string
	ResultText string
}

// StationResults is a slice with station results.
type StationResults struct {
	Results       []StationResult
	CurrentResult int
}
