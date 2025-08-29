package screen

import (
	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// NewGreetings returns a new instance of template.Image for greetings screen.
func NewGreetings(active chan bool, btnHandlersReg func(),
	storage pubsub.Storage, imagePath string) *template.Image {

	image := &entity.Image{ImagePath: imagePath}
	return template.NewImage(active, btnHandlersReg, storage, image, "greetings")
}
