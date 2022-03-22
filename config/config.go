package config

import (
	"baal/lib/file"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// GlobalConf represents a global config struct
type GlobalConf struct {
	MODE string
	PORT string
}

var (
	dir, _          = os.Getwd()
	defaultConfPath = path.Join(dir, "./configs/conf.default.yml")
	rootConfPath    = path.Join(dir, "./conf.yml")
)

// Module is used for `fx.provider` to inject dependencies
var Module fx.Option = fx.Options(fx.Provide(registration))

func registration() *GlobalConf {
	var Global GlobalConf
	confPath := rootConfPath
	if !file.IsExists(rootConfPath) {
		confPath = defaultConfPath
	}

	viper.SetConfigFile(confPath)
	viper.ReadInConfig()
	viper.AutomaticEnv()
	viper.Unmarshal(&Global)

	return &Global
}

// IsDev method is used to return whether it is currently in development mode
func (c *GlobalConf) IsDev() bool {
	return strings.ToUpper(c.MODE) == "DEBUG"
}
