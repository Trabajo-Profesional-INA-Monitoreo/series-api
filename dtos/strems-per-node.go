package dtos

type StreamsPerNodeResponse struct {
	Nodes []*StreamsPerNode
}

type StreamsPerNode struct {
	NodeName     string `gorm:"column:name"`
	NodeId       string `gorm:"column:node_id"`
	StreamsCount int    `gorm:"column:streams_count"`
	ErrorCount   int
}

type ErrorsOfNodes struct {
	NodeId     string `gorm:"column:node_id"`
	ErrorCount int    `gorm:"column:error_count"`
}
