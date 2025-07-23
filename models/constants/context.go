package constants

type ContextKey int

const (
	ContextKeyChannel ContextKey = iota
	ContextKeyCity
	ContextKeyEnabled
	ContextKeyFeed
	ContextKeyJob
	ContextKeyLanguage
	ContextKeyLevel
	ContextKeyOrder
	ContextKeyServer
)
