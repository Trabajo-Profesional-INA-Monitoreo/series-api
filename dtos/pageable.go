package dtos

import "math"

type Pageable struct {
	Pages         int
	TotalElements int
	Page          int
	PageSize      int
}

func NewPageable(totalElements int, page int, pageSize int) Pageable {
	return Pageable{Pages: int(math.Ceil(float64(totalElements) / float64(pageSize))), TotalElements: totalElements, Page: page, PageSize: pageSize}
}
