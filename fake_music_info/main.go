package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	port = ":8081"
)

func main() {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://127.0.0.1:8081",
			"http://localhost:8081",
		},
		AllowMethods:     []string{echo.GET},
		AllowCredentials: true,
	}))
	server.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	server.Use(recoveryMiddleware)

	publicAPI := server.Group("/api/v1")

	publicAPI.GET("/info", info)

	log.Println("Starting fake Info API at ", port)
	log.Fatal(server.Start(port))
}

func info(c echo.Context) error {
	group := c.QueryParam("group")
	if group == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "group is required")
	}

	song := c.QueryParam("song")
	if song == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "song is required")
	}

	fake := faker.New()

	if fake.BoolWithChance(5) {
		return echo.ErrInternalServerError
	}

	dt := fake.Time().Time(time.Now()).Format("02.01.2006")

	coupletsCount := fake.IntBetween(4, 8)
	text := strings.Join(fake.Lorem().Sentences(coupletsCount), "\n\n")

	return c.JSON(http.StatusOK, echo.Map{
		"release_date": dt,
		"text":         text,
		"link":         fake.YouTube().GenerateFullURL(),
	})
}

func recoveryMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				c.Error(echo.ErrInternalServerError)
			}
		}()
		return next(c)
	}
}
