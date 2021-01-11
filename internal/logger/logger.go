package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// Creates logger once and returns it.
func Get() *zap.Logger {
	once.Do(func() {
		logger, _ = zap.NewProduction()
	})

	return logger
}
