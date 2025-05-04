package password

import "testing"

func TestPassword(t *testing.T) {
	pass1 := "xvii"
	hash1, err := Generate(pass1)
	if err != nil {
		t.Fatalf("unexpected error while generating password 1: %v", err)
	}

	if err := Compare(hash1, pass1); err != nil {
		t.Fatalf("unexpected error while comparing password: %v", err)
	}

	pass2 := "zxc"
	if err := Compare(hash1, pass2); err == nil {
		t.Fatal("expected error while comparing different hashes, but got nothing")
	}
}
