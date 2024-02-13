package responses

import "time"

type LastForecast struct {
	RunId         uint64             `json:"cor_id"`
	CalibrationId uint64             `json:"cal_id"`
	ForecastDate  time.Time          `json:"forecast_date"`
	Streams       []ForecastedStream `json:"series"`
}

type ForecastedStream struct {
	StreamId  uint64     `json:"series_id"`
	Qualifier string     `json:"qualifier"`
	Forecasts [][]string `json:"pronosticos"`
}
