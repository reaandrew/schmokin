#!/usr/bin/env bash
schmokin "$URL" --jq '. | length' --le 6
