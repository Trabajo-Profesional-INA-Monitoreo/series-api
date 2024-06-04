package entities

import "time"

type DetectedError struct {
	ErrorId          uint64 `gorm:"primary_key;auto_increment"`
	StreamId         uint64
	Stream           *Stream            `gorm:"references:StreamId"`
	ConfiguredStream []ConfiguredStream `gorm:"many2many:configured_streams_errors;"` // It creates an intermediate table with AutoCreate
	DetectedDate     time.Time
	RequestId        string    `gorm:"type:varchar(254)"` // A unique identifier that we can get from the requests to INA
	ErrorType        ErrorType `gorm:"index"`
	ExtraInfo        string
}
