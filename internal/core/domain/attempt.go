package domain

const (
	AttemptUserAllow AttemptOperation = AttemptOperation(iota)
	AttemptValidate
	AttemptDenyUserByIp
	AttemptDenyUserByCountry
)

type AttemptOperation int

type AttemptRequest struct {
	UserIpTime
	Location *string
}
