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
}

func CreateRepositories(connectionData string) *Repositories {
	log.Infof("Attempting connection to DB")
	connection, err := gorm.Open(postgres.Open(connectionData), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
		&entities.Network{},
		&entities.ConfiguredStream{},
		&entities.Stream{},
		&entities.Configuration{},
		&entities.DetectedError{},
		&entities.ConfiguredMetric{},
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
	}
	log.Infof("Done creating repositories")
	return &repos
}
