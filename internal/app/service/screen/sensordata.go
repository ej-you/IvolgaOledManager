package screen

import (
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// NewSensTemp returns a new instance of template.SensorData for temperature sensor.
func NewSensTemp(active chan bool, btnHandlersReg BtnHandlersRegFunc,
	storage pubsub.Storage) *template.SensorData {

	return template.NewSensorData(active, btnHandlersReg, storage, repo.SensTempKey)
}

// NewSensHumid returns a new instance of template.SensorData for humidity sensor.
func NewSensHumid(active chan bool, btnHandlersReg BtnHandlersRegFunc,
	storage pubsub.Storage) *template.SensorData {

	return template.NewSensorData(active, btnHandlersReg, storage, repo.SensHumidKey)
}

// NewSensPress returns a new instance of template.SensorData for pressure sensor.
func NewSensPress(active chan bool, btnHandlersReg BtnHandlersRegFunc,
	storage pubsub.Storage) *template.SensorData {

	return template.NewSensorData(active, btnHandlersReg, storage, repo.SensPressKey)
}

// NewSensWindSpeed returns a new instance of template.SensorData for wind speed sensor.
func NewSensWindSpeed(active chan bool, btnHandlersReg BtnHandlersRegFunc,
	storage pubsub.Storage) *template.SensorData {

	return template.NewSensorData(active, btnHandlersReg, storage, repo.SensWindSpeedKey)
}

// NewSensWindDir returns a new instance of template.SensorData for wind direction sensor.
func NewSensWindDir(active chan bool, btnHandlersReg BtnHandlersRegFunc,
	storage pubsub.Storage) *template.SensorData {

	return template.NewSensorData(active, btnHandlersReg, storage, repo.SensWindDirKey)
}
