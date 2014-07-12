# LastSpotify #

Simple Last.fm radio player written in Go and QML, using Spotify as music source.

(Requires a Last.fm account and Spotify Premium)

## Installation ##
LastRadio is available for Ubuntu 14.04 LTS x86_64.


```
#!bash

sudo add-apt-repository ppa:martin-borho/lastradio
sudo apt-get update
sudo apt-get install lastradio
```

## Build


#### Develop with QtCreator/UbuntuSDK ####

LastRadio is developed by using the [UbuntuSDK](http://developer.ubuntu.com/apps/sdk/) with Go support.

* checkout the source
* open the lastradio.goproject file
* have only one package uncommented in lastradio.goproject, build and repeat with the next. [Bugreport](https://bugs.launchpad.net/qtcreator-plugin-go/+bug/1322853) for this.
* Run

#### From source ####

Install missing dependency:

```
#!bash
sudo apt install mercurial g++ golang-go portaudio19-dev \
     qtbase5-private-dev qtdeclarative5-private-dev \
     libqt5opengl5-dev qtdeclarative5-qtquick2-plugin \
     qtdeclarative5-window-plugin qtdeclarative5-localstorage-plugin \
     qtdeclarative5-controls-plugin
```

Download and install [Libspotify](https://developer.spotify.com/technologies/libspotify/):

```
#!bash

wget https://developer.spotify.com/download/libspotify/libspotify-12.1.51-Linux-x86_64-release.tar.gz
tar xvfz libspotify-12.1.51-Linux-x86_64-release.tar.gz
cd libspotify-12.1.51-Linux-x86_64-release
# install system wide
sudo make install prefix=/usr/local
```

checkout source and build:
```
#!bash

git clone git@gitorious.org:lastradio/lastradio.git
cd lastspotify
export GOPATH=/path/to/lastspotify/
go get gopkg.in/qml.v0
go get code.google.com/p/portaudio-go/portaudio
go get github.com/shkh/lastfm-go/lastfm
go get github.com/op/go-libspotify/spotify
go run src/lastradio/main.go
``` 


### License ###
Copyright 2014 Martin Borho <martin@borho.net> 

GPLv3 - see LICENSE for details
