package constants

import (
	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	commands "github.com/kaellybot/kaelly-commands"
)

// AnkamaGame ties a game to its command set and to the production application
// allowed to receive it. Registering replaces the application's whole command set,
// so a GAME that disagrees with the credentials would silently overwrite another
// bot's commands: the application is checked against ApplicationID first.
type AnkamaGame struct {
	AMQPGame amqp.Game
	// ApplicationID is the production Discord application of the game.
	ApplicationID string
	GetCommands   func() []*discordgo.ApplicationCommand
}

func GetGames() []AnkamaGame {
	return []AnkamaGame{
		{
			AMQPGame:      amqp.Game_DOFUS_GAME,
			ApplicationID: "202916641414184960",
			GetCommands:   commands.GetDofusCommands,
		},
		{
			AMQPGame:      amqp.Game_DOFUS_TOUCH,
			ApplicationID: "393925392618094612",
			GetCommands:   commands.GetTouchCommands,
		},
		{
			AMQPGame:      amqp.Game_DOFUS_RETRO,
			ApplicationID: "1239972754132832256",
			GetCommands:   commands.GetRetroCommands,
		},
	}
}

// GetGame returns the game registered for amqpGame, and whether there is one.
func GetGame(amqpGame amqp.Game) (AnkamaGame, bool) {
	for _, game := range GetGames() {
		if game.AMQPGame == amqpGame {
			return game, true
		}
	}

	return AnkamaGame{}, false
}
