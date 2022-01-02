package debug_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/celogeek/piwigo-cli/internal/debug"
	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
)

func TestHelloWorldDump(t *testing.T) {
	var test struct {
		Hello string `json:"hello"`
		World string `json:"world"`
	}
	test.Hello = "abc"
	test.World = "def"

	want := `{
  "hello": "abc",
  "world": "def"
}`
	received := debug.Dump(test)
	if received != want {
		t.Fatalf("Dump hello world failed!\nReceive:\n\"%s\"\nWant:\n\"%s\"\n", received, want)
	}
}

func TestDumpTimeResult(t *testing.T) {
	var test struct {
		CreatedAt *piwigotools.TimeResult
	}
	now := time.Now()
	tr := piwigotools.TimeResult(now)
	test.CreatedAt = &tr

	want := fmt.Sprintf(`{
  "CreatedAt": "%s"
}`, now.Format("2006-01-02 15:04:05"))

	received := debug.Dump(test)
	if received != want {
		t.Fatalf("Dump TimeResult failed!\nReceive:\n\"%s\"\nWant:\n\"%s\"\n", received, want)
	}
}

func TestDumpNullTimeResult(t *testing.T) {
	var test struct {
		CreatedAt *piwigotools.TimeResult
	}

	want := fmt.Sprint(`{
  "CreatedAt": null
}`)

	received := debug.Dump(test)
	if received != want {
		t.Fatalf("Dump TimeResult failed!\nReceive:\n\"%s\"\nWant:\n\"%s\"\n", received, want)
	}
}

func TestDumpError(t *testing.T) {
	test := math.Inf(1)
	want := ""
	received := debug.Dump(test)
	if received != want {
		t.Fatalf("Dump TimeResult failed!\nReceive:\n\"%s\"\nWant:\n\"%s\"\n", received, want)
	}
}
