SCRIPT=$(realpath $0)
BASEDIR=$(dirname $SCRIPT)
(cd $BASEDIR && GOPATH=$BASEDIR go run src/gopkg.in/qml.v1/cmd/genqrc/main.go share/lastradio/)
mv $BASEDIR/qrc.go $BASEDIR/src/lastradio/
