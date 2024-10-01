package info_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"songster/models"
)

const (
	basePath = "/api/v1"
	infoPath = "/info"
)

type InfoClient struct {
	addr   string
	client *http.Client
}

type infoResponse struct {
	Text        string `json:"text"`
	Link        string `json:"link"`
	ReleaseDate string `json:"release_date"`
}

func NewInfoClient(addr string) *InfoClient {
	return &InfoClient{
		addr:   addr,
		client: &http.Client{},
	}
}

func (c InfoClient) GetSong(group, song string) (*models.Song, error) {
	if group == "" || song == "" {
		return nil, errors.New("group and song are required")
	}

	baseURL, err := url.Parse(fmt.Sprintf("%s%s%s", c.addr, basePath, infoPath))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)

	baseURL.RawQuery = params.Encode()

	resp, err := c.client.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("info request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("info request failed with status code %d", resp.StatusCode)
	}

	info, err := parseInfoResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse info: %w", err)
	}

	songModel, err := fillSong(info, group, song)
	if err != nil {
		return nil, fmt.Errorf("failed to fill song: %w", err)
	}

	return songModel, nil
}

func parseInfoResponse(body io.ReadCloser) (*infoResponse, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var info infoResponse
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &info, nil
}

func fillSong(info *infoResponse, group, song string) (*models.Song, error) {
	dt, err := time.Parse(models.ReleaseDateFormat, info.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse release date: %w", err)
	}

	textParts := strings.Split(info.Text, "\n\n")
	couplets := make([]models.Couplet, len(textParts))
	for i, part := range textParts {
		couplets[i] = models.Couplet{
			Text: part,
		}
	}

	return &models.Song{
		Band:        group,
		Song:        song,
		Link:        info.Link,
		ReleaseDate: dt,
		Couplets:    couplets,
	}, nil
}
