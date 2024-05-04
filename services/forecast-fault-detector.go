package services

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

type forecastFaultDetectorService struct {
	configuredStreamsRepository repositories.DetectionConfiguredStreamsRepository
	errorsRepository            repositories.ErrorDetectionRepository
	inaApiClient                clients.InaAPiClient
}

func newForecastFaultDetectorService(configuredStreamsRepository repositories.DetectionConfiguredStreamsRepository, errorsRepository repositories.ErrorDetectionRepository, inaApiClient clients.InaAPiClient) *forecastFaultDetectorService {
	return &forecastFaultDetectorService{configuredStreamsRepository: configuredStreamsRepository, errorsRepository: errorsRepository, inaApiClient: inaApiClient}
}

func (f forecastFaultDetectorService) handleMissingForecast(stream entities.Stream, configuredStream entities.ConfiguredStream, res *responses.LastForecast) {
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

func (f forecastFaultDetectorService) handleStream(stream entities.Stream) {
	configuredStreams := f.configuredStreamsRepository.FindConfiguredStreamsWithCheckErrorsForStream(stream)
	// Una misma serie puede tener multiples calibrados, estos calibrados pertenecen a la configuracion
	for _, configuredStream := range configuredStreams {
		res, err := f.inaApiClient.GetLastForecast(configuredStream.CalibrationId)
		if err != nil {
			log.Errorf("Error performing check for stream %v with configuration %v: %v", stream.StreamId, configuredStream.ConfiguredStreamId, err)
			continue
		}
		f.handleMissingForecast(stream, configuredStream, res)

		forecast := res.GetForecastOfStream(stream.StreamId)
		if forecast.MainForecast.Forecasts != nil {
			f.checkMainForecastErrors(stream, forecast, configuredStream, res)
		}

		if shouldFetchObservedStream(forecast, configuredStream) {
			f.checkObservedValuesTogetherWithForecast(stream, configuredStream, forecast, res)
		}
	}
}

func (f forecastFaultDetectorService) checkMainForecastErrors(stream entities.Stream, forecast *responses.Forecast, configuredStream entities.ConfiguredStream, res *responses.LastForecast) {
	consecutiveOutliers := 0
	forecastedDays := 0
	for _, hourlyForecast := range forecast.MainForecast.Forecasts {
		timestamp, value, err := convertForecastStringData(hourlyForecast[0], hourlyForecast[2])
		if err != nil {
			log.Errorf("Error parsing forecast for stream %v with cal id %v - err: %v", stream.StreamId, configuredStream.CalibrationId, err)
			continue
		}
		isOutsideBoundaries := configuredStream.NormalLowerThreshold > value || configuredStream.NormalUpperThreshold < value
		if isOutsideBoundaries && consecutiveOutliers == 0 {
			f.handleValueOutsideBoundaries(stream, res, timestamp, configuredStream, value)
			consecutiveOutliers++
		}
		if !isOutsideBoundaries {
			consecutiveOutliers = 0
		}

		forecastedDays = addTimeToForecastedDays(*timestamp, res, forecastedDays)
	}
	if forecastedDays < MinForecastedDays {
		f.saveForecastWasNotCompleteError(stream, res, forecast, configuredStream)
	}
}

func (f forecastFaultDetectorService) checkObservedValuesTogetherWithForecast(stream entities.Stream, configuredStream entities.ConfiguredStream, forecast *responses.Forecast, res *responses.LastForecast) {
	log.Debugf("Performing check of values out of bands for forecasted (%v - cal %v) - observed (%v) streams", configuredStream.StreamId, configuredStream.CalibrationId, *configuredStream.ObservedRelatedStreamId)
	timeStart, timeEnd := getDateRangeOfForecast(forecast.MainForecast)
	observedData, err := getObservedDataFromStream(*configuredStream.ObservedRelatedStreamId, timeStart, timeEnd, f.inaApiClient)
	if err != nil {
		log.Errorf("Error performing out of bands check on configured stream: %v - Err: %v", configuredStream.ConfiguredStreamId, err)
		return
	}
	outsideBands := countObservedValuesOutsideErrorBands(forecast, observedData)
	if tooManyValuesOutsideBands(outsideBands, len(observedData)) {
		f.saveObservedValuesAreOutsideErrorBands(stream, res, configuredStream, outsideBands)
	}
}

func (f forecastFaultDetectorService) handleValueOutsideBoundaries(stream entities.Stream, res *responses.LastForecast, timestamp *time.Time, configuredStream entities.ConfiguredStream, value float64) {
	reqErrorId := fmt.Sprintf("%v-%v-%v-%v", res.RunId, res.CalibrationId, stream.StreamId, *timestamp)
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
			ExtraInfo:        fmt.Sprintf("Valor %v - Timestamp %v - Corrida %v - Fecha pronostico %v", value, *timestamp, res.RunId, res.ForecastDate),
		}
		f.errorsRepository.Create(detected)
	} else if !contains(detectedError.ConfiguredStream, configuredStream) {
		detectedError.ConfiguredStream = append(detectedError.ConfiguredStream, configuredStream)
		f.errorsRepository.Update(detectedError)
	}
}

func addTimeToForecastedDays(timestamp time.Time, res *responses.LastForecast, forecastedDays int) int {
	if timestamp.After(res.ForecastDate) {
		diff := timestamp.Sub(res.ForecastDate)
		if math.Floor(diff.Hours()/DayInHours) > float64(forecastedDays) {
			forecastedDays++
		}
	}
	return forecastedDays
}

func countObservedValuesOutsideErrorBands(forecast *responses.Forecast, observedData []responses.ObservedDataResponse) int {
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
	return outsideBands
}

func (f forecastFaultDetectorService) saveObservedValuesAreOutsideErrorBands(stream entities.Stream, res *responses.LastForecast, configuredStream entities.ConfiguredStream, outsideBands int) {
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

func (f forecastFaultDetectorService) saveForecastWasNotCompleteError(stream entities.Stream, res *responses.LastForecast, forecast *responses.Forecast, configuredStream entities.ConfiguredStream) {
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
