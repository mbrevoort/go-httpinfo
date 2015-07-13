package httpinfo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPInfo(t *testing.T) {
	durationLower := 10 * time.Millisecond
	durationUpper := 15 * time.Millisecond

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(durationLower)
		w.Write([]byte("test"))
	})
	info := New(h)
	ts := httptest.NewServer(info)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d statusCode but received %d", http.StatusOK, resp.StatusCode)
	}

	if info.Status() != http.StatusOK {
		t.Errorf("Expected %d statusCode but received %d", http.StatusOK, info.Status())
	}

	if info.Size() != 4 {
		t.Errorf("Expected Size() to be %d but was %d", 4, info.Size())
	}

	if info.Elapsed() < durationLower {
		t.Errorf("Expected Elapsed() to be greater than %d but was %d", durationLower, info.Elapsed())
	}

	if info.Elapsed() > durationUpper {
		t.Errorf("Expected Elapsed() to be less than %d but was %d", durationUpper, info.Elapsed())
	}

}
