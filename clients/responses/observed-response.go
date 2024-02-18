package responses

import "time"

type ObservedDataResponse struct {
	DataType   string    `json:"tipo"`
	StreamId   uint64    `json:"series_id"`
	TimeStart  time.Time `json:"timestart"`
	TimeEnd    time.Time `json:"timeend"`
	TimeUpdate time.Time `json:"timeupdate"`
	Value      float64   `json:"valor"`
	DataId     string    `json:"id"`
}
