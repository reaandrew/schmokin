#!/usr/bin/env bash
schmokin "$URL" --jq '.status' --eq "UP"
