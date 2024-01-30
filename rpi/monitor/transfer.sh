#!/bin/bash

BIN=monitor
TARGET_OS=linux
TARGET_ARCH=arm
ARM_VERSION=7
DEVICE_ADDR=192.168.0.101
REMOTE_USER=pi
REMOTE_PATH=/home/$REMOTE_USER/$BIN

echo "Building for OS: $TARGET_OS, architecture: $TARGET_ARCH ..."

GOOS=$TARGET_OS GOARCH=$TARGET_ARCH GOARM=$ARM_VERSION go build

if [ $? -eq 0 ]; then
    echo "Build successful"
else
    exit 1
fi

command -v scp >/dev/null 2>&1

if [ $? -ne 0 ]; then
    echo >&2 "scp not installed, aborting"
    rm $BIN
    exit 1
fi

echo "Securely transfering to $REMOTE_USER@$DEVICE_ADDR ..."

scp $BIN $REMOTE_USER@$DEVICE_ADDR:$REMOTE_PATH


if [ $? -eq 0 ]; then
    echo "Transfer complete"
else
    rm $BIN
    exit 1
fi

rm $BIN
echo "Done"
