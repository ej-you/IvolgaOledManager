package db

import "math/rand/v2"

var _ StationResultRepoDB = (*mockRepoDB)(nil)

// StationResultRepoDB implementation.
type mockRepoDB struct{}

func NewMockStationResultRepoDB() StationResultRepoDB {
	return &mockRepoDB{}
}

// GetTemperature returns random temperature value in the half-open interval [20.0,30.0).
func (r *mockRepoDB) GetTemperature() (float64, error) {
	temper := rand.Float64()*10 + 20
	return temper, nil
}
