package bot

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken string
	GuildID  string
	s        *discordgo.Session
)

func Run() {
	// create a session
	s, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		slog.Error("Invalid bot parameters", slog.Any("err", err.Error()))
		return
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info("Logged in as:", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(newMessage)
	// open session
	err = s.Open()
	if err != nil {
		slog.Error("Cannot open session", slog.Any("err", err.Error()))
		return
	}

	defer s.Close() // close session, after function termination

	// keep bot running until there is NO os interruption (ctrl + C)
	slog.Info("Bot running....")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	slog.Info("Press Ctrl + C to exit")
	<-stop

	slog.Info("Bot stopped")
}
