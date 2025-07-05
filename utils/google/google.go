package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"users-service/constants"
	"users-service/logger"
	"users-service/models"
	"users-service/utils/configs"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthService struct {
	config *oauth2.Config
	logger logger.Logger
}

func NewGoogleOAuthService() (*GoogleOAuthService, error) {
	log := logger.GetLoggerWithoutContext()

	// Get application config
	applicationConfig, err := configs.Get(constants.ApplicationConfig)
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to get application config")
		return nil, err
	}

	// Get Google OAuth credentials from config
	clientID := applicationConfig.GetString("google.client_id")
	clientSecret := applicationConfig.GetString("google.client_secret")
	redirectURL := applicationConfig.GetString("google.redirect_url")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("Google OAuth credentials not configured")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthService{
		config: config,
		logger: log,
	}, nil
}

// GetAuthURL returns the Google OAuth authorization URL
func (g *GoogleOAuthService) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCodeForToken exchanges authorization code for access token
func (g *GoogleOAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		g.logger.With(zap.Error(err)).Error("Failed to exchange code for token")
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}
	return token, nil
}

// GetUserInfo retrieves user information from Google
func (g *GoogleOAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*models.GoogleUserInfo, error) {
	// Make HTTP request to Google's userinfo endpoint
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		g.logger.With(zap.Error(err)).Error("Failed to create HTTP request")
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		g.logger.With(zap.Error(err)).Error("Failed to get user info from Google")
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		g.logger.With(zap.Error(err)).Error("Failed to decode user info response")
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Convert to our model
	googleUserInfo := &models.GoogleUserInfo{
		ID:            userInfo.ID,
		Email:         userInfo.Email,
		VerifiedEmail: userInfo.VerifiedEmail,
		Name:          userInfo.Name,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Picture:       userInfo.Picture,
		Locale:        userInfo.Locale,
	}

	return googleUserInfo, nil
}

// ValidateToken validates the Google access token
func (g *GoogleOAuthService) ValidateToken(ctx context.Context, token *oauth2.Token) (bool, error) {
	// Try to get user info to validate token
	_, err := g.GetUserInfo(ctx, token)
	if err != nil {
		return false, nil // Token is invalid
	}

	return true, nil
}

// RefreshToken refreshes the access token using refresh token
func (g *GoogleOAuthService) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := g.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		g.logger.With(zap.Error(err)).Error("Failed to refresh token")
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}
