package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/erotokritosVall/xmapp/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %+v", err)
	}

	server := server.New()

	server.Start(exitChannel)
}
