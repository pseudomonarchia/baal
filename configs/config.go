package configs

import (
	"baal/libs/utils"
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
	dir, _                 = os.Getwd()
	defaultConfPath string = path.Join(dir, "./configs/conf.default.yml")
	rootConfPath    string = path.Join(dir, "./conf.yml")
)

func registration() *GlobalConf {
	var Global GlobalConf
	confPath := rootConfPath
	if !utils.FileIsExists(rootConfPath) {
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

// Module is used for `fx.provider` to inject dependencies
var Module fx.Option = fx.Options(fx.Provide(registration))
