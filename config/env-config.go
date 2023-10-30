package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type ApiConfig struct {
	LogLevel   string
	ServerPort string
}

// initEnv Initializes the configuration properties from a config file and environment
func initEnv() (*viper.Viper, error) {
	v := viper.New()

	// Configure viper to read env variables with the CLI_ prefix
	v.AutomaticEnv()
	v.SetEnvPrefix("cli")
	// Use a replacer to replace env variables underscores with points. This let us
	// use nested configurations in the config file and at the same time define
	// env variables for the nested configurations
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Add env variables supported
	_ = v.BindEnv("log", "level")
	_ = v.BindEnv("server", "port")

	// Try to read configuration from config file. If config file
	// does not exist then ReadInConfig will fail but configuration
	// can be loaded from the environment variables, so we shouldn't
	// return an error in that case
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Warnf("Config | Warning Message | Configuration could not be read from config file. Using env variables instead")
	}

	return v, nil
}

func GetConfig() *ApiConfig {
	env, err := initEnv()
	if err != nil {
		log.Fatalf("Failed to read environment, exiting")
	}
	logLevel := env.GetString("log.level")
	if logLevel == "" {
		log.Warnf("Missing log level, using info")
		logLevel = "info"
	}
	serverPort := env.GetString("server.port")
	if serverPort == "" {
		log.Fatalf("Missing server port, exiting")
	}
	return &ApiConfig{
		LogLevel:   logLevel,
		ServerPort: serverPort,
	}
}
