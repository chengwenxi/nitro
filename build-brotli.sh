#!/bin/bash

set -e

mydir=`dirname $0`
cd "$mydir"

BUILD_WASM=false
BUILD_LOCAL=false
USE_DOCKER=false
TARGET_DIR=target/
SOURCE_DIR=brotli

usage(){
    echo "brotli builder for mantle"
    echo
    echo "usage: $0 [options]"
    echo
    echo "use one or more of:"
    echo " -w     build wasm (uses emscripten)"
    echo " -l     build local"
    echo
    echo "to avoid dependencies you might want:"
    echo " -d     build inside docker container"
    echo
    echo "Other options:"
    echo " -s     source dir default: $SOURCE_DIR"
    echo " -t     target dir default: $TARGET_DIR"
    echo " -h     help"
    echo
    echo "all relative paths are relative to script location"
}

while getopts "s:t:c:wldh" option; do
    case $option in
        h)
            usage
            exit
            ;;
        w)
            BUILD_WASM=true
            ;;
        l)
            BUILD_LOCAL=true
            ;;
        d)
            USE_DOCKER=true
            ;;
        t)
            TARGET_DIR="$OPTARG"
            ;;
        s)
            SOURCE_DIR="$OPTARG"
            ;;
    esac
done

if ! $BUILD_WASM && ! $BUILD_LOCAL; then
    usage
    exit
fi

if [ ! -d "$TARGET_DIR" ]; then
    mkdir -p "${TARGET_DIR}lib"
    ln -s "lib" "${TARGET_DIR}lib64" # Fedora build
fi
TARGET_DIR_ABS=`cd -P "$TARGET_DIR"; pwd`


if $USE_DOCKER; then
    if $BUILD_WASM; then
        DOCKER_BUILDKIT=1 docker build --target brotli-wasm-export -o type=local,dest="$TARGET_DIR_ABS" .
    fi
    if $BUILD_LOCAL; then
        DOCKER_BUILDKIT=1 docker build --target brotli-library-export -o type=local,dest="$TARGET_DIR_ABS" .
    fi
    exit 0
fi

cd "$SOURCE_DIR"
if $BUILD_WASM; then
    mkdir -p buildfiles/build-wasm
    mkdir -p buildfiles/install-wasm
    TEMP_INSTALL_DIR_ABS=`cd -P buildfiles/install-wasm; pwd`
    cd buildfiles/build-wasm
    cmake ../../ -DCMAKE_C_COMPILER=emcc -DCMAKE_BUILD_TYPE=Release -DCMAKE_C_FLAGS=-fPIC -DCMAKE_INSTALL_PREFIX="$TEMP_INSTALL_DIR_ABS" -DCMAKE_AR=`which emar` -DCMAKE_RANLIB=`which touch`
    make -j
    make install
    cp -rv "$TEMP_INSTALL_DIR_ABS/lib" "$TARGET_DIR_ABS/lib-wasm"
    cd ..
fi

if $BUILD_LOCAL; then
    mkdir -p buildfiles/build-local
    cd buildfiles/build-local
    cmake ../../ -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX="$TARGET_DIR_ABS"
    make -j
    make install
    cd ..
fi
