package models

type Couplet struct {
	ID     int64
	SongID int64
	Text   string
}

func (Couplet) TableName() string {
	return "couplets"
}
