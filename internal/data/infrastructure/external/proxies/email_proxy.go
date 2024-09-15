package proxies

import (
	"github.com/craftizmv/rewards/internal/app/contracts"
	"github.com/craftizmv/rewards/pkg/logger"
	"go.uber.org/zap"
)

type EmailProxy struct {
	mailer contracts.Mailer
	logger logger.ILogger
}

func NewEmailProxy(mailer contracts.Mailer, logger logger.ILogger) *EmailProxy {
	return &EmailProxy{
		mailer: mailer,
		logger: logger,
	}
}

func (e *EmailProxy) SendEmail(name string, addr string, data string) error {
	err := e.mailer.SendEmail(name, addr, data)
	if err != nil {
		e.logger.Error("error sending email", zap.Error(err))
		return err
	}
	return nil
}
