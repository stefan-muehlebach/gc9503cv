package main

import (
	"testing"
)

func TestOpenCursor(t *testing.T) {
	cursor := OpenCursor("cross", 32)
	if cursor.img == nil {
		t.Error("Couldn't open cursor")
	} else {
		t.Log("Cursor open")
	}
}
