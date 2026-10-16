package auth

import (
	"testing"
	"time"
)

func TestJWTRoundTrip(t *testing.T) {
	svc := NewJWTService("test-secret", time.Hour)

	token, err := svc.Issue(42)
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	userID, err := svc.Verify(token)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if userID != 42 {
		t.Errorf("Verify returned userID %d, want 42", userID)
	}
}

func TestJWTVerifyRejectsTamperedToken(t *testing.T) {
	svc := NewJWTService("test-secret", time.Hour)

	token, err := svc.Issue(1)
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	tampered := token[:len(token)-1] + "x"
	if tampered == token {
		t.Fatal("tampering did not change the token")
	}

	if _, err := svc.Verify(tampered); err == nil {
		t.Error("Verify accepted a tampered token")
	}
}

func TestJWTVerifyRejectsExpiredToken(t *testing.T) {
	svc := NewJWTService("test-secret", -time.Hour) // already expired by the time it's issued

	token, err := svc.Issue(1)
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	if _, err := svc.Verify(token); err == nil {
		t.Error("Verify accepted an expired token")
	}
}
