package main

import (
	"log"

	"sschmc/config"
	"sschmc/internal/display"
)

func main() {
	if err := startApp(); err != nil {
		log.Fatal(err)
	}
}

func startApp() error {
	// load config
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// init display
	oled, err := display.New(cfg.App.OledBus, cfg.App.GreetingsImgPath)
	if err != nil {
		return err
	}

	if err := oled.Greetings(); err != nil {
		return err
	}
	return nil
}
