# Loop over the variables to get any args after the --
whitespace="[[:space:]]"
args=$#                          # number of command line args
for (( i=1; i<=args; i+=1 ))    # loop from 1 to N (where N is number of args)
do
    arg="${!i}"

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
                    if [ "$i" -ge 1 ]
                    then
                        previous="${*:$((i-1)):1}"
                        if [ "$previous" == "-d" ] && echo $arg | grep -q "@";
                        then
                            filename="$(echo $arg | cut -c2-)"
                            newFilename="$(mktemp)"

                            cp "$filename" "$newFilename"

                            envsubst < "$filename" > "$newFilename"
                            CURL_ARGS+=("@$newFilename")
                            continue
                        fi
                    fi

                    CURL_ARGS+=("$arg")
                fi
            fi
            ;;
    esac
done

CURL_ARGS+=("-v")
CURL_ARGS+=("-s")
CURL_ARGS+=("-w @$targetDirectory/schmokin.format")
CURL_ARGS+=("-o /tmp/schmokin-response")
CURL_ARGS+=("> /tmp/schmokin-output")
CURL_ARGS+=("2>&1")
