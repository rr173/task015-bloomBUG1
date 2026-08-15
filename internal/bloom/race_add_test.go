package bloom

import (
	"fmt"
	"sync"
	"testing"
)

// TestConcurrentAddNoRace verifies that concurrent Add calls do not trigger
// data races and that all added elements are found by Test (no false negatives).
func TestConcurrentAddNoRace(t *testing.T) {
	f, err := New(10000, 0.01)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			item := fmt.Sprintf("item-%d", id)
			f.Add(item)
		}(i)
	}
	wg.Wait()

	// Verify no false negatives: every added element must be found.
	for i := 0; i < goroutines; i++ {
		item := fmt.Sprintf("item-%d", i)
		if !f.Test(item) {
			t.Errorf("false negative for %q after concurrent Add", item)
		}
	}
}
