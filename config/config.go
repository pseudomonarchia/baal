package config

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"

	//nolint
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

var (
	dir, _       = os.Getwd()
	rootConfPath = path.Join(dir, "./config.yml")
)

// GlobalConf represents a global config struct
type GlobalConf struct {
	DEBUG bool
	PORT  int
	HTTPS bool
}

type secretConf struct {
	Oauth struct {
		Google struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
		}
	}
	Database struct {
		Mysql struct {
			Username string
			Database string
			Addr     string
			Port     int
			Password string
		}
	}
}

// Global config infomation
var Global GlobalConf

// Secret secret infomation
var Secret secretConf

var once sync.Once

func init() {
	fmt.Println("DIR:", dir)

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(rootConfPath)
	viper.ReadInConfig()
	viper.Unmarshal(&Secret, viper.DecodeHook(decodeHook))

	port, debug, https :=
		viper.GetInt("PORT"),
		viper.GetBool("DEBUG"),
		viper.GetBool("HTTPS")

	if port == 0 {
		port = 8080
	}

	Global = GlobalConf{
		DEBUG: debug,
		PORT:  port,
		HTTPS: https,
	}
}

func decodeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() == reflect.String {
		stringData := data.(string)
		if strings.HasPrefix(stringData, "${") && strings.HasSuffix(stringData, "}") {
			envVarValue := os.Getenv(strings.TrimPrefix(strings.TrimSuffix(stringData, "}"), "${"))
			if len(envVarValue) > 0 {
				return envVarValue, nil
			}
		}
	}

	return data, nil
}

// Setup override default configuration
func Setup(conf GlobalConf) {
	once.Do(func() {
		Global = conf
	})
}

// IsDev method is used to return whether it is currently in development mode
func (c *GlobalConf) IsDev() bool {
	return c.DEBUG
}

// PROTOCOL get local service url protocol
func (c *GlobalConf) PROTOCOL() string {
	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}

	return protocol
}
