package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
)

type FilterService interface {
	GetProcedures() []dtos.ProcedureFilter
	GetStations() []dtos.StationFilter
	GetVariables() []dtos.VariableFilter
}

type filterService struct {
	repository repositories.FilterRepository
}

func (f filterService) GetProcedures() []dtos.ProcedureFilter {
	return f.repository.GetProcedures()
}

func (f filterService) GetStations() []dtos.StationFilter {
	return f.repository.GetStations()
}

func (f filterService) GetVariables() []dtos.VariableFilter {
	return f.repository.GetVariables()
}

func NewFilterService(repository repositories.FilterRepository) FilterService {
	return &filterService{repository: repository}
}
