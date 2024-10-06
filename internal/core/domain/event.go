package domain

type UserIpTime struct {
	Ip        string `json:"ip"`
	User      string `json:"user"`
	Timestamp int32  `json:"timestamp"`
}
