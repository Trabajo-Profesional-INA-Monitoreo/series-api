package dtos

type StreamCardsParameters struct {
	params map[string]interface{}
	//streamId        *uint64
	//configurationId uint64
	//timeStart       time.Time
	//timeEnd         time.Time
	//varId           *uint64
	//procId          *uint64
	//stationId       *uint64
	//streamType      *entities.StreamType
	//page            int
	//pageSize        int
}

func NewStreamCardsParameters() *StreamCardsParameters {
	return &StreamCardsParameters{params: make(map[string]interface{})}
}

func (s *StreamCardsParameters) AddParam(key string, value interface{}) {
	s.params[key] = value
}

func (s *StreamCardsParameters) AddParamIfFound(key string, value string, found bool) {
	if found {
		s.params[key] = value
	}
}

func (s *StreamCardsParameters) AddParamOrDefault(key string, value string, found bool, defaultValue interface{}) {
	if found {
		s.params[key] = value
	} else {
		s.params[key] = defaultValue
	}
}

func (s *StreamCardsParameters) Get(key string) interface{} {
	return s.params[key]
}

func (s *StreamCardsParameters) GetAsInt(key string) *int {
	value := s.params[key]
	if value == nil {
		return nil
	}
	aux := int(value.(uint64))
	return &aux
}
