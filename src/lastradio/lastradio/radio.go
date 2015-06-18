package lastradio

import (
	"errors"
	spotify "github.com/op/go-libspotify/spotify"
	"github.com/shkh/lastfm-go/lastfm"
	"log"
	"sync"
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
	wg           sync.WaitGroup
}

func (radio *RecommendedRadio) Load() error {
	radio.artistslist = make(chan *LastFmArtist)
	go getLastFmTrack(radio.lastfm, radio.lastfmTracks, radio.artistslist)
	radio.wg.Add(1)
	err := radio.loadArtists()
	radio.wg.Wait()
	if err != nil {
		return err
	}
	go radio.Run()
	return nil
}

func (radio *RecommendedRadio) loadArtists() error {
	defer radio.wg.Done()
	radio.page = radio.page + 1
	var params = lastfm.P{"limit": 50, "page": radio.page}
	artists, err := radio.lastfm.User.GetRecommendedArtists(params)
	if err != nil {
		return err
	}
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
	return err
}

func (radio *RecommendedRadio) Run() {
	for {
		var next *LastFmArtist
		if len(radio.artists) > 0 {
			next, radio.artists = getRandomArtist(radio.artists)
			radio.artistslist <- next
			if len(radio.artists) <= 49 {
				radio.wg.Add(1)
				log.Printf("########### load more artists, sum %d\n", len(radio.artists))
				go radio.loadArtists()
			}
			radio.wg.Wait()
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
	wg              sync.WaitGroup
}

func (radio *TopArtistsRadio) Load() error {
	radio.artistslist = make(chan *LastFmArtist)
	go getLastFmTrack(radio.lastfm, radio.lastfmTracks, radio.artistslist)

	radio.wg.Add(1)
	err := radio.loadArtists()
	radio.wg.Wait()
	if err != nil {
		return err
	}
	go radio.Run()
	return nil
}

func (radio *TopArtistsRadio) loadArtists() error {
	defer radio.wg.Done()
	radio.page = radio.page + 1
	var params = lastfm.P{"user": radio.currentUsername, "limit": 50, "page": radio.page}
	artists, err := radio.lastfm.User.GetTopArtists(params)
	if err != nil {
		return err
	}
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
	return err
}

func (radio *TopArtistsRadio) Run() {
	for {
		var next *LastFmArtist
		if len(radio.artists) > 0 {
			next, radio.artists = getRandomArtist(radio.artists)
			radio.artistslist <- next
			if len(radio.artists) <= 49 {
				radio.wg.Add(1)
				log.Printf("########### load more artists, sum %d\n", len(radio.artists))
				go radio.loadArtists()
			}
			radio.wg.Wait()
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
	wg              sync.WaitGroup
}

func (radio *TopTracksRadio) Load() error {
	radio.wg.Add(1)
	err := radio.loadTopTracks()
	radio.wg.Wait()
	if err != nil {
		return err
	}
	go radio.Run()
	return nil
}

func (radio *TopTracksRadio) loadTopTracks() error {
	defer radio.wg.Done()
	radio.page = radio.page + 1
	params := lastfm.P{"user": radio.currentUsername, "limit": 50, "page": radio.page}
	tracks, err := radio.lastfm.User.GetTopTracks(params)
	if err != nil {
		return err
	}
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
	return err
}

func (radio *TopTracksRadio) Run() {
	for {
		if len(radio.tracks) > 0 {
			var next *LastFmTrack
			next, radio.tracks = getRandomTrack(radio.tracks)
			radio.lastfmTracks <- next
			if len(radio.tracks) <= 49 {
				radio.wg.Add(1)
				log.Printf("########### load more tracks, sum %d\n", len(radio.tracks))
				go radio.loadTopTracks()
			}
			radio.wg.Wait()
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
	wg              sync.WaitGroup
}

func (radio *LovedTracksRadio) Load() error {
	radio.wg.Add(1)
	err := radio.loadLovedTracks()
	radio.wg.Wait()
	if err != nil {
		return err
	}
	go radio.Run()
	return nil
}

func (radio *LovedTracksRadio) loadLovedTracks() error {
	defer radio.wg.Done()
	radio.page = radio.page + 1
	params := lastfm.P{"user": radio.currentUsername, "limit": 50, "page": radio.page}
	tracks, err := radio.lastfm.User.GetLovedTracks(params)
	if err != nil {
		return err
	}
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
	return err
}

func (radio *LovedTracksRadio) Run() {
	for {
		if len(radio.tracks) > 0 {
			var next *LastFmTrack
			next, radio.tracks = getRandomTrack(radio.tracks)
			radio.lastfmTracks <- next
			if len(radio.tracks) <= 49 {
				radio.wg.Add(1)
				log.Printf("Loading more tracks, currently %d tracks.\n", len(radio.tracks))
				go radio.loadLovedTracks()
			}
			radio.wg.Wait()
		}
	}
}

type SimilarRadio struct {
	spotify      *spotify.Session
	lastfm       *lastfm.Api
	artists      []*LastFmArtist
	artistslist  chan *LastFmArtist
	lastfmTracks chan *LastFmTrack
	page         int
	lastFmUser   *LastFmUser
	bandName     string
	wg           sync.WaitGroup
}

func (radio *SimilarRadio) Load() error {
	radio.artistslist = make(chan *LastFmArtist)
	go getLastFmTrack(radio.lastfm, radio.lastfmTracks, radio.artistslist)
	radio.wg.Add(1)
	err := radio.loadSimilar()
	radio.wg.Wait()
	if err != nil {
		return err
	}
	go radio.Run()
	return nil
}

func (radio *SimilarRadio) loadSimilar() error {
	defer radio.wg.Done()
	radio.page = radio.page + 1
	params := lastfm.P{"artist": radio.bandName, "autocorrect": 1, "limit": 50, "page": radio.page}
	artists, err := radio.lastfm.Artist.GetSimilar(params)
	if err != nil {
		return err
	}
	if len(artists.Similars) < 1 {
		return errors.New("No artists found!")
	}
	for _, artist := range artists.Similars {
		lastFmArtist := &LastFmArtist{
			Name:  artist.Name,
			Image: artist.Images[2].Url,
		}
		radio.artists = append(radio.artists, lastFmArtist)
	}
	return err
}

func (radio *SimilarRadio) Run() {
	for {
		var next *LastFmArtist
		if len(radio.artists) > 0 {
			next, radio.artists = getRandomArtist(radio.artists)
			radio.artistslist <- next
			if len(radio.artists) <= 49 {
				radio.wg.Add(1)
				log.Printf("Loading more artists, currently %d\n", len(radio.artists))
				go radio.loadSimilar()
			}
			radio.wg.Wait()
		}
	}
}

type TagTracksRadio struct {
	spotify      *spotify.Session
	lastfm       *lastfm.Api
	tracks       []*LastFmTrack
	lastfmTracks chan *LastFmTrack
	page         int
	lastFmUser   *LastFmUser
	tagName      string
	wg           sync.WaitGroup
}

func (radio *TagTracksRadio) Load() error {
	radio.wg.Add(1)
	err := radio.loadTagTracks()
	radio.wg.Wait()
	if err != nil {
		return err
	}
	go radio.Run()
	return nil
}

func (radio *TagTracksRadio) loadTagTracks() error {
	defer radio.wg.Done()
	radio.page = radio.page + 1
	params := lastfm.P{"tag": radio.tagName, "limit": 50, "page": radio.page}
	tracks, err := radio.lastfm.Tag.GetTopTracks(params)
	if err != nil {
		return err
	}
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
	return err
}

func (radio *TagTracksRadio) Run() {
	for {
		if len(radio.tracks) > 0 {
			var next *LastFmTrack
			next, radio.tracks = getRandomTrack(radio.tracks)
			radio.lastfmTracks <- next
			if len(radio.tracks) <= 49 {
				radio.wg.Add(1)
				log.Printf("Loading more tracks, currently %d\n", len(radio.tracks))
				go radio.loadTagTracks()
			}
			radio.wg.Wait()
		}
	}
}
