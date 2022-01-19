#!/usr/bin/env bash

for arg in "$@";
do
    case "$arg" in
        -v|--version)
            echo "schmokin"
            echo 
            printf "%-20s: %s\n" "version" "$VERSION"
            printf "%-20s: %s\n" "sha1" "$CHECKSUM_SHA1"
            shift
            exit 0
        ;;
        -h|--help)
            echo "schmokin"
            echo
            echo "A wrapper for curl providing chainable assertions to create simple but powerful smoke tests all written in bash"
            echo
            echo "URL"
            echo
            echo "  https://schmok.in" 
            echo
            echo "Usage:"
            echo
            echo "  schmokin <url> [schmokin-options]... [--] [curl-options]..."
            echo "  schmokin -h | --help"
            echo "  schmokin -v | --version"
            echo
            echo "Options:"
            echo
            echo " Assertions:"
            echo
            printf "    %-25s: %s\n" "--eq" "equals comparison"
            printf "    %-25s: %s\n" "--gt" "greater than comparison"
            printf "    %-25s: %s\n" "--ge" "greater than or equals comparison"
            printf "    %-25s: %s\n" "--lt" "less than comparison"
            printf "    %-25s: %s\n" "--le" "less than or equals comparison"
            printf "    %-25s: %s\n" "--co" "contains comparison"
            echo
            echo " Extractors:"
            echo
            printf "    %-25s: %s\n" "--jq" "JQ expression extractor"
            printf "    %-25s: %s\n" "--req-header" "HTTP Request Header extractor"
            printf "    %-25s: %s\n" "--resp-header" "HTTP Response Header extractor"
            printf "    %-25s: %s\n" "--status" "HTTP status extractor"
            echo
            echo "Examples"
            echo
            printf "    %-25s: %s\n" "Assert on HTTP status" "schmokin $URL --status -eq 200"
            printf "    %-25s: %s\n" "Range Assertions" "schmokin $URL --status --gt 200 --lt 202"
            printf "    %-25s: %s\n" "Using curl args" "schmokin $URL --req-header 'X-FU' --eq 'BAR' -- -H 'X-FU: BAR'"
            exit 0
            ;;
    esac
done

