package bot

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

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

		newChannelPermissionsOverwrites := newCategoryPermissionOverrides
		newChannelReadonlyPermissionsOverwrites := newCategoryPermissionOverrides

		for _, player := range gamePlayers {
			newChannelPermissionsOverwrites = append(newChannelPermissionsOverwrites, playerPermissionsWrite(player)...)
		}

		for _, player := range gamePlayers {
			newChannelReadonlyPermissionsOverwrites = append(newChannelReadonlyPermissionsOverwrites, playerPermissionsRead(player)...)
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

		s.ChannelMessageSend(message.ChannelID, "Game created! Have fun!")
	}
}
