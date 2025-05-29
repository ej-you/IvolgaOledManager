// Package storage provides interfaces with key-value storage for each entity.
// It contains storage manager with all interfaces.
package storage

import (
	"sschmc/internal/app/entity"
	"sschmc/internal/pkg/storage"
)

const (
	_keyAppStatus   = "app-status" // key for app status
	_valueNone      = ""           // value for app status
	_valueGreetings = "greetings"  // value for app status
	_valueMenuMain  = "menu-main"  // value for app status and key for main menu struct
	_valueMenuLevel = "menu-level" // value for app status and key for level menu struct
	_valueMessage   = "message"    // value for app status and key for message struct
)

// AppRepoStorage contains general storage methods for app status.
type AppRepoStorage interface {
	SetNone()
	IsNone() bool

	SetGreetings()
	IsGreetings() bool

	SetMenuMain()
	IsMenuMain() bool

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
	GetLevel() *entity.Menu
	SetLevel(value *entity.Menu)
}

// MessageRepoStorage contains message entity methods.
type MessageRepoStorage interface {
	Get() *entity.Message
	Set(value *entity.Message)
}

// RepoStorageManager contains all storage repos.
type RepoStorageManager struct {
	App     AppRepoStorage
	Menu    MenuRepoStorage
	Message MessageRepoStorage
}

func NewRepoStorageManager(store storage.Storage) *RepoStorageManager {
	return &RepoStorageManager{
		App:     NewAppStorage(store),
		Menu:    NewMenuStorage(store),
		Message: NewMessageStorage(store),
	}
}
