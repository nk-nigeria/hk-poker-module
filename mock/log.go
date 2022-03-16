package mock

import (
	"github.com/heroiclabs/nakama-common/runtime"
	log "github.com/sirupsen/logrus"
)

type MockLog struct {
}

func (m MockLog) Debug(format string, v ...interface{}) {
	log.Debug(format, v)
}

func (m MockLog) Info(format string, v ...interface{}) {
	log.Infof(format, v)
}

func (m MockLog) Warn(format string, v ...interface{}) {
	log.Warn(format, v)
}

func (m MockLog) Error(format string, v ...interface{}) {
	log.Debugf(format, v)
}

func (m MockLog) WithField(key string, v interface{}) runtime.Logger {
	log.WithField(key, v)

	return m
}

func (m MockLog) WithFields(fields map[string]interface{}) runtime.Logger {
	log.WithFields(fields)
	return m
}

func (m MockLog) Fields() map[string]interface{} {
	return log.Fields{}
}
