package buttons

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"sschmc/internal/app/entity"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/system"
)

var (
	_mainMenuItems          = 2               // amount of main menu items
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

// screenNone sets "none" app-status and update render.
func (b *Buttons) screenNone() {
	b.store.App.SetNone()
	b.render <- struct{}{}
}

// screenGreetings sets "greetings" app-status and update render.
func (b *Buttons) screenGreetings() {
	b.store.App.SetGreetings()
	b.render <- struct{}{}
}

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

// screenMessage sets "message" app-status and update render.
// But if delete button is selected it deletes all messages with selected level.
func (b *Buttons) screenMessage() {
	levelMenu := b.store.Menu.GetLevel()
	selectedItem := levelMenu.Items[levelMenu.SelectedItem]
	// get selected message ID from level menu (or level if item is delete button)
	selectedMsgValue, ok := selectedItem.Value.(string)
	if !ok {
		errlog.Print(errors.New("message id is not string"))
		return
	}

	// check if delete button is selected
	if selectedItem.IsDeleteButton() {
		b.deleteWithLevel(selectedMsgValue)
		return
	}

	msg := &entity.Message{ID: selectedMsgValue}
	// get message by ID
	err := b.msgRepoDB.GetByID(msg)
	if err != nil {
		errlog.Print(err)
		return
	}

	msg.Format()
	b.store.Message.Set(msg)
	b.store.App.SetMessage()
	b.render <- struct{}{}
}

// screenSensor sets "sensor" app-status and update render.
func (b *Buttons) screenSensor() {
	stationMenu := b.store.Menu.GetStation()
	selectedItem := stationMenu.Items[stationMenu.SelectedItem]
	// get selected sensor from station menu
	selectedSensor, ok := selectedItem.Value.(*entity.StationSensor)
	if !ok {
		errlog.Print(errors.New("sensor value is not *entity.StationSensor"))
		return
	}

	b.store.Sensor.Set(selectedSensor)
	b.store.App.SetSensor()
	b.render <- struct{}{}
}

// deleteWithLevel deletes all messages with selected level and display previous menu.
func (b *Buttons) deleteWithLevel(level string) {
	err := b.msgRepoDB.DeleteAllWithLevel(level)
	if err != nil {
		errlog.Print(err)
		return
	}
	log.Printf("All messages with level %s was deleted successfully", level)
	b.screenMenuLogs()
}

// deleteMessage deletes opened message and display previous menu.
func (b *Buttons) deleteMessage() {
	msg := b.store.Message.Get()
	if !msg.IsDeleteButton() {
		return
	}
	err := b.msgRepoDB.DeleteByID(msg)
	if err != nil {
		errlog.Print(err)
		return
	}
	log.Printf("Message %q (level %s) was deleted successfully", msg.Header, msg.Level)
	b.screenMenuLevel()
}

// updateStation change selected sensor status, updates
// station config file and restart station service.
func (b *Buttons) updateStation() {
	// get all sensors and selected sensor
	stationSensors := b.store.Sensor.GetAll()
	sensor := b.store.Sensor.Get()
	// update selected sensor status and station config
	stationSensors[sensor.Idx].ChangeActive()
	if err := b.stationRepoFile.UpdateSensors(stationSensors); err != nil {
		errlog.Print(err)
	}
	// restart station service
	if err := system.RestartService(b.stationService); err != nil {
		errlog.Print(err)
	}
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

// messageScrollUp scrolls message text up for one line.
func (b *Buttons) messageScrollUp(msg *entity.Message) {
	// scroll up message
	msg.ScrollUp()
	// update render with new message view
	b.render <- struct{}{}
}

// messageScrollDown scrolls message text down for one line.
func (b *Buttons) messageScrollDown(msg *entity.Message) {
	// scroll up message
	msg.ScrollDown()
	// update render with new message view
	b.render <- struct{}{}
}
