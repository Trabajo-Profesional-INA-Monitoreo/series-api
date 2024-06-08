package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"gorm.io/gorm"
	"time"
)

type InputsRepository interface {
	GetTotalStreams(configurationId uint64) int
	GetTotalStations(configurationId uint64) int
	GetTotalStreamsByError(id uint64, start time.Time, end time.Time, value entities.ErrorType) (int, []int64)
}

type inputsRepository struct {
	connection *gorm.DB
}

func NewInputsRepository(connection *gorm.DB) InputsRepository {
	return &inputsRepository{connection}
}

func (db *inputsRepository) GetTotalStreams(configurationId uint64) int {
	var count int64
	db.connection.Model(
		&entities.Stream{},
	).Select("COUNT(streams.stream_id)").Joins(
		"JOIN configured_streams ON configured_streams.stream_id = streams.stream_id",
	).Where(
		"configured_streams.configuration_id = ?", configurationId,
	).Where(
		"configured_streams.deleted = false",
	).Find(&count)
	return int(count)
}

func (db *inputsRepository) GetTotalStations(configurationId uint64) int {
	var count int64
	db.connection.Model(
		&entities.Stream{},
	).Select("COUNT(streams.station_id)").Joins(
		"JOIN configured_streams ON configured_streams.stream_id = streams.stream_id",
	).Where(
		"configured_streams.configuration_id = ?", configurationId,
	).Where(
		"configured_streams.deleted = false",
	).Group("streams.station_id").Find(&count)
	return int(count)
}

func (db *inputsRepository) GetTotalStreamsByError(configId uint64, timeStart time.Time, timeEnd time.Time, error entities.ErrorType) (int, []int64) {
	var streams []int64
	db.connection.Model(
		&entities.ConfiguredStream{},
	).Select("streams.stream_id").Joins(
		"JOIN streams ON streams.stream_id = configured_streams.stream_id",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.configured_stream_configured_stream_id=configured_streams.configured_stream_id",
	).Joins(
		"JOIN detected_errors ON detected_errors.error_id = configured_streams_errors.detected_error_error_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"detected_errors.error_type = ?", error,
	).Where(
		"detected_errors.detected_date >= ? AND detected_errors.detected_date <= ?", timeStart, timeEnd,
	).Where(
		"streams.stream_type = ?", entities.Observed,
	).Where(
		"configured_streams.deleted = false",
	).Group("streams.stream_id").Find(&streams)
	return len(streams), streams
}
