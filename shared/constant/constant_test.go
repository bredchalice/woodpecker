package constant

import (
	"testing"
	"time"
)

func TestTaskTimeoutAllowsControlPlaneRecovery(t *testing.T) {
	const minimumGrace = 5 * time.Minute
	if TaskTimeout < minimumGrace {
		t.Fatalf("TaskTimeout = %s, want at least %s", TaskTimeout, minimumGrace)
	}
}
