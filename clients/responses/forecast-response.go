package responses

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type LastForecast struct {
	RunId         uint64              `json:"cor_id"`
	CalibrationId uint64              `json:"cal_id"`
	ForecastDate  time.Time           `json:"forecast_date"`
	Streams       []*ForecastedStream `json:"series"`
}

type ForecastedStream struct {
	StreamId  uint64     `json:"series_id"`
	Qualifier string     `json:"qualifier"`
	Forecasts [][]string `json:"pronosticos"`
}

type Forecast struct {
	MainForecast *ForecastedStream
	P05Forecast  *ForecastedStream
	P95Forecast  *ForecastedStream
	P01Forecast  *ForecastedStream
	P99Forecast  *ForecastedStream
}

func (f LastForecast) GetForecastOfStream(streamId uint64) *Forecast {
	forecast := &Forecast{}
	for _, stream := range f.Streams {
		if stream.StreamId == streamId {
			if stream.Qualifier == "main" {
				forecast.MainForecast = stream
			} else if stream.Qualifier == "p05" || stream.Qualifier == "error_band_05" {
				forecast.P05Forecast = stream
			} else if stream.Qualifier == "p95" || stream.Qualifier == "error_band_95" {
				forecast.P95Forecast = stream
			} else if stream.Qualifier == "p99" || stream.Qualifier == "error_band_99" {
				forecast.P99Forecast = stream
			} else if stream.Qualifier == "p01" || stream.Qualifier == "error_band_01" {
				forecast.P01Forecast = stream
			}
		}
	}
	if forecast.MainForecast == nil {
		log.Warnf("Requested forecast for stream %v but it was not found on request, check if config is correct", streamId)
	}
	return forecast
}

func ConvertToFloats(forecast [][]string) []float64 {
	var values []float64
	for _, forecast := range forecast {
		value, _ := strconv.ParseFloat(forecast[2], 64)
		values = append(values, value)
	}
	return values
}

func (f LastForecast) ConvertToCalibratedStreamsDataResponse(streamId uint64) dtos.CalibratedStreamsDataResponse {
	var P05Streams = convertToCalibratedStreamsData(f, "p05", streamId, "error_band_05")
	var MainStreams = convertToCalibratedStreamsData(f, "main", streamId, "")
	var P75Streams = convertToCalibratedStreamsData(f, "p75", streamId, "error_band_75")
	var P95Streams = convertToCalibratedStreamsData(f, "p95", streamId, "error_band_95")
	var P25Streams = convertToCalibratedStreamsData(f, "p25", streamId, "error_band_25")
	var P99Streams = convertToCalibratedStreamsData(f, "p99", streamId, "error_band_99")
	var P01Streams = convertToCalibratedStreamsData(f, "p01", streamId, "error_band_01")

	return dtos.CalibratedStreamsDataResponse{
		P05Streams,
		MainStreams,
		P75Streams,
		P95Streams,
		P25Streams,
		P99Streams,
		P01Streams,
	}
}

func convertToCalibratedStreamsData(f LastForecast, qualifier string, streamId uint64, alternativeQualifier string) []dtos.CalibratedStreamsData {
	var calibratedStreams []dtos.CalibratedStreamsData

	for _, stream := range f.Streams {
		if stream.StreamId == streamId && (stream.Qualifier == qualifier || stream.Qualifier == alternativeQualifier) {
			for _, forecast := range stream.Forecasts {
				value, _ := strconv.ParseFloat(forecast[2], 64)
				date, _ := time.Parse("2006-01-02T15:04:05Z07:00", forecast[0])
				calibratedStreams = append(calibratedStreams, dtos.CalibratedStreamsData{
					Time:  date,
					Value: value,
				})
			}
		}
	}

	return calibratedStreams
}
