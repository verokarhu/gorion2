package resource

import (
	"testing"
)

func Test_parseAnimationParams(t *testing.T) {
	emptykey, somekey := "", "something_f5432_d4367.png"

	if i, _ := parseAnimationParams(emptykey); i != 1 {
		t.Error("expected 1, got", i)
	}

	if _, i := parseAnimationParams(emptykey); i != 0 {
		t.Error("expected 0, got", i)
	}

	if i, _ := parseAnimationParams(somekey); i != 5432 {
		t.Error("expected 5432, got", i)
	}

	if _, i := parseAnimationParams(somekey); i != 4367 {
		t.Error("expected 5432, got", i)
	}
}
