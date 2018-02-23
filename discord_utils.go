package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	// Discord channel types are enums
	// GUILD_TEXT	0
	// DM	1
	// GUILD_VOICE	2
	// GROUP_DM	3
	// GUILD_CATEGORY	4

	// GuildText is the guild text channel
	GuildText = iota

	// DMType is the direct message channel type
	DMType

	// GuildVoice is the voice channel type
	GuildVoice

	// GroupDM is the group dm type
	GroupDM

	// GuildCategory is the category? type
	GuildCategory
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

func isDirectMessage(session *discordgo.Session, channelID string) (bool, error) {
	channel, err := session.Channel(channelID)
	if err != nil {
		return false, fmt.Errorf("Discord Error retrieving channel: %s", err)
	}

	return channel.Type == DMType, nil
}
