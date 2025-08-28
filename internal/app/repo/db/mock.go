package db

import (
	rand "math/rand/v2"

	"IvolgaOledManager/internal/app/repo"
)

const (
	_fromRand  = 20  // start random interval
	_toRand    = 30  // end random interval
	_fromRand2 = 200 // start random interval 2
	_toRand2   = 300 // end random interval 2
)

// Ensure sensors mock implements interface.
var _ repo.SensorRepoDB = (*SensorsMock)(nil)

// SensorsMock is a mock of repo for sensors data.
// It generates random data values.
type SensorsMock struct{}

// NewMockStationResultRepoDB returns a new instance of SensorsMock.
func NewMockStationResultRepoDB() *SensorsMock {
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
	return getRandomFloat64(_fromRand2, _toRand2), nil
}

// GetWindSpeed returns random wind speed value.
func (r *SensorsMock) GetWindSpeed() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetWindDirection returns random wind direction value.
func (r *SensorsMock) GetWindDirection() (float64, error) {
	return getRandomFloat64(_fromRand2, _toRand2), nil
}

// getRandomFloat64 returns random float64 in the half-open interval [from, to).
func getRandomFloat64(from, to float64) float64 {
	return rand.Float64()*(to-from) + from //nolint:gosec // crypto/rand is excessive
}
