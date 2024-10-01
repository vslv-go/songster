package handlers

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) PublicURLs() map[string]map[string]echo.HandlerFunc {
	return map[string]map[string]echo.HandlerFunc{
		"/songs": {
			"GET":  h.Songs,
			"POST": h.AddSong,
		},
		"/songs/:id": {
			"GET":    h.Song,
			"PUT":    h.UpdateSong,
			"DELETE": h.DeleteSong,
		},
	}
}
