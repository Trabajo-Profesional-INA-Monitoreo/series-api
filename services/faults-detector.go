package services

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"time"
)

const DayInHours = 24
const MinForecastedDays = 4

type FaultDetector interface {
	DetectFaults()
}

type faultDetectorService struct {
	streamsRepository           repositories.StreamRepository
	configuredStreamsRepository repositories.ConfiguredStreamsRepository
	errorsRepository            repositories.ErrorsRepository
	inaApiClient                clients.InaAPiClient
}

func NewFaultDetectorService(streamsRepository repositories.StreamRepository,
	configuredStreamsRepository repositories.ConfiguredStreamsRepository,
	errorsRepository repositories.ErrorsRepository,
	inaApiClient clients.InaAPiClient,
) FaultDetector {
	return &faultDetectorService{
		streamsRepository:           streamsRepository,
		configuredStreamsRepository: configuredStreamsRepository,
		errorsRepository:            errorsRepository,
		inaApiClient:                inaApiClient,
	}
}

func (f faultDetectorService) handleMissingForecast(stream entities.Stream, configuredStream entities.ConfiguredStream, res *responses.LastForecast) {
	now := time.Now()
	diff := now.Sub(res.ForecastDate)
	reqErrorId := fmt.Sprintf("%v", res.RunId)
	detectedError := f.errorsRepository.GetDetectedErrorForStreamWithIdAndType(stream.StreamId, reqErrorId, entities.ForecastMissing)
	detected := detectedError.RequestId == reqErrorId
	if diff.Minutes() > configuredStream.UpdateFrequency && !detected {
		// There should be a new forecast already
		// We save the detected error
		log.Debugf("Detected missing forecast for: %v", stream.StreamId)
		detected := entities.DetectedError{
			StreamId:         stream.StreamId,
			Stream:           &stream,
			ConfiguredStream: []entities.ConfiguredStream{configuredStream},
			DetectedDate:     time.Now(),
			RequestId:        reqErrorId,
			ErrorType:        entities.ForecastMissing,
			ExtraInfo:        fmt.Sprintf("Fecha ultimo pronostico %v - Tiempo transcurrido %v", res.ForecastDate, diff.String()),
		}
		f.errorsRepository.Create(detected)
	} else if detected && !contains(detectedError.ConfiguredStream, configuredStream) {
		// We already detected the error, we need to save the relationship to the current ConfiguredStream
		detectedError := f.errorsRepository.GetDetectedErrorForStreamWithIdAndType(stream.StreamId, reqErrorId, entities.ForecastMissing)
		detectedError.ConfiguredStream = append(detectedError.ConfiguredStream, configuredStream)
		f.errorsRepository.Update(detectedError)
	}
}

func (f faultDetectorService) handleObservedValues(data []responses.ObservedDataResponse, stream entities.Stream, configuredStreams []entities.ConfiguredStream) {
	for _, observed := range data {
		isNull := observed.Value == nil || observed.DataId == ""
		reqId := fmt.Sprintf("%v", observed.TimeStart)
		if isNull && !f.errorsRepository.AlreadyDetectedErrorForStreamWithIdAndType(stream.StreamId, reqId, entities.NullValue) {
			// We detected a new null value
			log.Debugf("Detected null value for: %v", stream.StreamId)
			detected := entities.DetectedError{
				StreamId:         stream.StreamId,
				Stream:           &stream,
				ConfiguredStream: configuredStreams,
				DetectedDate:     time.Now(),
				RequestId:        reqId,
				ErrorType:        entities.NullValue,
				ExtraInfo:        fmt.Sprintf("Timestamp del dato nulo %v", observed.TimeStart),
			}
			f.errorsRepository.Create(detected)
		}
	}

	for _, configuredStream := range configuredStreams {
		consecutiveOutliers := 0
		for _, observed := range data {
			isOutlier := valueIsAnOutlier(configuredStream, observed)
			if isOutlier && consecutiveOutliers == 0 {
				reqId := fmt.Sprintf("%v", observed.TimeStart)
				detectedError := f.errorsRepository.GetDetectedErrorForStreamWithIdAndType(stream.StreamId, reqId, entities.ObservedOutlier)
				detected := detectedError.RequestId == reqId
				if !detected {
					log.Debugf("Detected outlier value in observed stream for: %v", stream.StreamId)
					detected := entities.DetectedError{
						StreamId:         stream.StreamId,
						Stream:           &stream,
						ConfiguredStream: []entities.ConfiguredStream{configuredStream},
						DetectedDate:     time.Now(),
						RequestId:        reqId,
						ErrorType:        entities.ObservedOutlier,
						ExtraInfo:        fmt.Sprintf("Valor %v - Timestamp %v", observed.Value, observed.TimeStart),
					}
					f.errorsRepository.Create(detected)
				} else if !contains(detectedError.ConfiguredStream, configuredStream) {
					detectedError.ConfiguredStream = append(detectedError.ConfiguredStream, configuredStream)
					f.errorsRepository.Update(detectedError)
				}
				consecutiveOutliers++
			}
			if !isOutlier {
				consecutiveOutliers = 0
			}
		}
	}

}

func (f faultDetectorService) getObservedDataFromStream(streamId uint64, timeStart time.Time, timeEnd time.Time) ([]responses.ObservedDataResponse, error) {
	res, err := f.inaApiClient.GetObservedData(streamId, timeStart, timeEnd)
	if err != nil {
		return nil, fmt.Errorf("error performing check for stream %v: %v", streamId, err)
	}
	return res, nil
}

func (f faultDetectorService) handleObservedStreams(stream entities.Stream) {
	configuredStreams := f.configuredStreamsRepository.FindConfiguredStreamsWithCheckErrorsForStream(stream)
	data, err := f.getObservedDataFromStream(stream.StreamId, time.Now().Add(time.Duration(-24)*time.Hour), time.Now())
	if err != nil {
		log.Errorf("Error detecting observed errors: %v", err)
		return
	}
	f.handleObservedValues(data, stream, configuredStreams)
}

func (f faultDetectorService) handleForecastedStream(stream entities.Stream) {
	configuredStreams := f.configuredStreamsRepository.FindConfiguredStreamsWithCheckErrorsForStream(stream)
	// Una misma serie puede tener multiples calibrados, estos calibrados pertenecen a la configuracion
	for _, configuredStream := range configuredStreams {
		res, err := f.inaApiClient.GetLastForecast(configuredStream.CalibrationId)
		if err != nil {
			log.Errorf("Error performing check for stream %v with configuration %v", stream.StreamId, configuredStream.ConfiguredStreamId)
			return
		}
		f.handleMissingForecast(stream, configuredStream, res)

		forecast := res.GetForecastOfStream(stream.StreamId)

		if forecast.MainForecast.Forecasts != nil {
			consecutiveOutliers := 0
			forecastedDays := 0
			for _, hourlyForecast := range forecast.MainForecast.Forecasts {
				timestamp, err := time.Parse("2006-01-02T15:04:05.999Z", hourlyForecast[0])
				if err != nil {
					log.Errorf("Error parsing timestamp %v for stream %v with cal id %v - err: %v", hourlyForecast[0], stream.StreamId, configuredStream.CalibrationId, err)
					continue
				}
				value, err := strconv.ParseFloat(hourlyForecast[2], 64)
				if err != nil {
					log.Errorf("Error parsing forecast value %v for stream %v with cal id %v - err: %v", hourlyForecast[2], stream.StreamId, configuredStream.CalibrationId, err)
					continue
				}
				isOutsideBoundaries := configuredStream.NormalLowerThreshold > value || configuredStream.NormalUpperThreshold < value
				if isOutsideBoundaries && consecutiveOutliers == 0 {
					reqErrorId := fmt.Sprintf("%v-%v-%v-%v", res.RunId, res.CalibrationId, stream.StreamId, timestamp)
					detectedError := f.errorsRepository.GetDetectedErrorForStreamWithIdAndType(stream.StreamId, reqErrorId, entities.ForecastOutOfBounds)
					detected := detectedError.RequestId == reqErrorId
					if !detected {
						// We detected an outlier in the forecast
						log.Debugf("Detected outlier in forecast for: %v", stream.StreamId)
						detected := entities.DetectedError{
							StreamId:         stream.StreamId,
							Stream:           &stream,
							ConfiguredStream: []entities.ConfiguredStream{configuredStream},
							DetectedDate:     time.Now(),
							RequestId:        reqErrorId,
							ErrorType:        entities.ForecastOutOfBounds,
							ExtraInfo:        fmt.Sprintf("Valor %v - Timestamp %v - Corrida %v - Fecha pronostico %v", value, timestamp, res.RunId, res.ForecastDate),
						}
						f.errorsRepository.Create(detected)
					} else if !contains(detectedError.ConfiguredStream, configuredStream) {
						detectedError.ConfiguredStream = append(detectedError.ConfiguredStream, configuredStream)
						f.errorsRepository.Update(detectedError)
					}
					consecutiveOutliers++
				}
				if !isOutsideBoundaries {
					consecutiveOutliers = 0
				}

				if timestamp.After(res.ForecastDate) {
					diff := timestamp.Sub(res.ForecastDate)
					if math.Floor(diff.Hours()/DayInHours) > float64(forecastedDays) {
						forecastedDays++
					}
				}

			}
			if forecastedDays < MinForecastedDays {
				reqErrorId := fmt.Sprintf("%v-%v-%v", res.RunId, res.CalibrationId, stream.StreamId)
				detectedError := f.errorsRepository.GetDetectedErrorForStreamWithIdAndType(stream.StreamId, reqErrorId, entities.Missing4DaysHorizon)
				detected := detectedError.RequestId == reqErrorId
				lastForecastedDate, _ := time.Parse("2006-01-02T15:04:05.999Z", forecast.MainForecast.Forecasts[len(forecast.MainForecast.Forecasts)-1][0])
				if !detected {
					log.Debugf("Detected missing 4 days horizon forecast for: %v", stream.StreamId)
					detected := entities.DetectedError{
						StreamId:         stream.StreamId,
						Stream:           &stream,
						ConfiguredStream: []entities.ConfiguredStream{configuredStream},
						DetectedDate:     time.Now(),
						RequestId:        reqErrorId,
						ErrorType:        entities.Missing4DaysHorizon,
						ExtraInfo:        fmt.Sprintf("Corrida %v - Fecha pronostico %v - Ultima fecha pronosticada %v - Diferencia %v", res.RunId, res.ForecastDate, lastForecastedDate, lastForecastedDate.Sub(res.ForecastDate)),
					}
					f.errorsRepository.Create(detected)
				} else if !contains(detectedError.ConfiguredStream, configuredStream) {
					detectedError.ConfiguredStream = append(detectedError.ConfiguredStream, configuredStream)
					f.errorsRepository.Update(detectedError)
				}
			}
		}

		if shouldFetchObservedStream(forecast, configuredStream) {
			log.Debugf("Performing check of values out of bands for forecasted (%v - cal %v) - observed (%v) streams", configuredStream.StreamId, configuredStream.CalibrationId, *configuredStream.ObservedRelatedStreamId)
			timeStart, timeEnd := getDateRangeOfForecast(forecast.MainForecast)
			observedData, err := f.getObservedDataFromStream(*configuredStream.ObservedRelatedStreamId, timeStart, timeEnd)
			if err != nil {
				log.Errorf("Error performing out of bands check on configured stream: %v - Err: %v", configuredStream.ConfiguredStreamId, err)
				continue
			}
			outsideBands := 0
			forecastedIndex := 0
			mainForecast := forecast.MainForecast.Forecasts
			for observedIndex := 0; observedIndex < len(observedData) && forecastedIndex < len(forecast.MainForecast.Forecasts); {
				hourlyTime := parseForecastedDate(mainForecast[forecastedIndex][0])
				observedData := observedData[observedIndex]
				if observedData.TimeStart.Before(hourlyTime) || observedData.TimeStart.Equal(hourlyTime) {
					if observedIsOutsideErrorBands(forecast, forecastedIndex, observedData.Value) {
						outsideBands++
					}
					observedIndex++
				}
				forecastedIndex++
			}
			if tooManyValuesOutsideBands(outsideBands, len(observedData)) {
				reqErrorId := fmt.Sprintf("%v-%v-%v", res.RunId, res.CalibrationId, stream.StreamId)
				detectedError := f.errorsRepository.GetDetectedErrorForStreamWithIdAndType(stream.StreamId, reqErrorId, entities.OutsideOfErrorBands)
				detected := detectedError.RequestId == reqErrorId
				if !detected {
					log.Debugf("Detected too many values outside error bands for stream: %v - calibration: %v", stream.StreamId, configuredStream.CalibrationId)
					detected := entities.DetectedError{
						StreamId:         stream.StreamId,
						Stream:           &stream,
						ConfiguredStream: []entities.ConfiguredStream{configuredStream},
						DetectedDate:     time.Now(),
						RequestId:        reqErrorId,
						ErrorType:        entities.OutsideOfErrorBands,
						ExtraInfo:        fmt.Sprintf("Corrida %v - Fecha pronostico %v - Cantidad valores fuera de bandas %v", res.RunId, res.ForecastDate, outsideBands),
					}
					f.errorsRepository.Create(detected)
				} else if !contains(detectedError.ConfiguredStream, configuredStream) {
					detectedError.ConfiguredStream = append(detectedError.ConfiguredStream, configuredStream)
					f.errorsRepository.Update(detectedError)
				}
			}
		}
	}
}

func (f faultDetectorService) DetectFaults() {

	log.Debugf("FaultsDetector | Performing fault checks...")
	streams := f.streamsRepository.GetStreams()
	for _, stream := range streams {
		if stream.IsObserved() {
			f.handleObservedStreams(stream)
		} else if stream.IsForecasted() {
			f.handleForecastedStream(stream)
		}
	}
	log.Debugf("FaultsDetector | Fault check done")
}
