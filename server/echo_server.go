package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

const (
	ReadTimeout  = 15 * time.Second
	WriteTimeout = 15 * time.Second
)

type EchoServer struct {
	app  *echo.Echo
	conf *EchoConfig
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

func NewEchoServer(conf *EchoConfig) *EchoServer {
	e := echo.New()
	return &EchoServer{
		app:  e,
		conf: conf,
	}
}

func (s *EchoServer) Start() {
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())
	s.app.Server.ReadTimeout = ReadTimeout
	s.app.Server.WriteTimeout = WriteTimeout

	s.app.GET("v1/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// TODO: initialize any routers

	s.app.Logger.Fatal(s.app.Start(s.conf.Port))
}
