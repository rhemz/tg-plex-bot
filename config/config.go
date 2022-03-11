package config

import (
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

/*
type Configuration struct {
    Server ServerConfiguration
    Telegram TelegramConfiguration
}

type ServerConfiguration struct {
    Port int
}

type TelegramConfiguration struct {

}
*/

var (
	TelegramMessagesDisabled bool
	config                   *viper.Viper
)

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(file string) {
	var err error
	TelegramMessagesDisabled = false

	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(file)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		zap.S().Fatal("error on parsing configuration file")
	}

	config.BindEnv("telegram.botId", "TELEGRAM_BOT_ID")
	config.BindEnv("telegram.apiToken", "TELEGRAM_API_TOKEN")
	config.BindEnv("telegram.hookUrl", "TELEGRAM_HOOK_URL")
	config.BindEnv("ipinfo.apiToken", "IPINFO_API_TOKEN")
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

func GetConfig() *viper.Viper {
	return config
}
