package helper

import (
	"campyuk-api/config"
	"context"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleAPI interface {
	GetUrlAuth(state string) string
	GetToken(code string) (*oauth2.Token, error)
	CreateCalendar(token *oauth2.Token, detail CalendarDetail) error
}

type CalendarDetail struct {
	Summay   string
	Location string
	Start    string
	End      string
	Emails   []string
}

type googleAPI struct {
	conf *oauth2.Config
}

func NewOauth(cfg *config.AppConfig) GoogleAPI {
	conf := &oauth2.Config{
		RedirectURL:  cfg.GOOGLE_REDIRECT_CALLBACK,
		ClientID:     cfg.GOOGLE_CLIENT_ID,
		ClientSecret: cfg.GOOGLE_CLIENT_SECRET,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/calendar"},
		Endpoint:     google.Endpoint,
	}

	return &googleAPI{conf: conf}
}

func (ga *googleAPI) GetUrlAuth(state string) string {
	return ga.conf.AuthCodeURL(state)
}

func (ga *googleAPI) GetToken(code string) (*oauth2.Token, error) {
	token, err := ga.conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (ga *googleAPI) CreateCalendar(token *oauth2.Token, detail CalendarDetail) error {
	calendarService, err := calendar.NewService(oauth2.NoContext, option.WithTokenSource(ga.conf.TokenSource(oauth2.NoContext, token)))
	if err != nil {
		return err
	}

	attendees := []*calendar.EventAttendee{}
	for _, e := range detail.Emails {
		a := &calendar.EventAttendee{Email: e}
		attendees = append(attendees, a)
	}

	event := &calendar.Event{
		Summary:  detail.Summay,
		Location: detail.Location,
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Format(time.RFC3339), // Wajib format RFC3339
			TimeZone: "Asia/Jakarta",
		},
		End: &calendar.EventDateTime{
			// DateTime: time.Date(2023, 02, 10, 13, 20, 10, 10, time.Local).Format(time.RFC3339), // artinya YYYY-MM-DD HH-MM-SS-NS Location
			DateTime: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			TimeZone: "Asia/Jakarta",
		},
		Attendees: attendees,
	}

	event, err = calendarService.Events.Insert("primary", event).SendUpdates("all").Do()
	if err != nil {
		return err
	}

	return nil
}
