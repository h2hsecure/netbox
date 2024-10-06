package domain

type SessionCliam struct {
	UserId string `json:"userId"`
	Ip     string `json:"ip"`
}

func WithDefaultCliam(userId, ip string) SessionCliam {
	return SessionCliam{
		UserId: userId,
		Ip:     ip,
	}
}
