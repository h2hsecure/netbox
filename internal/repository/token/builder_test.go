package token_test

import (
	"testing"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/repository/token"
)

func Test_TokenValidation(t *testing.T) {
	tokenstr, err := token.CreateToken("testuserId", "1.1.1.1", time.Duration(3600))

	if err != nil {
		t.Fatal(err)
	}

	claim, err := token.VerifyToken(tokenstr)

	if err != nil {
		t.Fatal(err)
	}

	if claim == nil {
		t.Fatalf("claim is nil")
	}

}
