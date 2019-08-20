strace -f -y -e trace=write "$@" | grep boo
