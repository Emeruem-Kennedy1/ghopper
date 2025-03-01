package auth

import (
	"net/http"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// SpotifyAuthInterface defines the methods we use from SpotifyAuth
type SpotifyAuthInterface interface {
	AuthURL() string
	CallBack(r *http.Request) (*spotify.Client, error)
	GetUserInfo(client *spotify.Client) (*spotify.PrivateUser, error)
	GetAuthenticator() AuthenticatorInterface
}

type AuthenticatorInterface interface {
	AuthURL(state string) string
	Token(state string, r *http.Request) (*oauth2.Token, error)
	NewClient(token *oauth2.Token) spotify.Client
	SetAuthInfo(clientID, secretKey string)
}


// Ensure the SpotifyAuth and spotify.Authenticator implement our interfaces
var _ SpotifyAuthInterface = (*SpotifyAuth)(nil)
var _ AuthenticatorInterface = (*spotify.Authenticator)(nil)
