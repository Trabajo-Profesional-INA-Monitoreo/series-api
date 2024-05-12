package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
)

type FilterService interface {
	GetProcedures(configId uint64) []dtos.FilterValue
	GetStations(configId uint64) []dtos.FilterValue
	GetVariables(configId uint64) []dtos.FilterValue
	GetNodes(configId uint64) []dtos.FilterValue
}

type filterService struct {
	repository repositories.FilterRepository
}

func (f filterService) GetProcedures(configId uint64) []dtos.FilterValue {
	return f.repository.GetProcedures(configId)
}

func (f filterService) GetNodes(configId uint64) []dtos.FilterValue {
	return f.repository.GetNodes(configId)
}

func (f filterService) GetStations(configId uint64) []dtos.FilterValue {
	return f.repository.GetStations(configId)
}

func (f filterService) GetVariables(configId uint64) []dtos.FilterValue {
	return f.repository.GetVariables(configId)
}

func NewFilterService(repository repositories.FilterRepository) FilterService {
	return &filterService{repository: repository}
}
