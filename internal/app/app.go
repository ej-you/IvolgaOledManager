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

	"sschmc/config"
	"sschmc/internal/app/controller/buttons"
	"sschmc/internal/app/controller/renderer"
	repodb "sschmc/internal/app/repo/db"
	repostorage "sschmc/internal/app/repo/storage"
	"sschmc/internal/pkg/db"
	"sschmc/internal/pkg/storage"
)

const (
	_menuUpdateDuration = 300 * time.Millisecond // duration for update menu output to display
	_checkAliveTimeout  = time.Second            // duration for checking button is alive
)

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
	fmt.Println("cfg:", cfg)
	dbStorage, err := db.New(cfg.DB.DSN,
		db.WithTranslateError(),
		db.WithDisableColorful())
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
	fmt.Println("db:", a.dbStorage)
	// db, err := database.New(a.cfg.DB.DSN)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("db:", db)
	// rows, err := db.QueryContext(context.TODO(), `SELECT * FROM storage`)
	// if err != nil {
	// 	return err
	// }
	// defer rows.Close()
	// var id int
	// var level, header, content, createdAt string
	// for rows.Next() {
	// 	rows.Scan(&id, &level, &header, &content, &createdAt)
	// 	fmt.Println(id, level, header, content, createdAt)
	// }

	// return nil

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
	messageRepoDB := repodb.NewMessageRepoDB(a.dbStorage)

	// init renderer
	render, err := renderer.New(
		a.cfg.Hardware.Oled.Bus,
		a.cfg.App.GreetingsImgPath,
		_menuUpdateDuration,
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
		messageRepoDB,
		storageManager,
		updateDisplay,
	)
	if err != nil {
		return err
	}

	startControllers(appContext, btns, render)
	log.Println("App shutdown successfully!")
	return nil
}

// startControllers starts buttons and renderer in separately goroutines.
// This function is blocking. Context are used to stop controllers.
func startControllers(ctx context.Context, btns *buttons.Buttons, render *renderer.Renderer) {
	var wg sync.WaitGroup
	wg.Add(2)
	// start renderer
	go func() {
		defer wg.Done()
		render.StartWithShutdown(ctx)
	}()
	// start handle buttons rising/falling
	go func() {
		defer wg.Done()
		btns.HandleAll(ctx)
	}()
	wg.Wait()
}
