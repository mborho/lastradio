package lastradio

import (
	"math/rand"
	"time"
)

func randInt(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func getRandomArtist(artists []*LastFmArtist) (*LastFmArtist, []*LastFmArtist) {
	randIndex := randInt(len(artists))
	next := artists[randIndex]
	artists = append(artists[0:randIndex], artists[randIndex+1:]...)
	return next, artists
}

func getRandomTrack(tracks []*LastFmTrack) (*LastFmTrack, []*LastFmTrack) {
	randIndex := randInt(len(tracks))
	next := tracks[randIndex]
	tracks = append(tracks[0:randIndex], tracks[randIndex+1:]...)
	return next, tracks
}
