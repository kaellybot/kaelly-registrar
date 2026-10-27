package discord

import (
	"errors"
	"testing"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-registrar/models/constants"
)

func game(t *testing.T, amqpGame amqp.Game) constants.AnkamaGame {
	t.Helper()
	game, found := constants.GetGame(amqpGame)
	if !found {
		t.Fatalf("no game for %v", amqpGame)
	}
	return game
}

func TestCheckApplication_GameOwnApplicationIsAccepted(t *testing.T) {
	dofus := game(t, amqp.Game_DOFUS_GAME)

	if err := checkApplication(dofus, dofus.ApplicationID); err != nil {
		t.Errorf("checkApplication() = %v, want nil", err)
	}
}

// GAME=1 with another bot's credentials would replace that bot's whole command set.
func TestCheckApplication_AnotherApplicationIsRefused(t *testing.T) {
	err := checkApplication(game(t, amqp.Game_DOFUS_GAME), "123456789")

	if !errors.Is(err, errApplicationMismatch) {
		t.Errorf("checkApplication() = %v, want %v", err, errApplicationMismatch)
	}
}

func TestCheckApplication_GameWithoutApplicationIsRefused(t *testing.T) {
	err := checkApplication(constants.AnkamaGame{AMQPGame: amqp.Game_DOFUS_RETRO}, "")

	if !errors.Is(err, errNoApplication) {
		t.Errorf("checkApplication() = %v, want %v", err, errNoApplication)
	}
}
