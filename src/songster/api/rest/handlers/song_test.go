package handlers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateDate(t *testing.T) {
	testDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		date     string
		wantTime time.Time
		wantErr  bool
	}{
		{"", "", time.Time{}, true},
		{"", "test", time.Time{}, true},
		{"", "2000", time.Time{}, true},
		{"", "2000-01-01", time.Time{}, true},
		{"", "01.01.2000", testDate, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tm, err := validateDate(tt.date)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.wantTime, tm)
		})
	}
}

func TestValidateDates(t *testing.T) {
	testFrom := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	testTo := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		dates    string
		wantFrom time.Time
		wantTo   time.Time
		wantErr  bool
	}{
		{"", "", time.Time{}, time.Time{}, false},
		{"", "-", time.Time{}, time.Time{}, true},
		{"", "1-1", time.Time{}, time.Time{}, true},
		{"", "2000-01-01", time.Time{}, time.Time{}, true},
		{"", "2000-01-01-2010-01-01", time.Time{}, time.Time{}, true},
		{"", "01.01.2000", time.Time{}, time.Time{}, true},
		{"", "01.01.2000-", time.Time{}, time.Time{}, true},
		{"", "01.01.2000-01.01.2010", testFrom, testTo, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tmFrom, tmTo, err := validateDates(tt.dates)
			require.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				require.Equal(t, tt.wantFrom, tmFrom)
				require.Equal(t, tt.wantTo, tmTo)
			}
		})
	}
}
