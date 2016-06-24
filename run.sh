export GOPATH=`pwd`
SCRIPT=$(realpath $0)
BASEDIR=$(dirname $SCRIPT)
sh $BASEDIR/build_qrc.sh
GOPATH=$BASEDIR go build -race -x -v  -o bin/lastradio lastradio
#GOPATH=$BASEDIR go build -race -x -v -a -o bin/lastradio lastradio
#GOPATH=$BASEDIR GODEBUG=cgocheck=0 go install -x -v -a lastradio
$BASEDIR/bin/lastradio
