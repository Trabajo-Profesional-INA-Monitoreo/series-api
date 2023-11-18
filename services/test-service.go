package services

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"

type TestService interface {
	GetTest() string
	PostTest(dto dtos.GreetingTestDto)
}

type testService struct {
	greeting string
}

func NewTestService() TestService {
	return &testService{greeting: "Hello"}
}

func (s *testService) GetTest() string {
	return s.greeting
}

func (s *testService) PostTest(dto dtos.GreetingTestDto) {
	s.greeting = dto.Greeting
}
