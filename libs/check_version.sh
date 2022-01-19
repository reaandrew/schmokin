#!/usr/bin/env bash

targetDirectory=${targetDirectory:-~/.schmokin}

duration="$(($(date +%s) - $(cat "$targetDirectory"/timestamp)))"
if [ "$duration" -gt "172800" ];
then
    date +%s > "$targetDirectory"/timestamp
    LATEST_VERSION=$(curl -s https://github.com/reaandrew/schmokin/releases/latest 2>&1 | grep -oP "[0-9.]{2,}")
    set +e
    vercomp "$VERSION" "$LATEST_VERSION"
    VER_DIFF="$?"
    set -e
    if [ "$VER_DIFF" -eq "2" ];
    then
        echo -e "${NOTICE}Notice ${NC} A new version of schmokin is available. Visit https://schmok.in to download."
    fi
fi

