package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
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

func (service *Impl) RegisterCommands(game constants.AnkamaGame) error {
	clientID := viper.GetString(constants.ClientID)
	guildID := ""
	if viper.GetBool(constants.Production) {
		if err := checkApplication(game, clientID); err != nil {
			log.Error().Err(err).Str(constants.LogGame, game.AMQPGame.String()).
				Msgf("Commands would overwrite another application's, registration stopped")
			return err
		}
	} else {
		log.Info().Msgf("Development mode enabled, registering commands in dedicated development guild")
		guildID = constants.DevelopmentGuildID
	}

	log.Info().Str(constants.LogGame, game.AMQPGame.String()).Msgf("Registering commands...")
	appCommands, err := service.session.ApplicationCommandBulkOverwrite(clientID, guildID, game.GetCommands())
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

// checkApplication refuses to register a game's commands on any application but the
// game's own. The token must belong to CLIENT_ID for Discord to accept the call, so
// checking CLIENT_ID is enough to keep a GAME/credentials mismatch from replacing
// another bot's commands.
func checkApplication(game constants.AnkamaGame, clientID string) error {
	if game.ApplicationID == "" {
		return fmt.Errorf("%w: %v", errNoApplication, game.AMQPGame)
	}

	if clientID != game.ApplicationID {
		return fmt.Errorf("%w: %v expects %v, got %q", errApplicationMismatch,
			game.AMQPGame, game.ApplicationID, clientID)
	}

	return nil
}
