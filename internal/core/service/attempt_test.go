package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/h2hsecure/netbox/internal/core/domain"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
)

var mockUserIpTime = domain.UserIpTime{
	Ip:   "1.1.1.1",
	Path: "/index.html",
}

func TestNginxAttempt_simple_request(t *testing.T) {
	RegisterTestingT(t)

	op := testService.AccessAtempt(testCtx, "", domain.AttemptRequest{
		UserIpTime: mockUserIpTime,
	},
	)
	Expect(op).To(Equal(domain.AttemptValidate))

	count, err := mockCache.Get(testCtx, "c"+mockUserIpTime.Ip)
	Expect(err).To(BeNil())
	Expect(count).To(Equal("1"))
	mockCache.Clear()
}

func TestNginxAttempt_noop_country(t *testing.T) {
	RegisterTestingT(t)

	op := testService.AccessAtempt(testCtx, "", domain.AttemptRequest{
		UserIpTime: mockUserIpTime,
		Location:   lo.ToPtr("TR"),
	},
	)
	Expect(op).To(Equal(domain.AttemptUserAllow))
	mockCache.Clear()
}

func TestNginxAttempt_allow_country(t *testing.T) {
	RegisterTestingT(t)
	op := testService.AccessAtempt(testCtx, "", domain.AttemptRequest{
		UserIpTime: mockUserIpTime,
		Location:   lo.ToPtr("EN"),
	},
	)
	Expect(op).To(Equal(domain.AttemptValidate))
	mockCache.Clear()
}

func TestNginxAttempt_deny_country(t *testing.T) {
	RegisterTestingT(t)

	op := testService.AccessAtempt(testCtx, "", domain.AttemptRequest{
		UserIpTime: mockUserIpTime,
		Location:   lo.ToPtr("NL"),
	},
	)
	Expect(op).To(Equal(domain.AttemptDenyUserByCountry))
	mockCache.Clear()
}

func TestNginxAttempt_default_country(t *testing.T) {
	RegisterTestingT(t)

	op := testService.AccessAtempt(testCtx, "", domain.AttemptRequest{
		UserIpTime: mockUserIpTime,
		Location:   lo.ToPtr("CN"),
	},
	)
	Expect(op).To(Equal(domain.AttemptValidate))
	mockCache.Clear()
}

func TestNginxAttempt_deny_ip(t *testing.T) {
	RegisterTestingT(t)
	mockCache.Set(testCtx, "1.1.1.1", "block", 0*time.Second)

	op := testService.AccessAtempt(testCtx, "", domain.AttemptRequest{
		UserIpTime: mockUserIpTime,
		Location:   nil,
	},
	)
	Expect(op).To(Equal(domain.AttemptDenyUserByIp))
	mockCache.Clear()
}

func TestNginxAttempt_Validtoken(t *testing.T) {
	RegisterTestingT(t)
	userId := uuid.NewString()
	duration := time.Until(time.Now().AddDate(100, 0, 0))
	token, err := tokenService.CreateToken(userId, mockUserIpTime.Ip, duration)
	Expect(err).To(BeNil())

	op := testService.AccessAtempt(testCtx, token, domain.AttemptRequest{
		UserIpTime: mockUserIpTime,
		Location:   nil,
	},
	)
	Expect(op).To(Equal(domain.AttemptUserAllow))
	count, err := mockCache.Get(testCtx, "c"+mockUserIpTime.Ip)
	Expect(err).To(BeNil())
	Expect(count).To(Equal("1"))

	count, err = mockCache.Get(testCtx, "c"+userId)
	Expect(err).To(BeNil())
	Expect(count).To(Equal("1"))

	mockCache.Clear()
}

var someToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI4ZmZkMDVjNS1lMmQ5LTRiMzMtYmE5ZC0xNWJhMzc4NTM0OGMiLCJleHAiOjE3NDMwNDY3NDIsInVzZXJJZCI6IiIsImlwIjoiIn0.wbPBkyMZWANMxWgLECJu_wLvLw9BTJHciyC_GBVV1F0"

func Test_Some_Token(t *testing.T) {

	op := testService.AccessAtempt(testCtx, someToken, domain.AttemptRequest{
		UserIpTime: mockUserIpTime,
		Location:   nil,
	},
	)

	if op != domain.AttemptValidate {
		t.Fatalf("it should return alllow rather than: %d", op)
	}
}
