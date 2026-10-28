package cli

import (
	"bytes"
	"testing"
)

// Check that the version command prints the correct version string.
func TestVersionOutput(t *testing.T) {
	want := "alert-tune " + Version + "\n"
	for _, args := range [][]string{{"version"}, {"--version"}} {
		var out bytes.Buffer
		cmd := NewRootCmd()
		cmd.SetOut(&out)
		cmd.SetArgs(args)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		if got := out.String(); got != want {
			t.Errorf("%v: got %q, want %q", args, got, want)
		}
	}
}

// Check that the version command rejects extra arguments, since it doesn't accept any.
func TestVersionRejectsArgs(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"version", "extra"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected an error for an extra argument")
	}
}
