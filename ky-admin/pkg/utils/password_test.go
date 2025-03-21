package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	if hashedPassword == password {
		t.Error("Hashed password should not be equal to original password")
	}

	if !CheckPassword(hashedPassword, password) {
		t.Error("CheckPassword should return true for correct password")
	}

	if CheckPassword(hashedPassword, "wrongPassword") {
		t.Error("CheckPassword should return false for incorrect password")
	}
}

func TestCheckPassword(t *testing.T) {
	correctPassword := "correctPassword123"
	wrongPassword := "wrongPassword123"

	hashedPassword, err := HashPassword(correctPassword)
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	if !CheckPassword(hashedPassword, correctPassword) {
		t.Error("CheckPassword should return true for correct password")
	}

	if CheckPassword(hashedPassword, wrongPassword) {
		t.Error("CheckPassword should return false for incorrect password")
	}
}
