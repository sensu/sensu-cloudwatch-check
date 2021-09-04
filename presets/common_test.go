package presets

import (
	"log"
	"os"
	"testing"
)

var (
	enableQuiet = false
)

func quiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	if enableQuiet {
		os.Stdout = null
		os.Stderr = null
		log.SetOutput(null)
	}
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}

func TestPresetGetMeasurementString(t *testing.T) {

}
