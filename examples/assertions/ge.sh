#!/usr/bin/env bash
schmokin "$URL" --jq '. | length' --ge 6
