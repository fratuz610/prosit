package web

import (
	"testing"
)

func TestGenerateID(t *testing.T) {

	srcID := "node server.js"
	newID := generateID(srcID, 0)
	expectedID := "nodeserv-000"

	t.Logf("ID generated from '%s': '%s'", srcID, newID)

	if newID != expectedID {
		t.Fatalf("ID mismatch: expected '%s' and got '%s'", expectedID, newID)
	}
}
