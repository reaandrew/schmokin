# Schmokin

[![Build Status](https://travis-ci.org/reaandrew/schmokin.svg?branch=master)](https://travis-ci.org/reaandrew/schmokin)

A simple utility which wraps curl and allows you to make assertions on the HTTP requests, responses, timings and other metrics provided by curl.

[![asciicast](https://asciinema.org/a/u2mdeeHToo7mBdbnEBBjKAMqO.png)](https://asciinema.org/a/u2mdeeHToo7mBdbnEBBjKAMqO)

## Under the hood

- Made with bash.
- Uses [https://stedolan.github.io/jq/](jq) and the jq expression syntax to make JSON assertions, **jq** needs to be installed.
- Pretty printed JSON output is not yet supported.
- Uses shunit and python to test (check out `schmokin_test`)
    - Probably moving to a bash only web server for the fun of it.
- Should work with HTTP/1.1 and HTTP1.0 but not HTTP/2 yet.
- The Python Web Server is currently Flask and is return HTTP/1.0 output.
- Not tested with HTTPS.
- Currently requires either an Ubuntu or Debian installation.

## Testing

The tests first setup a simple webserver which is built in Python Flask.  Once all the tests have completed it uses the `shunit` `oneTimeTearDown` method to kill the test server.

```
make test
```

## Installation

```
curl -Ls https://github.com/reaandrew/schmokin/releases/download/0.1.0/schmokin_install | bash
```

## Running

Schmokin requires the 1st argument to be the url, followed by an optional set of Schmokin arguments, followed by a delimited of `--` and finally followed by any optional extra curl arguments.  Schmokin literallly proxies any arguments supplied after the `--` delimiter to curl.

```
schmokin <url> [schmokin-args] -- [curl-args]
```
Schmokin outputs a pretty format and returns either an exit code of 0 (PASSED) or 1 (FAILED).

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
schmokin $URL --resp-header "Content-Type" --eq 'application/json'
```

**Assert equals on a Request Header**

```
schmokin $URL --req-header "Content-Type" --eq 'application/json'
```

**Assert using contains**

```
schmokin $URL --resp-header "Server" --contains 'Python'
```

**Add additional CURL arguments**

```
./schmokin $ENDPOINT/array --req-header "X-FU" --eq 'BAR' -- -H "X-FU: BAR"
```
