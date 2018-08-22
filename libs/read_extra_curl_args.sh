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
                    CURL_ARGS+=(\"$arg\")
                else
                    CURL_ARGS+=($arg)
                fi
            fi
            ;;
    esac
done

EXTRA_ARGS=()
EXTRA_ARGS+=("-v")
EXTRA_ARGS+=("-s")
EXTRA_ARGS+=("-o /tmp/schmokin-response")
EXTRA_ARGS+=("> /tmp/schmokin-output")
EXTRA_ARGS+=("2>&1")

CURL_ARGS=( "${EXTRA_ARGS[@]}" "${CURL_ARGS[@]}" )

for i in "${!CURL_ARGS[@]}"; do 
      VALUE="${CURL_ARGS[$i]}"
      if [ "$VALUE" = "--next" ]; then
          CURL_ARGS=( "${CURL_ARGS[@]:0:$i}" "${CURL_ARGS[@]:$i}" "-o /tmp/schmokin-response" )
      fi
done

