package controllers

import (
	"baal/libs/logger"

	"go.uber.org/fx"
)

// Controllers represents a global controllers struct
type Controllers struct {
	Index *Index
}

func registration(log *logger.Logger) *Controllers {
	return &Controllers{
		&Index{log},
	}
}

// Module is used for `fx.provider` to inject dependencies
var Module fx.Option = fx.Options(fx.Provide(registration))
