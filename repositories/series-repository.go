package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StreamRepository interface {
}

type dbRepository struct {
	connection *gorm.DB
}

func NewDbRepository(connectionData string) StreamRepository {
	log.Infof("Attempting connection to DB")
	connection, err := gorm.Open(postgres.Open(connectionData), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	log.Infof("Connected to DB successfully")
	log.Infof("Executing auto migrate")
	err = connection.AutoMigrate(&entities.ConfiguredStream{}, &entities.Stream{}, &entities.Station{}, &entities.Network{})
	if err != nil {
		log.Fatalf("Failed to auto migrate model to DB: %v", err)
	}
	log.Infof("Executed automigrate successfully")
	return &dbRepository{connection}
}
