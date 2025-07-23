package discord

import (
	"github.com/bwmarrin/discordgo"
	commands "github.com/kaellybot/kaelly-commands"
	"github.com/kaellybot/kaelly-registrar/models/constants"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New(token string) (*Impl, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Error().Err(err).Msgf("Connecting to Discord gateway failed")
		return nil, err
	}

	return &Impl{session: dg}, nil
}

func (service *Impl) RegisterCommands() error {
	guildID := ""
	if !viper.GetBool(constants.Production) {
		log.Info().Msgf("Development mode enabled, registering commands in dedicated development guild")
		guildID = constants.DevelopmentGuildID
	}

	appCommands, err := service.session.ApplicationCommandBulkOverwrite(viper.GetString(constants.ClientID),
		guildID, commands.GetCommands())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create commands, registration stopped")
		return err
	}
	log.Info().Msgf("Commands successfully registered!")
	for _, command := range appCommands {
		log.Info().Msgf("%v (type %v): %v", command.Name, command.Type, command.ID)
	}

	return nil
}

func (service *Impl) DeleteCommands() error {
	commands, err := service.session.ApplicationCommands(viper.GetString(constants.ClientID), "")
	if err != nil {
		log.Error().Err(err).Msgf("Failed to retrieve registered commands, deletion stopped")
		return err
	}

	for _, cmd := range commands {
		errDel := service.session.ApplicationCommandDelete(viper.GetString(constants.ClientID), "", cmd.ID)
		if errDel != nil {
			log.Error().Err(errDel).Msgf("Failed to delete commands, deletion stopped")
			return errDel
		}
		log.Info().Msgf("Commands successfully deleted!")
	}

	return nil
}

func (service *Impl) Shutdown() {
	log.Info().Int(constants.LogShard, service.session.ShardID).Msgf("Closing Discord connections...")
	err := service.session.Close()
	if err != nil {
		log.Warn().Err(err).Msgf("Cannot close session and shutdown correctly")
	}
}
