#!/bin/bash

function destiny_dir()
{
    local ORIGIN_DIR="$1"
    local DESTINY_DIR=${ORIGIN_DIR/cbr/mobi}
    echo "$DESTINY_DIR"
}

function mkdir_destiny()
{
    DIR=$(destiny_dir "$1")
    mkdir -p "$DIR"
}

function convert2mobi()
{
    local ORIGIN_FILE="$1"
    local DESTINY_FILE=$(destiny_dir "$ORIGIN_FILE")
    local MOBI_FILE="${DESTINY_FILE/cbr/mobi}"
    ebook-convert "$ORIGIN_FILE" "$MOBI_FILE"
}
export -f destiny_dir
export -f mkdir_destiny
export -f convert2mobi

DEPTH=2
DEPTH_OUTPUT_FOLDER=1
ROOT_FOLDER=$1

# create output folder
find $ROOT_FOLDER -type d -mindepth $DEPTH_OUTPUT_FOLDER -maxdepth $DEPTH_OUTPUT_FOLDER \
    -exec bash -c 'mkdir_destiny "{}"' bash {} \;

# convert files
find $ROOT_FOLDER -type f -mindepth $DEPTH -maxdepth $DEPTH \
    -exec bash -c 'convert2mobi "{}"' bash {} \;
