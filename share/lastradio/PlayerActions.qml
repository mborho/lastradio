import QtQuick 2.0

Rectangle {
    id: playerActions
    height: 40
    visible: true
    color: "lightgrey"
    anchors {
        right: parent.right
        left: parent.left
        bottom: parent.bottom
    }
    function clickedAnimation(button) {
        clickAnim1.target = button
        clickAnim2.target = button
        clickAnim.start()
    }
    SequentialAnimation {
        id: clickAnim
        running: false
        NumberAnimation { id:clickAnim1; property: "scale"; to: 0.8; duration: 100}
        NumberAnimation { id:clickAnim2; property: "scale"; to: 1.0; duration: 100}
    }
    Row {
        spacing: 8
        anchors.centerIn: parent
        Image {
            id: playButton
            property bool paused: false
            source: "png/pause.png"
            width: 40
            height: 40
            MouseArea {
                anchors.fill:parent
                onClicked: {
                    playerActions.clickedAnimation(playButton)
                    if (!playButton.paused) {
                        player.pause()
                        playButton.paused = true
                        playButton.source = "png/play.png"
                        progressTimer.running = false
                    } else {
                        if (playButton.paused) {
                            progressTimer.running = true
                        }
                        player.play()
                        playButton.source = "png/pause.png"
                        playButton.paused = false
                    }
                }
            }
        }
        Image {
            id: skipButton
            source: "png/skip.png"
            width: 40
            height: 40
            MouseArea {
                anchors.fill:parent
                onClicked: {
                    playerActions.clickedAnimation(skipButton)
                    playerView.skipTrack()
                }
            }
        }
        Image {
            id: loveButton
            source: (track.isLoved) ? "png/loved.png" :"png/unloved.png"
            width: 20
            height: 20
            smooth: true
            anchors {
                verticalCenter: parent.verticalCenter
            }
            MouseArea {
                anchors.fill:parent
                onClicked: {
                    playerActions.clickedAnimation(loveButton)
                    if(track.isLoved) {
                        player.sendToLastfm("unlove")
                        track.isLoved = false
                    } else {
                        player.sendToLastfm("love")
                        track.isLoved = true
                    }
                }
            }
        }
        /*Image {
            id: banButton
            source: "png/dislike.png"
            width: 20
            height: 20
            smooth: true
            anchors {
                verticalCenter: parent.verticalCenter
            }
            MouseArea {
                anchors.fill:parent
                onClicked: {
                    playerActions.clickedAnimation(banButton)
                    player.sendToLastfm("ban")
                    playerView.skipTrack()
                }
            }
        }*/
        Image {
            id: radioButton
            source: "png/radio.png"
            width: 20
            height: 20
            smooth: true
            anchors {
                verticalCenter: parent.verticalCenter
                verticalCenterOffset: -1
            }
            MouseArea {
                anchors.fill:parent
                onClicked: {
                    playerActions.clickedAnimation(radioButton)
                    stack.pop()                    }
            }
        }
    }
}
