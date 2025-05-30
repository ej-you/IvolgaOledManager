package main

import (
	"log"

	"sschmc/config"
	"sschmc/internal/app"
)

func main() {
	if err := startApp(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func startApp() error {
	// load config
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// init app
	application, err := app.New(cfg)
	if err != nil {
		return err
	}
	// run app
	if err := application.Run(); err != nil {
		return err
	}
	return nil
}
