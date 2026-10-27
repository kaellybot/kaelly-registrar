package application

import (
	"errors"
	"testing"

	"github.com/kaellybot/kaelly-registrar/models/constants"
	"github.com/spf13/viper"
)

func TestNew_UnsetOrUnknownGameIsRefused(t *testing.T) {
	for _, value := range []any{nil, 0, 4, "retro"} {
		viper.Set(constants.Game, value)
		t.Cleanup(func() { viper.Set(constants.Game, nil) })

		if _, err := New(); !errors.Is(err, errUnknownGame) {
			t.Errorf("GAME=%v: New() = %v, want %v", value, err, errUnknownGame)
		}
	}
}

func TestNew_KnownGameIsKept(t *testing.T) {
	viper.Set(constants.Game, 3)
	t.Cleanup(func() { viper.Set(constants.Game, nil) })

	app, err := New()
	if err != nil {
		t.Fatalf("New() = %v, want nil", err)
	}

	if got := app.game.AMQPGame.String(); got != "DOFUS_RETRO" {
		t.Errorf("game = %v, want DOFUS_RETRO", got)
	}
}
