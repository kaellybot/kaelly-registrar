package constants

import (
	"testing"

	amqp "github.com/kaellybot/kaelly-amqp"
	commands "github.com/kaellybot/kaelly-commands"
)

func commandNames(game AnkamaGame) []string {
	names := make([]string, 0)
	for _, command := range game.GetCommands() {
		names = append(names, command.Name)
	}
	return names
}

func TestGetGame_EachGameRegistersItsOwnCommands(t *testing.T) {
	tests := []struct {
		game amqp.Game
		want int
	}{
		{game: amqp.Game_DOFUS_GAME, want: len(commands.GetDofusCommands())},
		{game: amqp.Game_DOFUS_TOUCH, want: len(commands.GetTouchCommands())},
		{game: amqp.Game_DOFUS_RETRO, want: len(commands.GetRetroCommands())},
	}

	for _, test := range tests {
		t.Run(test.game.String(), func(t *testing.T) {
			game, found := GetGame(test.game)
			if !found {
				t.Fatalf("GetGame(%v) not found", test.game)
			}

			if got := commandNames(game); len(got) != test.want {
				t.Errorf("%v registers %v, want %v commands", test.game, got, test.want)
			}
		})
	}

	retro, _ := GetGame(amqp.Game_DOFUS_RETRO)
	for _, name := range commandNames(retro) {
		if name == commands.ItemCommandName {
			t.Errorf("DOFUS Retro registers /%v, a DOFUS command", name)
		}
	}
}

// ANY_GAME is what an unset GAME reads as: it must not pick a command set.
func TestGetGame_AnyGameIsUnknown(t *testing.T) {
	if _, found := GetGame(amqp.Game_ANY_GAME); found {
		t.Error("GetGame(ANY_GAME) found, want not found")
	}
}

// Two games sharing an application would let either one overwrite the other's
// commands without the guard noticing.
func TestGetGames_EachGameHasItsOwnApplication(t *testing.T) {
	seen := make(map[string]amqp.Game)
	for _, game := range GetGames() {
		if game.ApplicationID == "" {
			t.Errorf("%v has no application", game.AMQPGame)
			continue
		}

		if other, found := seen[game.ApplicationID]; found {
			t.Errorf("%v and %v share application %v", other, game.AMQPGame, game.ApplicationID)
		}
		seen[game.ApplicationID] = game.AMQPGame
	}
}
