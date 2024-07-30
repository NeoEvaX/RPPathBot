package bot

import "github.com/bwmarrin/discordgo"

func gmPermissions(gmUser discordgo.Member) []*discordgo.PermissionOverwrite {
	permissions := []*discordgo.PermissionOverwrite{
		{
			ID:    gmUser.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionManageMessages,
			Deny:  0,
		},
		{
			ID:    gmUser.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionManageChannels,
			Deny:  0,
		},
		{
			ID:    gmUser.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionViewChannel,
			Deny:  0,
		},
		{
			ID:    gmUser.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionManageRoles,
			Deny:  0,
		},
		{
			ID:    gmUser.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionSendMessages,
			Deny:  0,
		},
		{
			ID:    gmUser.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionAllChannel,
			Deny:  0,
		},
	}
	return permissions
}

func playerPermissionsWrite(playerUser discordgo.Member) []*discordgo.PermissionOverwrite {
	permissions := []*discordgo.PermissionOverwrite{
		{
			ID:    playerUser.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionSendMessages,
			Deny:  0,
		},
		{
			ID:    playerUser.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionViewChannel,
			Deny:  0,
		},
	}
	return permissions
}

func playerPermissionsRead(playerUser discordgo.Member) []*discordgo.PermissionOverwrite {
	permissions := []*discordgo.PermissionOverwrite{
		{
			ID:    playerUser.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionViewChannel,
			Deny:  discordgo.PermissionSendMessages,
		},
	}
	return permissions
}
