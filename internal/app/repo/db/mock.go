package db

import "math/rand/v2"

const (
	_fromRand     = 20 // start random interval
	_intervalRand = 10 // len of random interval
)

var _ StationResultRepoDB = (*mockRepoDB)(nil)

// StationResultRepoDB implementation.
type mockRepoDB struct{}

func NewMockStationResultRepoDB() StationResultRepoDB {
	return &mockRepoDB{}
}

// GetTemperature returns random temperature value in the half-open interval [20.0,30.0).
func (r *mockRepoDB) GetTemperature() (float64, error) {
	temper := rand.Float64()*_intervalRand + _fromRand //nolint:gosec // crypto/rand is excessive
	return temper, nil
}
