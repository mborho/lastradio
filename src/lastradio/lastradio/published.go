package lastradio

import (
	"log"

	spotify "github.com/op/go-libspotify/spotify"
	"github.com/shkh/lastfm-go/lastfm"
	"gopkg.in/qml.v0"
)

type Player struct {
	paused       bool
	LoggedIn     bool
	loginLastFm  bool
	loginSpotify bool
	Radio        Radio
	Public       *PublishedData
	Spotify      *spotify.Session
	Lastfm       *lastfm.Api
	LastFmUser   *LastFmUser
	control      chan string
	lastfmTracks chan *LastFmTrack
	playQueue    chan *LastFmTrack
	spotifyQueue chan *LastFmTrack
	Exit         chan bool
}

func (p *Player) LoginToLastFm(user, password string) bool {
	lastFmUser := &LastFmUser{
		Username:  user,
		Password:  password,
		ApiKey:    "7efd4b4e0d365c1b2867b055abaf2ac3",
		ApiSecret: "46852d04b0c882f6aa905b3af719dc83",
	}
	api, err := loginToLastFm(lastFmUser)
	if err != nil {
		log.Print(err)
		return false
	}
	p.Lastfm = api
	p.LastFmUser = lastFmUser
	p.loginLastFm = true
	p.checkLoggedIn()
	return true
}

func (p *Player) LoginToSpotify(user, password string) bool {

	spotifyUser := &SpotifyUser{
		AppKeyPath: "share/lastradio/spotify_appkey.key",
		Username:   user,
		Password:   password,
		Remember:   false,
	}

	if p.Spotify == nil {
		session, sErr := getSpotifySession(spotifyUser)
		if sErr != nil {
			log.Print(sErr)
			return false
		}
		p.Spotify = session
	}

	err := loginToSpotifySession(p.Spotify, spotifyUser)
	if err != nil {
		log.Print(err)
		return false
	}
	p.loginSpotify = true
	go HandleSpotifySession(p.Spotify, p.Exit)
	p.checkLoggedIn()
	return true
}

func (p *Player) checkLoggedIn() {
	if p.loginSpotify == true && p.loginLastFm == true {
		p.LoggedIn = true
		qml.Changed(p, &p.LoggedIn)
	}
}

func (p *Player) Logout() {
	p.loginSpotify = false
	p.loginLastFm = false
	p.LoggedIn = false
	qml.Changed(p, &p.LoggedIn)
	p.checkLoggedIn()
}

func (p *Player) LoadRadio(name string) error {
	p.control = make(chan string)
	p.spotifyQueue = make(chan *LastFmTrack, 3)
	p.playQueue = make(chan *LastFmTrack, 3)
	p.lastfmTracks = make(chan *LastFmTrack, 3)
	_ = initializeAudioConsumer(p.Spotify)
	go getTrackInfo(p.Lastfm, p.LastFmUser.Username, p.lastfmTracks, p.spotifyQueue)
	go getSpotifyData(p.Spotify, p.spotifyQueue, p.playQueue, p.control)
	go p.Controller()
	err := p.setRadio(name)
	if err != nil {
		log.Print(err)
	}
	return err
}

func (p *Player) setRadio(mode string) error {
	switch mode {
	case "recommended":
		p.Radio = &RecommendedRadio{
			spotify:      p.Spotify,
			lastfm:       p.Lastfm,
			lastfmTracks: p.lastfmTracks,
			lastFmUser:   p.LastFmUser,
		}
	case "top":
		p.Radio = &TopTracksRadio{
			spotify:      p.Spotify,
			lastfm:       p.Lastfm,
			lastfmTracks: p.lastfmTracks,
			lastFmUser:   p.LastFmUser,
		}
	case "topartists":
		p.Radio = &TopArtistsRadio{
			spotify:      p.Spotify,
			lastfm:       p.Lastfm,
			lastfmTracks: p.lastfmTracks,
			lastFmUser:   p.LastFmUser,
		}
	case "loved":
		p.Radio = &LovedTracksRadio{
			spotify:      p.Spotify,
			lastfm:       p.Lastfm,
			lastfmTracks: p.lastfmTracks,
			lastFmUser:   p.LastFmUser,
		}
	}
	err := p.Radio.Load()
	return err
}

func (p *Player) Controller() {
	endOfTrack := p.Spotify.EndOfTrackUpdates()
	for {
		select {
		case command := <-p.control:
			if command == "next" {
				p.startNextTrack()
			}
		case <-endOfTrack:
			p.startNextTrack()
		}
	}
}

func (p *Player) startNextTrack() {
	track := <-p.playQueue
	p.Public.SetData(track)
	go playSpotifyTrack(p.Spotify, track.SpotifyLink)
}

func (p *Player) Pause() {
	p.paused = true
	p.Spotify.Player().Pause()
}

func (p *Player) Skip() {
	p.Play()
}

func (p *Player) Stop() {
	p.Spotify.Player().Unload()
}

func (p *Player) Play() {
	if p.paused {
		p.Spotify.Player().Play()
		p.paused = false
		return
	}
	go func() {
		p.control <- "next"
	}()
}

func (p *Player) SendToLastfm(mode string) {
	log.Print("To Last.fm: " + mode)
	if mode == "nowplaying" {
		go SendNowPlaying(p.Lastfm, p.Public.Artist, p.Public.Name, p.Public.Album)
	} else if mode == "scrobble" {
		go Scrobble(p.Lastfm, p.Public.Artist, p.Public.Name, p.Public.Album)
	} else if mode == "love" {
		go Love(p.Lastfm, p.Public.Artist, p.Public.Name)
	} else if mode == "ban" {
		go Ban(p.Lastfm, p.Public.Artist, p.Public.Name)
	} else if mode == "unlove" {
		go Unlove(p.Lastfm, p.Public.Artist, p.Public.Name)
	}
}

type PublishedData struct {
	Name         string
	Artist       string
	Album        string
	Year         int
	Image        string
	IsLoved      bool
	Duration     float64
	ScrobbleAt   int
	NowPlayingAt int
}

func (p *PublishedData) SetData(track *LastFmTrack) {
	p.Name = track.Name
	p.Artist = track.Artist.Name
	p.Album = track.Album
	p.Year = track.Year
	p.IsLoved = track.IsLoved
	// set title
	if track.Image != "" {
		p.Image = track.Image
	} else if track.Artist.Image != "" {
		p.Image = track.Artist.Image
	} else {
		p.Image = "png/dummy.png"
	}
	// set times when to scrobble etc
	p.Duration = track.Duration.Seconds()
	p.NowPlayingAt, p.ScrobbleAt = getTrackScrobbleTimes(p.Duration)
	// signal changes to Qml side of things
	qml.Changed(p, &p.Name)
	qml.Changed(p, &p.Artist)
	qml.Changed(p, &p.Album)
	qml.Changed(p, &p.Year)
	qml.Changed(p, &p.Image)
	qml.Changed(p, &p.Duration)
	qml.Changed(p, &p.IsLoved)
	qml.Changed(p, &p.NowPlayingAt)
	qml.Changed(p, &p.ScrobbleAt)
}

func getTrackScrobbleTimes(duration float64) (int, int) {
	scrobble_time := duration / 2
	if scrobble_time > 240 {
		scrobble_time = 240
	} else if scrobble_time < 30 {
		scrobble_time = 30
	}
	return 30, int(scrobble_time)
}
