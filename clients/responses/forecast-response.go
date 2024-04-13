package responses

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
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
}

func (f LastForecast) GetForecastOfStream(streamId uint64) *Forecast {
	forecast := &Forecast{}
	for _, stream := range f.Streams {
		if stream.StreamId == streamId {
			if stream.Qualifier == "main" {
				forecast.MainForecast = stream
			} else if stream.Qualifier == "p05" {
				forecast.P05Forecast = stream
			} else if stream.Qualifier == "p95" {
				forecast.P95Forecast = stream
			}
		}
	}
	return forecast
}

func (f LastForecast) GetMainForecast() []float64 {
	for _, stream := range f.Streams {
		if stream.Qualifier == "main" {
			var values []float64
			for _, forecast := range stream.Forecasts {
				value, _ := strconv.ParseFloat(forecast[2], 64)
				values = append(values, value)
			}
			return values
		}
	}
	return []float64{}
}

func (f LastForecast) ConvertToCalibratedStreamsDataResponse() dtos.CalibratedStreamsDataResponse {
	var P05Streams = convertToCalibratedStreamsData(f, "p05")
	var MainStreams = convertToCalibratedStreamsData(f, "main")
	var P75Streams = convertToCalibratedStreamsData(f, "p75")
	var P95Streams = convertToCalibratedStreamsData(f, "p95")
	var P25Streams = convertToCalibratedStreamsData(f, "p25")

	return dtos.CalibratedStreamsDataResponse{
		P05Streams,
		MainStreams,
		P75Streams,
		P95Streams,
		P25Streams,
	}
}

func convertToCalibratedStreamsData(f LastForecast, qualifier string) []dtos.CalibratedStreamsData {
	var calibratedStreams []dtos.CalibratedStreamsData

	for _, stream := range f.Streams {
		if stream.Qualifier == qualifier {
			for _, forecast := range stream.Forecasts {
				value, _ := strconv.ParseFloat(forecast[2], 64)
				date, _ := time.Parse("2006-01-02T15:04:05Z07:00", forecast[0])
				calibratedStreams = append(calibratedStreams, dtos.CalibratedStreamsData{
					Time:      date,
					Value:     value,
					Qualifier: forecast[3],
				})
			}
		}
	}

	return calibratedStreams
}
