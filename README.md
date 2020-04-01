# aaronvb/logrequest
Package `aaronvb/logrequest` is a Go middleware log output inspired by the Ruby on Rails log out for requests. Example output:

```sh
Started GET "/" 127.0.0.1:12345 HTTP/1.1
Completed 200 in 3.7455ms
```

The output can directly sent to `log.Logger` or to a map[string]string with the key `started` and `completed`.

## Install
```sh
go get -u github.com/aaronvb/logrequest
```