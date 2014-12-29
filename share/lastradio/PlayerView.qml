import QtQuick 2.0
import QtQuick.Controls 1.1

Rectangle {
    id: playerView
    visible: false
    color: "transparent"

    signal skipTrack
    onSkipTrack: {
        progressDuration.text = "0:00"
        progressTimer.running = false
        player.skip()
    }

    signal stop()
    onStop: {
        player.stop()
        stack.pop()
    }

    Rectangle {
        id: trackBox
        height: 360
        color: "transparent"
        anchors {
            right: parent.right
            left: parent.left
            top: parent.top
            bottom: playerActions.top
        }
        //Column {
        Rectangle {
            id: trackTitle
            color: "transparent"
            height: childrenRect.height
            anchors {
                margins: 10
                bottomMargin: 0
                top: parent.top
                left: parent.left
                right: parent.right
            }
            Text {
                width: parent.width
                text: track.name
                font.pointSize: 17
                color: "#FFF"
                font.bold: true
                wrapMode: Text.WordWrap
                elide: Text.ElideRight
                maximumLineCount: 2
                onTextChanged: {
                    trackProgress.value = 0
                    if (track.duration > 0) {
                        progressTimer.running = true
                        loveButton.enabled = true
                    }
                }
            }
        }
        Rectangle {
            id: trackAlbum
            height: childrenRect.height
            visible: track.album
            color: "transparent"
            anchors {
                margins: 10
                top: trackTitle.bottom
                topMargin: 5
                left: parent.left
                right: parent.right
            }
            Text {
                width: parent.width
                text: track.album
                wrapMode: Text.WordWrap
                elide: Text.ElideRight
                maximumLineCount: 2
                font.pointSize: 14
                color: "#FFF"
            }
        }
        Rectangle {
            id: trackYear
            height: childrenRect.height
            visible: track.year
            color: "transparent"
            anchors {
                margins: 10
                top: trackAlbum.bottom
                topMargin: 5
                left: parent.left
                right: parent.right
            }
            Text {
                text: track.year
                font.pointSize: 10
                color: "#FFF"
            }
        }
        Image {
            visible: track.image
            source: track.image
            height: trackBox.height
            width: trackBox.width
            z: -1
            opacity: 0.5
        }
        Rectangle {
            height: childrenRect.height
            color: "transparent"
            anchors {
                margins: 10
                bottom: parent.bottom
                bottomMargin: 30
                left: parent.left
                right: parent.right
            }
            Text {
                width: parent.width
                text: track.artist
                wrapMode: Text.WordWrap
                elide: Text.ElideRight
                maximumLineCount: 2
                font.pointSize: 17
                color: "#FFF"
                font.bold: true
            }
        }
    }
    Rectangle {
        id: progressBox
        visible: (track.duration)
        height: 10
        color: "transparent"
        anchors {
            left: parent.left
            right: parent.right
            bottom: trackBox.bottom
            bottomMargin: 10
        }
        function formatDuration(duration) {
            var minSecs = duration % 60
            minSecs = (minSecs < 10) ? "0" + minSecs : minSecs
            return parseInt(duration / 60) + ":" + minSecs
        }
        Timer {
            id: progressTimer
            interval: 1000
            running: false
            repeat: true
            onTriggered: {
                trackProgress.value = trackProgress.value + 1
                progressDuration.text = progressBox.formatDuration(
                            trackProgress.value)
                if (trackProgress.value === track.nowPlayingAt) {
                    player.sendToLastfm("nowplaying")
                } else if (trackProgress.value === track.scrobbleAt) {
                    player.sendToLastfm("scrobble")
                }
            }
        }
        Row {
            anchors.fill: parent
            Rectangle {
                id: progressDone
                width: 40
                height: parent.height
                color: "transparent"
                Text {
                    id: progressDuration
                    text: "0:00"
                    horizontalAlignment: Text.AlignHCenter
                    anchors.fill: parent
                    color: "#FFF"
                    font.bold: true
                }
            }
            ProgressBar {
                id: trackProgress
                width: parent.width - progressDone.width - progressTotal.width
                height: 10
                maximumValue: track.duration
                minimumValue: 0
                value: 0
            }
            Rectangle {
                id: progressTotal
                width: 40
                height: parent.height
                color: "transparent"
                Text {
                    text: (track.duration > 0) ? progressBox.formatDuration(
                                                     track.duration) : "-"
                    horizontalAlignment: Text.AlignHCenter
                    anchors.fill: parent
                    color: "#FFF"
                    font.bold: true
                }
            }
        }
    }

    PlayerActions {
        id: playerActions
    }

}
