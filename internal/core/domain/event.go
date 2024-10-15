package domain

type UserIpTime struct {
	Ip        string `json:"ip"`
	User      string `json:"user"`
	Path      string `json:"path"`
	Timestamp int32  `json:"timestamp"`
}
