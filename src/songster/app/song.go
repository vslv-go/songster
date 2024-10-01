package app

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"songster/models"
)

type songRepo interface {
	SaveSong(ctx context.Context, song *models.Song) error
	Song(ctx context.Context, id int64, offset, limit int) (*models.Song, int64, error)
	Songs(ctx context.Context, filter models.SongFilter, offset, limit int) ([]models.Song, int64, error)
	DeleteSong(ctx context.Context, id int64) error
}

func (app *App) CreateSong(ctx context.Context, group, song string) (*models.Song, error) {
	if group == "" || song == "" {
		return nil, NewBadRequestError(nil, "both group and song are required")
	}

	songModel, err := app.infoClient.GetSong(group, song)
	if err != nil {
		return nil, NewInternalError(err, "failed to get song info")
	}

	err = app.repo.SaveSong(ctx, songModel)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, NewBadRequestError(nil, fmt.Sprintf("song %q of band %q already exists", song, group))
		}
		return nil, NewInternalError(err, "failed to save song")
	}

	return songModel, nil
}

func (app *App) UpdateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	if song == nil || song.ID == 0 {
		return nil, NewBadRequestError(nil, "song is required")
	}

	found, _, err := app.repo.Song(ctx, song.ID, 0, 0)
	if err != nil {
		return nil, NewInternalError(err, "failed to get song")
	}

	if found == nil || found.ID == 0 {
		return nil, NewNotFoundError(nil, "song not found")
	}

	err = app.repo.SaveSong(ctx, song)
	if err != nil {
		return nil, NewInternalError(err, "failed to save song")
	}

	return song, nil
}

func (app *App) Song(ctx context.Context, id int64, offset, limit int) (*models.Song, int64, error) {
	if id == 0 {
		return nil, 0, NewBadRequestError(nil, "id is required")
	}

	song, totalCouplets, err := app.repo.Song(ctx, id, offset, limit)
	if err != nil {
		return nil, 0, NewInternalError(err, "failed to get song")
	}

	if song == nil || song.ID == 0 {
		return nil, 0, NewNotFoundError(nil, "song not found")
	}

	return song, totalCouplets, nil
}

func (app *App) Songs(ctx context.Context, filter models.SongFilter, offset, limit int) ([]models.Song, int64, error) {
	songs, total, err := app.repo.Songs(ctx, filter, offset, limit)
	if err != nil {
		return nil, 0, NewInternalError(err, "failed to get songs")
	}
	return songs, total, nil
}

func (app *App) DeleteSong(ctx context.Context, id int64) error {
	if id < 1 {
		return NewBadRequestError(nil, "id is required")
	}

	err := app.repo.DeleteSong(ctx, id)
	if err != nil {
		return NewInternalError(err, "failed to delete song")
	}

	return nil
}
