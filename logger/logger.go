package logger

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Entry

var once sync.Once

func InitLogger() {
	once.Do(
		func() {
			log := logrus.New()
			log.SetFormatter(&logrus.TextFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
				FullTimestamp:   true,
			})
			log.SetLevel(logrus.DebugLevel)
			Log = log.WithTime(time.Now())
		})
}

func LogErrorResponce(code int, message string) *logrus.Entry {
	return Log.WithFields(logrus.Fields{
		"code":    code,
		"message": message,
	})
}
