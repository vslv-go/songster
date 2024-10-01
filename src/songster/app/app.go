package app

import (
	"songster/models"
)

type Repo interface {
	songRepo
}

type InfoClient interface {
	GetSong(group, song string) (*models.Song, error)
}

type App struct {
	repo       Repo
	infoClient InfoClient
}

func New(repo Repo, client InfoClient) *App {
	return &App{
		repo:       repo,
		infoClient: client,
	}
}
