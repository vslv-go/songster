package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"songster/app"
)

const (
	defaultPerPageCount = 10
)

type Handler struct {
	app *app.App
}

// NewAPIHandler
// @title Songster API
// @version 1.0
// @description This is a sample API
// @host localhost:8080
// @BasePath /api/v1
func NewAPIHandler(app *app.App) *Handler {
	return &Handler{
		app: app,
	}
}

func handleError(err error) error {
	switch {
	case app.IsBadRequestError(err):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case app.IsNotFoundError(err):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	default:
		log.WithError(err).Error()
		return echo.ErrInternalServerError
	}
}

func parseOffsetLimit(pageParam, countParam string) (offset, limit int) {
	page, _ := strconv.Atoi(pageParam)
	count, _ := strconv.Atoi(countParam)

	if count == -1 {
		return 0, -1
	}

	if count < 1 {
		count = defaultPerPageCount
	}

	if page < 1 {
		page = 1
	}

	offset = (page - 1) * count
	limit = count

	return
}
