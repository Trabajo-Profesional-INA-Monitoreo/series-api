package clients

import (
	"encoding/json"
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type InaAPiClient interface {
	GetLastForecast(calibrationId uint64) (*responses.LastForecast, error)
}

type inaApiClientImpl struct {
	baseUrl    string
	authHeader string
}

func NewInaApiClientImpl(apiConfig *config.ApiConfig) InaAPiClient {
	return &inaApiClientImpl{baseUrl: apiConfig.InaBaseUrl, authHeader: fmt.Sprintf("Bearer %v", apiConfig.InaToken)}
}

func (i inaApiClientImpl) GetLastForecast(calibrationId uint64) (*responses.LastForecast, error) {
	url := fmt.Sprintf("%v/sim/calibrados/%v/corridas/last", i.baseUrl, calibrationId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return nil, err
	}
	req.Header.Add("Authorization", i.authHeader)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error making request: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("forecast response error: got %v", res.StatusCode)
	}
	decodedBody := &responses.LastForecast{}
	err = json.NewDecoder(res.Body).Decode(decodedBody)
	if err != nil {
		log.Errorf("Error decoding response: %v", err)
		return nil, err
	}
	return decodedBody, nil
}
