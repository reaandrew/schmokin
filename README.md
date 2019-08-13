# reqit - Declarative file based HTTP client

This is a simple library which uses files to declare the HTTP request and to output information from the response.  The file is a multi-document yaml file with the first document being the HTTP request properties and the second being the content that is to be sent, no yaml properties are need for this section just the raw data e.g. JSON, XML, Form, Strings etc...

An example file:

```yaml
---
method: POST
url: https://postman-echo.com/post?foo1=bar1&foo2=bar2
headers:
  X-SOMETHING: Boom
---
{
  "a":1
}
```

An example response:

```shell
Response Info
                        Method: POST
                        Status: 200
                  Elapsed time: 396 ms
                          Size: 359 Bytes
                      Encoding: utf-8
Request Headers
                    User-Agent: python-requests/2.22.0
               Accept-Encoding: gzip, deflate
                        Accept: */*
                    Connection: keep-alive
                   X-SOMETHING: Boom
                Content-Length: 7
Response Headers
              Content-Encoding: gzip
                  Content-Type: application/json; charset=utf-8
                          Date: Thu, 13 Jun 2019 08:16:04 GMT
                          ETag: W/"167-JoBy2tgB/RL3o7FkGQZAQ0RUsao"
                        Server: nginx
                    set-cookie: sails.sid=s%3AFCviHXQukVCG0k9xJGOz5WTwNjqgN6I9.7ZRiCnsRX63zZ%2F7ZTJNyJ63Rmdr%2F1gv0LVh6xkSIDFs; Path=/; HttpOnly
                          Vary: Accept-Encoding
                Content-Length: 244
                    Connection: keep-alive
Cookies
                     sails.sid: s%3AFCviHXQukVCG0k9xJGOz5WTwNjqgN6I9.7ZRiCnsRX63zZ%2F7ZTJNyJ63Rmdr%2F1gv0LVh6xkSIDFs
Data
b'{"args":{"foo1":"bar1","foo2":"bar2"},"data":{},"files":{},"form":{},"headers":{"x-forwarded-proto":"https","host":"postman-echo.com","content-length":"7","accept":"*/*","accept-encoding":"gzip, deflate","user-agent":"python-requests/2.22.0","x-something":"Boom","x-forwarded-port":"443"},"json":null,"url":"https://postman-echo.com/post?foo1=bar1&foo2=bar2"}''
```
