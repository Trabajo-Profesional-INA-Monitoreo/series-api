package entities

import "time"

type DetectedError struct {
	ErrorId            uint64 `gorm:"primary_key;auto_increment"`
	StreamId           uint64
	Stream             *Stream `gorm:"references:StreamId"`
	ConfiguredStreamId uint64
	ConfiguredStream   *ConfiguredStream `gorm:"references:ConfiguredStreamId"`
	DetectedDate       time.Time
	RequestId          string `gorm:"type:varchar(254)"` // A unique identifier that we can get from the requests to INA
	ErrorType          ErrorType
}
