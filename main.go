package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"santhor.tv/internal/handler"
	"santhor.tv/internal/youtube"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("ERRO: Falha ao carregar o arquivo .env: %v", err)
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func main() {
	loadEnv()

	port := getEnv("PORT", "8080")
	apiKey := getEnv("YOUTUBE_API_KEY", "")
	channelID := getEnv("YOUTUBE_CHANNEL_ID", "")

	if apiKey == "" || channelID == "" || port == "" {
		log.Fatal("ERRO: As variáveis de ambiente 'YOUTUBE_API_KEY', 'YOUTUBE_CHANNEL_ID' e 'SERVER_PORT' são obrigatórias.")
	}

	ytClient := youtube.NewClient(apiKey)
	redirectHandler := handler.NewRedirectHandler(ytClient, channelID)

	mux := http.NewServeMux()
	mux.Handle("/", redirectHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Servidor iniciado em http://localhost:%s \nPronto para redirecionar https://santhor.tv", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("ERRO: Falha ao iniciar o servidor: %v", err)
	}
}
