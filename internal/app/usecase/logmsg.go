package usecase

import (
	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"context"
	"fmt"
	"strconv"
	"strings"
)

var (
	_levelsAmount = 6 // amount of log levels
	_levelName    = map[int]string{
		0: "trace",
		1: "debug",
		2: "info",
		3: "warn",
		4: "error",
		5: "fatal",
	}
)

// Ensure sensors' config usecase implements interface.
var _ LogMsgUsecase = (*LogMsgUC)(nil)

// LogMsgUC represents a usecase for log messages.
type LogMsgUC struct {
	logMsgRepoDB repo.LogMsgRepoDB
}

// NewLogMsgUsecase returns a new instance of LogMsgUC.
func NewLogMsgUsecase(logMsgRepoDB repo.LogMsgRepoDB) *LogMsgUC {
	return &LogMsgUC{
		logMsgRepoDB: logMsgRepoDB,
	}
}

// GetLevelCount returns menu with level numbers and amount of messages with this levels.
func (l *LogMsgUC) GetLevelCountAsMenu() (*entity.Menu, error) {
	// get levels count from DB
	levels, err := l.logMsgRepoDB.GetLevelCount()
	if err != nil {
		return nil, fmt.Errorf("get level count as menu: %w", err)
	}

	// create slice for all levels count and fill it with zeros
	allLevels := make([]entity.LogMsgLvlCount, _levelsAmount)
	for idx := range allLevels {
		allLevels[idx].Level = idx
	}
	// update level count values according to gotten DB data
	for _, lvl := range levels {
		allLevels[lvl.Level].Count = lvl.Count
	}

	// create empty menu
	menu := &entity.Menu{
		Title: "Уровни логов",
		Items: make([]*entity.MenuItem, 0, len(_levelName)),
	}
	// append menu items
	var itemName string
	var itemCtx context.Context
	for _, levelCount := range allLevels {
		itemName = fmt.Sprintf("%s (%d)", _levelName[levelCount.Level], levelCount.Count)
		itemCtx = context.WithValue(context.Background(), entity.LogLvlCountCtxKey, levelCount)
		menu.Items = append(menu.Items, entity.NewMenuItem(itemCtx, itemName))
	}
	return menu, nil
}

// GetWithLevel returns menu with messages with given level ordered by created datetime.
func (l *LogMsgUC) GetWithLevelAsMenu(level string) (*entity.Menu, error) {
	// get messages with selected level from DB
	levelMessages, err := l.logMsgRepoDB.GetWithLevel(level)
	if err != nil {
		return nil, fmt.Errorf("get with level as menu: %w", err)
	}

	// convert level from string into int
	levelInt, err := strconv.Atoi(level)
	if err != nil {
		return nil, fmt.Errorf("level to int: %w", err)
	}

	// create title for level menu
	levelLower := _levelName[levelInt]
	title := strings.ToTitle(levelLower[:1]) + levelLower[1:] + " логи"
	// create level menu
	menu := &entity.Menu{
		Title: title,
		Items: make([]*entity.MenuItem, 0, len(levelMessages)),
	}
	// append menu items
	var itemName string
	var itemCtx context.Context
	for _, msg := range levelMessages {
		itemName = msg.Header
		itemCtx = context.WithValue(context.Background(), entity.LogMsgWithLvlCtxKey, msg)
		menu.Items = append(menu.Items, entity.NewMenuItem(itemCtx, itemName))
	}
	return menu, nil
}

// GetByID returns message with given ID. ID field must be presented.
func (l *LogMsgUC) GetByID(msg *entity.LogMsg) error {
	return l.logMsgRepoDB.GetByID(msg)
}

// DeleteByID deletes message record by its ID. ID field must be presented.
func (l *LogMsgUC) DeleteByID(id string) error {
	return l.logMsgRepoDB.DeleteByID(id)
}

// DeleteAllWithLevel deletes all message records with given level.
func (l *LogMsgUC) DeleteAllWithLevel(level string) error {
	return l.logMsgRepoDB.DeleteAllWithLevel(level)
}
