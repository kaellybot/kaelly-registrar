package i18n

import (
	"embed"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
)

type Language struct {
	Locale          discordgo.Locale
	AMQPLocale      amqp.Language
	TranslationFile string
}

const (
	frenchFile     = "fr.json"
	englishFile    = "en.json"
	spanishFile    = "es.json"
	germanFile     = "de.json"
	portugueseFile = "pt.json"

	DefaultLocale = discordgo.EnglishGB
)

//go:embed *.json
var Folder embed.FS

func GetLanguages() []Language {
	return []Language{
		{
			Locale:          discordgo.French,
			TranslationFile: frenchFile,
			AMQPLocale:      amqp.Language_FR,
		},
		{
			Locale:          discordgo.EnglishGB,
			TranslationFile: englishFile,
			AMQPLocale:      amqp.Language_EN,
		},
		{
			Locale:          discordgo.EnglishUS,
			TranslationFile: englishFile,
			AMQPLocale:      amqp.Language_EN,
		},
		{
			Locale:          discordgo.SpanishES,
			TranslationFile: spanishFile,
			AMQPLocale:      amqp.Language_ES,
		},
		{
			Locale:          discordgo.German,
			TranslationFile: germanFile,
			AMQPLocale:      amqp.Language_DE,
		},
		{
			Locale:          discordgo.PortugueseBR,
			TranslationFile: portugueseFile,
			AMQPLocale:      amqp.Language_PT,
		},
	}
}

func MapDiscordLocale(locale discordgo.Locale) amqp.Language {
	for _, language := range GetLanguages() {
		if language.Locale == locale {
			return language.AMQPLocale
		}
	}

	return amqp.Language_ANY
}

func MapAMQPLocale(locale amqp.Language) discordgo.Locale {
	for _, language := range GetLanguages() {
		if language.AMQPLocale == locale {
			return language.Locale
		}
	}

	return DefaultLocale
}
