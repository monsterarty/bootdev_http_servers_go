package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID
}

func TestHashPassword_NotEmpty(t *testing.T) {
	password := "Asdfmovie1234"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if hashedPassword == "" {
		t.Errorf(`HashPassword(password) = %q, returns empty password`, hashedPassword)
	}
}
func TestHashPassword_IsDifferent(t *testing.T) {
	password := "Asdfmovie1234"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if hashedPassword == password {
		t.Errorf(`HashPassword(password) did not hash password`)
	}
}

func TestCheckPasswordHash_CorrectPassword(t *testing.T) {
	password := "Asdfmovie1234"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	match, err := CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Errorf(`CheckPasswordHash(password) does not match`)
	}
}

func TestCheckPasswordHash_IncorrectPassword(t *testing.T) {
	password := "Asdfmovie1234"
	incorrectPassword := "AsdfMovie12345"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	match, err := CheckPasswordHash(incorrectPassword, hashedPassword)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Errorf(`CheckPasswordHash(password) accepts wrong password`)
	}
}

func TestMakeJWT_MakeAndValidate(t *testing.T) {
	userID := uuid.New()
	secret := "unittest-SECRET"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf(`MakeJWT failed: %v`, err)
	}
	userIDfromToken, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf(`ValidateJWT failed: %v`, err)
	}

	if userIDfromToken != userID {
		t.Errorf("Expected userID %v, got %v", userID, userIDfromToken)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "unittest-SECRET"
	wrongSecret := "UT-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf(`MakeJWT failed: %v`, err)
	}
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatalf(`Expected incorrect secret to fail validation: %v`, err)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "unittest-SECRET"

	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf(`MakeJWT failed: %v`, err)
	}
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf(`Expected expired token to fail validation: %v`, err)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	invalidToken := "unittest-INVALID-TOKEN"
	secret := "unittest-SECRET"

	_, err := ValidateJWT(invalidToken, secret)
	if err == nil {
		t.Fatalf("Expected invalid token to fail validation: %v", err)
	}
}
