package main

import (
	"log"
	"os"
	"os/signal"

	"gopkg.in/qml.v1"

	"lastradio/lastradio"
)

func main() {
	if err := qml.Run(run); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {

	defer lastradio.Close()

	// handle signals and libspotify
	exit := make(chan bool)
	go handleSignals(exit)

	engine := qml.NewEngine()

	component, err := engine.LoadFile("qrc:///share/lastradio/main.qml")
	if err != nil {
		return err
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
	return nil
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
