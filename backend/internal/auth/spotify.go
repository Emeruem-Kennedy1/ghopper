package auth

import (
	"fmt"
	"net/http"

	"github.com/Emeruem-Kennedy1/ghopper/config"
	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/pkg/utils"
	"github.com/zmb3/spotify"
)

// Make the function a variable so it can be swapped in tests
var CreateOrUpdateUserFromSpotifyDataFunc = CreateOrUpdateUserFromSpotifyDataImpl

type SpotifyAuth struct {
	authenticator AuthenticatorInterface
	state         string
	config        *config.Config
}

func NewSpotifyAuth(cfg *config.Config) (*SpotifyAuth, error) {
	auth := spotify.NewAuthenticator(
		cfg.SpotifyRedirectURI,
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserReadEmail,
		spotify.ScopeUserTopRead,
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopePlaylistModifyPublic,
	)
	auth.SetAuthInfo(cfg.SpotifyClientID, cfg.SpotifyClientSecret)

	state, err := utils.GenerateRandomString(10)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate random string: %v", err)
	}

	return &SpotifyAuth{
		authenticator: &auth,
		state:         state,
		config:        cfg,
	}, nil
}

func (sa *SpotifyAuth) AuthURL() string {
	return sa.authenticator.AuthURL(sa.state)
}

func (sa *SpotifyAuth) CallBack(r *http.Request) (*spotify.Client, error) {
	if st := r.FormValue("state"); st != sa.state {
		return nil, fmt.Errorf("state mismatch: %s != %s", st, sa.state)
	}

	tok, err := sa.authenticator.Token(sa.state, r)
	if err != nil {
		return nil, fmt.Errorf("couldn't get token: %v", err)
	}

	client := sa.authenticator.NewClient(tok)
	return &client, nil
}

func (sa *SpotifyAuth) GetUserInfo(client *spotify.Client) (*spotify.PrivateUser, error) {
	user, err := client.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("couldn't get user: %v", err)
	}

	return user, nil
}

func CreateOrUpdateUserFromSpotifyDataImpl(userRepo repository.UserRepositoryInterface, spotifyUser spotify.PrivateUser) (*models.User, string, error) {
	user := &models.User{
		ID:          spotifyUser.ID,
		DisplayName: spotifyUser.DisplayName,
		Email:       spotifyUser.Email,
		Country:     spotifyUser.Country,
		SpotifyURI:  string(spotifyUser.URI),
	}

	if len(spotifyUser.Images) > 0 {
		user.ProfileImage = spotifyUser.Images[0].URL
	}

	err := userRepo.UpsertUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create or update user: %v", err)
	}

	token, err := GenerateTokenFunc(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return user, token, nil

}

func CreateOrUpdateUserFromSpotifyData(userRepo repository.UserRepositoryInterface, spotifyUser spotify.PrivateUser) (*models.User, string, error) {
	return CreateOrUpdateUserFromSpotifyDataFunc(userRepo, spotifyUser)
}

func (sa *SpotifyAuth) GetAuthenticator() AuthenticatorInterface {
	return sa.authenticator
}
