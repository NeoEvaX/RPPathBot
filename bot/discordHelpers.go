package bot

import "github.com/bwmarrin/discordgo"

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
