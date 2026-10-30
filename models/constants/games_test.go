package constants

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	commands "github.com/kaellybot/kaelly-commands"
	"github.com/kaellybot/kaelly-registrar/models/i18n"
	di18n "github.com/kaysoro/discordgo-i18n"
)

// Discord refuses the whole command set when a single description is out of bounds.
const (
	descriptionMinLength = 1
	descriptionMaxLength = 100
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

// Game names are inserted into the translations, so a description fitting for
// DOFUS can overflow for DOFUS Touch or DOFUS Retro.
func TestGetGames_DescriptionsFitDiscordLimits(t *testing.T) {
	di18n.SetDefault(i18n.DefaultLocale)
	for _, language := range i18n.GetLanguages() {
		if err := di18n.LoadBundleFS(language.Locale, i18n.Folder, language.TranslationFile); err != nil {
			t.Fatalf("loading %v translations: %v", language.Locale, err)
		}
	}

	for _, game := range GetGames() {
		t.Run(game.AMQPGame.String(), func(t *testing.T) {
			for _, command := range game.GetCommands() {
				if command.Type != discordgo.ChatApplicationCommand {
					continue
				}
				checkDescriptions(t, "/"+command.Name, command.Description,
					command.DescriptionLocalizations, command.Options)
			}
		})
	}
}

func checkDescriptions(t *testing.T, path, description string,
	localizations *map[discordgo.Locale]string, options []*discordgo.ApplicationCommandOption) {
	t.Helper()
	checkDescription(t, path, "default", description)
	if localizations != nil {
		for locale, localized := range *localizations {
			checkDescription(t, path, string(locale), localized)
		}
	}

	for _, option := range options {
		optionPath := fmt.Sprintf("%v %v", path, option.Name)
		checkDescriptions(t, optionPath, option.Description,
			&option.DescriptionLocalizations, option.Options)
		if option.Type == discordgo.ApplicationCommandOptionInteger && option.MaxValue != 0 {
			checkMaxValue(t, optionPath, option)
		}
	}
}

func checkDescription(t *testing.T, path, locale, description string) {
	t.Helper()
	length := utf8.RuneCountInString(description)
	if length < descriptionMinLength || length > descriptionMaxLength {
		t.Errorf("%v (%v) description is %v characters long, want %v to %v: %q",
			path, locale, length, descriptionMinLength, descriptionMaxLength, description)
	}

	if strings.Contains(description, "<no value>") {
		t.Errorf("%v (%v) description has an unfilled placeholder: %q", path, locale, description)
	}
}

// The maximum differs between games, so a description stating it must follow the option's.
func checkMaxValue(t *testing.T, path string, option *discordgo.ApplicationCommandOption) {
	t.Helper()
	maxValue := strconv.Itoa(int(option.MaxValue))
	descriptions := map[string]string{"default": option.Description}
	for locale, localized := range option.DescriptionLocalizations {
		descriptions[string(locale)] = localized
	}

	for locale, description := range descriptions {
		if !strings.Contains(description, maxValue) {
			t.Errorf("%v (%v) description does not state the maximum %v: %q",
				path, locale, maxValue, description)
		}
	}
}
