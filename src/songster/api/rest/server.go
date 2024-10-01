package rest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/echo-swagger"

	_ "songster/api/rest/docs"
	"songster/api/rest/handlers"
	"songster/app"
)

const (
	basePath = "/api/v1"
)

type Server struct {
	echo *echo.Echo
	app  *app.App
	port string
}

func New(app *app.App, port string) *Server {
	return &Server{
		app:  app,
		port: port,
	}
}

func (s *Server) Init() {
	s.echo = echo.New()
	s.echo.HideBanner = true
	s.echo.HidePort = true

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://127.0.0.1:8080",
			"http://localhost:8080",
		},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))
	s.echo.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	s.echo.Use(recoveryMiddleware)

	apiHandler := handlers.NewAPIHandler(s.app)

	publicAPI := s.echo.Group(basePath)

	addUrls(publicAPI, apiHandler.PublicURLs())

	s.echo.GET("/docs/*", echoSwagger.WrapHandler)
}

func (s *Server) Run() error {
	return s.echo.Start(":" + s.port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func addUrls(group *echo.Group, urls map[string]map[string]echo.HandlerFunc) {
	for path := range urls {
		for method, handler := range urls[path] {
			group.Add(method, path, handler)
		}
	}
}

func recoveryMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
				c.Error(echo.ErrInternalServerError)
			}
		}()
		return next(c)
	}
}
