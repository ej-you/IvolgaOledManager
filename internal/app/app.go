// Package app provides struct with Run method to start full application.
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
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/display"
	"IvolgaOledManager/internal/app/service/sensordata"
	"IvolgaOledManager/internal/app/usecase"
	"IvolgaOledManager/internal/pkg/db"
	"IvolgaOledManager/internal/pkg/logger"
	"IvolgaOledManager/internal/pkg/pubsub"
)

var _ Service = (*button.Buttons)(nil)
var _ Service = (*display.Display)(nil)

// App service interface.
type Service interface {
	StartWithShutdown(ctx context.Context) error
}

type App struct {
	cfg      *config.Config
	services []Service
}

// New returns new app instance.
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
	_, err = db.New(cfg.DB.DSN,
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
	sensorRepoDB := repodb.NewMockStationResultRepoDB()
	// init usecases
	sensorUC := usecase.NewSensorDataUsecase(sensorRepoDB, storage)

	// init buttons services
	btns, err := button.NewButtons(cfg)
	if err != nil {
		return nil, fmt.Errorf("create buttons service: %w", err)
	}
	// init display service
	displ, err := display.NewDisplay(cfg)
	if err != nil {
		return nil, fmt.Errorf("create display service: %w", err)
	}
	// init updater services
	tempUpdater := sensordata.NewTemperatureUpdater(cfg, storage, sensorUC)
	humidUpdater := sensordata.NewHumidityUpdater(cfg, storage, sensorUC)
	pressUpdater := sensordata.NewPressureUpdater(cfg, storage, sensorUC)
	windSpeedUpdater := sensordata.NewWindSpeedUpdater(cfg, storage, sensorUC)
	windDirUpdater := sensordata.NewWindDirUpdater(cfg, storage, sensorUC)

	return &App{
		cfg: cfg,
		services: []Service{btns, displ,
			tempUpdater, humidUpdater, pressUpdater, windSpeedUpdater, windDirUpdater},
	}, nil
}

// Run starts all services. This function is blocking.
// It waits for os signal to gracefully shutdown all services.
// Or it waits for fall down one of the services and stop other services.
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
	var wg sync.WaitGroup // nolint:varnamelen // generally accepted name
	serviceErr := make(chan error, 1)
	for _, service := range a.services {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := service.StartWithShutdown(appContext); err != nil {
				serviceErr <- err
			}
		}()
	}
	logrus.Info("all services were started successfully")

	select {
	case handledSignal := <-quitSig:
		cancel()
		logrus.Infof("Got %s signal. Shutdown services...", handledSignal.String())
	case err := <-serviceErr:
		cancel()
		appErr = fmt.Errorf("service: %w", err)
		logrus.Info("One of the services fell down. Shutdown other services...")
	case <-appContext.Done():
		appErr = appContext.Err()
		logrus.Info("Context canceled. Shutdown app...")
	}

	// wait for all services
	wg.Wait()
	logrus.Info("All services was stopped. Shutdown app")
	return appErr
}
