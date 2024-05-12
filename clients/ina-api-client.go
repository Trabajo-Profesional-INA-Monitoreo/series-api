package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const directMeasurementProcedure = 1
const waterLevelMeasurment = 2

type InaAPiClient interface {
	GetLastForecast(calibrationId uint64) (*responses.LastForecast, error)
	GetObservedData(streamId uint64, timeStart time.Time, timeEnd time.Time) ([]responses.ObservedDataResponse, error)
	GetStream(streamId uint64) (*responses.InaStreamResponse, error)
	GetMainStreamFromStation(stationId uint64) (*responses.MainStreamFromStation, error)
}

type inaApiClientImpl struct {
	baseUrl    string
	authHeader string
}

func closeReaderAndPrintError(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		log.Errorf("Error closing reader: %v", err)
	}
}

func NewInaApiClientImpl(apiConfig *config.ServiceConfigurationData) InaAPiClient {
	return &inaApiClientImpl{baseUrl: apiConfig.InaBaseUrl, authHeader: fmt.Sprintf("Bearer %v", apiConfig.InaToken)}
}

func (i inaApiClientImpl) GetMainStreamFromStation(stationId uint64) (*responses.MainStreamFromStation, error) {
	url := fmt.Sprintf("%v/obs/puntual/series?estacion_id=%v", i.baseUrl, stationId)

	log.Debugf("Performing stream request: %v", url)
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
	defer closeReaderAndPrintError(res.Body)
	if res.StatusCode != 200 {
		return nil, errors.Join(exceptions.MapCodeToError(res.StatusCode), fmt.Errorf("stations response error: got %v", res.StatusCode))
	}
	var decodedBody []responses.InaStreamResponse
	err = json.NewDecoder(res.Body).Decode(&decodedBody)
	if err != nil {
		log.Errorf("Error decoding response: %v", err)
		return nil, err
	}

	for _, stream := range decodedBody {
		if stream.Procedure.Id == directMeasurementProcedure && stream.Variable.Id == waterLevelMeasurment {
			return responses.NewMainStreamFromStation(stream), nil
		}
	}

	return nil, nil
}

func (i inaApiClientImpl) GetLastForecast(calibrationId uint64) (*responses.LastForecast, error) {
	url := fmt.Sprintf("%v/sim/calibrados/%v/corridas/last", i.baseUrl, calibrationId)
	log.Debugf("Performing forecast request: %v", url)
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
	defer closeReaderAndPrintError(res.Body)
	if res.StatusCode != 200 {
		return nil, errors.Join(exceptions.MapCodeToError(res.StatusCode), fmt.Errorf("forecast response error: got %v", res.StatusCode))
	}
	decodedBody := &responses.LastForecast{}
	err = json.NewDecoder(res.Body).Decode(decodedBody)
	if err != nil {
		log.Errorf("Error decoding response: %v", err)
		return nil, err
	}
	return decodedBody, nil
}

func (i inaApiClientImpl) GetObservedData(streamId uint64, timeStart time.Time, timeEnd time.Time) ([]responses.ObservedDataResponse, error) {
	url := fmt.Sprintf("%v/obs/puntual/series/%v/observaciones?timestart=%v&timeend=%v", i.baseUrl, streamId, timeStart.Format(time.RFC3339Nano), timeEnd.Format(time.RFC3339Nano))
	log.Debugf("Performing observed request: %v", url)
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
	defer closeReaderAndPrintError(res.Body)
	if res.StatusCode != 200 {
		return nil, errors.Join(exceptions.MapCodeToError(res.StatusCode), fmt.Errorf("observed response error: got %v", res.StatusCode))
	}
	var decodedBody []responses.ObservedDataResponse
	err = json.NewDecoder(res.Body).Decode(&decodedBody)
	if err != nil {
		log.Errorf("Error decoding response: %v", err)
		return nil, err
	}
	return decodedBody, nil
}

func (i inaApiClientImpl) GetStream(streamId uint64) (*responses.InaStreamResponse, error) {
	url := fmt.Sprintf("%v/obs/puntual/series/%v", i.baseUrl, streamId)

	log.Debugf("Performing stream request: %v", url)
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
	defer closeReaderAndPrintError(res.Body)
	if res.StatusCode != 200 {
		return nil, errors.Join(exceptions.MapCodeToError(res.StatusCode), fmt.Errorf("Stream response error: got %v", res.StatusCode))
	}
	var decodedBody responses.InaStreamResponse
	err = json.NewDecoder(res.Body).Decode(&decodedBody)
	if err != nil {
		log.Errorf("Error decoding response: %v", err)
		return nil, err
	}
	return &decodedBody, nil
}
