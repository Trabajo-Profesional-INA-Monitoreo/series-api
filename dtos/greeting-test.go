package dtos

type GreetingTestDto struct {
	Greeting string `json:"greeting" binding:"required,min=2,max=100"`
}
