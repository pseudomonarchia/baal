package config

import (
	"baal/lib/file"
	"os"
	"path"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// GlobalConf represents a global config struct
type GlobalConf struct {
	DEBUG bool
	PORT  string
}

var (
	dir, _          = os.Getwd()
	defaultConfPath = path.Join(dir, "./config/conf.default.yml")
	rootConfPath    = path.Join(dir, "./conf.yml")
)

// Module is used for `fx.provider` to inject dependencies
var Module fx.Option = fx.Options(fx.Provide(registration))

func registration() *GlobalConf {
	var global GlobalConf
	confPath := defaultConfPath
	if file.IsExists(rootConfPath) {
		confPath = rootConfPath
	}

	viper.SetConfigFile(confPath)
	viper.ReadInConfig()
	viper.Unmarshal(&global)

	global.DEBUG = os.Getenv("DEBUG") == "true"
	global.PORT = os.Getenv("PORT")

	return &global
}

// IsDev method is used to return whether it is currently in development mode
func (c *GlobalConf) IsDev() bool {
	return c.DEBUG
}
