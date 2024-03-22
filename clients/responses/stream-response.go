package responses

type InaStreamResponse struct {
	Type      string    `json:"tipo"`
	Id        uint      `json:"id"`
	DateRange DateRange `json:"date_range"`
	Unit      Unit      `json:"unidades"`
	Station   Station   `json:"estacion"`
	Variable  Variable  `json:"var"`
	Procedure Procedure `json:"procedimiento"`
}

type DateRange struct {
	TimeStart *string `json:"timestart"`
	TimeEnd   *string `json:"timeend"`
	Count     *string `json:"count"`
}

type Unit struct {
	Name      string `json:"nombre"`
	Abrev     string `json:"abrev"`
	UnitsId   uint   `json:"UnitsID"`
	UnitsType string `json:"UnitsType"`
	Id        int    `json:"id"`
}

type Procedure struct {
	Name        string `json:"nombre"`
	Abrev       string `json:"abrev"`
	Description string `json:"descripcion"`
	Id          int    `json:"id"`
}

type Variable struct {
	Var             string `json:"var"`
	Name            string `json:"nombre"`
	Abrev           string `json:"abrev"`
	Type            string `json:"type"`
	Datatype        string `json:"datatype"`
	Valuetype       string `json:"valuetype"`
	GeneralCategory string `json:"GeneralCategory"`
	VariableName    string `json:"VariableName"`
	SampleMedium    string `json:"SampleMedium"`
	DefUnitId       int    `json:"def_unit_id"`
	//TimeSupport     struct {} `json:"timeSupport"`
	//DefHoraCorte interface{} `json:"def_hora_corte"`
	Id int `json:"id"`
}

type Station struct {
	Name            string      `json:"nombre"`
	ExternalId      string      `json:"id_externo"`
	Table           string      `json:"tabla"`
	Province        string      `json:"provincia"`
	Country         string      `json:"pais"`
	River           string      `json:"rio"`
	HasObservations bool        `json:"has_obs"`
	Type            string      `json:"tipo"`
	Automatic       bool        `json:"automatica"`
	Enabled         bool        `json:"habilitar"`
	Owner           string      `json:"propietario"`
	Abreviatura     interface{} `json:"abreviatura"`
	Locality        interface{} `json:"localidad"`
	Real            bool        `json:"real"`
	AlertLevel      int         `json:"nivel_alerta"`
	EvacuationLevel float64     `json:"nivel_evacuacion"`
	LowWaterLevel   float64     `json:"nivel_aguas_bajas"`
	Altitude        interface{} `json:"altitud"`
	Public          bool        `json:"public"`
	CeroIgn         float64     `json:"cero_ign"`
	Id              int         `json:"id"`
	Geom            Geolocation `json:"geom"`
}

type Geolocation struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
