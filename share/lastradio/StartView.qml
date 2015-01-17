import QtQuick 2.0
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.1
import QtQuick.Controls.Styles 1.1

Rectangle {
    id: startView
    visible: false
    color: "lightgrey"

    signal startRadio(string name, string label, string username)
    onStartRadio: {
        stack.push(playerView)
        if(radioName !== label || name === "similar") {
            player.stop()
            started = true
            radioName = label
            currentUsername = username
            player.loadRadio(name, username)
            player.play()
        }
    }

    signal show()
    onShow: { }

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
        ListElement {
            name: "similar"
            label: "Similar Artists"
        }
    }

    Rectangle {
        id: tools
        height:40
        color: "#F2F1F0"
        anchors {
            top: parent.top
            left: parent.left
            right: parent.right
        }

        RowLayout {
            anchors.verticalCenter: parent.verticalCenter
            Button {
                visible: started
                text: "Back to radio"
                onClicked: {
                    stack.push(playerView)
                }
                /*style: ButtonStyle {
                        background: Rectangle {
                            implicitWidth: 100
                            implicitHeight: 25
                            border.width: control.activeFocus ? 2 : 1
                            border.color: "#888"
                            radius: 4
                            gradient: Gradient {
                                GradientStop { position: 0 ; color: control.pressed ? "#ccc" : "#eee" }
                                GradientStop { position: 1 ; color: control.pressed ? "#aaa" : "#ccc" }
                            }
                        }
                    }*/
            }
            Button {
                text: "Reset accounts"
                onClicked: {
                    clearSetting("lastfmUser");
                    clearSetting("lastfmPwd");
                    clearSetting("spotifyUser");
                    clearSetting("spotifyPwd");
                    started = false
                    playerView.stop()
                    player.logout()
                    stack.push(loginView)
                }
            }
        }
    }
    ListView {
        id: radioListView
        anchors {
            left: parent.left
            right: parent.right
            top: tools.bottom
            bottom: parent.bottom
        }

        anchors.margins: 40
        model: radioList
        property string username: ""
        header: Item {
            id: radioListHeader
            width: parent.width
            height: childrenRect.height

            Text {
                id: listHeaderText
                text: "Choose your station"
                font.bold: true
                font.pointSize: 17
                //color: "#FFF"
            }
            /*Rectangle {
                id: userChoice
                height: childrenRect.height
                color: "transparent"
                anchors {
                    top: listHeaderText.bottom
                    right:parent.right
                    left:parent.left
                }
                Row {
                    height: 50
                    anchors {
                        right: parent.right
                        left: parent.left
                    }
                    Text {
                        id: userInputLabel
                        text: "for "
                        //width: parent.width/2
                        font.pointSize: 14
                        //color: "#FFF"
                    }
                    TextField {
                        id: lastfmUserInput
                        width: parent.width/2
                        anchors.verticalCenter: userInputLabel.verticalCenter
                        placeholderText: "you"
                        onEditingFinished: {
                            radioListView.username = text.trim()
                        }

                        style: TextFieldStyle {
                            background: Rectangle {
                                color: "#FFFFFF"
                                border.color: "#333"
                                border.width: 1
                                radius: 5
                            }
                        }
                    }
                }
            }*/
        }

        delegate: Item {
            width: parent.width
            height: 30
            Text {
                id: radioLabel
                text: label
                font.bold: false
                font.pointSize: 17
                //color: "#FFF"
            }
            MouseArea {
                anchors.fill: parent
                onClicked: {
                    if(name !== "similar") {
                        radioListView.currentIndex = index
                        startView.startRadio(name, label, radioListView.username)
                    }
                }
            }
            Image {
                source: "png/online.png"
                width: 26
                height: 26
                anchors.right: parent.right
                visible: (started && radioListView.currentIndex === index)
            }
            Rectangle {
                id: bandChoice
                visible: (name === "similar")
                height: childrenRect.height
                color: "transparent"
                anchors {
                    top: radioLabel.bottom
                    right:parent.right
                    left:parent.left
                }
                Row {
                    height: 50
                    width: parent.width
                    Text {
                        id: bandInputLabel
                        width: parent.width/7
                        text: "to "
                        font.pointSize: 12
                    }
                    TextField {
                        id: bandUserInput
                        width: parent.width/2
                        anchors.verticalCenter: bandInputLabel.verticalCenter
                        placeholderText: "artist name"
                        style: TextFieldStyle {
                            background: Rectangle {
                                color: "#FFFFFF"
                                border.color: "#333"
                                border.width: 1
                                radius: 5
                            }
                        }
                    }
                    Button {
                        id: relatedsStart
                        width: parent.width/5
                        anchors.verticalCenter: bandInputLabel.verticalCenter
                        text: "Start"
                        onClicked: {
                            radioListView.currentIndex = index
                            startView.startRadio(name, label, bandUserInput.text)
                        }
                    }
                }
            }
        }
    }
}
