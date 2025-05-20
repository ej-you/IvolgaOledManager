package main

import (
	"log"
	"time"

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
	defer oled.Close()

	if err := oled.Greetings(); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	if err := oled.Clear(); err != nil {
		return err
	}

	return nil
}
