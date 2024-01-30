package scans

import (
	"context"
	"time"

	"github.com/Lutz-Pfannenschmidt/yagll"
	"github.com/Ullaakut/nmap/v3"
	"github.com/google/uuid"
)

type Scan struct {
	Result    *nmap.Run   `json:"result"`
	StartTime int64       `json:"startTime"`
	EndTime   int64       `json:"endTime"`
	Config    *ScanConfig `json:"config"`
}

type ScanConfig struct {
	Targets  string `json:"targets"`
	Ports    string `json:"ports"`
	OSScan   bool   `json:"osScan"`
	TopPorts bool   `json:"topPorts"`
}

type ScanManager struct {
	Scans map[uuid.UUID]*Scan `json:"scans"`
}

func NewScanManager() *ScanManager {
	return &ScanManager{
		Scans: make(map[uuid.UUID]*Scan),
	}
}

func (sm *ScanManager) StartScan(config *ScanConfig, callback func(uuid.UUID)) uuid.UUID {
	id := uuid.New()

	sm.Scans[id] = &Scan{
		StartTime: time.Now().Unix(),
		EndTime:   0,
		Config:    config,
	}

	go func() {
		scanner, err := nmap.NewScanner(
			context.Background(),

			nmap.WithTargets("scanme.nmap.org"),
			nmap.WithPorts("100"),
			nmap.WithVerbosity(3),
			nmap.WithServiceInfo(),
			nmap.WithFilterHost(func(h nmap.Host) bool {
				return h.Status.State != "down"
			}),
		)
		if err != nil {
			yagll.Errorf("Error creating scanner: %s", err.Error())
		}

		result, warnings, err := scanner.Run()
		if err != nil {
			yagll.Errorf("Error running scan: %s", err.Error())
		}
		if len(*warnings) > 0 {
			yagll.Debugf("Scan finished with warnings: %s", *warnings)
		}

		sm.Scans[id].Result = result
		sm.Scans[id].EndTime = time.Now().Unix()

		callback(id)

	}()

	return id

}
