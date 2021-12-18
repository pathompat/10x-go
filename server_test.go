package main

import (
	"testing"
)

func TestInput1ShouldBeDisplay1(t *testing.T) {
	v := testFunc(1)
	if 1 != v {
		t.Error("test of 1 should be 1 but have", v)
	}
}
