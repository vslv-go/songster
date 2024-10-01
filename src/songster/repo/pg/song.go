package pg

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"songster/models"
)

func (db *DB) SaveSong(ctx context.Context, song *models.Song) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := tx.Where("song_id = ?", song.ID).Delete(&models.Couplet{}).Error
	if err != nil {
		return fmt.Errorf("failed to clear couplets: %w", err)
	}

	err = tx.Save(song).Error
	if err != nil {
		return fmt.Errorf("failed to save song: %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (db *DB) Song(ctx context.Context, id int64, offset, limit int) (*models.Song, int64, error) {
	var (
		song  models.Song
		total int64
	)

	err := db.WithContext(ctx).Model(&models.Couplet{}).Where("song_id = ?", id).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get couplets count: %w", err)
	}

	err = db.WithContext(ctx).
		Preload("Couplets", func(db *gorm.DB) *gorm.DB {
			return db.Order("id").Offset(offset).Limit(limit)
		}).
		Find(&song, id).
		Limit(1).
		Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get song: %w", err)
	}

	return &song, total, err
}

func (db *DB) Songs(ctx context.Context, filter models.SongFilter, offset, limit int) ([]models.Song, int64, error) {
	var (
		songs []models.Song
		total int64
	)

	q := db.WithContext(ctx).Model(&models.Song{})

	if filter.Band != "" {
		q.Where("band ilike ?", "%"+filter.Band+"%")
	}

	if filter.Song != "" {
		q.Where("song ilike ?", "%"+filter.Song+"%")
	}

	if filter.Link != "" {
		q.Where("link ilike ?", "%"+filter.Link+"%")
	}

	if !filter.ReleaseDateFrom.IsZero() {
		q.Where("release_date >= ?", filter.ReleaseDateFrom)
	}

	if !filter.ReleaseDateTo.IsZero() {
		q.Where("release_date <= ?", filter.ReleaseDateTo)
	}

	err := q.Count(&total).
		Offset(offset).
		Limit(limit).
		Find(&songs).
		Error

	return songs, total, err
}

func (db *DB) DeleteSong(ctx context.Context, id int64) error {
	return db.WithContext(ctx).Delete(&models.Song{}, id).Error
}
