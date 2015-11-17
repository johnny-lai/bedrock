package itest

import (
	"net/http"
	"os"
	"testing"
)

func host() string {
	return os.Getenv("TEST_HOST")
}

func TestHealth(t *testing.T) {
	url := host() + "/api/v1/health"

	r, err := http.Get(url)
	if err != nil {
		t.Fail()
		return
	}

	if r.StatusCode != 200 {
		t.Fail()
		return
	}
}

func BenchmarkHealth(b *testing.B) {
	url := host() + "/api/v1/health"

	for n := 0; n < b.N; n++ {
		r, err := http.Get(url)
		if err != nil {
			b.Fail()
			return
		}

		if r.StatusCode != 200 {
			b.Fail()
			return
		}
	}
}
