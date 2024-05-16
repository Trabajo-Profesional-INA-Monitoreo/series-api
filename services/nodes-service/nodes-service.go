package nodes_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	metrics_service "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/metrics-service"
	log "github.com/sirupsen/logrus"
	"time"
)

type NodesService interface {
	GetNodes(start time.Time, end time.Time, id uint64) dtos.StreamsPerNodeResponse
}

type nodesServiceImpl struct {
	nodesRepository repositories.NodeRepository
	inaClient       clients.InaAPiClient
}

func NewNodesService(repositories *config.Repositories, inaClient clients.InaAPiClient) NodesService {
	return &nodesServiceImpl{repositories.NodeRepository, inaClient}
}

func (s nodesServiceImpl) GetNodes(timeStart time.Time, timeEnd time.Time, configId uint64) dtos.StreamsPerNodeResponse {
	nodes := s.nodesRepository.GetStreamsPerNodeById(configId)
	errorsPerNode := s.nodesRepository.GetErrorsOfNodes(configId, timeStart, timeEnd)
	for _, node := range nodes {
		if node.MainStreamId != nil {
			stream, err := s.inaClient.GetStream(*node.MainStreamId)
			if err != nil {
				log.Errorf("Error getting stream: %v for node summary", err)
				continue
			}
			node.LastUpdate = stream.DateRange.TimeEnd
			levels, err := s.inaClient.GetObservedData(*node.MainStreamId, timeStart, timeEnd)
			if err != nil {
				log.Errorf("Error getting levels: %v for node summary", err)
				continue
			}
			calculator := metrics_service.NewCalculatorOfWaterLevels(stream.Station.AlertLevel, stream.Station.EvacuationLevel, stream.Station.LowWaterLevel)
			totalValues := uint64(0)
			for _, level := range levels {
				if level.Value != nil {
					calculator.Compute(*level.Value)
					totalValues++
				}
			}
			node.AlertWaterLevels = calculator.GetAlertsCount()
			node.TotalWaterLevels = totalValues
		}
	}
	for _, errors := range errorsPerNode {
		for _, node := range nodes {
			if node.NodeId == errors.NodeId {
				node.ErrorCount = errors.ErrorCount
				break
			}
		}
	}
	return dtos.StreamsPerNodeResponse{Nodes: nodes}
}
