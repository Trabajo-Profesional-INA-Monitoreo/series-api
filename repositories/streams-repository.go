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
	GetStreamCards(parameters dtos.StreamCardsParameters) (*dtos.StreamCardsResponse, error)
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
	).Joins("Network").Joins("Station").Joins("Procedure").Joins("Unit").Joins("Variable").Where(
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

func (db *streamsRepository) GetStreamCards(parameters dtos.StreamCardsParameters) (*dtos.StreamCardsResponse, error) {
	var streamCards []dtos.StreamCard
	result := db.connection.Model(
		&entities.ConfiguredStream{},
	).Joins("Stream").Joins("Station").Joins("Procedure").Joins("Variable").Where(
		"configured_streams.configured_stream_id = ?", parameters.ConfigurationId,
	).Where(
		"(? IS NULL OR streams.stream_id = ?)", parameters.StreamId,
	).Where(
		"(? IS NULL OR stations.station_id = ?) ", parameters.StationId,
	).Where(
		"(? IS NULL OR procedures.procedure_id = ?)", parameters.ProcId,
	).Where(
		"(? IS NULL OR streams.variable_id = ?)", parameters.VarId,
	).Where(
		"(? IS NULL OR streams.stream_type = ?)", parameters.StreamType,
	).Limit(parameters.PageSize).Offset(parameters.Page * parameters.PageSize).Find(&streamCards)

	var totalElements int
	result = db.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"count(configured_streams.configured_stream_id)").Joins("Stream").Where(
		"configured_streams.configured_stream_id = ?", parameters.ConfigurationId,
	).Where(
		"(? IS NULL OR streams.stream_id = ?)", parameters.StreamId,
	).Where(
		"(? IS NULL OR stations.station_id = ?) ", parameters.StationId,
	).Where(
		"(? IS NULL OR procedures.procedure_id = ?)", parameters.ProcId,
	).Where(
		"(? IS NULL OR streams.variable_id = ?)", parameters.VarId,
	).Where(
		"(? IS NULL OR streams.stream_type = ?)", parameters.StreamType,
	).Find(&totalElements)

	if result.Error != nil {
		log.Errorf("Error executing GetStreamWithAssociatedData query: %v", result.Error)
		return nil, result.Error
	}

	return dtos.NewStreamCardsResponse(streamCards, dtos.NewPageable(totalElements, parameters.Page, parameters.PageSize)), nil
}
