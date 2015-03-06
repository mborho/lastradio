package lastradio

import (
	spotify "github.com/op/go-libspotify/spotify"
	"github.com/shkh/lastfm-go/lastfm"
	"gopkg.in/qml.v1"
	"log"
	"time"
)

var (
	control chan string
)

func init() {
	control = make(chan string)
}

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

func (p *Player) LoadRadio(name string, username string) error {
	p.spotifyQueue = make(chan *LastFmTrack, 3)
	p.playQueue = make(chan *LastFmTrack, 3)
	p.lastfmTracks = make(chan *LastFmTrack, 3)

	go getTrackInfo(p.Lastfm, p.LastFmUser.Username, p.lastfmTracks, p.spotifyQueue)
	go getSpotifyData(p.Spotify, p.spotifyQueue, p.playQueue, control)
	go p.Controller()
	err := p.setRadio(name, username)

	if err != nil {
		log.Print(err)
	}
	return err
}

func (p *Player) setRadio(mode string, username string) error {
	if username == "" {
		username = p.LastFmUser.Username
	}
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
			spotify:         p.Spotify,
			lastfm:          p.Lastfm,
			lastfmTracks:    p.lastfmTracks,
			lastFmUser:      p.LastFmUser,
			currentUsername: username,
		}
	case "topartists":
		p.Radio = &TopArtistsRadio{
			spotify:         p.Spotify,
			lastfm:          p.Lastfm,
			lastfmTracks:    p.lastfmTracks,
			lastFmUser:      p.LastFmUser,
			currentUsername: username,
		}
	case "loved":
		p.Radio = &LovedTracksRadio{
			spotify:         p.Spotify,
			lastfm:          p.Lastfm,
			lastfmTracks:    p.lastfmTracks,
			lastFmUser:      p.LastFmUser,
			currentUsername: username,
		}
	case "similar":
		p.Radio = &SimilarRadio{
			spotify:      p.Spotify,
			lastfm:       p.Lastfm,
			lastfmTracks: p.lastfmTracks,
			lastFmUser:   p.LastFmUser,
			bandName:     username,
		}
	case "tag":
		p.Radio = &TagTracksRadio{
			spotify:      p.Spotify,
			lastfm:       p.Lastfm,
			lastfmTracks: p.lastfmTracks,
			lastFmUser:   p.LastFmUser,
			tagName:      username,
		}
	}
	err := p.Radio.Load()
	return err
}

func (p *Player) Controller() {
	started := make(chan *LastFmTrack)
	endTrackTime := time.Now()
	endOfTrack := p.Spotify.EndOfTrackUpdates()
	streamingErrors := p.Spotify.StreamingErrors()
	//logMessages := p.Spotify.LogMessages()
	for {
		select {
		case command := <-control:
			log.Print("CONTROL: ", command)
			if command == "next" {
				p.startNextTrack(started)
			}
		case <-endOfTrack:
			log.Print("SIGNAL: endOfTrack")
			duration := time.Since(endTrackTime)
			if duration.Seconds() > 5 {
				// handle multiple endOfTrack signals with a threshold of 5 secs
				p.startNextTrack(started)
			} else {
				log.Print("EndOfTrack threshold not reached")
			}
			endTrackTime = time.Now()
		case track := <-started:
			p.Public.SetData(track)
		case err := <-streamingErrors:
			log.Print("SIGNAL: streaming error")
			log.Print(err)
			/*case msg := <-logMessages:
			  log.Print("LOG MESSAGE: ", msg)*/
		}
	}
}

func (p *Player) startNextTrack(started chan *LastFmTrack) {
	go p.playSpotifyTrack(started)
}

func (p *Player) playSpotifyTrack(started chan *LastFmTrack) { //nextTrack *LastFmTrack, started chan *LastFmTrack) {
	// Parse the track
	nextTrack := <-p.playQueue
	log.Print("Loading Track: ", nextTrack.Artist.Name, nextTrack.Name)
	link, err := p.Spotify.ParseLink(nextTrack.SpotifyLink)
	if err != nil {
		log.Fatal(err)
	}
	track, err := link.Track()
	if err != nil {
		log.Fatal(err)
	}

	// Load the track and play it
	track.Wait()
	player := p.Spotify.Player()
	if err := player.Load(track); err != nil {
		log.Fatal(err)
	}
	player.Seek(1000000)
	player.Play()
	started <- nextTrack
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
		control <- "next"
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
