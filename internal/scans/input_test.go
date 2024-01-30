package scans_test

import (
	"testing"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/scans"
)

func TestTransformPortRange(t *testing.T) {
	in := "100~4, 2 0   0~2,3,2 0 -22"
	expected := "100-104,200-202,3,20-22"
	actual := scans.TransformPortRange(in)
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
		return
	}
}
