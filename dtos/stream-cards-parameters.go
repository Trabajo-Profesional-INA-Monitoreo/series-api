package dtos

type QueryParameters struct {
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

func NewQueryParameters() *QueryParameters {
	return &QueryParameters{params: make(map[string]interface{})}
}

func (s *QueryParameters) AddParam(key string, value interface{}) {
	s.params[key] = value
}

func (s *QueryParameters) AddParamIfFound(key string, value string, found bool) {
	if found {
		s.params[key] = value
	}
}

func (s *QueryParameters) AddParamOrDefault(key string, value string, found bool, defaultValue interface{}) {
	if found {
		s.params[key] = value
	} else {
		s.params[key] = defaultValue
	}
}

func (s *QueryParameters) Get(key string) interface{} {
	return s.params[key]
}

func (s *QueryParameters) GetAsInt(key string) *int {
	value := s.params[key]
	if value == nil {
		return nil
	}
	aux := int(value.(uint64))
	return &aux
}
