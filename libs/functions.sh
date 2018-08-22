initialize_schmokin_files(){
    if [ ! -f "$targetDirectory/timestamp" ]; then
        date +%s > "$targetDirectory/timestamp"
    fi

    touch "$targetDirectory/context"
}

# https://stackoverflow.com/questions/4023830/how-to-compare-two-strings-in-dot-separated-version-format-in-bash
vercomp () {
    if [[ "$1" == "$2" ]]
    then
        return 0
    fi
    local IFS=.
    local i ver1=("$1") ver2=("$2")
    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done
    for ((i=0; i<${#ver1[@]}; i++))
    do
        if [[ -z ${ver2[i]} ]]
        then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]}))
        then
            return 1
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 2
        fi
    done
    return 0
}

HEADING(){
    printf "${BOLD}%s %s${NC}\\n" "$@"
    echo ""
}

PASS(){
    printf "${GREEN}%-6s${NC}: %s \\n" "${2:-PASS}" "$1" 
}

FAIL(){
    printf "${RED}%-6s${NC}: %s \\n" "${2:-FAIL}" "$1" 
    FAILED=1
}
