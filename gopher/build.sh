#!/bin/bash

set -e

PROJECT_NAME=`ls src/`
GOANDROID="$HOME/src/goandroid/go/bin"

# Build the Android binary
echo "Build the Android binary..."
mkdir -p android/libs/armeabi-v7a
mkdir -p android/obj/local/armeabi-v7a
CC="$NDK_ROOT/bin/arm-linux-androideabi-gcc"
CC=$CC GOPATH="`pwd`:$GOPATH" GOROOT="" GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 $GOANDROID/go install $GOFLAGS -v -ldflags="-android -shared -extld $CC -extldflags '-march=armv7-a -mfloat-abi=softfp -mfpu=vfpv3-d16'" -tags android $@ $PROJECT_NAME
cp bin/linux_arm/$PROJECT_NAME android/libs/armeabi-v7a/lib$PROJECT_NAME.so
cp bin/linux_arm/$PROJECT_NAME android/obj/local/armeabi-v7a/lib$PROJECT_NAME.so

# Build the host binary
echo "Build the host binary..."
GOPATH="`pwd`:$GOPATH" go install $@ $PROJECT_NAME

