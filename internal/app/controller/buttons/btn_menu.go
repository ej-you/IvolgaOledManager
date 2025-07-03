package buttons

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/errlog"
)

var (
	_mainMenuItemLogsIdx    = 0               // first item index
	_mainMenuItemLogs       = "Log messages"  // first item title
	_mainMenuItemStationIdx = 1               // second item index
	_mainMenuItemStation    = "Sensors setup" // second item title

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

// screenMenuMain sets "menu-main" app-status and update render.
func (b *Buttons) screenMenuMain() {
	// create empty menu
	mainMenu := &entity.Menu{
		Title: "Main menu",
		Items: []*entity.MenuItem{
			entity.NewMenuItem(_mainMenuItemLogs, nil),
			entity.NewMenuItem(_mainMenuItemStation, nil),
		},
	}
	b.store.Menu.SetMain(mainMenu)
	b.store.App.SetMenuMain()
	b.render <- struct{}{}
}

// screenMenuLogsOrStation uses screenMenuLogs or screenMenuStation handler
// according to selected item in main menu.
func (b *Buttons) screenMenuLogsOrStation() {
	mainMenu := b.store.Menu.GetMain()

	switch mainMenu.SelectedItem {
	case _mainMenuItemLogsIdx:
		b.screenMenuLogs()
	case _mainMenuItemStationIdx:
		b.screenMenuStation()
	}
}

// screenMenuLogs sets "menu-logs" app-status and update render
func (b *Buttons) screenMenuLogs() {
	// get levels count from DB
	levelsCount, err := b.msgRepoDB.GetLevelsCount()
	if err != nil {
		errlog.Print(err)
		return
	}

	// create slice for all levels count and fill it with zeros
	allLevels := make([]entity.MessageLevelCount, _levelsAmount)
	for idx := range allLevels {
		allLevels[idx].Level = idx
	}
	// update level count values according to gotten DB data
	for _, levelCount := range levelsCount {
		allLevels[levelCount.Level].Count = levelCount.Count
	}

	// create empty menu
	logsMenu := &entity.Menu{
		Title: "Logs menu",
		Items: make([]*entity.MenuItem, 0, len(_levelName)),
	}
	// append menu items
	var name string
	for _, levelCount := range allLevels {
		name = fmt.Sprintf("%s (%d)", _levelName[levelCount.Level], levelCount.Count)
		logsMenu.Items = append(logsMenu.Items, entity.NewMenuItem(name, levelCount))
	}

	b.store.Menu.SetLogs(logsMenu)
	b.store.App.SetMenuLogs()
	b.render <- struct{}{}
}

// screenMenuStation sets "menu-station" app-status and update render.
func (b *Buttons) screenMenuStation() {
	allSensors, err := b.stationRepoFile.ParseSensors()
	if err != nil {
		errlog.Print(err)
		return
	}

	// create empty menu
	sensorsMenu := &entity.Menu{
		Title: "Station sensors",
		Items: make([]*entity.MenuItem, 0, len(allSensors)),
	}
	// append menu items
	for _, sensor := range allSensors {
		sensorsMenu.Items = append(sensorsMenu.Items, entity.NewMenuItem(sensor.Name, sensor))
	}

	b.store.Sensor.SetAll(allSensors)
	b.store.Menu.SetStation(sensorsMenu)
	b.store.App.SetMenuStation()
	b.render <- struct{}{}
}

// screenMenuLevel sets "menu-level" app-status and update render.
func (b *Buttons) screenMenuLevel() {
	// get selected level from main menu
	logsMenu := b.store.Menu.GetLogs()
	selectedItem, ok := logsMenu.Items[logsMenu.SelectedItem].Value.(entity.MessageLevelCount)
	if !ok {
		errlog.Print(errors.New("menu item is not entity.MessageLevelCount"))
		return
	}
	// skip button handling if messages amount with selected level is zero.
	if selectedItem.Count == 0 {
		return
	}

	// get messages with selected level from DB
	levelMessages, err := b.msgRepoDB.GetWithLevel(strconv.Itoa(selectedItem.Level))
	if err != nil {
		errlog.Print(err)
		return
	}

	// create title for level menu
	levelLower := _levelName[selectedItem.Level]
	title := strings.ToTitle(levelLower[0:1]) + levelLower[1:] + " messages"
	// create level menu
	levelMenu := &entity.Menu{
		Title: title,
		Items: make([]*entity.MenuItem, 0, len(levelMessages)+1),
	}
	// append menu items
	var header string
	for _, msg := range levelMessages {
		header = msg.Header
		levelMenu.Items = append(levelMenu.Items, entity.NewMenuItem(header, msg.ID))
	}
	// button for deleting all messages with selected level
	deleteWithLevelBtn := entity.NewMenuItemDeleteButton(strconv.Itoa(selectedItem.Level))
	levelMenu.Items = append(levelMenu.Items, deleteWithLevelBtn)

	b.store.Menu.SetLevel(levelMenu)
	b.store.App.SetMenuLevel()
	b.render <- struct{}{}
}

// menuScrollUp selects the previous item in given menu.
func (b *Buttons) menuScrollUp(menu *entity.Menu) {
	// scroll down menu
	menu.SelectPrevious()
	// update render with new menu view
	b.render <- struct{}{}
}

// menuScrollDown selects the next item in given menu.
func (b *Buttons) menuScrollDown(menu *entity.Menu) {
	// scroll down menu
	menu.SelectNext()
	// update render with new menu view
	b.render <- struct{}{}
}
