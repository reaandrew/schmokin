# Schmokin

A simple utility which takes the output from curl and allows you to make assertion on the headers or the body.

[![asciicast](https://asciinema.org/a/u2mdeeHToo7mBdbnEBBjKAMqO.png)](https://asciinema.org/a/u2mdeeHToo7mBdbnEBBjKAMqO)

## Under the hood

- Made with bash.
- Uses [https://stedolan.github.io/jq/](jq) and the jq expression syntax to make JSON assertions.
- Pretty printed JSON output is not supported.
- Uses shunit and python to test (check out `schmokin_test`)
- Should work with HTTP/1.1 and HTTP1.0 but not HTTP/2 yet.
- The Python Web Server is currently Flask and is return HTTP/1.0 output
- Not tested with HTTPS

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

In order for schmokin to get access to all the curl output, it is required that you use a base curl of:

```
curl -vs $URL 2>&1 | schmokin 
```

- `-s` to suppress the progress bar.
- `-v` see request headers, response headers, connection info, body etc...
- `2>&1` redirect standard error to standard output.

Schmokin outputs a pretty format and returns either an exit code of 0 (PASSED) or 1 (FAILED).

## Examples

**Assert on status code**

```
curl -vs $URL 2>&1 | schmokin -s 200
```

**Assert equals on JSON output**

```
curl -vs $URL 2>&1 | schmokin --jq-expr '.status' --eq "UP"
```

**Assert greater than on JSON output**

```
curl -vs $URL 2>&1 | schmokin --jq-expr '. | length' --gt 4
```

**Assert greater than or equal on JSON output**

```
curl -vs $URL 2>&1 | schmokin --jq-expr '. | length' --ge 5
```

**Assert less than on JSON output**

```
curl -vs $URL 2>&1 | schmokin --jq-expr '. | length' --lt 6
```

**Assert less than or equal on JSON output**

```
curl -vs $URL 2>&1 | schmokin --jq-expr '. | length' --le 5
```

**Expressions can be combined and evaluated in order**

```
curl -vs $URL 2>&1 | schmokin --jq-expr '. | length' --gt 4 \ -s 200
```

**Assert equals on a Response Header**

```
curl -vs $URL 2>&1 | schmokin --resp-header "Content-Type" --eq 'application/json'
```

**Assert equals on a Request Header**

```
curl -vs $URL 2>&1 | schmokin --req-header "Content-Type" --eq 'application/json'
```

**Assert using contains**

```
curl -vs $URL 2>&1 | schmokin --resp-header "Server" --contains 'Python'
```
