# aaronvb/logrequest
Package `aaronvb/logrequest` is a Go middleware log output inspired by the Ruby on Rails log output for requests. Example output:

```sh
Started GET "/" 127.0.0.1:12345 HTTP/1.1
Completed 200 in 3.7455ms
```

The output can directly sent to `log.Logger` or to a map[string]string with the key `started` and `completed`.

## Install
```sh
go get -u github.com/aaronvb/logrequest
```

## Middleware Example (using [gorilla/mux](https://github.com/gorilla/mux))
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aaronvb/logrequest"

	"github.com/gorilla/mux"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", ":8080")
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/foobar", app.foobar).Methods("GET")

	// Middleware
	r.Use(app.logRequest)

	return r
}

func (app *application) foobar(w http.ResponseWriter, r *http.Request) {
	time.Sleep(300 * time.Millisecond)
	fmt.Fprintln(w, "Hello world")
}

// Middleware

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lr := logrequest.LogRequest{Request: r, Writer: w, Handler: next}
		lr.ToLogger(app.infoLog)
	})
}
```

```sh
> go run main.go
INFO	2020/03/31 22:40:09 Starting server on :8080
INFO	2020/03/31 22:40:13 Started GET "/foobar" [::1]:55044 HTTP/1.1
INFO	2020/03/31 22:40:13 Completed 200 in 300.131639ms
INFO	2020/03/31 22:40:18 Started GET "/foobar" [::1]:55044 HTTP/1.1
INFO	2020/03/31 22:40:18 Completed 200 in 302.047625ms
```

## Showing Parameters
```sh
INFO	2020/03/31 22:40:13 Started GET "/foobar" [::1]:55044 HTTP/1.1
INFO	2020/03/31 22:40:13 Parameters: {"foo" => "bar"}
INFO	2020/03/31 22:40:13 Completed 200 in 300.131639ms
```
Check out my other middleware package to output incoming parameters, which is also influenced by the Ruby on Rails logger:  [https://github.com/aaronvb/logparams](https://github.com/aaronvb/logparams)
