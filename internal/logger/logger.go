package logger

import "go.uber.org/zap"

func New(level string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	if err := cfg.Level.UnmarshalText([]byte(level)); err != nil {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	l, _ := cfg.Build()
	return l
}

func ReplaceGlobals(l *zap.Logger) { zap.ReplaceGlobals(l) }
