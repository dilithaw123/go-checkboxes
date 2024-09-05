package util

import (
	"slices"
	"testing"
)

func TestEncodeSelection(t *testing.T) {
	var index uint64 = 123
	var set bool = true
	bytes := EncodeSelection(index, set)
	expected := []byte{1, 123, 0, 0, 0, 0, 0, 0, 0}
	if len(bytes) != len(expected) || !slices.Equal(bytes, expected) {
		t.Errorf("Expected %v, got %v", expected, bytes)
	}
}

func TestDecodeSelection(t *testing.T) {
	bytes := []byte{1, 123, 0, 0, 0, 0, 0, 0, 0}
	index, set := DecodeSelection(bytes)
	if index != 123 || !set {
		t.Errorf("Expected 123, true, got %d, %t", index, set)
	}
}
