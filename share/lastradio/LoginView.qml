import QtQuick 2.0
import QtQuick.Controls 1.2
import QtQuick.Controls.Styles 1.2

Rectangle {
    id: loginView
    color:"lightgrey"
    /*anchors {
        fill:parent
    }*/

    property bool loggedIn: player.loggedIn
    onLoggedInChanged: {
        if(!player.loggedIn) {
            clearLoginForm();
        } else {
            stack.push(startView)
        }
    }

    function clearLoginForm() {
        loginMessage.text = ""
        lastFmLoginButton.visible = true
        spotifyLoginButton.visible = true
        stack.push(startView)
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
        visible: !player.loggedIn
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
        visible: !player.loggedIn
        anchors {
            //horizontalCenter: parent.horizontalCenter
            top: parent.top
            topMargin: 100
            right: parent.right
            rightMargin: 30
            left:parent.left
            leftMargin: 30
        }
        Text {
            id: loginMessage
            text: ""
        }

        GroupBox {
            title: "LastFm-Login"            
            anchors {

            }
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
