# Loop over the variables to get any args after the --
whitespace="[[:space:]]"
for arg in "$@";
do
    case "$arg" in
        "--")
            EXTRA_PARAMS=1
            ;;
        *)
            if [ "$EXTRA_PARAMS" -eq 1 ];
            then
                if [[ $arg =~ $whitespace ]]
                then
                    CURL_ARGS+=(\""$arg"\")
                else
                    CURL_ARGS+=("$arg")
                fi
            fi
            ;;
    esac
done

HOMEPATH="$(cd ~/.schmokin && pwd)"

CURL_ARGS+=("-v")
CURL_ARGS+=("-s")
CURL_ARGS+=("-w @$HOMEPATH/schmokin.format")
CURL_ARGS+=("-o /tmp/schmokin-response")
CURL_ARGS+=("> /tmp/schmokin-output")
CURL_ARGS+=("2>&1")
