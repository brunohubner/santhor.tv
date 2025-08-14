package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type searchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
	} `json:"items"`
}

type CachedVideoURL struct {
	videoURL   string
	expiration time.Time
}

type Client struct {
	apiKey         string
	httpClient     *http.Client
	cachedVideoURL *CachedVideoURL
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cachedVideoURL: nil,
	}
}

func (c *Client) GetLatestVideoURL(ctx context.Context, channelID string) (string, error) {
	if c.cachedVideoURL != nil && time.Now().Before(c.cachedVideoURL.expiration) {
		return c.cachedVideoURL.videoURL, nil
	}

	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%s&order=date&maxResults=1&key=%s",
		channelID,
		c.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("ERRO: Falha ao criar requisição para a API do YouTube: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ERRO: Falha ao chamar a API do YouTube: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ERRO: API do YouTube retornou status inesperado: %s", resp.Status)
	}

	var responsePayload searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return "", fmt.Errorf("ERRO: Falha ao decodificar resposta do YouTube: %w", err)
	}

	if len(responsePayload.Items) == 0 {
		return "", fmt.Errorf("ERRO: Nenhum vídeo encontrado para o canal %s", channelID)
	}

	videoID := responsePayload.Items[0].ID.VideoID
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	cachedVideoURL := &CachedVideoURL{
		expiration: time.Now().Add(time.Minute),
		videoURL:   videoURL,
	}
	c.cachedVideoURL = cachedVideoURL

	return videoURL, nil
}
