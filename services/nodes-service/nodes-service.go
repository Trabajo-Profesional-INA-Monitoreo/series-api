package nodes_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	metrics_service "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/metrics-service"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type NodesService interface {
	GetNodes(*dtos.QueryParameters) dtos.StreamsPerNodeResponse
}

type nodesServiceImpl struct {
	nodesRepository repositories.NodeRepository
	inaClient       clients.InaAPiClient
}

func NewNodesService(repositories *config.Repositories, inaClient clients.InaAPiClient) NodesService {
	return &nodesServiceImpl{repositories.NodeRepository, inaClient}
}

func (s nodesServiceImpl) GetNodes(parameters *dtos.QueryParameters) dtos.StreamsPerNodeResponse {
	configId := parameters.Get("configurationId").(uint64)
	page := *parameters.GetAsInt("page")
	pageSize := *parameters.GetAsInt("pageSize")
	timeStart := parameters.Get("timeStart").(time.Time)
	timeEnd := parameters.Get("timeEnd").(time.Time)
	nodes, pageable := s.nodesRepository.GetStreamsPerNodeById(configId, page, pageSize)
	var wg sync.WaitGroup
	for _, node := range nodes {
		if node.MainStreamId != nil {
			wg.Add(1)
			go s.getLastUpdateAndLevelOfNode(node, timeStart, timeEnd, &wg)
		}
	}

	var nodeIds []uint64
	for _, node := range nodes {
		id, _ := strconv.ParseUint(node.NodeId, 10, 64)
		nodeIds = append(nodeIds, id)
	}
	errorsPerNode := s.nodesRepository.GetErrorsOfNodes(configId, timeStart, timeEnd, nodeIds)
	for _, errors := range errorsPerNode {
		for _, node := range nodes {
			if node.NodeId == errors.NodeId {
				node.ErrorCount = errors.ErrorCount
				break
			}
		}
	}
	wg.Wait()
	return dtos.StreamsPerNodeResponse{Nodes: nodes, Pageable: pageable}
}

func (s nodesServiceImpl) getLastUpdateAndLevelOfNode(node *dtos.StreamsPerNode, timeStart time.Time, timeEnd time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	stream, err := s.inaClient.GetStream(*node.MainStreamId)
	if err != nil {
		log.Errorf("Error getting stream: %v for node summary", err)

	}
	node.LastUpdate = stream.DateRange.TimeEnd
	levels, err := s.inaClient.GetObservedData(*node.MainStreamId, timeStart, timeEnd)
	if err != nil {
		log.Errorf("Error getting levels: %v for node summary", err)

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
