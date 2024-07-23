package handlers

import (
	"testing"

	"github.com/jasonflorentino/go-http-server/src/lib"
)

func TestHandleEchoEmpty(t *testing.T) {
	status, body := Echo(lib.Request{Target: []string{"echo"}})
	if status != 400 || body != nil {
		t.Fatalf(`Echo with no path should return status 400 and body nil, got %d and %v`, status, body)
	}
}

func TestHandleEchoFull(t *testing.T) {
	want := "hello"
	status, body := Echo(lib.Request{Target: []string{"echo", want}})
	v, ok := body.(string)
	if !ok {
		t.Fatalf(`Echo didn't return string body, got %v of type %T`, body, body)
	}
	if status != 200 && string(v) != want {
		t.Fatalf(`Echo with %s should return status 200 and body %s, got %d and %v`, want, want, status, string(v))
	}
}
