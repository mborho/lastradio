SCRIPT=$(realpath $0)
BASEDIR=$(dirname $SCRIPT)
sh $BASEDIR/build_qrc.sh
GOPATH=$BASEDIR go install -x -v lastradio
$BASEDIR/bin/lastradio
