package db

import (
	"math/rand/v2"

	"IvolgaOledManager/internal/app/repo"
)

const (
	_fromRand = 20 // start random interval
	_toRand   = 30 // end random interval
)

var _ repo.SensorRepoDB = (*SensorsMock)(nil)

// SensorsMock is a repo.SensorsRepoDB implementation.
type SensorsMock struct{}

func NewMockStationResultRepoDB() repo.SensorRepoDB {
	return &SensorsMock{}
}

// GetTemperature returns random temperature value.
func (r *SensorsMock) GetTemperature() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetHumidity returns random humidity value.
func (r *SensorsMock) GetHumidity() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetPressure returns random pressure value.
func (r *SensorsMock) GetPressure() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetWindSpeed returns random wind speed value.
func (r *SensorsMock) GetWindSpeed() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetWindDirection returns random wind direction value.
func (r *SensorsMock) GetWindDirection() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// getRandomFloat64 returns random float64 in the half open interval [from, to).
func getRandomFloat64(from, to float64) float64 {
	return rand.Float64()*(from-to) + from //nolint:gosec // crypto/rand is excessive
}
