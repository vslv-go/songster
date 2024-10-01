-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs
(
    id           BIGSERIAL PRIMARY KEY,
    band         VARCHAR(255) NOT NULL,
    song         VARCHAR(255) NOT NULL,
    link         VARCHAR(255) NOT NULL,
    release_date timestamptz NOT NULL,

    CONSTRAINT u_songs UNIQUE (band, song)
);

CREATE TABLE couplets
(
    id      BIGSERIAL PRIMARY KEY,
    song_id BIGINT REFERENCES songs (id) ON DELETE CASCADE,
    text    TEXT NOT NULL
);

CREATE INDEX idx_couplets_song_id ON couplets (song_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE couplets, songs CASCADE;
-- +goose StatementEnd
