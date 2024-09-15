package server

import (
	"github.com/craftizmv/rewards/internal/app/handlers/http"
	"github.com/craftizmv/rewards/internal/app/usecase"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

const (
	ReadTimeout  = 15 * time.Second
	WriteTimeout = 15 * time.Second
)

type EchoServer struct {
	app     *echo.Echo
	conf    *EchoConfig
	log     logger.ILogger
	useCase usecase.RewardUseCase
}

type EchoConfig struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"`
	Host                string   `mapstructure:"host"`
}

func NewEchoServer(conf *EchoConfig, log logger.ILogger, useCase usecase.RewardUseCase) *EchoServer {
	e := echo.New()
	return &EchoServer{
		app:     e,
		conf:    conf,
		log:     log,
		useCase: useCase,
	}
}

func (s *EchoServer) Start() {
	// using middleware to recover and log
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())

	s.app.Server.ReadTimeout = ReadTimeout
	s.app.Server.WriteTimeout = WriteTimeout

	s.app.GET("v1/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// init http handlers
	s.initRewardHttpHandler(s.useCase)

	s.app.Logger.Fatal(s.app.Start(s.conf.Port))
}

func (s *EchoServer) initRewardHttpHandler(usecase usecase.RewardUseCase) {

	// create usecase

	rewardHandler := http.NewRewardHandler(usecase, s.log)

	// routers
	rewardRouter := s.app.Group("v1/rewards")
	rewardRouter.POST("eligibility", rewardHandler.CheckRewardEligibility)

}
