package server

import (
	"context"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/labstack/echo/v4"
)

type Server interface {
	RunHttpServer(ctx context.Context, echo *echo.Echo, log logger.ILogger, cfg *EchoConfig) error
}
