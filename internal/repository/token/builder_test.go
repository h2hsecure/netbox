package token_test

import (
	"testing"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/repository/token"
)

func Test_TokenValidation(t *testing.T) {
	ts := token.NewTokenService("testSecret", time.Hour*12)
	tokenstr, err := ts.CreateToken("testuserId", "1.1.1.1", time.Duration(3600))

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
