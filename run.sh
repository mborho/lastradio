export GOPATH=`pwd`
SCRIPT=$(realpath $0)
BASEDIR=$(dirname $SCRIPT)
sh $BASEDIR/build_qrc.sh
GOPATH=$BASEDIR go install -race -x -v lastradio
$BASEDIR/bin/lastradio
