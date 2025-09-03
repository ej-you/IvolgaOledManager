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
var _ repo.SensordataRepoDB = (*SensordataMock)(nil)

// SensordataMock is a mock of repo for sensors data.
// It generates random data values.
type SensordataMock struct{}

// NewMockSensordataRepoDB returns a new instance of SensorsMock.
func NewMockSensordataRepoDB() *SensordataMock {
	return &SensordataMock{}
}

// GetTemperature returns random temperature value.
func (r *SensordataMock) GetTemperature() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetHumidity returns random humidity value.
func (r *SensordataMock) GetHumidity() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetPressure returns random pressure value.
func (r *SensordataMock) GetPressure() (float64, error) {
	return getRandomFloat64(_fromRand2, _toRand2), nil
}

// GetWindSpeed returns random wind speed value.
func (r *SensordataMock) GetWindSpeed() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetWindDirection returns random wind direction value.
func (r *SensordataMock) GetWindDirection() (float64, error) {
	return getRandomFloat64(_fromRand2, _toRand2), nil
}

// getRandomFloat64 returns random float64 in the half-open interval [from, to).
func getRandomFloat64(from, to float64) float64 {
	return rand.Float64()*(to-from) + from //nolint:gosec // crypto/rand is excessive
}
