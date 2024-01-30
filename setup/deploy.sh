#!/bin/bash

SRVR_DIR="$HOME/.spot"
LOGS_DIR="$SRVR_DIR/logs"
CONFFILE="$SRVR_DIR/spot.json"
VERSION="v0.1"
DEV_ADDR=
DEV_PORT=8000
MONGODB_PORT=27017

mkdir -p $SRVR_DIR
mkdir -p $LOGS_DIR

echo "Writing config to $CONFFILE"
source setup/config.sh
cat $CONFFILE

EXE=$GOPATH/bin/spot
GOBIN=$GOPATH/bin go install setup/spot.go
if [ $? -eq 0 ]; then
    echo "Install successful"
else
    exit 1
fi

echo "Deploying $VERSION"
PARAM="-config=$CONFFILE"
set -x
$EXE $PARAM &
{ set +x; } 2>/dev/null

# wait for the server to start up
sleep 1
echo "Done"
