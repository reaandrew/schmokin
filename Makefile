.PHONY: test
test: shunit2-2.1.7/
	pip install --user -r requirements.txt
	shellcheck schmokin
	./schmokin_test

shunit2-2.1.7/:
	curl -L "https://github.com/kward/shunit2/archive/v2.1.7.tar.gz" | tar zx

