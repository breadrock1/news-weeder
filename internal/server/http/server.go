package http

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"news-weeder/internal/server"
	"news-weeder/internal/weeder"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "news-weeder/docs"
)

type HttpServer struct {
	config *server.Config
	server *echo.Echo
	weeder *weeder.Weeder
}

func Init(conf *server.Config, weeder *weeder.Weeder) *server.Server {
	httpServer := &HttpServer{
		config: conf,
		server: echo.New(),
		weeder: weeder,
	}

	return &server.Server{Server: httpServer}
}

func (s *HttpServer) setupServer() {
	s.server = echo.New()

	s.server.Use(middleware.CORS())
	s.server.Use(middleware.Recover())
	s.server.Use(InitLogger(s.config))

	_ = s.CreateSearchGroup()

	s.server.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (s *HttpServer) Start(_ context.Context) error {
	s.setupServer()
	return s.server.Start(s.config.Address)
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
