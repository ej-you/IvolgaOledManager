// Package db provide *gorm.DB for interaction with the database through GORM methods.
package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	_logLevelError = "error" // error log level tag
	_logLevelWarn  = "warn"  // warn log level tag
)

// Logger is an internal interface compatible with a logger.Writer.
// It is used to configure a custom DB logger.
type Logger interface {
	Printf(format string, args ...any)
}

// dbSettings is a db settings for *gorm.DB with custom options.
// It is used when creating a new *gorm.DB object.
type dbSettings struct {
	customLogger    Logger
	logLevel        logger.LogLevel
	translateError  bool
	ignoreNotFound  bool
	disableColorful bool
}

// Option represents an option for DB struct initializing.
type Option func(*dbSettings)

// New returns new DB instance with connection to given DSN.
// Options can be set with "WithSmth" funcs.
func New(dsn string, options ...Option) (*gorm.DB, error) {
	dbStorage := &dbSettings{
		customLogger:    log.Default(),
		logLevel:        logger.Info,
		translateError:  false,
		ignoreNotFound:  false,
		disableColorful: false,
	}

	// apply all options to customize DB struct
	for _, opt := range options {
		opt(dbStorage)
	}

	gormDB, err := gorm.Open(
		withConn(dsn),
		&gorm.Config{
			// set UTC time zone
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
			Logger: logger.New(
				dbStorage.customLogger,
				logger.Config{
					LogLevel:                  dbStorage.logLevel,
					IgnoreRecordNotFoundError: dbStorage.ignoreNotFound,
					Colorful:                  !dbStorage.disableColorful,
				},
			),
			TranslateError: dbStorage.translateError,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("open db connection: %w", err)
	}

	dbStorage.customLogger.Printf("successfully connected to DB")
	return gormDB, nil
}

// WithLogger sets custom logger for DB. Optional.
func WithLogger(customLogger Logger) Option {
	return func(d *dbSettings) {
		d.customLogger = customLogger
	}
}

// WithLogLevel sets log level for DB. Accepted values: "info", "warn", "error". Optional.
func WithLogLevel(logLevel string) Option {
	// default level is info level
	level := logger.Info
	switch logLevel {
	case _logLevelWarn:
		level = logger.Warn
	case _logLevelError:
		level = logger.Error
	}

	return func(d *dbSettings) {
		d.logLevel = level
	}
}

// WithErrorLogLevel sets error log level for DB. Optional.
func WithErrorLogLevel() Option {
	return func(d *dbSettings) {
		d.logLevel = logger.Error
	}
}

// WithWarnLogLevel sets warn log level for DB. Optional.
func WithWarnLogLevel() Option {
	return func(d *dbSettings) {
		d.logLevel = logger.Warn
	}
}

// WithTranslateError sets translate error parameter true. Optional.
func WithTranslateError() Option {
	return func(d *dbSettings) {
		d.translateError = true
	}
}

// WithIgnoreNotFound sets ignore record not found error parameter true. Optional.
func WithIgnoreNotFound() Option {
	return func(d *dbSettings) {
		d.ignoreNotFound = true
	}
}

// WithDisableColorful sets colorful log output false. Optional.
func WithDisableColorful() Option {
	return func(d *dbSettings) {
		d.disableColorful = true
	}
}

// withConn sets connection for DB. Required.
// There is the MyySQL is used as DB.
func withConn(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}
