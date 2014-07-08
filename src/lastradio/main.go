package main

import (
	"log"
	"os"
	"os/signal"

	"gopkg.in/qml.v0"

	"lastradio/lastradio"
)

func main() {

	defer lastradio.Close()

	// handle signals and libspotify
	exit := make(chan bool)
	go handleSignals(exit)

	qml.Init(nil)
	engine := qml.NewEngine()

	component, err := engine.LoadFile("share/lastradio/main.qml")
	if err != nil {
		panic(err)
	}

	currentTrack := &lastradio.PublishedData{}
	player := &lastradio.Player{
		Public: currentTrack,
		Exit:   exit,
	}
	context := engine.Context()
	context.SetVar("player", player)
	context.SetVar("track", currentTrack)

	window := component.CreateWindow(nil)

	window.Show()
	window.Wait()
}

func handleSignals(exit chan bool) {
	log.Print("TODO: signal handling")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	for _ = range signals {
		select {
		case exit <- true:
		default:
		}
	}
}
