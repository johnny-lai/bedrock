package itest

import (
  "os"
	"testing"
)

func host() string {
  return os.Getenv("TEST_HOST")
}

func TestHealth(t *testing.T) {
  url := host() + "/health"

  r, err := makeRequest("GET", url, nil)
  if err != nil {
    t.Fail()
    return
  }

  err = processResponse(r, 200)
  if err != nil {
    t.Fail()
    return
  }
}

func BenchmarkHealth(b *testing.B) {
  url := host() + "/health"

  for n := 0; n < b.N; n++ {
    r, err := makeRequest("GET", url, nil)
    if err != nil {
      b.Fail()
      return
    }

    err = processResponse(r, 200)
    if err != nil {
      b.Fail()
      return
    }
  }
}
