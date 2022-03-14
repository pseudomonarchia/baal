package controllers

import (
	"baal/configs"
	"baal/libs/logger"

	"go.uber.org/fx"
)

// Controllers represents a global controllers struct
type Controllers struct {
	Index *Index
}

// ControllerInjection represents a inject controller struct
type ControllerInjection struct {
	Log  *logger.Logger
	Conf *configs.GlobalConf
}

func registration(log *logger.Logger, conf *configs.GlobalConf) *Controllers {
	injection := &ControllerInjection{
		log,
		conf,
	}

	return &Controllers{
		&Index{injection},
	}
}

// Module is used for `fx.provider` to inject dependencies
var Module fx.Option = fx.Options(fx.Provide(registration))
