// Package storage provides interfaces with key-value storage for each entity.
// It contains storage manager with all interfaces.
package storage

import (
	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/storage"
)

const (
	_keyAppStatus     = "app-status"   // key for app status
	_valueNone        = ""             // value for app status
	_valueGreetings   = "greetings"    // value for app status
	_valueMenuMain    = "menu-main"    // value for app status and key for main menu struct
	_valueMenuLogs    = "menu-logs"    // value for app status and key for menu-logs struct
	_valueMenuLevel   = "menu-level"   // value for app status and key for level menu struct
	_valueMessage     = "message"      // value for app status and key for message struct
	_valueMenuStation = "menu-station" // value for app status and key for station menu struct
	_valueSensor      = "sensor"       // value for app status and key for sensor struct
	_valueSensors     = "sensors"      // key for station sensors slice
)

// AppRepoStorage contains general storage methods for app status.
type AppRepoStorage interface {
	SetNone()
	IsNone() bool

	SetGreetings()
	IsGreetings() bool

	SetMenuMain()
	IsMenuMain() bool

	// first main menu branch

	SetMenuStation()
	IsMenuStation() bool

	SetSensor()
	IsSensor() bool

	// second main menu branch

	SetMenuLogs()
	IsMenuLogs() bool

	SetMenuLevel()
	IsMenuLevel() bool

	SetMessage()
	IsMessage() bool

	IsMenuAny() bool
}

// MenuRepoStorage contains menu entity methods.
type MenuRepoStorage interface {
	GetMain() *entity.Menu
	SetMain(value *entity.Menu)
	// first main menu branch
	GetLogs() *entity.Menu
	SetLogs(value *entity.Menu)
	GetLevel() *entity.Menu
	SetLevel(value *entity.Menu)
	// second main menu branch
	GetStation() *entity.Menu
	SetStation(value *entity.Menu)
}

// MessageRepoStorage contains message entity methods.
type MessageRepoStorage interface {
	Get() *entity.Message
	Set(value *entity.Message)
}

// SensorRepoStorage contains station sensor entity methods.
type SensorRepoStorage interface {
	GetAll() entity.StationSensors
	SetAll(value entity.StationSensors)
	Get() *entity.StationSensor
	Set(value *entity.StationSensor)
}

// RepoStorageManager contains all storage repos.
type RepoStorageManager struct {
	App     AppRepoStorage
	Menu    MenuRepoStorage
	Message MessageRepoStorage
	Sensor  SensorRepoStorage
}

func NewRepoStorageManager(store storage.Storage) *RepoStorageManager {
	return &RepoStorageManager{
		App:     NewAppStorage(store),
		Menu:    NewMenuStorage(store),
		Message: NewMessageStorage(store),
		Sensor:  NewSensorStorage(store),
	}
}
