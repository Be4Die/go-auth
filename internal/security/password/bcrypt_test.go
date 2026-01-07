package password

import "testing"

func TestHashAndCompare_Success(t *testing.T) {
	s := New()
	hash, err := s.Hash("Password123!")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if err := s.Compare(hash, "Password123!"); err != nil {
		t.Fatalf("compare should succeed: %v", err)
	}
}

func TestCompare_Failure(t *testing.T) {
	s := New()
	hash, _ := s.Hash("Password123!")
	if err := s.Compare(hash, "wrong"); err == nil {
		t.Fatalf("compare should fail for wrong password")
	}
}
