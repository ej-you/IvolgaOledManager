package buttons

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"sschmc/internal/app/entity"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/gpiobutton"
)

var (
	_levelsAmount = 6
	_levelName    = map[int]string{
		0: "trace",
		1: "debug",
		2: "info",
		3: "warn",
		4: "error",
		5: "fatal",
	}
)

// BtnEntRisingHandler handles all cases of ENT button rising.
func (b *Buttons) BtnEntRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.btnAllGreetings()
		case b.store.App.IsGreetings():
			b.btnEntMenuMain()
		case b.store.App.IsMenuMain():
			b.btnEntMenuLevel()
		case b.store.App.IsMenuLevel():
			b.btnEntMessage()
		default:
			log.Println("*** ENTER pressed ***")
		}
	}
}

// btnEntMenuMain update status to menu-main and creates render task.
func (b *Buttons) btnEntMenuMain() {
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

	// create empty menu and append menu items
	mainMenu := &entity.Menu{
		Title: "Main menu",
		Items: make([]*entity.MenuItem, 0, len(_levelName)),
	}

	var name string
	for _, levelCount := range allLevels {
		name = fmt.Sprintf("%s (%d)", _levelName[levelCount.Level], levelCount.Count)
		mainMenu.Items = append(mainMenu.Items, entity.NewMenuItem(name, levelCount))
	}

	// save menu to storage
	b.store.Menu.SetMain(mainMenu)
	b.store.App.SetMenuMain()
	// update render according to new app-status
	b.render <- struct{}{}
}

// btnEntMenuLevel update status to menu-level and creates render task.
func (b *Buttons) btnEntMenuLevel() {
	// get selected level from main menu
	mainMenu := b.store.Menu.GetMain()
	selectedItem, ok := mainMenu.Items[mainMenu.SelectedItem].Value.(entity.MessageLevelCount)
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
		Items: make([]*entity.MenuItem, 0, len(levelMessages)),
	}
	// append menu items
	var header string
	for _, msg := range levelMessages {
		header = msg.Header
		levelMenu.Items = append(levelMenu.Items, entity.NewMenuItem(header, msg.ID))
	}

	b.store.Menu.SetLevel(levelMenu)
	b.store.App.SetMenuLevel()
	// update render according to new app-status
	b.render <- struct{}{}
}

// btnEntMessage clears rendered data and updates app-status in storage to none.
func (b *Buttons) btnEntMessage() {
	// get selected message ID from level menu
	levelMenu := b.store.Menu.GetLevel()
	selectedMsgID, ok := levelMenu.Items[levelMenu.SelectedItem].Value.(string)
	if !ok {
		errlog.Print(errors.New("message id is not string"))
		return
	}

	msg := &entity.Message{
		ID: selectedMsgID,
	}

	// get message by ID
	err := b.msgRepoDB.GetByID(msg)
	if err != nil {
		errlog.Print(err)
		return
	}

	msg.Format()
	b.store.Message.Set(msg)
	b.store.App.SetMessage()
	// update render according to new app-status
	b.render <- struct{}{}
}
