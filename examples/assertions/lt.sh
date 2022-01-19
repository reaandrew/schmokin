#!/usr/bin/env bash
schmokin "$URL" --jq '. | length' --lt 4
