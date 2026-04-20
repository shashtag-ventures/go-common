package executil

import (
	"context"
	"testing"
)

func TestCommandContext(t *testing.T) {
	// Let's test a common utility that exists in /bin or /usr/bin.
	cmd := CommandContext(context.Background(), "echo", "hello")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected echo to run, but got error: %v", err)
	}

	if string(output) != "hello\n" {
		t.Errorf("expected output 'hello\\n', got %q", string(output))
	}
}
