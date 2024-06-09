package entities

import (
	"fmt"
	"time"
)

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

func (error DetectedError) ToString() string {
	var errorType = error.ErrorType.Translate()
	return fmt.Sprintf("Se detectó un error de tipo: %v \nEl error se detectó a las: %v en la serie %v ", errorType, error.DetectedDate.Format("2006-01-02 15:04:05"), error.Stream.StreamId)
}
