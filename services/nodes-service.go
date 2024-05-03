package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"time"
)

type NodesService interface {
	GetNodes(start time.Time, end time.Time, id uint64) dtos.StreamsPerNodeResponse
}

type nodesServiceImpl struct {
	repository      repositories.StreamRepository
	nodesRepository repositories.NodeRepository
}

func NewNodesService(repositories *config.Repositories) NodesService {
	return &nodesServiceImpl{repositories.StreamsRepository, repositories.NodeRepository}
}

func (s nodesServiceImpl) GetNodes(timeStart time.Time, timeEnd time.Time, configId uint64) dtos.StreamsPerNodeResponse {
	nodes := s.nodesRepository.GetStreamsPerNodeById(configId)
	errorsPerNode := s.repository.GetErrorsOfNodes(configId, timeStart, timeEnd)

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
