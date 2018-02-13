package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func extractNameAndDiscriminator(userID string) (string, string) {
	if userID == "" {
		return "", ""
	}

	splitUserID := strings.Split(userID, "#")
	if len(splitUserID) == 1 {
		return splitUserID[0], ""
	}
	return splitUserID[0], splitUserID[1]
}

func findUserID(session *discordgo.Session, guildID, userName, userDiscriminator string) (string, error) {
	members, err := session.GuildMembers(guildID, "", 1000)
	if err != nil {
		return "", fmt.Errorf("Discord error while trying to search guild members: %s", err)
	}

	// case where they did not provide discriminator
	if userDiscriminator == "" {
		for _, member := range members {
			if fmt.Sprintf("%s", member.User.Username) == userName {
				return member.User.ID, nil
			}
		}
	}
	for _, member := range members {
		if fmt.Sprintf("%s#%s", member.User.Username, member.User.Discriminator) == fmt.Sprintf("%s#%s", userName, userDiscriminator) {
			return member.User.ID, nil
		}
	}
	return "", nil
}
