// Package app provides App interface for run full application.
package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gorm.io/gorm"
	"periph.io/x/host/v3"

	"IvolgaOledManager/config"
	"IvolgaOledManager/internal/app/controller/buttons"
	"IvolgaOledManager/internal/app/controller/renderer"
	repodb "IvolgaOledManager/internal/app/repo/db"
	repostorage "IvolgaOledManager/internal/app/repo/storage"
	"IvolgaOledManager/internal/app/service"
	"IvolgaOledManager/internal/app/usecase"

	"IvolgaOledManager/internal/pkg/db"
	"IvolgaOledManager/internal/pkg/storage"
)

const (
	_renderUpdateDuration = 5 * time.Second // duration for update rendered display output
	_checkAliveTimeout    = time.Second     // duration for checking button is alive
)

var (
	_ Service = (*buttons.Buttons)(nil)
	_ Service = (*renderer.Renderer)(nil)
	_ Service = (*service.SensorsDataUpdate)(nil)
)

// App service interface.
type Service interface {
	StartWithShutdown(ctx context.Context)
}

var _ App = (*app)(nil)

type App interface {
	Run() error
}

// App implementation.
type app struct {
	cfg       *config.Config
	store     storage.Storage
	dbStorage *gorm.DB
}

// New returns App interface.
func New(cfg *config.Config) (App, error) {
	// connect to DB
	dbStorage, err := db.New(cfg.DB.DSN,
		db.WithTranslateError(),
		db.WithDisableColorful(),
		db.WithWarnLogLevel())
	if err != nil {
		return nil, err
	}

	// initialise all relevant drivers
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("init drivers: %w", err)
	}
	return &app{
		cfg:       cfg,
		store:     storage.NewMap(),
		dbStorage: dbStorage,
	}, nil
}

// Run starts full application.
func (a app) Run() error {
	// ctx for app
	appContext, appCancel := context.WithCancel(context.Background())
	// channel to update display
	updateDisplay := make(chan struct{})

	// handle shutdown process signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	// create gracefully shutdown task
	go func() {
		defer appCancel()
		handledSignal := <-quit
		log.Printf("Get %q signal. Shutdown app...", handledSignal.String())
	}()

	// init repos
	storageManager := repostorage.NewRepoStorageManager(a.store)
	stationResultRepoDB := repodb.NewSensorsRepoDB(a.dbStorage)
	// init usecases
	stationUC := usecase.NewStationResultUsecase(stationResultRepoDB, storageManager.StationResults)

	// init temperature update
	tempUpdate := service.NewSensorsDataUpdate(storageManager, stationUC, updateDisplay)

	// init renderer
	render, err := renderer.New(
		a.cfg.Hardware.Oled.Bus,
		a.cfg.App.GreetingsImgPath,
		_renderUpdateDuration,
		storageManager,
		updateDisplay,
	)
	if err != nil {
		return fmt.Errorf("init renderer: %w", err)
	}

	// init buttons
	btns, err := buttons.New(
		a.cfg.Hardware.Buttons.Escape,
		a.cfg.Hardware.Buttons.Up,
		a.cfg.Hardware.Buttons.Down,
		a.cfg.Hardware.Buttons.Enter,
		_checkAliveTimeout,
		storageManager,
		updateDisplay,
	)
	if err != nil {
		return err
	}

	startServices(appContext, btns, render, tempUpdate)
	log.Println("App shutdown successfully!")
	return nil
}

// startServices starts all services in separately goroutines.
// This function is blocking. Context are used to stop services.
func startServices(ctx context.Context, services ...Service) {
	// start all services
	var wg sync.WaitGroup // nolint:varnamelen // generally accepted name
	for _, service := range services {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.StartWithShutdown(ctx)
		}()
	}

	wg.Wait()
}
