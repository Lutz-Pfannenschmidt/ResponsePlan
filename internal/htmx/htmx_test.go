package htmx_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/htmx"
)

func TestParseBody(t *testing.T) {
	in := "foo=bar&baz=qux"
	expected := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}
	actual, err := htmx.ParseBody(io.NopCloser(strings.NewReader(in)))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
		return
	}
}
