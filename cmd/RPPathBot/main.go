package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	bot "nex-verse.com/RolePlayPathBot/bot"
)

func main() {
	// Set up the logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	bot.BotToken = os.Getenv("DISCORD_TOKEN")
	bot.Run() // call the run function of bot/bot.go
}
