package helpers_test

import (
	"testing"

	"github.com/thlacroix/goadvent/helpers"
)

func TestGCD(t *testing.T) {
	if gcd := helpers.GCD(30, 20); gcd != 10 {
		t.Errorf("GCD of 30 and 20 should be 10, not %d", gcd)
	}
	if gcd := helpers.GCD(20, 30); gcd != 10 {
		t.Errorf("GCD of 20 and 30 should be 10, not %d", gcd)
	}
	if gcd := helpers.GCD(52, 0); gcd != 52 {
		t.Errorf("GCD of 52 and 0 should be 52, not %d", gcd)
	}
	if gcd := helpers.GCD(0, 52); gcd != 52 {
		t.Errorf("GCD of 0 and 52 should be 52, not %d", gcd)
	}
}

func TestAbs(t *testing.T) {
	if a := helpers.Abs(10); a != 10 {
		t.Errorf("Abs of 10 should be 10, not %d", a)
	}
	if a := helpers.Abs(-10); a != 10 {
		t.Errorf("Abs of -10 should be 10, not %d", a)
	}
	if a := helpers.Abs(0); a != 0 {
		t.Errorf("Abs of 0 should be 0, not %d", a)
	}
}
