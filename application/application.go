package application

import (
	"fmt"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-registrar/models/constants"
	"github.com/kaellybot/kaelly-registrar/services/discord"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	game, found := constants.GetGame(amqp.Game(viper.GetInt32(constants.Game)))
	if !found {
		return nil, fmt.Errorf("%w, got %q", errUnknownGame, viper.GetString(constants.Game))
	}

	discordService, err := discord.New(viper.GetString(constants.Token))
	if err != nil {
		return nil, err
	}

	return &Impl{game: game, discordService: discordService}, nil
}

func (app *Impl) Run() error {
	err := app.discordService.RegisterCommands(app.game)
	if err != nil {
		return err
	}

	return nil
}

func (app *Impl) Shutdown() {
	app.discordService.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
