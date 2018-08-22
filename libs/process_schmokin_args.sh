while [ -n "$1" ]; do
    case "$1" in
    --status)
       msg=$(printf "HTTP Status")
       RESULT=$(echo "$DATA" | sed -rn 's/.*HTTP\/[0-9.]+.*([0-9]{3}).*/\1/p' | sed 's/\r//g')
       ;;
    --jq)
        msg=$2
        RESULT=$(cat < "/tmp/schmokin-response" | jq "$2" | sed 's/\"//g')
        shift
       ;;
    --resp-body)
        msg=$2
        RESULT=$(cat < "/tmp/schmokin-response")
       ;;
    --eq)
        statement="expected ${msg:0:60} = $2 actual $RESULT"
        if [ "$RESULT" = "$2" ];
        then
         PASS "$statement" "PASS"
        else
         FAIL "$statement" "FAIL"
        fi
        shift
       ;;
    --gt)
        statement="expected ${msg:0:60} > $2 actual $RESULT"
        if [ "$RESULT" -gt "$2" ];
        then
         PASS "$statement" "PASS"
        else
         FAIL "$statement" "FAIL"
        fi
        shift
       ;;
    --ge)
        statement="expected ${msg:0:60} >= $2 actual $RESULT"
        if [ "$RESULT" -ge "$2" ];
        then
         PASS "$statement" "PASS"
        else
         FAIL "$statement" "FAIL"
        fi
        shift
       ;;
    --lt)
        statement="expected ${msg:0:60} < $2 actual $RESULT"
        if [ "$RESULT" -lt "$2" ];
        then
         PASS "$statement" "PASS"
        else
         FAIL "$statement" "FAIL"
        fi
        shift
       ;;
    --le)
        statement="expected ${msg:0:60} <= $2 actual $RESULT"
        if [ "$RESULT" -le "$2" ];
        then
         PASS "$statement" "PASS"
        else
         FAIL "$statement" "FAIL"
        fi
        shift
       ;;
    --resp-header)
        msg="response header $2"
        EXPECTED=$2
        RESULT=$(echo "$DATA" \
            | tr -d ' ' | grep "<$EXPECTED.*" | cut -d: -f2 | sed 's/\r//g')
        shift
        ;;
    --req-header)
        msg="request header $2"
        EXPECTED=$2
        RESULT=$(echo "$DATA" \
            | tr -d ' ' | grep ">$EXPECTED.*" | cut -d: -f2 | sed 's/\r//g')
        shift
        ;;
    --co)
        statement="expected ${RESULT:0:30} to contain ${2:0:30}"
        if  echo "$RESULT" | grep -q "$2";
        then
         PASS "$statement" "PASS"
        else
         FAIL "$statement" "FAIL"
        fi
        shift
       ;;
    --export)
        name=$2
        message="export $name=\"$RESULT\""
        echo "$message" >> $targetDirectory/context
        PASS "$message" "PASS"
        shift
        ;;
    --debug)
        echo "$DATA"
        ;;
    * )
        ;;
    esac
    shift
done

