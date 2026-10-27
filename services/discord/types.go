package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/kaellybot/kaelly-registrar/models/constants"
)

var (
	errNoApplication       = errors.New("no production application known for this game")
	errApplicationMismatch = errors.New("CLIENT_ID is not the production application of this game")
)

type Service interface {
	RegisterCommands(game constants.AnkamaGame) error
	DeleteCommands() error
	Shutdown()
}

type Impl struct {
	session *discordgo.Session
}
