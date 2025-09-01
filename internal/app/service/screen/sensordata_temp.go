package screen

// import (
// 	"IvolgaOledManager/internal/app/repo"
// 	"IvolgaOledManager/internal/app/service/button"
// 	"IvolgaOledManager/internal/pkg/pubsub"
// )

// // NewSensTemp returns a new instance of SensorData for temperature sensor.
// func NewSensTemp(active <-chan bool, btns button.Buttons, storage pubsub.Storage) *SensorData {
// 	handlers := button.Handlers{
// 		button.ButtonEsc: sensTempEsc,
// 	}
// 	btnHandlersReg := func() {
// 		btns.SetRisingHandlers(handlers)
// 	}
// 	return NewSensorData(active, btnHandlersReg, storage, repo.SensTempKey)
// }

// func sensTempEsc() {
// 	sensTempCh <- false
// 	greetCh <- true
// }
