package dtos

import "math"

type Pageable struct {
	Pages         int `json:"Pages"`
	TotalElements int `json:"TotalElements"`
	Page          int `json:"Page"`
	PageSize      int `json:"PageSize"`
}

func NewPageable(totalElements int, page int, pageSize int) Pageable {
	return Pageable{Pages: int(math.Ceil(float64(totalElements) / float64(pageSize))), TotalElements: totalElements, Page: page, PageSize: pageSize}
}
