package application

import (
	"errors"

	"github.com/kaellybot/kaelly-registrar/models/constants"
	"github.com/kaellybot/kaelly-registrar/services/discord"
)

var errUnknownGame = errors.New("GAME must be 1 (DOFUS), 2 (DOFUS Touch) or 3 (DOFUS Retro)")

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	game           constants.AnkamaGame
	discordService discord.Service
}
