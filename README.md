# Schmokin ([https://schmok.in](https://schmok.in))

A simple utility which wraps curl and allows you to make assertions on the HTTP requests, responses, timings and other metrics provided by curl.

## Installation

```shell
curl -Ls https://github.com/reaandrew/schmokin/releases/latest/download/schmokin_install | bash
```

## Running

Schmokin requires the 1st argument to be the url, followed by an optional set of Schmokin arguments, followed by a delimited of `--` and finally followed by any optional extra curl arguments.  Schmokin literallly proxies any arguments supplied after the `--` delimiter to curl.

```shell
schmokin <url> [schmokin-args] -- [curl-args]
```
Schmokin outputs a pretty format and returns either an exit code of 0 (PASSED) or 1 (FAILED).

## Under the hood

- Made with bash.
- Uses [https://stedolan.github.io/jq/](jq) and the jq expression syntax to make JSON assertions, **jq** needs to be installed.
- Uses shunit, shellcheck and python to test (check out `schmokin_test`).
- Tested with HTTP/1.0, HTTP1.1 and HTTP/2. 
- The Python Web Server is currently Flask and is return HTTP/1.0 output.
- Currently only Linux is supported (changing use of sed to awk should make this available for OSX too).
- Assertions on the request is not supported, only the response.
- The curl argument `--next` is not currently supported.

## Testing

The tests first setup a simple webserver which is built in Python Flask.  Once all the tests have completed it uses the `shunit` `oneTimeTearDown` method to kill the test server.

```shell
make test
```

## Linting

```shell
make lint
```

## Examples

**Assert on status code**

```
schmokin $URL --status --eq 200
```

**Assert equals on JSON output**

```
schmokin $URL --jq '.status' --eq "UP"
```

**Assert greater than on JSON output**

```
schmokin $URL --jq '. | length' --gt 4
```

**Assert greater than or equal on JSON output**

```
schmokin $URL --jq '. | length' --ge 5
```

**Assert less than on JSON output**

```
schmokin $URL --jq '. | length' --lt 6
```

**Assert less than or equal on JSON output**

```
schmokin $URL --jq '. | length' --le 5
```

**Expressions can be combined and evaluated in order**

```
schmokin $URL --jq '. | length' --gt 4 -s 200
```

**Assert equals on a Response Header**

```
schmokin $URL --res-header "Content-Type" --eq 'application/json'
```

**Assert equals on a Request Header**

```
schmokin $URL --req-header "Content-Type" --eq 'application/json'
```

**Assert using contains**

```
schmokin $URL --res-header "Server" --contains 'Python'
```

**Add additional CURL arguments**

```
./schmokin $ENDPOINT/array --req-header "X-FU" --eq 'BAR' -- -H "X-FU: BAR"
```

**Export Utility**

```shell
# Store the status message of the application in a variable named appStatus
schmokin $ENDPOINT/simple --jq '.status' --export appStatus 

# Post the appStatus value and assert on the value
schmokin $ENDPOINT/echo --jq '.message' --eq 'UP' -- -X POST -d '{\"message\":\"$appStatus\"}'
```

