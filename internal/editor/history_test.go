package editor_test

import (
	"ash/internal/editor"
	"testing"
)

func TestHistory_BasicFlow(t *testing.T) {
	h := editor.NewHistory()

	// Initial state check
	if h.Len() != 0 {
		t.Fatalf("expected Len() 0, got %d", h.Len())
	}
	if got := h.GetHistory(); got != nil {
		t.Fatalf("expected nil GetHistory on empty history, got %s", string(got))
	}

	// Save items
	h.SaveHistory([]rune("first"))
	h.SaveHistory([]rune("second"))

	if h.Len() != 2 {
		t.Fatalf("expected Len() 2, got %d", h.Len())
	}

	// MoveUp (Backward in time)
	if got := string(h.MoveUp()); got != "second" {
		t.Errorf("expected 'second', got '%s'", got)
	}
	if got := string(h.MoveUp()); got != "first" {
		t.Errorf("expected 'first', got '%s'", got)
	}

	// MoveUp at upper boundary (should stay at oldest item)
	if got := string(h.MoveUp()); got != "first" {
		t.Errorf("expected boundary hold at 'first', got '%s'", got)
	}

	// MoveDown (Forward in time)
	if got := string(h.MoveDown()); got != "second" {
		t.Errorf("expected 'second', got '%s'", got)
	}

	// MoveDown to the active prompt (past newest entry)
	if got := h.MoveDown(); got != nil {
		t.Errorf("expected nil when reaching active prompt, got '%s'", string(got))
	}
}

func TestHistory_Deduplication(t *testing.T) {
	h := editor.NewHistory()

	h.SaveHistory([]rune("cmd"))
	h.SaveHistory([]rune("cmd")) // Duplicate, should be ignored

	if h.Len() != 1 {
		t.Fatalf("expected Len() 1 after consecutive duplicate, got %d", h.Len())
	}
}

func TestHistory_Immutability(t *testing.T) {
	h := editor.NewHistory()
	original := []rune("hello")

	h.SaveHistory(original)
	original[0] = 'X' // Mutate source slice

	got := h.MoveUp()
	if string(got) == "Xello" {
		t.Fatal("SaveHistory did not clone the input slice")
	}

	got[0] = 'Y' // Mutate returned slice
	again := h.GetHistory()
	if string(again) == "Yello" {
		t.Fatal("GetHistory/MoveUp exposed direct reference to internal slice")
	}
}

// FuzzHistory randomly executes methods to catch out-of-bounds panics and state corruption.
func FuzzHistory(f *testing.F) {
	// Seed with a random sequence of operations
	f.Add([]byte{0, 1, 2, 3, 4, 1, 2})

	f.Fuzz(func(t *testing.T, ops []byte) {
		h := editor.NewHistory()
		var shadow [][]rune // Ground truth model to verify invariant checks

		for _, op := range ops {
			action := op % 4 // 4 distinct operations

			switch action {
			case 0: // SaveHistory
				// Generate a simple payload based on byte value
				payload := []rune{'c', 'm', 'd', rune(op)}
				h.SaveHistory(payload)

				// Update shadow state
				if len(payload) > 0 {
					if len(shadow) == 0 || string(shadow[len(shadow)-1]) != string(payload) {
						entry := make([]rune, len(payload))
						copy(entry, payload)
						shadow = append(shadow, entry)
					}
				}

			case 1: // MoveUp
				res := h.MoveUp()
				if len(shadow) == 0 && res != nil {
					t.Fatalf("MoveUp returned non-nil on empty history")
				}

			case 2: // MoveDown
				res := h.MoveDown()
				if len(shadow) == 0 && res != nil {
					t.Fatalf("MoveDown returned non-nil on empty history")
				}

			case 3: // GetHistory
				h.GetHistory()
			}

			// Invariant assertions that must hold true after every operation:
			if h.Len() != len(shadow) {
				t.Fatalf("Len mismatch: got %d, expected %d", h.Len(), len(shadow))
			}

			// Check for out-of-bounds panics via GetHistory
			current := h.GetHistory()
			if current != nil && len(shadow) == 0 {
				t.Fatalf("GetHistory returned non-nil when history is empty")
			}
		}
	})
}
