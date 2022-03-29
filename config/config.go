package config

import (
	"baal/lib/file"
	"fmt"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/spf13/viper"
)

// GlobalConf represents a global config struct
type GlobalConf struct {
	DEBUG bool
	PORT  int
	HOST  string
	HTTPS bool
}

var (
	dir, _          = os.Getwd()
	defaultConfPath = path.Join(dir, "./config/conf.default.yml")
	rootConfPath    = path.Join(dir, "./conf.yml")
)

// Global config infomation
var Global *GlobalConf
var once sync.Once

func init() {
	var global GlobalConf

	confPath := defaultConfPath
	if file.IsExists(rootConfPath) {
		confPath = rootConfPath
	}

	viper.SetConfigFile(confPath)
	viper.ReadInConfig()
	viper.Unmarshal(&global)

	Global = &global
}

// Setup override default configuration
func Setup(conf *GlobalConf) {
	once.Do(func() {
		Global = conf
	})
}

// IsDev method is used to return whether it is currently in development mode
func (c *GlobalConf) IsDev() bool {
	return c.DEBUG
}

// URL get local service address
func (c *GlobalConf) URL() string {
	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}

	return fmt.Sprintf(
		"%s://%s:%s",
		protocol,
		c.HOST,
		strconv.Itoa(c.PORT),
	)
}
