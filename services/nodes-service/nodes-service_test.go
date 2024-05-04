package nodes_service

import (
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

func TestShouldReturnTheNodesSummary(t *testing.T) {
	mockRepository := nodesRepositoryMock{}
	mockRepository.streamsPerNode = []*dtos.StreamsPerNode{{"Test", "1", 1, 0}}
	mockRepository.errorsPerNode = []dtos.ErrorsOfNodes{{"1", 1}}
	nodesService := &nodesServiceImpl{mockRepository}
	res := nodesService.GetNodes(time.Now(), time.Now(), 1)
	assert.Equal(t, 1, len(res.Nodes))
	assert.Equal(t, 1, res.Nodes[0].StreamsCount)
	assert.Equal(t, 1, res.Nodes[0].ErrorCount)
	assert.Equal(t, "Test", res.Nodes[0].NodeName)
	assert.Equal(t, "1", res.Nodes[0].NodeId)
}
