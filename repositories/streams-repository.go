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
	var streamCards []*dtos.StreamCard
	configId := parameters.GetAsInt("configurationId")
	streamId := parameters.GetAsInt("streamId")
	stationId := parameters.GetAsInt("stationId")
	procId := parameters.GetAsInt("procId")
	varId := parameters.GetAsInt("varId")
	streamType := parameters.GetAsInt("streamType")
	pageSize := *parameters.GetAsInt("pageSize")
	page := *parameters.GetAsInt("page")

	tx := db.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"configured_streams.configured_stream_id as configured_stream_id",
		"configured_streams.stream_id as stream_id",
		"configured_streams.check_errors as check_errors",
		"streams.procedure_id as procedure_id",
		"procedures.name as procedure_name",
		"streams.variable_id as variable_id",
		"variables.name as variable_name",
		"streams.station_id as station_id",
		"stations.name as station_name",
	).Joins(
		"JOIN streams ON streams.stream_id=configured_streams.stream_id",
	).Joins(
		"JOIN stations ON stations.station_id=streams.station_id",
	).Joins(
		"JOIN procedures ON procedures.procedure_id=streams.procedure_id",
	).Joins(
		"JOIN variables ON variables.variable_id=streams.variable_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	)

	var totalElements int
	countTx := db.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"count(configured_streams.configured_stream_id)",
	).Joins(
		"JOIN streams ON configured_streams.stream_id=streams.stream_id",
	).Where("configured_streams.configuration_id = ?", configId)

	if streamId != nil {
		tx.Where("streams.stream_id = ?", streamId)
		countTx.Where("streams.stream_id = ?", streamId)
	}
	if stationId != nil {
		tx.Where("stations.station_id = ? ", stationId)
		countTx.Where("streams.station_id = ?", stationId)
	}
	if procId != nil {
		tx.Where("procedures.procedure_id = ?", procId)
		countTx.Where("streams.procedure_id = ?", procId)
	}
	if varId != nil {
		tx.Where("streams.variable_id = ?", varId)
		countTx.Where("streams.variable_id = ?", varId)
	}
	if streamType != nil {
		tx.Where("streams.stream_type = ?", streamType)
		countTx.Where("streams.stream_type = ?", streamType)
	}
	tx.Limit(pageSize).Offset(page * pageSize).Find(&streamCards)

	if tx.Error != nil {
		log.Errorf("Error executing GetStreamCards query: %v", tx.Error)
		return nil, tx.Error
	}

	countTx.Find(&totalElements)

	if countTx.Error != nil {
		log.Errorf("Error executing GetStreamCards Count query: %v", countTx.Error)
		return nil, countTx.Error
	}

	return dtos.NewStreamCardsResponse(streamCards, dtos.NewPageable(totalElements, page, pageSize)), nil
}
