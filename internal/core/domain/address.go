package domain

import (
	"fmt"
	"strings"
)

type ConnectionItem struct {
	id       string
	hostname string
	raftPort string
	grpcPort string
}

func (c *ConnectionItem) GetId() string {
	return c.id
}

func (c *ConnectionItem) GrpcAddress() string {
	return c.hostname + ":" + c.grpcPort
}

func (c *ConnectionItem) RaftAddress() string {
	return c.hostname + ":" + c.raftPort
}

func ParseAddress(strs string) ([]ConnectionItem, error) {

	strList := strings.Split(strs, ",")

	var items []ConnectionItem
	for _, str := range strList {
		parseData := strings.Split(str, ":")
		if len(parseData) != 4 {
			return nil, fmt.Errorf("parsing address data failed: %s", str)
		}

		items = append(items, ConnectionItem{
			id:       parseData[0],
			hostname: parseData[1],
			raftPort: parseData[2],
			grpcPort: parseData[3],
		})
	}

	return items, nil
}
