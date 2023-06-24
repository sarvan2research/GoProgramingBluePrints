package trace

import (
	"bytes"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Return from new should not be nil")
	} else {
		tracer.Trace("Hello from trace package.")
		log.Println(buf.String())
		if buf.String() != "Hello from trace package.\n" {
			t.Errorf("Trace should not write '%s' .", buf.String())
		}
	}

}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("Something test")
}
