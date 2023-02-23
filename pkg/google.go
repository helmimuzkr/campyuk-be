package pkg

import (
	"campyuk-api/config"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Google Config
func NewGoogleConf(cfg *config.AppConfig) *oauth2.Config {
	conf := &oauth2.Config{
		RedirectURL:  cfg.GOOGLE_REDIRECT_CALLBACK,
		ClientID:     cfg.GOOGLE_CLIENT_ID,
		ClientSecret: cfg.GOOGLE_CLIENT_SECRET,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/calendar",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return conf
}

// Google API
type googleAPI struct {
	conf *oauth2.Config
}

func NewGoogleAPI(conf *oauth2.Config) *googleAPI {
	googleApi := &googleAPI{
		conf: conf,
	}

	return googleApi
}

func (g *googleAPI) GetEmail(accessToken string) (string, error) {
	// Create request to get user info.
	request, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)
	if err != nil {
		return "", err
	}
	bearer := "Bearer " + accessToken
	request.Header.Set("Authorization", bearer)
	// Hook the response.
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", errors.New("could not retrieve user")
	}
	// Decode body response to map.
	var resBody map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&resBody); err != nil {
		return "", err
	}

	email := resBody["email"].(string)

	return email, nil
}

func (g *googleAPI) CreateEvent(detailEvent map[string]string) (string, error) {
	ctx := context.TODO()

	// Get client.
	client, err := g.getClient(ctx)
	if err != nil {
		return "", err
	}

	// Create calendar service.
	calendarSrv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", err
	}

	// Setup calendar request.
	startTime, err := time.Parse("2006-01-02", detailEvent["start"])
	if err != nil {
		return "", errors.New("error parsing time in create event service")
	}
	endTime, err := time.Parse("2006-01-02", detailEvent["end"])
	if err != nil {
		return "", errors.New("error parsing time in create event service")
	}
	detailEvent["start"] = startTime.Format(time.RFC3339)
	detailEvent["end"] = endTime.Format(time.RFC3339)

	event := &calendar.Event{
		Summary:  detailEvent["summary"],
		Location: detailEvent["location"],
		Start: &calendar.EventDateTime{
			DateTime: detailEvent["start"], // Format time must be RFC3339.
			TimeZone: "Asia/Jakarta",
		},
		End: &calendar.EventDateTime{
			// DateTime: time.Date(2023, 02, 10, 13, 20, 10, 10, time.Local).Format(time.RFC3339), // it means YYYY-MM-DD HH-MM-SS-NS Location.
			DateTime: detailEvent["end"],
			TimeZone: "Asia/Jakarta",
		},
		Attendees: []*calendar.EventAttendee{{Email: detailEvent["email"]}},
	}

	// Execute calendar insert events.
	event, err = calendarSrv.Events.Insert("primary", event).SendUpdates("all").Do()
	if err != nil {
		return "", err
	}

	return event.HtmlLink, nil
}

func (g *googleAPI) getClient(ctx context.Context) (*http.Client, error) {
	// Read a token from a local file.
	f, err := os.Open(config.TokenPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Stream decoder, decode token json to struct oauth token.
	var token *oauth2.Token
	if err := json.NewDecoder(f).Decode(&token); err != nil {
		return nil, err
	}

	// You can create google client without token access,
	// only using refresh token you can make the client.
	// But, in this case i want to create a new proper token only using the refresh token with TokenSource method.

	// You can skip this create and update token step.
	// Create a new token using the refresh token.
	tokenSource := g.conf.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	// If the new access token different from the old token,
	// update the token in local storage with a new token.
	if newToken.AccessToken != token.AccessToken {
		newFile, err := os.Create(config.TokenPath)
		if err != nil {
			return nil, err
		}
		if err := json.NewEncoder(newFile).Encode(newToken); err != nil {
			return nil, err
		}
		log.Println("Token has been updated")
	}

	// Create the client.
	client := g.conf.Client(ctx, token)

	return client, nil
}
