package lastradio

import (
	"fmt"
	"log"
	"time"

	spotify "github.com/op/go-libspotify/spotify"
	"github.com/shkh/lastfm-go/lastfm"
)

type LastFmArtist struct {
	Name  string
	Image string
}

type LastFmTrack struct {
	Artist      *LastFmArtist
	Name        string
	SpotifyLink string
	Duration    time.Duration
	Album       string
	Year        int
	Image       string
	IsLoved     bool
}

func getLastFmTrack(api *lastfm.Api, trackList chan *LastFmTrack, artists chan *LastFmArtist) {
	log.Print("getLastFmTrack")
	for artist := range artists {
		tracks, err := api.Artist.GetTopTracks(lastfm.P{"artist": artist.Name, "limit": 5})
		if err != nil {
			log.Println(err)
			continue
		}
		index := randInt(len(tracks.Tracks))
		track := tracks.Tracks[index]
		trackList <- &LastFmTrack{
			Artist: artist,
			Name:   track.Name,
		}
	}
}

func getNextArtist(recommended *lastfm.UserGetRecommendedArtists) string {
	nextTrack := recommended.Artists[len(recommended.Artists)-1]
	recommended.Artists = recommended.Artists[:len(recommended.Artists)-1]
	return nextTrack.Name
}

func getTrackInfo(api *lastfm.Api, username string, trackList chan *LastFmTrack, spotifyQueue chan *LastFmTrack) {
	log.Print("getTrackInfo")
	for track := range trackList {
		params := lastfm.P{"username": username, "artist": track.Artist.Name, "track": track.Name}
		trackInfo, err := api.Track.GetInfo(params)
		if err != nil {
			log.Println(err)
			continue
		}
		if trackInfo.UserLoved == "1" {
			track.IsLoved = true
		}
		// get image
		imgLen := len(trackInfo.Album.Images)
		if imgLen > 0 {
			track.Image = trackInfo.Album.Images[imgLen-1].Url
		}
		spotifyQueue <- track
	}
}

func getSpotifyData(session *spotify.Session, tracklist, playlist chan *LastFmTrack, control chan string) error {
	log.Print("getSpotifyData")
	for track := range tracklist {
		term := track.Artist.Name + " " + track.Name
		log.Print("Received: ", term)
		query := prepareSpotifyTerm(term)
		if term != query {
			log.Print("  Prepared: " + query)
		}
		spec := spotify.SearchSpec{0, 1}
		var sOpts = &spotify.SearchOptions{
			Tracks: spec,
		}

		search, err := session.Search(query, sOpts)
		if err != nil {
			return err
		}
		search.Wait()
		for i := 0; i < search.Tracks(); i++ {
			spotifyLink, duration, album, year := getTrackData(search.Track(i))
			if spotifyLink != "" {
				track.SpotifyLink = spotifyLink
				track.Duration = duration
				track.Album = album
				track.Year = year
				log.Println("-> adding to play queue")
				playlist <- track
				break
			}
		}
	}
	return nil
}

func getTrackData(track *spotify.Track) (string, time.Duration, string, int) {
	track.Wait()
	album := track.Album()
	album.Wait()
	return fmt.Sprintf("%s", track.Link()), track.Duration(), album.Name(), album.Year()
}

func SendNowPlaying(api *lastfm.Api, artist, track, album string) {
	p := lastfm.P{"artist": artist, "track": track, "album": album}
	_, err := api.Track.UpdateNowPlaying(p)
	if err != nil {
		log.Print(err)
	}
}

func Scrobble(api *lastfm.Api, artist, track, album string) {
	start := time.Now().Unix()
	p := lastfm.P{"artist": artist, "album": album, "track": track, "timestamp": start, "chosenByUser": 0}
	_, err := api.Track.Scrobble(p)
	if err != nil {
		log.Print(err)
	}
}

func Love(api *lastfm.Api, artist, track string) {
	p := lastfm.P{"track": track, "artist": artist}
	err := api.Track.Love(p)
	if err != nil {
		log.Print(err)
	}
}

func Unlove(api *lastfm.Api, artist, track string) {
	p := lastfm.P{"artist": artist, "track": track}
	err := api.Track.UnLove(p)
	if err != nil {
		log.Print(err)
	}
}

func Ban(api *lastfm.Api, artist, track string) {
	p := lastfm.P{"artist": artist, "track": track}
	err := api.Track.Ban(p)
	if err != nil {
		log.Print(err)
	}
}
