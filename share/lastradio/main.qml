import QtQuick 2.0
import QtQuick.Controls 1.1
import QtQuick.LocalStorage 2.0

ApplicationWindow {
    title: "LastRadio"+((radioName) ? " - "+radioName : "")+((currentUsername) ? " ("+currentUsername+")" : "")
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
    property string currentUsername: ""

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

    LoginView {
        id:loginView
    }

    StartView {
        id: startView
    }

    PlayerView {
        id: playerView
    }

    StackView {
        id: stack
        initialItem: loginView
    }

}
