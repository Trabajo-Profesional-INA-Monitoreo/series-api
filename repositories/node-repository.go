package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"gorm.io/gorm"
)

type NodeRepository interface {
	Create(node *entities.Node) error
}

type nodeRepository struct {
	connection *gorm.DB
}

func (n nodeRepository) Create(node *entities.Node) error {
	result := n.connection.Create(&node)
	return result.Error
}

func NewNodeRepository(connection *gorm.DB) NodeRepository {
	return &nodeRepository{connection}
}
