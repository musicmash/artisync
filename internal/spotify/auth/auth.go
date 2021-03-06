package auth

import (
	"context"
	"fmt"

	"github.com/musicmash/artisync/internal/config"
	"golang.org/x/oauth2/clientcredentials"
)

func ValidateConfig(ctx context.Context, spotify *config.Spotify) error {
	conf := &clientcredentials.Config{
		ClientID:     spotify.ClientID,
		ClientSecret: spotify.ClientSecret,
		TokenURL:     spotify.TokenURL,
	}

	_, err := conf.Token(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get token: %w", err)
	}

	return nil
}
