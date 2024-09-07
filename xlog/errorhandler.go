package xlog

import (
	"context"
	"go.uber.org/zap"
)

type logErrorHandler struct {
}

func NewZapLogErrorHandler() *logErrorHandler {
	return &logErrorHandler{}
}

func (h *logErrorHandler) Handle(ctx context.Context, err error) {
	Error("Error handler log", zap.Error(err))
}
