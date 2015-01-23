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
        if(radioName !== label || name === "similar" || name === "tag") {
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
            label: "Similar"
        }
        ListElement {
            name: "tag"
            label: "Tag Radio"
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

        anchors.margins: 20
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
        }

        delegate: Item {
            width: parent.width            
            height: 45
            property bool active: (started && radioListView.currentIndex === index)
            Rectangle {
                anchors {
                    fill: parent
                    topMargin:3
                    bottomMargin:3
                }
                color: (active) ? "#E1E1E1" : "transparent"
                border.color:  (active) ? "#0D3D08": "#BBBBBB"
                radius: 6
                Rectangle {
                    anchors {
                        fill: parent
                        leftMargin: 10
                        rightMargin: 10
                    }
                    color: "transparent"
                    Text {
                        id: radioLabel
                        text: label + ((name === "similar") ?  " to" : "")
                        font.bold: false
                        font.pointSize: 17
                        anchors.verticalCenter: parent.verticalCenter
                    }
                    MouseArea {
                        anchors.fill: parent
                        onClicked: {
                            if(name !== "similar" && name !== "tag") {
                                radioListView.currentIndex = index
                                startView.startRadio(name, label, radioListView.username)
                            }
                        }
                    }
                    TextField {
                        id: radioTermInput
                        visible: (name === "similar" || name == "tag")
                        width: parent.width/2
                        anchors {
                            verticalCenter: parent.verticalCenter
                            right: extraStart.left
                            rightMargin: 10
                        }
                        placeholderText: (name === "similar") ? "artist name" : "tag name"
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
                        id: extraStart
                        visible: (name === "similar" || name == "tag")
                        width: parent.width/7
                        anchors {
                            right: parent.right
                            rightMargin: 0
                            verticalCenter: parent.verticalCenter
                        }
                        text: "Start"
                        onClicked: {
                            radioListView.currentIndex = index
                            startView.startRadio(name, label, radioTermInput.text)
                        }
                    }
                }
            }
        }
    }
}
