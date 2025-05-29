package entity

import (
	"log"
	"os"
	"testing"
	"unicode/utf8"
)

var _menu *Menu

func TestMain(m *testing.M) {
	// create menu
	_menu = &Menu{
		Title: "Main menu",
		Items: []*MenuItem{
			NewMenuItem("Какого hell ты делаешь, dude?", ""),
			NewMenuItem("Всё good", ""),
		},
	}

	logger := log.New(os.Stdout, "    ", log.Lshortfile)
	// print out letters of titles of menu items
	for idx, item := range _menu.Items {
		logger.Printf("-- Item %d --", idx)
		for idx, symbol := range item.Title {
			logger.Printf("symbol: %q | idx: %d | runeLen: %d\n", symbol, idx, utf8.RuneLen(symbol))
		}
	}
	// run tests
	os.Exit(m.Run())
}

func TestLongItem(t *testing.T) {
	t.Log("Scroll long menu item")
	item := _menu.Items[0]

	for range 20 {
		item.Scroll()
	}
	for range 10 {
		item.Scroll()
		t.Logf("(from %02d to %02d): %q", item.firstSymbol, item.lastSymbol, item.FormattedTitle())
	}
}

func TestShortItem(t *testing.T) {
	t.Log("Scroll short menu item")
	item := _menu.Items[1]

	for range 3 {
		item.Scroll()
		t.Logf("(from %02d to %02d): %q", item.firstSymbol, item.lastSymbol, item.FormattedTitle())
	}
}
