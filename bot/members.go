package bot

import (
	"errors"
	"log/slog"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func findMemberFromName(s *discordgo.Session, userId string) (member *discordgo.Member, err error) {
	userId = strings.Replace(userId, "<@", "", -1)
	userId = strings.Replace(userId, ">", "", -1)
	member, err = s.GuildMember(GuildID, userId)
	if err != nil {
		slog.Error("Error finding member", slog.String("userId", userId), slog.Any("err", err))
		return
	}
	return
}

func parseGamemasters(s *discordgo.Session, message []string) ([]discordgo.Member, error) {
	var gameMasters []discordgo.Member
	// start looping through after GM tag
	for _, a := range message[3:] {
		if a == "Players" {
			return gameMasters, nil
		}
		newGameMaster, err := findMemberFromName(s, a)
		if err != nil {
			slog.Error("Error finding member", slog.String("userId", a), slog.Any("err", err))
			return nil, err
		}
		gameMasters = append(gameMasters, *newGameMaster)
	}
	return gameMasters, nil
}

func parsePlayers(s *discordgo.Session, message []string) ([]discordgo.Member, error) {
	var players []discordgo.Member
	// start looping through after GMs, but skipping till its Players
	started := false
	for _, a := range message[4:] {
		if a == "Players" {
			started = true
			continue

		}
		if !started {
			continue
		}
		newPlayer, err := findMemberFromName(s, a)
		if err != nil {
			return nil, err
		}
		players = append(players, *newPlayer)
	}
	return players, nil
}

func findRole(s *discordgo.Session, roleName string) (discordgo.Role, error) {
	// Get the @everyone role
	Roles, _ := s.GuildRoles(GuildID)
	slog.Info("Roles", slog.Any("Roles", Roles))
	idx := slices.IndexFunc(Roles, func(r *discordgo.Role) bool { return r.Name == roleName })
	if idx == -1 {
		return discordgo.Role{}, errors.New("Could not find role")
	}
	return *Roles[idx], nil
}
