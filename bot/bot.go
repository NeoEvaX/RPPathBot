package bot

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		slog.Error("Error message")
	}
}

func Run() {
	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(newMessage)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running until there is NO os interruption (ctrl + C)
	slog.Info("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Hello World😃, %s", message.Content))
	case strings.Contains(message.Content, ".setupgame"):
		permissions, err := discord.UserChannelPermissions(message.Author.ID, message.ChannelID)
		if err != nil {
			slog.Error(err.Error())
		}
		if permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {

			slog.Info("Creating Category: setup")
			newCategory, err := discord.GuildChannelCreate(message.GuildID, "setup", discordgo.ChannelTypeGuildCategory)
			slog.Info("Created Category: setup")

			if err == nil {
				slog.Info("Creating Channel: channel")
				channelData := discordgo.GuildChannelCreateData{Name: "channel", Type: discordgo.ChannelTypeGuildText, ParentID: newCategory.ID}
				discord.GuildChannelCreateComplex(message.GuildID, channelData)
				slog.Info("Created Channel: channel")
			}
		} else {
			slog.Warn("User is not admin")
			slog.Warn("Permissions", slog.Int64("Permissions:", permissions))
		}
	case strings.Contains(message.Content, "!bye"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
		// add more cases if required
	}
}
