#!/usr/bin/env bash
schmokin "$URL" --jq '. | length' --gt 6
