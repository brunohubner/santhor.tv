package handler

import (
	"log"
	"net/http"

	"santhor.tv/internal/youtube"
)

type RedirectHandler struct {
	youtubeClient *youtube.Client
	channelID     string
}

func NewRedirectHandler(ytClient *youtube.Client, channelID string) *RedirectHandler {
	return &RedirectHandler{
		youtubeClient: ytClient,
		channelID:     channelID,
	}
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	videoURL, err := h.youtubeClient.GetLatestVideoURL(ctx, h.channelID)
	if err != nil {
		log.Printf("ERRO: Falha ao obter a URL do último vídeo: %v", err)
		fallbackURL := "https://www.youtube.com/@SanthorTV"
		http.Redirect(w, r, fallbackURL, http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, videoURL, http.StatusTemporaryRedirect)
}
