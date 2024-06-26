package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type ErrorDetectionRepository interface {
	AlreadyDetectedErrorForStreamWithIdAndType(uint64, string, entities.ErrorType) bool
	GetDetectedErrorForStreamWithIdAndType(streamId uint64, requestId string, errorType entities.ErrorType) entities.DetectedError
	Create(detectedError entities.DetectedError)
	Update(detectedError entities.DetectedError)
}

type ErrorMetricsRepository interface {
	GetErrorsPerDay(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorsCountPerDayAndType
	GetErrorIndicators(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorIndicator
	GetRelatedStreamsToError(parameters *dtos.QueryParameters) ([]dtos.ErrorRelatedStream, error)
	GetErrorsOfConfiguredStream(parameters *dtos.QueryParameters) (*dtos.DetectedErrorsOfStream, error)
	GetAverageDelayPerDay(start time.Time, end time.Time, configuredStreamId uint64) []*dtos.DelayPerDay
}

type ErrorsRepository interface {
	ErrorDetectionRepository
	ErrorMetricsRepository
}

type errorsRepository struct {
	connection *gorm.DB
}

func NewErrorsRepository(connection *gorm.DB) ErrorsRepository {
	return &errorsRepository{connection}
}

func (e errorsRepository) AlreadyDetectedErrorForStreamWithIdAndType(streamId uint64, requestId string, errorType entities.ErrorType) bool {
	var detectedError entities.DetectedError
	err := e.connection.Model(&entities.DetectedError{}).Where("stream_id = ? AND request_id = ? AND error_type = ?", streamId, requestId, errorType).First(&detectedError)
	if err.Error != nil {
		log.Errorf("Error executing AlreadyDetectedErrorForStreamWithIdAndType query: %v", err)
	}
	return detectedError.RequestId == requestId
}

func (e errorsRepository) Create(detectedError entities.DetectedError) {
	err := e.connection.Create(&detectedError)
	if err.Error != nil {
		log.Errorf("Error executing Create query: %v", err)
	}
}

func (e errorsRepository) Update(detectedError entities.DetectedError) {
	err := e.connection.Save(&detectedError)
	if err.Error != nil {
		log.Errorf("Error executing Update query: %v", err)
	}
}

func (e errorsRepository) GetDetectedErrorForStreamWithIdAndType(streamId uint64, requestId string, errorType entities.ErrorType) entities.DetectedError {
	var detectedError entities.DetectedError
	e.connection.Model(&entities.DetectedError{}).Where("stream_id = ? AND request_id = ? AND error_type = ?", streamId, requestId, errorType).First(&detectedError)
	return detectedError
}

func (e errorsRepository) GetErrorsPerDay(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorsCountPerDayAndType {
	var results []*dtos.ErrorsCountPerDayAndType
	err := e.connection.Model(
		&entities.DetectedError{},
	).Select(
		"detected_errors.error_type as error_type",
		"COUNT(detected_errors.error_type) as total",
		"DATE(detected_errors.detected_date) as date",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.detected_error_error_id = detected_errors.error_id",
	).Joins(
		"JOIN configured_streams ON configured_streams.configured_stream_id = configured_streams_errors.configured_stream_configured_stream_id",
	).Where(
		"detected_date >= ? AND detected_date <= ?", timeStart, timeEnd,
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"configured_streams.deleted = false",
	).Group(
		"DATE(detected_errors.detected_date)",
	).Group(
		"detected_errors.error_type",
	).Find(&results)
	if err.Error != nil {
		log.Errorf("Error executing GetErrorsPerDay query: %v", err)
	}
	return results
}

func (e errorsRepository) GetErrorIndicators(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorIndicator {
	var results []*dtos.ErrorIndicator
	err := e.connection.Model(
		&entities.DetectedError{},
	).Select(
		"detected_errors.error_type as error_type",
		"COUNT(detected_errors.error_id) as count",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.detected_error_error_id=detected_errors.error_id",
	).Joins(
		"JOIN configured_streams ON configured_streams.configured_stream_id=configured_streams_errors.configured_stream_configured_stream_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"configured_streams.deleted = false",
	).Where(
		"detected_date >= ? AND detected_date <= ?", timeStart, timeEnd,
	).Group(
		"detected_errors.error_type",
	).Find(&results)
	if err.Error != nil {
		log.Errorf("Error executing GetErrorIndicators query: %v", err)
	}
	return results
}

func (e errorsRepository) GetRelatedStreamsToError(parameters *dtos.QueryParameters) ([]dtos.ErrorRelatedStream, error) {
	timeStart := parameters.Get("timeStart").(time.Time)
	timeEnd := parameters.Get("timeEnd").(time.Time)
	configurationId := parameters.Get("configurationId").(uint64)
	errorId := parameters.Get("errorType").(uint64)
	var results []dtos.ErrorRelatedStream
	tx := e.connection.Model(
		&entities.DetectedError{},
	).Select(
		"streams.stream_id as stream_id",
		"streams.station_id as station_id",
		"stations.name as station_name",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.detected_error_error_id=detected_errors.error_id",
	).Joins(
		"JOIN configured_streams ON configured_streams.configured_stream_id = configured_streams_errors.configured_stream_configured_stream_id",
	).Joins(
		"JOIN streams ON streams.stream_id = configured_streams.stream_id",
	).Joins(
		"JOIN stations ON stations.station_id = streams.station_id",
	).Where(
		"configured_streams.configuration_id = ?", configurationId,
	).Where(
		"configured_streams.deleted = false",
	).Where(
		"detected_errors.error_type = ?", errorId,
	).Where(
		"detected_errors.detected_date >= ? AND detected_errors.detected_date <= ?", timeStart, timeEnd,
	).Group("streams.stream_id").Group("streams.station_id").Group("stations.name").Find(&results)
	if tx.Error != nil {
		log.Errorf("Error executing GetRelatedStreamsToError query: %v", tx.Error)
		return nil, tx.Error
	}
	return results, nil
}

func (e errorsRepository) GetErrorsOfConfiguredStream(parameters *dtos.QueryParameters) (*dtos.DetectedErrorsOfStream, error) {
	timeStart := parameters.Get("timeStart").(time.Time)
	timeEnd := parameters.Get("timeEnd").(time.Time)
	configStreamId := parameters.Get("configStreamId").(uint64)
	pageSize := *parameters.GetAsInt("pageSize")
	page := *parameters.GetAsInt("page")

	var errors []*dtos.ErrorDto
	tx := e.connection.Model(
		&entities.DetectedError{},
	).Select(
		"detected_errors.error_id as error_id",
		"detected_errors.detected_date as detected_date",
		"detected_errors.error_type as error_type_id",
		"detected_errors.extra_info as extra_info",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.detected_error_error_id=detected_errors.error_id",
	).Where(
		"configured_streams_errors.configured_stream_configured_stream_id = ?", configStreamId,
	).Where(
		"detected_errors.detected_date >= ? AND detected_errors.detected_date <= ?", timeStart, timeEnd,
	).Limit(pageSize).Offset((page - 1) * pageSize).Order("detected_errors.detected_date desc").Find(&errors)

	if tx.Error != nil {
		log.Errorf("Error executing GetErrorsOfConfiguredStream Count query: %v", tx.Error)
		return nil, tx.Error
	}

	var totalElements int
	countTx := e.connection.Model(
		&entities.DetectedError{},
	).Select(
		"count(detected_errors.error_id)",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.detected_error_error_id=detected_errors.error_id",
	).Where(
		"configured_streams_errors.configured_stream_configured_stream_id = ?", configStreamId,
	).Where(
		"detected_errors.detected_date >= ? AND detected_errors.detected_date <= ?", timeStart, timeEnd,
	).Find(&totalElements)

	if countTx.Error != nil {
		log.Errorf("Error executing GetErrorsOfConfiguredStream Count query: %v", countTx.Error)
		return nil, countTx.Error
	}

	return dtos.NewDetectedErrorsOfStream(errors, dtos.NewPageable(totalElements, page, pageSize)), nil
}

func (e errorsRepository) GetAverageDelayPerDay(timeStart time.Time, timeEnd time.Time, configuredStreamId uint64) []*dtos.DelayPerDay {
	var results []*dtos.DelayPerDay

	err := e.connection.Model(
		&entities.DetectedError{},
	).Select(
		"AVG(to_number(REGEXP_REPLACE(detected_errors.extra_info, 'Delay total en minutos: ', ''), '99999999D999S')) as average",
		"DATE(detected_errors.detected_date) as date",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.detected_error_error_id = detected_errors.error_id",
	).Joins(
		"JOIN configured_streams ON configured_streams.configured_stream_id = configured_streams_errors.configured_stream_configured_stream_id",
	).Where(
		"detected_date >= ? AND detected_date <= ?", timeStart, timeEnd,
	).Where(
		"configured_streams.configured_stream_id = ?", configuredStreamId,
	).Where(
		"detected_errors.error_type = ?", entities.Delay,
	).Where(
		"configured_streams.deleted = false",
	).Group(
		"DATE(detected_errors.detected_date)",
	).Order(
		"DATE(detected_errors.detected_date)",
	).Find(&results)

	if err.Error != nil {
		log.Errorf("Error executing GetAverageDelayPerDay query: %v", err)
	}

	return results
}
