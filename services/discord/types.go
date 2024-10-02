package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Service interface {
	RegisterCommands() error
	DeleteCommands() error
	Shutdown()
}

type Impl struct {
	session *discordgo.Session
}
