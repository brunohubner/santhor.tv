package main

import (
	"log"
	"net/http"
	"os"

	"santhor.tv/internal/handler"
	"santhor.tv/internal/youtube"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	port := getEnv("SERVER_PORT", "8080")
	apiKey := getEnv("YOUTUBE_API_KEY", "")
	channelID := getEnv("YOUTUBE_CHANNEL_ID", "")

	if apiKey == "" || channelID == "" {
		log.Fatal("As variáveis de ambiente YOUTUBE_API_KEY e YOUTUBE_CHANNEL_ID são obrigatórias.")
	}

	ytClient := youtube.NewClient(apiKey)
	redirectHandler := handler.NewRedirectHandler(ytClient, channelID)

	mux := http.NewServeMux()
	mux.Handle("/", redirectHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Servidor iniciado na porta %s. Pronto para redirecionar santhor.tv!", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
