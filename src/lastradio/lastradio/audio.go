package lastradio

import (
	"log"

	"code.google.com/p/portaudio-go/portaudio"
	spotify "github.com/op/go-libspotify/spotify"
)

func initializeAudioConsumer(session *spotify.Session, control chan string) bool {
	// setup audio consumer
	portaudio.Initialize()
	pa := newPortAudio()
	go pa.Player(control)
	session.SetAudioConsumer(pa)
	return true
}

func playSpotifyTrack(session *spotify.Session, uri string) {
	log.Print("playing ", uri)
	// Parse the track
	link, err := session.ParseLink(uri)
	if err != nil {
		log.Fatal(err)
	}
	track, err := link.Track()
	if err != nil {
		log.Fatal(err)
	}

	// Load the track and play it
	track.Wait()
	player := session.Player()
	if err := player.Load(track); err != nil {
		log.Fatal(err)
	}

	player.Play()
}

type audio struct {
	format spotify.AudioFormat
	frames []byte
}

type audio2 struct {
	format spotify.AudioFormat
	frames []int16
}

type portAudio struct {
	buffer chan *audio
}

func newPortAudio() *portAudio {
	return &portAudio{
		buffer: make(chan *audio, 8),
	}
}

func (pa *portAudio) WriteAudio(format spotify.AudioFormat, frames []byte) int {
	audio := &audio{format, frames}
	//log.Print("audio", len(frames), len(frames)/2)

	if len(frames) == 0 {
		//log.Print("no frames")
		return 0
	}

	select {
	case pa.buffer <- audio:
		//log.Print("return", len(frames))
		return len(frames)
	default:
		//log.Print("buffer full")
		return 0
	}
}

func (pa *portAudio) Player(control chan string) {
	log.Print("PORTADIO PLAYER")
	out := make([]int16, 2048*2)

	stream, err := portaudio.OpenDefaultStream(
		0,
		2,     // audio.format.Channels,
		44100, // float64(audio.format.SampleRate),
		len(out),
		&out,
	)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	stream.Start()
	defer stream.Stop()

	// Decode the incoming data which is expected to be 2 channels and
	// delivered as int16 in []byte, hence we need to convert it.
	waiting := true
	for audio := range pa.buffer {
		if len(audio.frames) != 2048*2*2 {
			// most probably end of song
			if waiting != true {
				control <- "next"
			}
			waiting = true
			continue
		} else {
			waiting = false
		}
		j := 0
		for i := 0; i < len(audio.frames); i += 2 {
			out[j] = int16(audio.frames[i]) | int16(audio.frames[i+1])<<8
			j++
		}

		stream.Write()
	}
}
