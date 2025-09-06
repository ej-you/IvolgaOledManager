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
var _ repo.SensdataRepoDB = (*SensdataMock)(nil)

// SensdataMock is a mock of repo for sensors data.
// It generates random data values.
type SensdataMock struct{}

// NewMockSensdataRepoDB returns a new instance of SensdataMock.
func NewMockSensdataRepoDB() *SensdataMock {
	return &SensdataMock{}
}

// GetTemperature returns random temperature value.
func (r *SensdataMock) GetTemperature() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetHumidity returns random humidity value.
func (r *SensdataMock) GetHumidity() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetPressure returns random pressure value.
func (r *SensdataMock) GetPressure() (float64, error) {
	return getRandomFloat64(_fromRand2, _toRand2), nil
}

// GetWindSpeed returns random wind speed value.
func (r *SensdataMock) GetWindSpeed() (float64, error) {
	return getRandomFloat64(_fromRand, _toRand), nil
}

// GetWindDirection returns random wind direction value.
func (r *SensdataMock) GetWindDirection() (float64, error) {
	return getRandomFloat64(_fromRand2, _toRand2), nil
}

// getRandomFloat64 returns random float64 in the half-open interval [from, to).
func getRandomFloat64(from, to float64) float64 {
	return rand.Float64()*(to-from) + from //nolint:gosec // crypto/rand is excessive
}
