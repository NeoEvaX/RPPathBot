package bot

import (
	"fmt"
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

func newMessage(s *discordgo.Session, message *discordgo.MessageCreate) {
	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == s.State.User.ID {
		return
	}

	brokenMessage := splitMessage(message.Content)

	switch {
	case brokenMessage[0] == ".setupgame":
		permissions, err := s.UserChannelPermissions(message.Author.ID, message.ChannelID)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		if permissions&discordgo.PermissionAdministrator != discordgo.PermissionAdministrator {
			slog.Warn("User is not admin")
			slog.Warn("Permissions", slog.Int64("Permissions:", permissions))
			return
		}

		slog.Info("BrokenMessage", slog.Any("Message", brokenMessage))
		if !verifySetupGameMessage(brokenMessage) {
			s.ChannelMessageSend(
				message.ChannelID,
				"Command formatted poorly. Try again: "+
					"`"+".setupgame"+"`"+
					"\"Name of game\" GM @GM @GM2(optional) Players @Player1 @Player2 etc",
			)

			slog.Error("Badly formatted command")
			return
		}
		gameName := brokenMessage[1]
		gameMasters, err := parseGamemasters(s, brokenMessage)
		if err != nil {
			return
		}

		gamePlayers, err := parsePlayers(s, brokenMessage)
		if err != nil {
			return
		}

		slog.Info("Found Gamemasters", slog.Any("gameMasters", gameMasters))
		slog.Info("Found Players", slog.Any("gamePlayers", gamePlayers))

		everyoneRole, err := findRole(s, "@everyone")
		if err != nil {
			slog.Error("Error finding @everyone role", slog.Any("err", err))
			return
		}

		rodRole, err := findRole(s, "Rod")
		if err != nil {
			slog.Error("Error finding @Rod role", slog.Any("err", err))
			return
		}

		channels, err := s.GuildChannels(GuildID)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		if doesCategoryExist(channels, gameName) {
			s.ChannelMessageSend(message.ChannelID, "Game already exists! Try again.")

			slog.Error("Game already exists")
			return
		}

		slog.Info(fmt.Sprintf("Creating Category: %s", gameName))
		newCategoryPermissionOverrides := []*discordgo.PermissionOverwrite{
			{ID: everyoneRole.ID, Type: discordgo.PermissionOverwriteTypeRole, Allow: 0, Deny: discordgo.PermissionViewChannel},
			{ID: rodRole.ID, Type: discordgo.PermissionOverwriteTypeRole, Allow: discordgo.PermissionViewChannel, Deny: 0},
		}

		for _, gm := range gameMasters {
			newCategoryPermissionOverrides = append(newCategoryPermissionOverrides, gmPermissions(gm)...)
		}

		slog.Info("Creating Category", slog.Any("Permissions", newCategoryPermissionOverrides))
		newCategoryData := discordgo.GuildChannelCreateData{Name: gameName, Type: discordgo.ChannelTypeGuildCategory, PermissionOverwrites: newCategoryPermissionOverrides}
		newCategory, err := s.GuildChannelCreateComplex(message.GuildID, newCategoryData)
		if err != nil {
			slog.Error("Error creating category", slog.Any("err", err))
			return
		}

		newChannelPermissionsOverwrites := []*discordgo.PermissionOverwrite{}
		newChannelReadonlyPermissionsOverwrites := []*discordgo.PermissionOverwrite{}

		for _, player := range gamePlayers {
			newChannelPermissionsOverwrites = append(newCategoryPermissionOverrides, playerPermissionsWrite(player)...)
		}

		for _, player := range gamePlayers {
			newChannelReadonlyPermissionsOverwrites = append(newCategoryPermissionOverrides, playerPermissionsRead(player)...)
		}

		slog.Info("Creating Game Channels",
			slog.Any("Write Permissions", newChannelPermissionsOverwrites),
			slog.Any("Read Permissions", newChannelReadonlyPermissionsOverwrites))

		channelData := discordgo.GuildChannelCreateData{
			Name:                 "story",
			Type:                 discordgo.ChannelTypeGuildText,
			ParentID:             newCategory.ID,
			PermissionOverwrites: newChannelPermissionsOverwrites,
		}
		s.GuildChannelCreateComplex(message.GuildID, channelData)

		channelData = discordgo.GuildChannelCreateData{
			Name:                 "ooc",
			Type:                 discordgo.ChannelTypeGuildText,
			ParentID:             newCategory.ID,
			PermissionOverwrites: newChannelPermissionsOverwrites,
		}
		s.GuildChannelCreateComplex(message.GuildID, channelData)

		channelData = discordgo.GuildChannelCreateData{
			Name:                 "maps",
			Type:                 discordgo.ChannelTypeGuildText,
			ParentID:             newCategory.ID,
			PermissionOverwrites: newChannelReadonlyPermissionsOverwrites,
		}
		s.GuildChannelCreateComplex(message.GuildID, channelData)

		channelData = discordgo.GuildChannelCreateData{
			Name:                 "combat-stats",
			Type:                 discordgo.ChannelTypeGuildText,
			ParentID:             newCategory.ID,
			PermissionOverwrites: newChannelReadonlyPermissionsOverwrites,
		}
		s.GuildChannelCreateComplex(message.GuildID, channelData)
	}
}

func doesCategoryExist(channels []*discordgo.Channel, newName string) bool {
	for _, c := range channels {
		// Check if channel is a guild text channel and not a voice or DM channel
		if c.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		if c.Name == newName {
			return true
		}
	}
	return false
}
