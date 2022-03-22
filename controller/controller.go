package controller

import (
	"baal/config"
	"baal/lib/logger"

	"go.uber.org/fx"
)

// Controllers represents a global controllers struct
type Controllers struct {
	Index *Index
}

// C represents a inject controller struct
type C struct {
	Log  *logger.Logger
	Conf *config.GlobalConf
}

var _ = (*Controllers)(nil)
var _ = (*C)(nil)

func registration(log *logger.Logger, conf *config.GlobalConf) *Controllers {
	injection := &C{
		log,
		conf,
	}

	return &Controllers{
		&Index{injection},
	}
}

// Module is used for `fx.provider` to inject dependencies
var Module fx.Option = fx.Options(fx.Provide(registration))
