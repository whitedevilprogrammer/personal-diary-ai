package services

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// var GoogleOAuthConfig = &oauth2.Config{
// 	RedirectURL:  "http://localhost:8080/auth/google/callback",
// 	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
// 	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
// 	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
// 	Endpoint:     google.Endpoint,
// }

func GetGoogleOAuthConfig() *oauth2.Config {
    return &oauth2.Config{
        RedirectURL:  "http://localhost:8080/auth/google/callback",
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        Scopes:       []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
        Endpoint: google.Endpoint,
    }
}


