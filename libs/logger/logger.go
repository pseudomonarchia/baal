package logger

import (
	"baal/configs"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Logger represents a global logger struct
type Logger struct {
	*zap.Logger
}

func registration(conf *configs.GlobalConf) *Logger {
	var log *zap.Logger
	if conf.IsDev() {
		log, _ = zap.NewDevelopment()
	} else {
		log, _ = zap.NewProduction()
	}

	defer log.Sync()
	return &Logger{log}
}

// Module is used for `fx.provider` to inject dependencies
var Module = fx.Option(fx.Provide(registration))
