package domain

import (
	"fmt"
	"strings"
)

type (
	CountryPolicy          string
	CountryPolicyOperation int
)

const (
	CountryPolicyAll = "*"

	CountryPolicyOperationNoop CountryPolicyOperation = CountryPolicyOperation(iota)
	CountryPolicyOperationAllow
	CountryPolicyOperationDeny
)

func (cp CountryPolicy) Parse() (map[string]CountryPolicyOperation, error) {
	retMap := make(map[string]CountryPolicyOperation)

	// if there is no policy load empty list
	if len(cp) == 0 {
		retMap[CountryPolicyAll] = CountryPolicyOperationAllow
		return retMap, nil
	}
	parts := strings.Split(string(cp), ",")

	for _, part := range parts {
		localAndOp := strings.Split(part, ":")
		if len(localAndOp) < 2 {
			return nil, fmt.Errorf("format of country policy is wrong: %s", cp)
		}
		switch localAndOp[1] {
		case "noop":
			retMap[localAndOp[0]] = CountryPolicyOperationNoop
		case "allow":
			retMap[localAndOp[0]] = CountryPolicyOperationAllow
		case "deny":
			retMap[localAndOp[0]] = CountryPolicyOperationDeny
		}
	}

	return retMap, nil
}
