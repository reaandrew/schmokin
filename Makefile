.PHONY: test
test: shunit2-2.1.7/ shellcheck-v0.5.0/
	pip install -q --user -r requirements.txt
	./shellcheck-v0.5.0/shellcheck schmokin 
	SCHMOKIN_TEST=1 ./schmokin_test

.PHONY: test_osx
test_osx: shunit2-2.1.7/ shell_check_osx
	brew install jq
	pip install -q --user -r requirements.txt
	shellcheck schmokin 
	SCHMOKIN_TEST=1 ./schmokin_test

shunit2-2.1.7/:
	curl -s -L "https://github.com/kward/shunit2/archive/v2.1.7.tar.gz" | tar zx

shellcheck-v0.5.0/:
	curl -s -L -o /tmp/shellcheck-v0.5.0.linux.x86_64.tar.xz https://storage.googleapis.com/shellcheck/shellcheck-v0.5.0.linux.x86_64.tar.xz
	tar -xJf /tmp/shellcheck-v0.5.0.linux.x86_64.tar.xz -C ./

.PHONY: shell_check_osx
shell_check_osx:
	brew install shellcheck
