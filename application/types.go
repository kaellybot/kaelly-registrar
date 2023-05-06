package application

import (
	"github.com/kaellybot/kaelly-registrar/services/discord"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	discordService discord.Service
}
