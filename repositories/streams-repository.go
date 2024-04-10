package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const directMeasurementProcedure = 1
const waterLevelVariable = 2

type StreamRepository interface {
	GetNetworks() []dtos.StreamsPerNetwork
	GetStations(configId uint64) *[]*dtos.StreamsPerStation
	GetErrorsOfStations(configId uint64, timeStart time.Time, timeEnd time.Time) []dtos.ErrorsOfStations
	GetTotalStreams(configurationId uint64) int
	//GetTotalNetworks(configurationId uint64) int
	GetTotalStations(configurationId uint64) int
	GetStreams() []entities.Stream
	GetStreamWithAssociatedData(streamId uint64) (entities.Stream, error)
	CreateUnit(entity entities.Unit) error
	CreateStation(entity entities.Station) error
	CreateProcedure(entity entities.Procedure) error
	CreateVariable(entity entities.Variable) error
	CreateStream(entity entities.Stream) error
	GetStreamCards(parameters dtos.QueryParameters) (*dtos.StreamCardsResponse, error)
	GetStreamsForOutputMetrics(configId uint64) ([]dtos.BehaviourStream, error)
	GetTotalStreamsWithNullValues(configId uint64, timeStart time.Time, timeEnd time.Time) int
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

func (db *streamsRepository) GetStations(configId uint64) *[]*dtos.StreamsPerStation {
	var stations *[]*dtos.StreamsPerStation

	tx := db.connection.Model(
		&entities.Stream{},
	).Select(
		"stations.name as station_name",
		"stations.station_id as station_id",
		"count(streams.stream_id) as streams_count",
	).Joins(
		"JOIN stations ON  stations.station_id = streams.station_id",
	).Joins(
		"JOIN configured_streams ON configured_streams.stream_id = streams.stream_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Group(
		"stations.name, stations.station_id",
	).Scan(&stations)

	if tx.Error != nil {
		log.Errorf("Error executing GetStations query: %v", tx.Error)
	}

	log.Debugf("Get stations query result: %v", stations)
	return stations
}

func (db *streamsRepository) GetErrorsOfStations(configId uint64, timeStart time.Time, timeEnd time.Time) []dtos.ErrorsOfStations {
	var errorsPerStation []dtos.ErrorsOfStations

	tx := db.connection.Model(
		&entities.Stream{},
	).Select(
		"streams.station_id as station_id",
		"count(detected_errors.error_id) as error_count",
	).Joins(
		"JOIN configured_streams ON configured_streams.stream_id = streams.stream_id",
	).Joins(
		"JOIN detected_errors ON detected_errors.stream_id = streams.stream_id ",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"detected_errors.detected_date >= ? AND detected_errors.detected_date <= ?", timeStart, timeEnd,
	).Group("streams.station_id").Scan(&errorsPerStation)

	if tx.Error != nil {
		log.Errorf("Error executing GetErrorsOfStations query: %v", tx.Error)
	}

	return errorsPerStation
}
func (db *streamsRepository) GetTotalStreams(configurationId uint64) int {
	var count int64
	db.connection.Model(
		&entities.Stream{},
	).Select("COUNT(streams.stream_id)").Joins(
		"JOIN configured_streams ON configured_streams.stream_id = streams.stream_id",
	).Where(
		"configured_streams.configuration_id = ?", configurationId,
	).Find(&count)
	return int(count)
}

func (db *streamsRepository) GetTotalStations(configurationId uint64) int {
	var count int64
	db.connection.Model(
		&entities.Stream{},
	).Select("COUNT(streams.station_id)").Joins(
		"JOIN configured_streams ON configured_streams.stream_id = streams.stream_id",
	).Where(
		"configured_streams.configuration_id = ?", configurationId,
	).Group("streams.station_id").Find(&count)
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

func (db *streamsRepository) GetTotalStreamsWithNullValues(configId uint64, timeStart time.Time, timeEnd time.Time) int {
	var count int64
	db.connection.Model(
		&entities.ConfiguredStream{},
	).Select("COUNT(detected_errors.error_id)").Joins(
		"JOIN streams ON streams.stream_id = configured_streams.stream_id",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.configured_stream_configured_stream_id=configured_streams.configured_stream_id",
	).Joins(
		"JOIN detected_errors ON detected_errors.error_id = configured_streams_errors.detected_error_error_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"detected_errors.error_type = ?", entities.NullValue,
	).Where(
		"detected_errors.detected_date >= ? AND detected_errors.detected_date <= ?", timeStart, timeEnd,
	).Where(
		"streams.stream_type = ?", entities.Observed,
	).Find(&count)
	return int(count)
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

func (db *streamsRepository) GetStreamsForOutputMetrics(configId uint64) ([]dtos.BehaviourStream, error) {
	var streams []dtos.BehaviourStream
	tx := db.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"streams.stream_id as stream_id",
		"stations.alert_level as alert_level",
		"stations.evacuation_level as evacuation_level",
		"stations.low_water_level as low_water_level",
	).Joins(
		"JOIN streams ON streams.stream_id=configured_streams.stream_id",
	).Joins(
		"JOIN stations ON stations.station_id=streams.station_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"streams.variable_id = ?", waterLevelVariable, // TODO validar
	).Where(
		"streams.procedure_id = ?", directMeasurementProcedure, // TODO validar
	).Where(
		"streams.stream_type = ?", entities.Observed, // TODO validar
	).Find(&streams)

	if tx.Error != nil {
		log.Errorf("Error executing GetStreamsForOutputMetrics query: %v", tx.Error)
		return nil, tx.Error
	}

	return streams, nil
}
