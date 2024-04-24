package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getEnvUint(env *viper.Viper, key string) uint {
	value := env.GetUint(key)
	if value == 0 {
		log.Fatalf("Missing value in configuration: %v", key)
	}
	return value
}

func getEnvString(env *viper.Viper, key string) string {
	value := env.GetString(key)
	if value == "" {
		log.Fatalf("Missing value in configuration: %v", key)
	}
	return value
}

func getEnvBool(env *viper.Viper, key string) bool {
	return env.GetBool(key)
}
