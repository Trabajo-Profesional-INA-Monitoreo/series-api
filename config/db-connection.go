package config

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repositories struct {
	StreamsRepository          repositories.StreamRepository
	ConfigurationRepository    repositories.ConfigurationRepository
	ConfiguredStreamRepository repositories.ConfiguredStreamsRepository
	ErrorsRepository           repositories.ErrorsRepository
	NodeRepository             repositories.NodeRepository
	FilterRepository           repositories.FilterRepository
	MetricsRepository          repositories.MetricsRepository
	RedundancyRepository       repositories.RedundancyRepository
	InputsRepository           repositories.InputsRepository
	StationsRepository         repositories.StationsRepository
	OutputsRepository          repositories.OutputsRepository
}

func loggerStringToGormLogger(loggerConfig string) logger.LogLevel {
	switch loggerConfig {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Info
	}
}

func CreateRepositories(connectionData string, loggerConfig string) *Repositories {
	log.Infof("Attempting connection to DB")
	connection, err := gorm.Open(postgres.Open(connectionData), &gorm.Config{
		Logger: logger.Default.LogMode(loggerStringToGormLogger(loggerConfig)),
	})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	log.Infof("Connected to DB successfully")
	log.Infof("Executing auto migrate")
	err = connection.AutoMigrate(
		&entities.Unit{},
		&entities.Procedure{},
		&entities.Variable{},
		&entities.Station{},
		&entities.ConfiguredStream{},
		&entities.Stream{},
		&entities.Configuration{},
		&entities.DetectedError{},
		&entities.ConfiguredMetric{},
		&entities.Node{},
		&entities.Redundancy{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate model to DB: %v", err)
	}
	log.Infof("Executed automigrate successfully")

	log.Infof("Creating repositories...")
	repos := Repositories{
		StreamsRepository:          repositories.NewStreamRepository(connection),
		ConfigurationRepository:    repositories.NewConfigurationRepository(connection),
		ConfiguredStreamRepository: repositories.NewConfiguredStreamsRepository(connection),
		ErrorsRepository:           repositories.NewErrorsRepository(connection),
		NodeRepository:             repositories.NewNodeRepository(connection),
		FilterRepository:           repositories.NewFilterRepository(connection),
		MetricsRepository:          repositories.NewMetricsRepository(connection),
		RedundancyRepository:       repositories.NewRedundancyRepository(connection),
		InputsRepository:           repositories.NewInputsRepository(connection),
		StationsRepository:         repositories.NewStationsRepository(connection),
		OutputsRepository:          repositories.NewOutputsRepository(connection),
	}
	log.Infof("Done creating repositories")
	return &repos
}
