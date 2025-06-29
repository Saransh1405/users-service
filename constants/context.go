// In constants/context.go
package constants

import "context"

type contextKey string

const (
	ClientNameKey contextKey = "user_id"
	LanguageKey   contextKey = "language"
)

// Helper functions for type-safe context operations
func WithClientName(ctx context.Context, clientName string) context.Context {
	return context.WithValue(ctx, ClientNameKey, clientName)
}

func GetClientName(ctx context.Context) (string, bool) {
	clientName, ok := ctx.Value(ClientNameKey).(string)
	return clientName, ok
}

func WithLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, LanguageKey, language)
}

func GetLanguage(ctx context.Context) (string, bool) {
	language, ok := ctx.Value(LanguageKey).(string)
	return language, ok
}
