package domain

import "fmt"

var (
	ErrNotFound      = fmt.Errorf("not found")
	ErrTokenExperied = fmt.Errorf("token expried")
	ErrTokenInvalid  = fmt.Errorf("token invalid")
)
