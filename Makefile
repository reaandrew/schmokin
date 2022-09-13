CHEAT_ARG := $(shell :> context)
OSFLAG 				:=
ifeq ($(OS),Windows_NT)
	OSFLAG += -D WIN32
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		OSFLAG += -D AMD64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		OSFLAG += -D IA32
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OSFLAG += -D LINUX
		SHELLCHECK_INSTALL=curl -s -L -o /tmp/shellcheck-v0.8.0.linux.x86_64.tar.xz https://github.com/koalaman/shellcheck/releases/download/v0.8.0/shellcheck-v0.8.0.linux.x86_64.tar.xz && tar -xJf /tmp/shellcheck-v0.8.0.linux.x86_64.tar.xz -C ./
	endif
	ifeq ($(UNAME_S),Darwin)
		OSFLAG += -D OSX
		SHELLCHECK_INSTALL=curl -s -L -o /tmp/shellcheck-v0.8.0.darwin.x86_64.tar.xz  https://github.com/koalaman/shellcheck/releases/download/v0.8.0/shellcheck-v0.8.0.darwin.x86_64.tar.xz && tar -xJf /tmp/shellcheck-v0.8.0.darwin.x86_64.tar.xz -C ./
	endif
		UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		OSFLAG += -D AMD64
	endif
		ifneq ($(filter %86,$(UNAME_P)),)
	OSFLAG += -D IA32
		endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		OSFLAG += -D ARM
	endif
endif

.PHONY: test
test: shunit2-2.1.7/ lint
	pip3 install -q --user -r requirements.txt
	SCHMOKIN_TEST=1 ./schmokin_test

.PHONY: lint
lint: shellcheck-v0.8.0/
	find ./ -name *.sh -and -not -path "*shunit*" -exec ./shellcheck-v0.8.0/shellcheck -x {} \;
	

shunit2-2.1.7/:
	curl -s -L "https://github.com/kward/shunit2/archive/v2.1.7.tar.gz" | tar zx

shellcheck-v0.8.0/:
	$(SHELLCHECK_INSTALL)

.PHONY: compress
compress:
	tar -czvf schmokin.tar.gz libs schmokin.format schmokin
