package models

import (
	"errors"
	"fmt"
	"time"
)

const (
	ReleaseDateFormat = "02.01.2006"
)

type Song struct {
	ID          int64
	Band        string
	Song        string
	Link        string
	ReleaseDate time.Time

	Couplets []Couplet `gorm:"foreignKey:song_id;references:ID"`
}

func (*Song) TableName() string {
	return "songs"
}

func (s *Song) Validate() error {
	if s.Band == "" {
		return errors.New("band is required")
	}

	if s.Song == "" {
		return errors.New("song is required")
	}

	if s.Link == "" {
		return errors.New("link is required")
	}

	if s.ReleaseDate.IsZero() {
		return errors.New("release date is required")
	}

	if s.ReleaseDate.After(time.Now()) {
		return errors.New("release date is in the future")
	}

	return nil
}

type SongFilter struct {
	Band, Song, Link string
	ReleaseDateFrom  time.Time
	ReleaseDateTo    time.Time
}

func (f *SongFilter) Validate() error {
	if !f.ReleaseDateFrom.IsZero() && f.ReleaseDateFrom.After(time.Now()) {
		return fmt.Errorf("release date 'from' must be before now")
	}

	if !f.ReleaseDateTo.IsZero() && f.ReleaseDateTo.After(time.Now()) {
		return fmt.Errorf("release date 'to' must be before now")
	}

	if !f.ReleaseDateFrom.IsZero() && !f.ReleaseDateTo.IsZero() && f.ReleaseDateFrom.After(f.ReleaseDateTo) {
		return fmt.Errorf("release date 'from' must be before 'to'")
	}

	return nil
}
