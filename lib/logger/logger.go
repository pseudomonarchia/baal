package logger

import (
	"baal/config"
	"sync"

	"go.uber.org/zap"
)

type logger struct {
	*zap.Logger
}

// Log is Global default logger
var Log *logger
var once sync.Once

// Setup use log
func Setup() {
	once.Do(func() {
		Log = newLogger()
	})
}

func newLogger() *logger {
	var log *zap.Logger
	if config.Global.IsDev() {
		log, _ = zap.NewDevelopment()
	} else {
		log, _ = zap.NewProduction()
	}

	return &logger{log}
}
