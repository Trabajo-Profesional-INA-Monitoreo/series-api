package dtos

import (
	"time"
)

type DelayPerDay struct {
	Average float64   `gorm:"column:average" json:"Average"`
	Date    time.Time `gorm:"column:date" json:"Date"`
}
