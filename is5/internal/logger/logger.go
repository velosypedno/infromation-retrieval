package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.SetOutput(os.Stdout)
	} else {
		multiWriter := io.MultiWriter(os.Stdout, file)
		Log.SetOutput(multiWriter)
	}

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	Log.SetLevel(logrus.DebugLevel)
}
