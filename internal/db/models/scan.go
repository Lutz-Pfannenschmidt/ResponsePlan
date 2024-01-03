package models

import "github.com/Ullaakut/nmap/v3"

type Scan struct {
	Subnets   *map[string]SubnetScan `json:"subnets"`
	OSScan    bool                   `json:"osScan"`
	PortScan  bool                   `json:"portScan"`
	StartTime int64                  `json:"startTime"`
	EndTime   int64                  `json:"endTime"`
}

type SubnetScan struct {
	Subnet string    `json:"subnet"`
	Result *nmap.Run `json:"result"`
}
