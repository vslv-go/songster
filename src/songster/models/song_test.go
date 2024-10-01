package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSongFilterValidate(t *testing.T) {
	now := time.Now()
	dateAfter := now.Add(time.Hour)
	dateBefore := now.Add(-time.Hour)

	tests := []struct {
		name    string
		filter  SongFilter
		wantErr bool
	}{
		{"", SongFilter{}, false},
		{"", SongFilter{ReleaseDateFrom: dateAfter}, true},
		{"", SongFilter{ReleaseDateTo: dateAfter}, true},
		{"", SongFilter{ReleaseDateFrom: now, ReleaseDateTo: dateBefore}, true},
		{"", SongFilter{ReleaseDateFrom: dateBefore, ReleaseDateTo: now}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.filter.Validate()
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestSongValidate(t *testing.T) {
	now := time.Now()
	dateAfter := now.Add(time.Hour)
	dateBefore := now.Add(-time.Hour)

	tests := []struct {
		name    string
		song    Song
		wantErr bool
	}{
		{"", Song{}, true},
		{"", Song{Band: "1"}, true},
		{"", Song{Song: "1"}, true},
		{"", Song{Link: "1"}, true},
		{"", Song{ReleaseDate: dateBefore}, true},
		{"", Song{Band: "1", Song: "1"}, true},
		{"", Song{Band: "1", Song: "1", Link: "1"}, true},
		{"", Song{Band: "1", Song: "1", Link: "1", ReleaseDate: dateAfter}, true},
		{"", Song{Band: "1", Song: "1", Link: "1", ReleaseDate: dateBefore}, false},
		{"", Song{Band: "1", Song: "1", Link: "1", ReleaseDate: now}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.song.Validate()
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
