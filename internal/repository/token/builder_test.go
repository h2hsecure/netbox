package token_test

import (
	"errors"
	"testing"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/repository/token"
)

const wrong_token_str = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoyNjg5MzM0NTI0fQ.vhfyD_Hx7PSf-IavHnfJvCK27A3m_StLRsO08k7FVfo"

func Test_TokenValidation(t *testing.T) {
	ts := token.NewTokenService("testSecret", time.Hour*12)
	tokenstr, err := ts.CreateToken("testuserId", "1.1.1.1", 3600*time.Second)

	if err != nil {
		t.Fatal(err)
	}

	claim, err := ts.VerifyToken(tokenstr)

	if err != nil {
		t.Fatal(err)
	}

	if claim == nil {
		t.Fatalf("claim is nil")
	}

	if claim.Subject != "testuserId" {
		t.Fatalf("subject is different than %s", "testuserId")
	}
}

func Test_TokenInValidation(t *testing.T) {
	ts := token.NewTokenService("testSecret", time.Hour*12)
	tokenstr, err := ts.CreateToken("testuserId", "1.1.1.1", 10)

	if err != nil {
		t.Fatal(err)
	}

	_, err = ts.VerifyToken(tokenstr)

	if !errors.Is(err, domain.ErrTokenExperied) {
		t.Fatalf("error should be: %v not %v", domain.ErrTokenExperied, err)
	}

	_, err = ts.VerifyToken(tokenstr + "invalid")

	if !errors.Is(err, domain.ErrTokenInvalid) {
		t.Fatalf("error should be: %v not %v", domain.ErrTokenInvalid, err)
	}

	_, err = ts.VerifyToken(wrong_token_str)

	if !errors.Is(err, domain.ErrTokenInvalid) {
		t.Fatalf("error should be: %v not %v", domain.ErrTokenInvalid, err)
	}
}
