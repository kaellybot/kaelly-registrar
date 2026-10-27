package constants

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// Discord Client ID.
	ClientID = "CLIENT_ID"

	// Discord Bot Token.
	Token = "TOKEN"

	// Game whose commands are registered: 1 (DOFUS), 2 (DOFUS Touch) or 3 (DOFUS Retro).
	Game = "GAME"

	// Metric port.
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic].
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"

	// Default values.
	defaultClientID   = ""
	defaultToken      = ""
	defaultGame       = 0 // No default game: each run must name the one it registers.
	defaultShardID    = 0
	defaultShardCount = 1
	defaultMetricPort = 2112
	defaultLogLevel   = zerolog.InfoLevel
	defaultProduction = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		ClientID:   defaultClientID,
		Token:      defaultToken,
		Game:       defaultGame,
		MetricPort: defaultMetricPort,
		LogLevel:   defaultLogLevel.String(),
		Production: defaultProduction,
	}
}
