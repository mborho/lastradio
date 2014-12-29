import QtQuick 2.0

Rectangle {
    id: startView
    visible: false
    color: "lightgrey"

    signal startRadio(string name, string label, string username)
    onStartRadio: {
        stack.push(playerView)
        if(radioName !== label) {
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
    }

    ListView {
        id: radioListView
        anchors.fill: parent
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
                text: label
                font.bold: false
                font.pointSize: 17
                //color: "#FFF"
            }
            MouseArea {
                anchors.fill: parent
                onClicked: {
                    radioListView.currentIndex = index
                    startView.startRadio(name, label, radioListView.username)
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
