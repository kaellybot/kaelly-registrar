package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Service interface {
	RegisterCommands() error
	Shutdown()
}

type Impl struct {
	session *discordgo.Session
}
