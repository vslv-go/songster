package handlers

type Song struct {
	ID          int64  `json:"id"`
	Band        string `json:"band"`
	Song        string `json:"song"`
	Link        string `json:"link"`
	ReleaseDate string `json:"release_date" example:"02.01.2006"`
}

type SongFilter struct {
	Band  string `query:"band"`
	Song  string `query:"song"`
	Link  string `query:"link"`
	Dates string `query:"dates" json:"dates" example:"02.01.2006-02.01.2007"`
}

type Couplet struct {
	ID     int64  `json:"id"`
	SongID int64  `json:"song_id"`
	Text   string `json:"text"`
}

type SongsResponse struct {
	Songs []Song `json:"songs"`
	Total int64  `json:"total"`
}

type SongCoupletsResponse struct {
	Song     Song      `json:"song"`
	Couplets []Couplet `json:"couplets"`
	Total    int64     `json:"total"`
}

type SongResponse struct {
	Song Song `json:"song"`
}

type AddSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}
