package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type ServiceConfigurationData struct {
	LogLevel                    string
	SqlLogLevel                 string
	ServerPort                  string
	DbUrl                       string
	FaultCronTime               string
	InaToken                    string
	InaBaseUrl                  string
	NotificationApiBaseUrl      string
	ForecastMaxWaitingTimeHours float64
	SecurityEnabled             bool
	KeycloakConfig              *KeycloakConfiguration
	DetectionMaxThreads         int
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
	_ = v.BindEnv("log", "sql")
	_ = v.BindEnv("server", "port")
	_ = v.BindEnv("datasource", "connection")
	_ = v.BindEnv("faults-detector", "cron")
	_ = v.BindEnv("faults-detector", "max-threads")
	_ = v.BindEnv("ina-client", "token")
	_ = v.BindEnv("ina-client", "base-url")
	_ = v.BindEnv("notifications-api-client", "base-url")
	_ = v.BindEnv("security", "enabled")
	_ = v.BindEnv("keycloak", "url")
	_ = v.BindEnv("keycloak", "realm")
	_ = v.BindEnv("keycloak", "client")
	_ = v.BindEnv("keycloak", "secret")

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

func GetConfig() *ServiceConfigurationData {
	env, err := initEnv()
	if err != nil {
		log.Fatalf("Failed to read environment, exiting")
	}

	logLevel := getEnvString(env, "log.level")
	sqlLogLevel := getEnvString(env, "log.sql")
	serverPort := getEnvString(env, "server.port")
	dbConnection := getEnvString(env, "datasource.connection")
	faultsDetectorCron := getEnvString(env, "faults-detector.cron")
	detectionMaxThreads := getEnvUint(env, "faults-detector.max-threads")
	inaBaseUrl := getEnvString(env, "ina-client.base-url")
	inaToken := getEnvString(env, "ina-client.token")
	notificationApiBaseUrl := getEnvString(env, "notifications-api-client.base-url")
	securityEnabled := getEnvBool(env, "security.enabled")
	kcConfig := getKeycloakConfig(env, securityEnabled)

	return &ServiceConfigurationData{
		LogLevel:               logLevel,
		SqlLogLevel:            sqlLogLevel,
		ServerPort:             serverPort,
		DbUrl:                  dbConnection,
		FaultCronTime:          faultsDetectorCron,
		InaBaseUrl:             inaBaseUrl,
		InaToken:               inaToken,
		NotificationApiBaseUrl: notificationApiBaseUrl,
		SecurityEnabled:        securityEnabled,
		KeycloakConfig:         kcConfig,
		DetectionMaxThreads:    int(detectionMaxThreads),
	}
}
