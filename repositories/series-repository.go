package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type StreamRepository interface {
	GetNetworks() []dtos.StreamsPerNetwork
	GetStations() []dtos.StreamsPerStation
}

type dbRepository struct {
	connection *gorm.DB
}

func NewDbRepository(connectionData string) StreamRepository {
	log.Infof("Attempting connection to DB")
	connection, err := gorm.Open(postgres.Open(connectionData), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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

func (db *dbRepository) GetNetworks() []dtos.StreamsPerNetwork {
	var networks []dtos.StreamsPerNetwork

	db.connection.Model(
		&entities.Stream{},
	).Select(
		"networks.name as networkname",
		"networks.network_id as networkid",
		"count(streams.stream_id) as streamscount",
	).Joins("JOIN networks ON streams.network_id = networks.network_id").Group("networks.name, networks.network_id").Scan(&networks)
	log.Debugf("Get network query result: %v", networks)
	return networks
}

func (db *dbRepository) GetStations() []dtos.StreamsPerStation {
	var stations []dtos.StreamsPerStation

	db.connection.Model(
		&entities.Stream{},
	).Select(
		"stations.name as stationname",
		"stations.station_id as stationid",
		"count(streams.stream_id) as streamscount",
	).Joins("JOIN stations ON streams.station_id = stations.station_id").Group("stations.name, stations.station_id").Scan(&stations)
	log.Debugf("Get stations query result: %v", stations)
	return stations
}
