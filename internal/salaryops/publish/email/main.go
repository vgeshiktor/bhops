package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	MS_GRAPH_BASE_URL = "https://graph.microsoft.com/v1.0"
)

func getAccessToken(
	appId string, secret string, scopes []string) (string, error) {

	// Create a confidential client application
	cred, err := confidential.
		NewCredFromSecret(secret)
	if err != nil {
		log.Error().Msgf("Error creating credential: %v", err)
	}

	authority:= "https://login.microsoftonline.com/consumers"
	clientApp, err := confidential.
		New(authority, appId, cred)
	if err != nil {
		log.Error().Msgf("Error creating confidential client application: %v", err)
	}

	   // Acquire a token
    ctx := context.Background()
    result, err := clientApp.AcquireTokenSilent(
		ctx, scopes)
    if err != nil {
        result, err = clientApp.AcquireTokenByCredential(ctx, scopes)
        if err != nil {
            log.Error().Msgf("Failed to acquire token: %v", err)
        }
    }

	return result.AccessToken, err
}

func main() {
	// load .env file
	// .env.local takes precedence (if present)
	err := godotenv.Load()
	if err != nil {
		log.Error().Msg("Error loading .env")
	}

	applicationId := os.Getenv("APPLICATION_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	// scopes := []string{"User.Read", "Mail.ReadWrite"}
	scopes := []string{"https://graph.microsoft.com/.default"}

	endpoint := MS_GRAPH_BASE_URL + "/me/messages"

	accessToken, err := getAccessToken(
		applicationId, clientSecret, scopes)

	if err != nil {
		log.Error().Msgf("Error getting access token: %v", err)
	}

	log.Info().Msgf("Access token: %v", accessToken)

	// Create the authorization header
	authHeader := "Bearer " + accessToken

	// Example usage in an HTTP request
	req, err := http.NewRequest(
		"GET", endpoint, nil)
	if err != nil {
		log.Error().Msgf("Error creating HTTP request: %v", err)
	}

	req.Header.Set("Authorization", authHeader)

	// Log the authorization header for debugging purposes
	log.Info().Msgf("Authorization Header: %v", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msgf("Error executing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("Received non-OK HTTP status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("Error reading response body: %v", err)
	}

	log.Info().Msgf("Response JSON: %v", string(body))
}
