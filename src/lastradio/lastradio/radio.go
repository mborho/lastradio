package lastradio

import (
	"errors"
	_ "log"

	spotify "github.com/op/go-libspotify/spotify"
	"github.com/shkh/lastfm-go/lastfm"
)

type Radio interface {
	Load() error
	Run()
}

type RecommendedRadio struct {
	spotify      *spotify.Session
	lastfm       *lastfm.Api
	artists      []*LastFmArtist
	artistslist  chan *LastFmArtist
	lastfmTracks chan *LastFmTrack
	lastFmUser   *LastFmUser
	page         int
}

func (radio *RecommendedRadio) Load() error {
	radio.artistslist = make(chan *LastFmArtist)
	go getLastFmTrack(radio.lastfm, radio.lastfmTracks, radio.artistslist)
	err := radio.loadArtists()
	if err == nil {
		go radio.Run()
	}
	return err
}

func (radio *RecommendedRadio) loadArtists() error {
	radio.page = radio.page + 1
	var params = lastfm.P{"limit": 50, "page": radio.page}
	artists, err := radio.lastfm.User.GetRecommendedArtists(params)
	if err == nil {
		if artists.Total < 1 {
			return errors.New("No artists found!")
		}
		for _, artist := range artists.Artists {
			lastFmArtist := &LastFmArtist{
				Name:  artist.Name,
				Image: artist.Images[2].Url,
			}
			radio.artists = append(radio.artists, lastFmArtist)
		}

	}
	return err
}

func (radio *RecommendedRadio) Run() {
	for {
		var next *LastFmArtist
		if len(radio.artists) > 0 {
			next, radio.artists = getRandomArtist(radio.artists)
			radio.artistslist <- next
			if len(radio.artists) <= 49 {
				go radio.loadArtists()
			}
		}
	}
}

type TopArtistsRadio struct {
	spotify         *spotify.Session
	lastfm          *lastfm.Api
	artists         []*LastFmArtist
	artistslist     chan *LastFmArtist
	lastfmTracks    chan *LastFmTrack
	lastFmUser      *LastFmUser
	page            int
	currentUsername string
}

func (radio *TopArtistsRadio) Load() error {
	radio.artistslist = make(chan *LastFmArtist)
	go getLastFmTrack(radio.lastfm, radio.lastfmTracks, radio.artistslist)
	err := radio.loadArtists()
	if err == nil {
		go radio.Run()
	}
	return err
}

func (radio *TopArtistsRadio) loadArtists() error {
	radio.page = radio.page + 1
	var params = lastfm.P{"user": radio.currentUsername, "limit": 50, "page": radio.page}
	artists, err := radio.lastfm.User.GetTopArtists(params)
	if err == nil {
		if artists.Total < 1 {
			return errors.New("No artists found!")
		}
		for _, artist := range artists.Artists {
			lastFmArtist := &LastFmArtist{
				Name:  artist.Name,
				Image: artist.Images[2].Url,
			}
			radio.artists = append(radio.artists, lastFmArtist)
		}
	}
	return err
}

func (radio *TopArtistsRadio) Run() {
	for {
		var next *LastFmArtist
		if len(radio.artists) > 0 {
			next, radio.artists = getRandomArtist(radio.artists)
			radio.artistslist <- next
			if len(radio.artists) <= 49 {
				go radio.loadArtists()
			}
		}
	}
}

type TopTracksRadio struct {
	spotify         *spotify.Session
	lastfm          *lastfm.Api
	tracks          []*LastFmTrack
	lastfmTracks    chan *LastFmTrack
	lastFmUser      *LastFmUser
	page            int
	currentUsername string
}

func (radio *TopTracksRadio) Load() error {
	err := radio.loadTopTracks()
	if err == nil {
		go radio.Run()
	}
	return err
}

func (radio *TopTracksRadio) loadTopTracks() error {
	radio.page = radio.page + 1
	params := lastfm.P{"user": radio.currentUsername, "limit": 50, "page": radio.page}
	tracks, err := radio.lastfm.User.GetTopTracks(params)
	if err == nil {
		if tracks.Total < 1 {
			return errors.New("No tracks found!")
		}
		for _, track := range tracks.Tracks {
			lastFmArtist := &LastFmArtist{
				Name: track.Artist.Name,
			}
			lastFmTrack := &LastFmTrack{
				Artist: lastFmArtist,
				Name:   track.Name,
			}
			radio.tracks = append(radio.tracks, lastFmTrack)
		}
	}
	return err
}

func (radio *TopTracksRadio) Run() {
	for {
		if len(radio.tracks) > 0 {
			var next *LastFmTrack
			next, radio.tracks = getRandomTrack(radio.tracks)
			radio.lastfmTracks <- next
			if len(radio.tracks) <= 49 {
				go radio.loadTopTracks()
			}
		}
	}
}

type LovedTracksRadio struct {
	spotify         *spotify.Session
	lastfm          *lastfm.Api
	tracks          []*LastFmTrack
	lastfmTracks    chan *LastFmTrack
	page            int
	lastFmUser      *LastFmUser
	currentUsername string
}

func (radio *LovedTracksRadio) Load() error {
	err := radio.loadLovedTracks()
	if err == nil {
		go radio.Run()
	}
	return err
}

func (radio *LovedTracksRadio) loadLovedTracks() error {
	radio.page = radio.page + 1
	params := lastfm.P{"user": radio.currentUsername, "limit": 50, "page": radio.page}
	tracks, err := radio.lastfm.User.GetLovedTracks(params)
	if err == nil {
		if tracks.Total < 1 {
			return errors.New("No tracks found!")
		}
		for _, track := range tracks.Tracks {
			lastFmArtist := &LastFmArtist{
				Name: track.Artist.Name,
			}
			// get image
			image := ""
			imgLen := len(track.Images)
			if imgLen > 0 {
				image = track.Images[imgLen-1].Url
			}
			lastFmTrack := &LastFmTrack{
				Artist: lastFmArtist,
				Name:   track.Name,
				Image:  image,
			}
			radio.tracks = append(radio.tracks, lastFmTrack)
		}
	}
	return err
}

func (radio *LovedTracksRadio) Run() {
	for {
		if len(radio.tracks) > 0 {
			var next *LastFmTrack
			next, radio.tracks = getRandomTrack(radio.tracks)
			radio.lastfmTracks <- next
			if len(radio.tracks) <= 49 {
				go radio.loadLovedTracks()
			}
		}
	}
}
