package resource

import (
	"testing"
)

func Test_parseNumframes(t *testing.T) {
	if i := parseNumframes(""); i != 1 {
		t.Error("expected 1, got", i)
	}

	if i := parseNumframes("something_f5432"); i != 5432 {
		t.Error("expected 5432, got", i)
	}
}
