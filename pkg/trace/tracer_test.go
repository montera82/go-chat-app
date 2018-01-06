package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {

	var buff bytes.Buffer

	trace := New(&buff)

	if trace == nil {
		t.Error("trace should not be nil")
	} else {
		trace.Trace("Hello Trace.")
		if buff.String() != "Hello Trace.\n" {
			t.Errorf("Trace should not have value %s", buff.String())
		}
	}
}
