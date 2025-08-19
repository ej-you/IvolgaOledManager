// App binary starts full application.
package main

import (
	"github.com/sirupsen/logrus"

	"IvolgaOledManager/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		logrus.Fatal(err)
	}
	if err := application.Run(); err != nil {
		logrus.Fatal(err)
	}
}
