package id_test

import (
	"testing"
	"time"

	"github.com/neopilot-ai/neo/pkg/id"
)

func run(t *testing.T, genFunc func() (string, error), compareFunc func(string, string) bool, order string) {
	ids := make([]string, 100)
	for i := 0; i < 100; i++ {
		id, err := genFunc()
		if err != nil {
			t.Fatalf("Failed to generate ID at index %d: %v", i, err)
		}
		ids[i] = id
		time.Sleep(time.Millisecond)
	}

	for i := range ids {
		if i == 0 {
			continue
		}
		if !compareFunc(ids[i-1], ids[i]) {
			t.Errorf("IDs not in %s order at index %d. Previous: %s, Current: %s", order, i, ids[i-1], ids[i])
			t.FailNow()
		}
		if len(ids[i]) != id.LENGTH {
			t.Errorf("ID at index %d has incorrect length. Expected: %d, Got: %d", i, id.LENGTH, len(ids[i]))
			t.FailNow()
		}
	}
}

func TestAscending(t *testing.T) {
	run(t, id.Ascending, func(a, b string) bool { return a <= b }, "ascending")
}

func TestDescending(t *testing.T) {
	run(t, id.Descending, func(a, b string) bool { return a >= b }, "descending")
}

