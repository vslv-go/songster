package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"songster/models"
)

// Songs
// @Router /songs [get]
// @Summary Songs list
// @Description Get list of songs
// @Tags Songs
// @Produce json
// @Param filter query SongFilter false "Fields filter"
// @Param page query int false "Number of page"
// @Param count query int false "Songs per page"
// @Success 200 {object} SongsResponse
// @Failure 400
// @Failure 500
func (h *Handler) Songs(c echo.Context) error {
	ctx := c.Request().Context()

	var filterData SongFilter
	err := c.Bind(&filterData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	songFilter, err := songFilterToModel(filterData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = songFilter.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	offset, limit := parseOffsetLimit(c.QueryParam("page"), c.QueryParam("count"))

	songs, total, err := h.app.Songs(ctx, *songFilter, offset, limit)
	if err != nil {
		return handleError(err)
	}

	return c.JSON(http.StatusOK, SongsResponse{
		Songs: songsToResponse(songs),
		Total: total,
	})
}

// Song
// @Router /songs/{id} [get]
// @Summary Song couplets
// @Description Get song couplets
// @Tags Songs
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Number of page"
// @Param count query int false "Couplets per page"
// @Success 200 {object} SongCoupletsResponse
// @Failure 400
// @Failure 500
func (h *Handler) Song(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Song ID")
	}

	offset, limit := parseOffsetLimit(c.QueryParam("page"), c.QueryParam("count"))

	song, totalCouplets, err := h.app.Song(ctx, id, offset, limit)
	if err != nil {
		return handleError(err)
	}

	return c.JSON(http.StatusOK, SongCoupletsResponse{
		Song:     songToResponse(*song),
		Couplets: coupletsToResponse(song.Couplets),
		Total:    totalCouplets,
	})
}

// DeleteSong
// @Router /songs/{id} [delete]
// @Summary Delete song
// @Description Delete song
// @Tags Songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200
// @Failure 400
// @Failure 500
func (h *Handler) DeleteSong(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Song ID")
	}

	err := h.app.DeleteSong(ctx, id)
	if err != nil {
		return handleError(err)
	}

	return c.NoContent(http.StatusOK)
}

// UpdateSong
// @Router /songs/{id} [put]
// @Summary Update song
// @Description Update song data
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body Song true "Song params"
// @Success 200 {object} SongResponse
// @Failure 400
// @Failure 500
func (h *Handler) UpdateSong(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Song ID")
	}

	var songData Song
	err := c.Bind(&songData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	song, err := songToModel(id, songData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = song.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	updated, err := h.app.UpdateSong(ctx, song)
	if err != nil {
		return handleError(err)
	}

	return c.JSON(http.StatusOK, SongResponse{Song: songToResponse(*updated)})
}

// AddSong
// @Router /songs [post]
// @Summary Add song
// @Description Add song data
// @Tags Songs
// @Accept json
// @Produce json
// @Param song body AddSongRequest true "Song params"
// @Success 200 {object} SongResponse
// @Failure 400
// @Failure 500
func (h *Handler) AddSong(c echo.Context) error {
	ctx := c.Request().Context()

	var req AddSongRequest

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.Group == "" || req.Song == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Both Group and Song must be provided")
	}

	song, err := h.app.CreateSong(ctx, req.Group, req.Song)
	if err != nil {
		return handleError(err)
	}

	return c.JSON(http.StatusOK, SongResponse{
		Song: songToResponse(*song),
	})
}

func songsToResponse(songs []models.Song) []Song {
	converted := make([]Song, len(songs))
	for i, song := range songs {
		converted[i] = songToResponse(song)
	}
	return converted
}

func songToResponse(song models.Song) Song {
	return Song{
		ID:          song.ID,
		Band:        song.Band,
		Song:        song.Song,
		Link:        song.Link,
		ReleaseDate: song.ReleaseDate.Format("02.01.2006"),
	}
}

func coupletsToResponse(couplets []models.Couplet) []Couplet {
	converted := make([]Couplet, len(couplets))
	for i, couplet := range couplets {
		converted[i] = coupletToResponse(couplet)
	}
	return converted
}

func coupletToResponse(couplet models.Couplet) Couplet {
	return Couplet{
		ID:     couplet.ID,
		SongID: couplet.SongID,
		Text:   couplet.Text,
	}
}

func songToModel(id int64, song Song) (*models.Song, error) {
	dt, err := validateDate(song.ReleaseDate)
	if err != nil {
		return nil, err
	}
	return &models.Song{
		ID:          id,
		Band:        song.Band,
		Song:        song.Song,
		Link:        song.Link,
		ReleaseDate: dt,
	}, nil
}

func songFilterToModel(f SongFilter) (*models.SongFilter, error) {
	from, to, err := validateDates(f.Dates)
	if err != nil {
		return nil, err
	}

	return &models.SongFilter{
		Band:            f.Band,
		Song:            f.Song,
		Link:            f.Link,
		ReleaseDateFrom: from,
		ReleaseDateTo:   to,
	}, nil
}

func validateDates(s string) (from, to time.Time, err error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return
	}

	dates := strings.Split(s, "-")
	if len(dates) != 2 {
		err = errors.New("invalid range format")
		return
	}

	from, err = validateDate(dates[0])
	if err != nil {
		return
	}

	to, err = validateDate(dates[1])

	return
}

func validateDate(s string) (time.Time, error) {
	dt, err := time.Parse(models.ReleaseDateFormat, s)
	if err != nil {
		err = fmt.Errorf("failed to parse date %q: %w", s, err)
	}
	return dt, err
}
