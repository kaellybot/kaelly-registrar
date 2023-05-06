package application

import (
	"github.com/kaellybot/kaelly-registrar/models/constants"
	"github.com/kaellybot/kaelly-registrar/services/discord"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	discordService, err := discord.New(viper.GetString(constants.Token))
	if err != nil {
		return nil, err
	}

	return &Impl{discordService: discordService}, nil
}

func (app *Impl) Run() error {
	err := app.discordService.RegisterCommands()
	if err != nil {
		return err
	}

	return nil
}

func (app *Impl) Shutdown() {
	app.discordService.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
