package dtos

type StreamsPerNetworkResponse struct {
	Networks []StreamsPerNetwork
}

//func NewStreamsPerNetworkResponse(networks []*StreamsPerNetwork) *StreamsPerNetworkResponse {
//	return &StreamsPerNetworkResponse{networks: networks}
//}

type StreamsPerNetwork struct {
	NetworkName  string `gorm:"column:networkname"`
	NetworkId    string `gorm:"column:networkid"`
	StreamsCount int    `gorm:"column:streamscount"`
}
