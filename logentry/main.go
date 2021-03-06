// Package logentry describes the format of log entries.
package logentry

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const timeFormat = "2006-01-02T15:04:05.000000-07:00"

const (
	HTime = iota
	HFetchType
	HOp
	HType
	HID
)

type Attachment struct {
	discordgo.MessageAttachment
	MessageID string
}

type Reaction struct {
	discordgo.MessageReaction
	Count int
}

type Embed struct {
	discordgo.MessageEmbed
	MessageID string
}

func idsFromUsers(users []*discordgo.User) (ids []string) {
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	return
}

func formatBool(name string, variable bool) string {
	if variable {
		return name
	} else {
		return ""
	}
}

func formatChannelType(t discordgo.ChannelType) string {
	switch t {
	case discordgo.ChannelTypeGuildText:
		return "text"
	case discordgo.ChannelTypeGuildVoice:
		return "voice"
	case discordgo.ChannelTypeGuildCategory:
		return "category"
	case discordgo.ChannelTypeDM:
		return "dm"
	case discordgo.ChannelTypeGroupDM:
		return "groupdm"
	case discordgo.ChannelTypeGuildNews:
		return "news"
	case discordgo.ChannelTypeGuildStore:
		return "store"
	default:
		log.Panicf("unsupported channel type %v", t)
		return "invalid"
	}
}

func Timestamp() string {
	return time.Now().Format(timeFormat)
}

func Type(v interface{}) string {
	switch v.(type) {
	case *discordgo.Message:
		return "message"
	case *Attachment:
		return "attachment"
	case *Reaction:
		return "reaction"
	case *Embed:
		return "embed"
	case *discordgo.Guild:
		return "guild"
	case *discordgo.Member:
		return "member"
	case *discordgo.Role:
		return "role"
	case *discordgo.Channel:
		return "channel"
	case *discordgo.PermissionOverwrite:
		return "permoverwrite"
	case *discordgo.Emoji:
		return "emoji"
	default:
		panic("unsupported type")
	}
}

func Make(ftype, op string, v interface{}) []string {
	var row []string

	switch v := v.(type) {
	case *discordgo.Message:
		row = []string{
			v.ID,
			v.Author.ID,
			string(v.EditedTimestamp),
			formatBool("tts", v.Tts),
			v.Content,
			formatBool("webhook", v.WebhookID != ""),
			v.Author.Username,
			v.Author.Avatar,
		}
		// only webhooks can override username/avatar at the moment
		if v.WebhookID == "" {
			row[6] = ""
			row[7] = ""
		}
	case *Attachment:
		row = []string{v.ID, v.MessageID, v.Filename}
	case *Reaction:
		row = []string{
			v.UserID,
			v.MessageID,
			v.Emoji.APIName(),
			strconv.Itoa(v.Count),
		}
	case *Embed:
		j, err := json.Marshal(v.MessageEmbed)
		if err != nil {
			panic(err)
		}

		row = []string{v.MessageID, string(j)}
	case *discordgo.Guild:
		row = []string{
			v.ID,
			v.Name,
			v.Icon,
			v.Splash,
			v.OwnerID,
			v.AfkChannelID,
			strconv.Itoa(v.AfkTimeout),
			formatBool("embeddable", v.EmbedEnabled),
			v.EmbedChannelID,
		}
	case *discordgo.Member:
		sort.StringSlice(v.Roles).Sort()
		row = []string{
			v.User.ID,
			v.User.Username,
			v.User.Discriminator,
			v.User.Avatar,
			v.Nick,
			strings.Join(v.Roles, ","),
		}
	case *discordgo.Role:
		row = []string{
			v.ID,
			v.Name,
			strconv.Itoa(v.Color),
			strconv.Itoa(v.Position),
			strconv.Itoa(v.Permissions),
			formatBool("hoist", v.Hoist),
		}
	case *discordgo.Channel:
		row = []string{
			v.ID,
			formatChannelType(v.Type),
			strconv.Itoa(v.Position),
			v.Name,
			v.Topic,
			formatBool("nsfw", v.NSFW),
			v.ParentID,
			strings.Join(idsFromUsers(v.Recipients), ","),
			v.Icon,
		}
	case *discordgo.PermissionOverwrite:
		row = []string{
			v.ID,
			v.Type,
			strconv.Itoa(v.Allow),
			strconv.Itoa(v.Deny),
		}
	case *discordgo.Emoji:
		row = []string{
			v.ID,
			v.Name,
			formatBool("nocolons", !v.RequireColons),
		}
	default:
		panic("unsupported type")
	}

	return append([]string{Timestamp(), ftype, op, Type(v)}, row...)
}
