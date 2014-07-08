import QtQuick 2.0
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.1
import QtQuick.LocalStorage 2.0
import QtQuick.Controls.Styles 1.1

ApplicationWindow {
    title: "LastRadio"+((radioName) ? " - "+radioName : "")
    id: root
    width: 400
    height: 400
    minimumHeight: 400
    minimumWidth: 400
    maximumHeight: 400
    maximumWidth: 400
    color: "#000000"
    property bool started: false
    property string radioName: ""
    property var db: null

    function openDB(name, value) {
        var identifier = "lastspotify";
        db = LocalStorage.openDatabaseSync(identifier, "", "", 1000);
        if (db.version === "") {
            db.changeVersion("", "0.1",
                function(tx) {
                    tx.executeSql('CREATE TABLE IF NOT EXISTS settings(key TEXT UNIQUE, value TEXT)');
                    console.log('Database created');
                });
            // reopen database with new version number
            db = LocalStorage.openDatabaseSync(identifier, "", "", 1000);
        }
    }

    function saveSetting(key, value) {
        openDB();
         db.transaction( function(tx){
             tx.executeSql('INSERT OR REPLACE INTO settings VALUES(?, ?)', [key, value]);
         });
    }

    function getSettings(callback) {
        openDB();
        var settings = {};
        db.readTransaction(
            function(tx){
                var rs = tx.executeSql('SELECT key, value FROM Settings');
                for(var i = 0; i < rs.rows.length; i++) {
                    var row = rs.rows.item(i);
                    settings[row.key] = row.value;
                }
                callback(settings);
            }
        );
    }

    function clearSetting(name) {
        openDB();
        db.transaction(function(tx){
            tx.executeSql('DELETE FROM Settings WHERE key = ?', [name]);
        });
    }

    toolBar: ToolBar {
        visible: startView.visible
        RowLayout {
            ToolButton {
                visible: started
                text: "Back to radio"
                onClicked: {
                    startView.visible = false
                    playerView.visible = true
                }
            }
            ToolButton {
                text: "Reset accounts"
                onClicked: {
                    clearSetting("lastfmUser");
                    clearSetting("lastfmPwd");
                    clearSetting("spotifyUser");
                    clearSetting("spotifyPwd");
                    started = false
                    playerView.stop()
                    player.logout()
                }
            }
        }
    }

    Rectangle {
        id: loginView
        anchors.fill:parent
        color:"lightgrey"
        visible: !player.loggedIn
        onVisibleChanged: {
            if(!player.loggedIn) {
                clearLoginForm();
            } else {
                startView.visible = true
            }
        }

        function clearLoginForm() {
            loginMessage.text = ""
            lastFmLoginButton.visible = true
            spotifyLoginButton.visible = true
            startView.visible = false
        }

        Component.onCompleted: {
            getSettings(function(settings) {
                if(settings["spotifyUser"] && settings["lastfmUser"]) {
                    var loginLF = player.loginToLastFm(settings["lastfmUser"], settings["lastfmPwd"]);
                    var loginSP = player.loginToSpotify(settings["spotifyUser"], settings["spotifyPwd"])
                }
            })
        }

        Text {
            anchors {
                top:parent.top
                topMargin: 30
                right: parent.right
                rightMargin: 30
                left: parent.left
                leftMargin: 30
            }
            font.bold: false
            font.pointSize: 14
            text: "Please specify your account\ncredentials!"
        }
        Column {
            anchors {
                horizontalCenter: parent.horizontalCenter
                top: parent.top
                topMargin: 100
            }
            Text {
                id: loginMessage
                text: ""
            }

            GroupBox {
                title: "LastFm-Login"
                Row {
                    TextField {
                        anchors.rightMargin: 10
                        id:lastfmUser
                        readOnly: false
                        style: TextFieldStyle {
                            textColor: "black"
                        }
                    }
                    TextField {
                        id:lastfmPwd
                        echoMode: TextInput.Password
                        readOnly: false
                        style: TextFieldStyle {
                            textColor: "black"
                        }
                    }
                    Button {
                        id: lastFmLoginButton
                        text: "Login"
                        onClicked: {
                            var login = player.loginToLastFm(lastfmUser.text, lastfmPwd.text)
                            if(!login) {
                                loginMessage.text = "Invalid Last.fm credentials!"
                            } else {
                                lastFmLoginButton.visible = false
                                lastfmUser.readOnly = true
                                lastfmPwd.readOnly = true
                                saveSetting("lastfmUser", lastfmUser.text)
                                saveSetting("lastfmPwd", lastfmPwd.text)
                            }
                        }
                    }
                }
            }

            GroupBox {
                title: "Spotify-Login (Premium-Account)"
                Row {
                    TextField {
                        id:spotifyUser
                        style: TextFieldStyle {
                            textColor: "black"
                        }
                    }
                    TextField {
                        id:spotifyPwd
                        echoMode: TextInput.Password
                        style: TextFieldStyle {
                            textColor: "black"
                        }
                    }
                    Button {
                        id: spotifyLoginButton
                        text: "Login"
                        onClicked: {
                            var login = player.loginToSpotify(spotifyUser.text, spotifyPwd.text)
                            if(!login) {
                                loginMessage.text = "Invalid Spotify credentials!"
                            } else {
                                spotifyLoginButton.visible = false
                                spotifyUser.readOnly = true
                                spotifyPwd.readOnly = true
                                saveSetting("spotifyUser", spotifyUser.text)
                                saveSetting("spotifyPwd", spotifyPwd.text)
                            }
                        }
                    }
                }
            }
        }
    }

    Rectangle {
        id: startView
        visible: false
        anchors.fill: parent
        color: "lightgrey"

        signal startRadio(string name, string label)
        onStartRadio: {
            startView.visible = false
            playerView.visible = true
            if(radioName !== label) {
                player.stop()
                started = true
                radioName = label
                player.loadRadio(name)
                player.play()
            }
        }

        signal show()
        onShow: {
            visible = true
            playerView.visible = false
        }

        ListModel {
            id: radioList
            ListElement {
                name: "top"
                label: "Top Tracks"
            }
            ListElement {
                name: "loved"
                label: "Loved Tracks"
            }
            ListElement {
                name: "recommended"
                label: "Recommended Artists"
            }
            ListElement {
                name: "topartists"
                label: "Top Artists"
            }
        }

        ListView {
            id: radioListView
            anchors.fill: parent
            anchors.margins: 40
            model: radioList
            delegate: Item {
                width: parent.width
                height: 30
                Text {
                    text: label
                    font.bold: false
                    font.pointSize: 17
                    //color: "#FFF"
                }
                MouseArea {
                    anchors.fill: parent
                    onClicked: {
                        radioListView.currentIndex = index
                        startView.startRadio(name, label)
                    }
                }
                Image {
                    source: "png/online.png"
                    width: 26
                    height: 26
                    anchors.right: parent.right
                    visible: (started && radioListView.currentIndex === index)
                }
            }
        }
    }

    Rectangle {
        id: playerView
        anchors.fill: parent
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
            playerView.visible = false
            startView.visible = true
        }

        Rectangle {
            id: trackBox
            height: 360
            color: "transparent"
            anchors {
                right: parent.right
                left: parent.left
                top: parent.top
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
                    text: track.name
                    font.pointSize: 17
                    color: "#FFF"
                    font.bold: true
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
                    text: track.album
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
                height: parent.height
                width: parent.width
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
                    text: track.artist
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

        Rectangle {
            id: playerActions
            visible: (!startView.visible && !loginView.visible)
            height: 40
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
                Image {
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
                }
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
                            startView.show()
                        }
                    }
                }
            }
        }
    }
}
