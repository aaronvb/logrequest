package logrequest

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestToLogger(t *testing.T) {
	tables := []struct {
		statusCode               int
		method                   string
		path                     string
		expectedStartedResults   string
		expectedCompletedResults string
	}{
		{http.StatusOK, http.MethodGet, "/foo", fmt.Sprintf(`Started %s "%s"`, http.MethodGet, "/foo"), fmt.Sprintf("Completed %d in", http.StatusOK)},
		{http.StatusUnauthorized, http.MethodPost, "/bar/create", fmt.Sprintf(`Started %s "%s"`, http.MethodPost, "/bar/create"), fmt.Sprintf("Completed %d in", http.StatusUnauthorized)},
		{http.StatusNotFound, http.MethodGet, "/hello/world", fmt.Sprintf(`Started %s "%s"`, http.MethodGet, "/hello/world"), fmt.Sprintf("Completed %d in", http.StatusNotFound)},
		{http.StatusInternalServerError, http.MethodGet, "/", fmt.Sprintf(`Started %s "%s"`, http.MethodGet, "/"), fmt.Sprintf("Completed %d in", http.StatusInternalServerError)},
		{http.StatusServiceUnavailable, http.MethodPut, "/foo/update", fmt.Sprintf(`Started %s "%s"`, http.MethodPut, "/foo/update"), fmt.Sprintf("Completed %d in", http.StatusServiceUnavailable)},
	}

	for _, table := range tables {
		var str bytes.Buffer
		var logger = log.Logger{}
		logger.SetOutput(&str)

		app := &application{
			infoLog: &logger,
		}
		req, err := http.NewRequest(table.method, table.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		testHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(table.statusCode)
		})

		rr := httptest.NewRecorder()
		handler := app.logRequestToLogger(testHandler)
		handler.ServeHTTP(rr, req)

		if strings.Contains(str.String(), table.expectedStartedResults) == false {
			t.Errorf("Expected output was incorrect, %s does not contain %s", str.String(), table.expectedStartedResults)
		}

		if strings.Contains(str.String(), table.expectedCompletedResults) == false {
			t.Errorf("Expected output was incorrect, %s does not contain %s", str.String(), table.expectedCompletedResults)
		}
	}
}

func TestToString(t *testing.T) {
	tables := []struct {
		statusCode               int
		method                   string
		path                     string
		expectedStartedResults   string
		expectedCompletedResults string
	}{
		{http.StatusOK, http.MethodGet, "/foo", fmt.Sprintf(`Started %s "%s"`, http.MethodGet, "/foo"), fmt.Sprintf("Completed %d in", http.StatusOK)},
		{http.StatusUnauthorized, http.MethodPost, "/bar/create", fmt.Sprintf(`Started %s "%s"`, http.MethodPost, "/bar/create"), fmt.Sprintf("Completed %d in", http.StatusUnauthorized)},
		{http.StatusNotFound, http.MethodGet, "/hello/world", fmt.Sprintf(`Started %s "%s"`, http.MethodGet, "/hello/world"), fmt.Sprintf("Completed %d in", http.StatusNotFound)},
		{http.StatusInternalServerError, http.MethodGet, "/", fmt.Sprintf(`Started %s "%s"`, http.MethodGet, "/"), fmt.Sprintf("Completed %d in", http.StatusInternalServerError)},
		{http.StatusServiceUnavailable, http.MethodPut, "/foo/update", fmt.Sprintf(`Started %s "%s"`, http.MethodPut, "/foo/update"), fmt.Sprintf("Completed %d in", http.StatusServiceUnavailable)},
	}

	for _, table := range tables {
		var str bytes.Buffer
		var logger = log.Logger{}
		logger.SetOutput(&str)

		app := &application{
			infoLog: &logger,
		}
		req, err := http.NewRequest(table.method, table.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		testHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(table.statusCode)
		})

		rr := httptest.NewRecorder()
		handler := app.logRequestToString(testHandler)
		handler.ServeHTTP(rr, req)

		if strings.Contains(str.String(), table.expectedStartedResults) == false {
			t.Errorf("Expected output was incorrect, %s does not contain %s", str.String(), table.expectedStartedResults)
		}

		if strings.Contains(str.String(), table.expectedCompletedResults) == false {
			t.Errorf("Expected output was incorrect, %s does not contain %s", str.String(), table.expectedCompletedResults)
		}
	}
}

// Helpers

type application struct {
	infoLog *log.Logger
}

func (app *application) logRequestToLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lr := LogRequest{Request: r, Writer: w, Handler: next}
		lr.ToLogger(app.infoLog)
	})
}

func (app *application) logRequestToString(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lr := LogRequest{Request: r, Writer: w, Handler: next}
		app.infoLog.Println(lr.ToString()["started"])
		app.infoLog.Println(lr.ToString()["completed"])
	})
}
