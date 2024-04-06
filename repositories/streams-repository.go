package repositories

import (
	"errors"
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StreamRepository interface {
	GetNetworks() []dtos.StreamsPerNetwork
	GetStations() []dtos.StreamsPerStation

	GetTotalStreams() int
	GetTotalNetworks() int
	GetTotalStations() int
	GetStreams() []entities.Stream
	GetStreamWithAssociatedData(streamId uint64) (entities.Stream, error)
	CreateUnit(entity entities.Unit) error
	CreateStation(entity entities.Station) error
	CreateProcedure(entity entities.Procedure) error
	CreateVariable(entity entities.Variable) error
	CreateStream(entity entities.Stream) error
}

type streamsRepository struct {
	connection *gorm.DB
}

func NewStreamRepository(connection *gorm.DB) StreamRepository {
	return &streamsRepository{connection}
}

func (db *streamsRepository) GetNetworks() []dtos.StreamsPerNetwork {
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

func (db *streamsRepository) GetStations() []dtos.StreamsPerStation {
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

func (db *streamsRepository) GetTotalStreams() int {
	var count int64
	db.connection.Model(
		&entities.Stream{},
	).Count(&count)
	return int(count)
}

func (db *streamsRepository) GetTotalStations() int {
	var count int64
	db.connection.Model(
		&entities.Station{},
	).Count(&count)
	return int(count)
}

func (db *streamsRepository) GetTotalNetworks() int {
	var count int64
	db.connection.Model(
		&entities.Network{},
	).Count(&count)
	return int(count)
}

func (db *streamsRepository) GetStreams() []entities.Stream {
	var streams []entities.Stream

	db.connection.Model(
		&entities.Stream{},
	).Find(&streams)

	return streams
}

func (db *streamsRepository) GetStreamWithAssociatedData(streamId uint64) (entities.Stream, error) {
	var stream entities.Stream

	result := db.connection.Model(
		&entities.Stream{},
	).Joins("Station").Joins("Procedure").Joins("Unit").Joins("Variable").Where(
		"streams.stream_id = ?", streamId,
	).Find(&stream)

	if result.Error != nil {
		log.Errorf("Error executing GetStreamWithAssociatedData query: %v", result.Error)
		return stream, result.Error
	}

	if result.RowsAffected == 0 {
		return stream, errors.Join(exceptions.NewNotFound(), fmt.Errorf("stream with id %v not found", streamId))
	}

	return stream, nil
}

func (db *streamsRepository) CreateUnit(unit entities.Unit) error {
	result := db.connection.Save(&unit)
	return result.Error
}

func (db *streamsRepository) CreateStation(station entities.Station) error {
	result := db.connection.Save(&station)
	return result.Error
}

func (db *streamsRepository) CreateProcedure(procedure entities.Procedure) error {
	result := db.connection.Save(&procedure)
	return result.Error
}

func (db *streamsRepository) CreateVariable(variable entities.Variable) error {
	result := db.connection.Save(&variable)
	return result.Error
}

func (db *streamsRepository) CreateStream(stream entities.Stream) error {
	result := db.connection.Save(&stream)
	return result.Error
}
