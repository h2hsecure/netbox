package service_test

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestOpenSession_ok(t *testing.T) {
	RegisterTestingT(t)

	token, err := testService.OpenSession(testCtx, mockUserIpTime)

	Expect(err).To(BeNil())
	Expect(token).To(HaveLen(192))
}
