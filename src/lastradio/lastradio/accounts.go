package lastradio

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"code.google.com/p/portaudio-go/portaudio"
	spotify "github.com/op/go-libspotify/spotify"
	"github.com/shkh/lastfm-go/lastfm"
)

var (
	LOCAL_APP_DIR = os.Getenv("HOME") + "/.local/share/lastradio"
)

type LastFmUser struct {
	Username  string
	Password  string
	ApiKey    string
	ApiSecret string
}

func (u *LastFmUser) GetUsername() string {
	return u.Username
}

type SpotifyUser struct {
	AppKeyPath string
	Username   string
	Password   string
	Remember   bool
	Debug      bool
}

func LoginLastFm(user *LastFmUser) (api *lastfm.Api, err error) {
	api = lastfm.New(user.ApiKey, user.ApiSecret)
	err = api.Login(user.Username, user.Password)
	return api, err
}

func loginToLastFm(lastFmUser *LastFmUser) (*lastfm.Api, error) {
	// login and get api client
	lastFmApi, err := LoginLastFm(lastFmUser)
	if err != nil {
		return nil, err
	}
	return lastFmApi, err
}

func getSpotifySession(spotifyUser *SpotifyUser) (*spotify.Session, error) {
	appKey, err := ioutil.ReadFile(spotifyUser.AppKeyPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Printf("libspotify %s", spotify.BuildId())

	session, err := spotify.NewSession(&spotify.Config{
		ApplicationKey:   appKey,
		ApplicationName:  "test",
		CacheLocation:    LOCAL_APP_DIR,
		SettingsLocation: LOCAL_APP_DIR,

		// Disable playlists to make playback faster
		DisablePlaylistMetadataCache: true,
		InitiallyUnloadPlaylists:     true,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return session, err
}

func loginToSpotifySession(session *spotify.Session, spotifyUser *SpotifyUser) error {
	var err error
	if len(spotifyUser.Password) > 0 {
		credentials := spotify.Credentials{
			Username: spotifyUser.Username,
			Password: spotifyUser.Password,
		}
		if err = session.Login(credentials, spotifyUser.Remember); err != nil {
			log.Fatal(err)
			return err
		}
	} else {
		if err = session.Relogin(); err != nil {
			log.Fatal(err)
			return err
		}
	}

	// Wait for login and expect it to go fine
	select {
	case err = <-session.LoginUpdates():
		if err != nil {
			log.Print(err)
			return err
		}
		return nil
	}
}

func HandleSpotifySession(session *spotify.Session, exit <-chan bool) {
	exitAttempts := 0
	running := true
	log.Print("TODO: spotify signal handling")
	for running {
		select {
		case <-session.LogoutUpdates():
			running = false
		case <-exit:
			if exitAttempts >= 3 {
				os.Exit(42)
			}
			exitAttempts++
			session.Logout()
		case <-time.After(5 * time.Second):
		}
	}
	session.Close()
	os.Exit(32)
}

func Close() {
	portaudio.Terminate()
}

func init() {
	createLocalDir(LOCAL_APP_DIR)
}

func createLocalDir(appDir string) {
	exists, err := localDirExists(appDir)
	if err != nil {
		log.Fatal(err)
	}
	if exists == false {
		log.Print("Create local dir")
		err := os.Mkdir(appDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func localDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
