package logger

import (
	"baal/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Logger represents a global logger struct
type Logger struct {
	*zap.Logger
}

var _ = (*Logger)(nil)

// Module is used for `fx.provider` to inject dependencies
var Module = fx.Option(fx.Provide(registration))

func registration(conf *config.GlobalConf) *Logger {
	var log *zap.Logger
	if conf.IsDev() {
		log, _ = zap.NewDevelopment()
	} else {
		log, _ = zap.NewProduction()
	}

	defer log.Sync()
	return &Logger{log}
}
