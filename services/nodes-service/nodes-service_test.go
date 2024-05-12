package nodes_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type nodesRepositoryMock struct {
	streamsPerNode []*dtos.StreamsPerNode
	errorsPerNode  []dtos.ErrorsOfNodes
}

func (n nodesRepositoryMock) Create(node *entities.Node) (uint64, error) {
	return 0, nil
}

func (n nodesRepositoryMock) Update(node *entities.Node) error {
	return nil
}

func (n nodesRepositoryMock) GetNodesById(id uint64) []*dtos.Node {
	return nil
}

func (n nodesRepositoryMock) GetStreamsPerNodeById(formatUint uint64) []*dtos.StreamsPerNode {
	return n.streamsPerNode
}

func (n nodesRepositoryMock) MarkAsDeletedOldNodes(id uint64, ids []uint64) {
}

func (n nodesRepositoryMock) DeleteByConfig(configId uint64) {
}

func (n nodesRepositoryMock) GetErrorsOfNodes(configId uint64, timeStart time.Time, timeEnd time.Time) []dtos.ErrorsOfNodes {
	return n.errorsPerNode
}

type MockInaApiClient struct {
	ObservedData []responses.ObservedDataResponse
	Error        error
	Stream       responses.InaStreamResponse
}

func (m MockInaApiClient) GetLastForecast(calibrationId uint64) (*responses.LastForecast, error) {
	return nil, nil
}

func (m MockInaApiClient) GetObservedData(streamId uint64, timeStart time.Time, timeEnd time.Time) ([]responses.ObservedDataResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.ObservedData, nil
}

func (m MockInaApiClient) GetStream(streamId uint64) (*responses.InaStreamResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return &m.Stream, nil
}

func (m MockInaApiClient) GetMainStreamFromStation(stationId uint64) (*responses.MainStreamFromStation, error) {
	return nil, nil
}

func TestShouldReturnTheNodesSummary(t *testing.T) {
	mockRepository := nodesRepositoryMock{}
	mainStream := uint64(26)
	mockRepository.streamsPerNode = []*dtos.StreamsPerNode{{"Test", "1", 1, 0, &mainStream, 0, 0, nil}}
	mockRepository.errorsPerNode = []dtos.ErrorsOfNodes{{"1", 1}}
	mockClient := MockInaApiClient{}
	alertLevel := 4.0
	evacLevel := 5.0
	lowLevel := 1.0
	mockClient.ObservedData = []responses.ObservedDataResponse{{Value: &alertLevel}}
	mockClient.Stream = responses.InaStreamResponse{DateRange: responses.DateRange{TimeEnd: &time.Time{}}, Station: responses.Station{AlertLevel: &alertLevel, EvacuationLevel: &evacLevel, LowWaterLevel: &lowLevel}}
	nodesService := &nodesServiceImpl{mockRepository, mockClient}
	res := nodesService.GetNodes(time.Now(), time.Now(), 1)
	assert.Equal(t, 1, len(res.Nodes))
	assert.Equal(t, 1, res.Nodes[0].StreamsCount)
	assert.Equal(t, 1, res.Nodes[0].ErrorCount)
	assert.Equal(t, uint64(1), res.Nodes[0].AlertWaterLevels)
	assert.Equal(t, uint64(1), res.Nodes[0].TotalWaterLevels)
	assert.Equal(t, "Test", res.Nodes[0].NodeName)
	assert.Equal(t, "1", res.Nodes[0].NodeId)
}
