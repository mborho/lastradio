package lastradio

import (
	"testing"
)

func Test_prepareSpotifyTerm(t *testing.T) {
	s := "The Skull Defekts 02 King of Misinformation (feat. Daniel Higgs)"
	prepared := prepareSpotifyTerm(s)
	if prepared != "The Skull Defekts King of Misinformation" {
		t.Error(s + " not stripped")
	}

	s = "Wu-Tang Clan '96 Recreation (Demo)"
	prepared = prepareSpotifyTerm(s)
	if prepared != "Wu-Tang Clan Recreation" {
		t.Error(s + " not stripped")
	}

	s = "Sleaford Mods 09 Showboat"
	prepared = prepareSpotifyTerm(s)
	if prepared != "Sleaford Mods Showboat" {
		t.Error(s + " not stripped")
	}
}
