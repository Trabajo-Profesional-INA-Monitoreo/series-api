package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type StationsRepository interface {
	GetStations(configId uint64, page int, pageSize int) (*[]*dtos.StreamsPerStation, dtos.Pageable)
	GetErrorsOfStations(configId uint64, timeStart time.Time, timeEnd time.Time, stationIds []uint64) []dtos.ErrorsOfStations
}

type stationsRepository struct {
	connection *gorm.DB
}

func NewStationsRepository(connection *gorm.DB) StationsRepository {
	return &stationsRepository{connection}
}

func (db *stationsRepository) GetStations(configId uint64, page int, pageSize int) (*[]*dtos.StreamsPerStation, dtos.Pageable) {
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
	).Where(
		"configured_streams.deleted = false",
	).Group(
		"stations.name, stations.station_id",
	).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&stations)

	if tx.Error != nil {
		log.Errorf("Error executing GetStations query: %v", tx.Error)
	}

	var totalElements int
	countTx := db.connection.Model(
		&entities.Stream{},
	).Select(
		"COUNT(DISTINCT stations.station_id)",
	).Joins(
		"JOIN stations ON  stations.station_id = streams.station_id",
	).Joins(
		"JOIN configured_streams ON configured_streams.stream_id = streams.stream_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"configured_streams.deleted = false",
	).Find(&totalElements)

	if countTx.Error != nil {
		log.Errorf("Error executing GetStations Count query: %v", countTx.Error)
	}

	log.Debugf("Get stations query result: %v", stations)
	return stations, dtos.NewPageable(totalElements, page, pageSize)
}

func (db *stationsRepository) GetErrorsOfStations(configId uint64, timeStart time.Time, timeEnd time.Time, stationIds []uint64) []dtos.ErrorsOfStations {
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
	).Where(
		"configured_streams.deleted = false",
	).Where(
		"streams.station_id IN ?", stationIds,
	).Group("streams.station_id").Scan(&errorsPerStation)

	if tx.Error != nil {
		log.Errorf("Error executing GetErrorsOfStations query: %v", tx.Error)
	}

	return errorsPerStation
}
