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

const directMeasurementProcedure = 1
const waterLevelVariable = 2

type FaultDetectionStreamsRepository interface {
	GetStreams() []entities.Stream
}

type ManagerStreamsRepository interface {
	GetStreamWithAssociatedData(streamId uint64) (entities.Stream, error)
	CreateUnit(entity entities.Unit) error
	CreateStation(entity entities.Station) error
	CreateProcedure(entity entities.Procedure) error
	CreateVariable(entity entities.Variable) error
	CreateStream(entity entities.Stream) error
	GetStreamCards(parameters dtos.QueryParameters) (*dtos.StreamCardsResponse, error)
	GetRedundancies(configuredStreamId uint64) []int
}

type StreamRepository interface {
	FaultDetectionStreamsRepository
	ManagerStreamsRepository
}

type streamsRepository struct {
	connection *gorm.DB
}

func NewStreamRepository(connection *gorm.DB) StreamRepository {
	return &streamsRepository{connection}
}

func (db *streamsRepository) GetStreams() []entities.Stream {
	var streams []entities.Stream

	db.connection.Model(
		&entities.Stream{},
	).Where("streams.stream_type != ?", entities.Curated).Find(&streams)

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

func (db *streamsRepository) GetStreamCards(parameters dtos.QueryParameters) (*dtos.StreamCardsResponse, error) {
	var streamCards []*dtos.StreamCard
	configId := parameters.Get("configurationId").(uint64)
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
		"configured_streams.calibration_id as calibration_id",
		"configured_streams.observed_related_stream_id as observed_related_stream_id",
		"streams.procedure_id as procedure_id",
		"streams.stream_type as stream_type",
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
	).Where(
		"configured_streams.deleted = false",
	)

	var totalElements int
	countTx := db.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"count(configured_streams.configured_stream_id)",
	).Joins(
		"JOIN streams ON configured_streams.stream_id=streams.stream_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"configured_streams.deleted = false",
	)

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
	tx.Limit(pageSize).Offset((page - 1) * pageSize).Find(&streamCards)

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

func (db *streamsRepository) GetRedundancies(configuredStreamId uint64) []int {
	var redundancies []int

	tx := db.connection.Model(
		&entities.Redundancy{},
	).Select(
		"redundancies.redundancy_id as redundancy_id",
	).Where(
		"redundancies.configured_stream_id = ?", configuredStreamId,
	).Scan(&redundancies)

	if tx.Error != nil {
		log.Errorf("Error executing GetRedundancies query: %v", tx.Error)
	}

	log.Debugf("Get redundancies query result: %v", redundancies)
	return redundancies
}
