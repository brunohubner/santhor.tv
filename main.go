package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"santhor.tv/internal/handler"
	"santhor.tv/internal/youtube"
)

func main() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	channelID := os.Getenv("YOUTUBE_CHANNEL_ID")

	if apiKey == "" || channelID == "" || port == "" {
		log.Fatal("As variáveis de ambiente 'YOUTUBE_API_KEY', 'YOUTUBE_CHANNEL_ID' e 'SERVER_PORT' são obrigatórias.")
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
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
