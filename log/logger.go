package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"syscall"
	"time"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&nested.Formatter{
		NoColors: true,
	})
	logFile, err := os.OpenFile(time.Now().Format("./data/log/2006-Jan-150405.log"),
		syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0777)
	if err != nil {
		Logger.Panicln(err)
	}
	Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}
