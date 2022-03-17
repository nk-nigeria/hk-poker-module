package mock

import (
	log "github.com/ciaolink-game-platform/cgp-chinese-poker-module/pkg/log"
	"github.com/heroiclabs/nakama-common/runtime"
)

type MockLog struct {
}

func (m MockLog) Debug(format string, v ...interface{}) {
	log.GetLogger().Debug(format, v)
}

func (m MockLog) Info(format string, v ...interface{}) {
	log.GetLogger().Info(format, v)
}

func (m MockLog) Warn(format string, v ...interface{}) {
	log.GetLogger().Warn(format, v)
}

func (m MockLog) Error(format string, v ...interface{}) {
	log.GetLogger().Debug(format, v)
}

func (m MockLog) WithField(key string, v interface{}) runtime.Logger {
	log.GetLogger().WithField(key, v)

	return m
}

func (m MockLog) WithFields(fields map[string]interface{}) runtime.Logger {
	log.GetLogger().WithFields(fields)
	return m
}

func (m MockLog) Fields() map[string]interface{} {
	return nil
}
