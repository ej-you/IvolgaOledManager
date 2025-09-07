// Package app provides object with Run method to start full application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"periph.io/x/host/v3"

	"IvolgaOledManager/config"
	repodb "IvolgaOledManager/internal/app/repo/db"
	repofs "IvolgaOledManager/internal/app/repo/fs"
	"IvolgaOledManager/internal/app/service/button"
	displayservice "IvolgaOledManager/internal/app/service/display"
	"IvolgaOledManager/internal/app/service/render"
	"IvolgaOledManager/internal/app/service/screen"
	"IvolgaOledManager/internal/app/service/sensdata"
	"IvolgaOledManager/internal/app/usecase"
	"IvolgaOledManager/internal/pkg/db"
	"IvolgaOledManager/internal/pkg/display"
	"IvolgaOledManager/internal/pkg/logger"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// Ensure GPIO-button implements interface.
var _ Service = (*button.Buttons)(nil)

// Ensure OLED-display implements interface.
var _ Service = (*display.Service)(nil)

// Ensure render implements interface.
var _ Service = (*render.Render)(nil)

// Ensure sensor data updater implements interface.
var _ Service = (*sensdata.Updater)(nil)

// Ensure screen implements interface.
var _ Service = screen.Screen(nil)

// Service describes an app service.
type Service interface {
	// StartWithShutdown starts service and wait for context cancellation to shutdown it.
	StartWithShutdown(ctx context.Context) error
	// Ready returns true if service was completely started and is ready-to-use now.
	Ready() <-chan struct{}
}

// App is a main object that starts full app.
type App struct {
	cfg      *config.Config
	storage  pubsub.Storage
	services []Service
}

// New returns a new instance of App.
func New() (*App, error) {
	// initialise all relevant drivers
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("init drivers: %w", err)
	}

	// load config
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("create config: %w", err)
	}
	// setup logger
	logger.InitLogrus(cfg.App.LogLevel, cfg.App.LogFormat)

	// connect to DB
	dbStorage, err := db.New(cfg.DB.DSN,
		db.WithTranslateError(),
		db.WithIgnoreNotFound(),
		db.WithDisableColorful(),
		db.WithLogLevel(cfg.App.LogLevel),
		db.WithLogger(logrus.StandardLogger()))
	if err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}
	// create storage
	storage := pubsub.NewKeyValueStorage()

	// init repos
	sensdataRepoDB := repodb.NewMockSensdataRepoDB()
	logMsgRepoDB := repodb.NewLogMsgRepoDB(dbStorage)
	sensconfRepoFS, err := repofs.NewSensconfRepoFS(cfg.Other.Station.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("sensorconf repo fs: %w", err)
	}
	// init usecases
	sensdataUC := usecase.NewSensdataUsecase(sensdataRepoDB, storage)
	logMsgUC := usecase.NewLogMsgUsecase(logMsgRepoDB)
	sensconfUC := usecase.NewSensconfUsecase(sensconfRepoFS, cfg.Other.Station.ServiceName)

	// init buttons services
	btns, err := button.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("create buttons services: %w", err)
	}
	// init display service
	disp, err := displayservice.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("create display service: %w", err)
	}
	// init render service
	rend := render.New(disp, storage)
	// init updater services
	tempUpdater := sensdata.NewTemperatureUpdater(cfg, storage, sensdataUC)
	humidUpdater := sensdata.NewHumidityUpdater(cfg, storage, sensdataUC)
	pressUpdater := sensdata.NewPressureUpdater(cfg, storage, sensdataUC)
	windSpeedUpdater := sensdata.NewWindSpeedUpdater(cfg, storage, sensdataUC)
	windDirUpdater := sensdata.NewWindDirUpdater(cfg, storage, sensdataUC)
	// init screens services
	screenManager := screen.NewManager(cfg, btns, rend, storage, sensconfUC, logMsgUC)

	return &App{
		cfg:     cfg,
		storage: storage,
		services: []Service{btns, disp, rend, screenManager,
			tempUpdater, humidUpdater, pressUpdater, windSpeedUpdater, windDirUpdater},
	}, nil
}

// Run starts all services. This function is blocking.
// It waits for os signal to gracefully shutdown all services.
// Or it waits for fall down one of the services and stops other services.
func (a *App) Run() error {
	var appErr error

	// ctx for app
	appContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle shutdown process signals
	quitSig := make(chan os.Signal, 1)
	signal.Notify(quitSig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// start all services
	var wgRunning sync.WaitGroup
	var wgReady sync.WaitGroup
	serviceErr := make(chan error, 1)
	for _, service := range a.services {
		wgRunning.Add(1)
		wgReady.Add(1)
		// start service
		go func() {
			defer wgRunning.Done()
			if err := service.StartWithShutdown(appContext); err != nil {
				serviceErr <- err
			}
		}()
		// wait for service is ready
		go func() {
			defer wgReady.Done()
			<-service.Ready()
		}()
	}
	// wait for all services until they are ready
	wgReady.Wait()
	logrus.Info("all services were started successfully")

	select {
	case handledSignal := <-quitSig:
		cancel()
		logrus.Infof("got %s signal. Shutdown services...", handledSignal.String())
	case err := <-serviceErr:
		cancel()
		appErr = fmt.Errorf("service: %w", err)
		logrus.Info("one of the services fell down. Shutdown other services...")
	case <-appContext.Done():
		appErr = appContext.Err()
		logrus.Info("context canceled. Shutdown app...")
	}

	// wait for all services
	wgRunning.Wait()
	logrus.Info("all services was stopped. Shutdown app")
	return appErr
}
