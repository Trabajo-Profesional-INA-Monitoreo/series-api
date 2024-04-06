package dtos

type Pageable struct {
	Pages         int
	TotalElements int
	Page          int
	PageSize      int
}

func NewPageable(totalElements int, page int, pageSize int) Pageable {
	return Pageable{Pages: totalElements / pageSize, TotalElements: totalElements, Page: page, PageSize: pageSize}
}
