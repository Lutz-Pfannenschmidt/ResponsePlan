package scanner

import (
	"context"
	"fmt"

	"github.com/Ullaakut/nmap/v3"
)

type scanManager struct {
	scanner *nmap.Scanner
}

func NewScanManager(options ...nmap.Option) (*scanManager, error) {
	scanner, err := nmap.NewScanner(
		context.Background(),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scanner: %w", err)
	}

	return &scanManager{
		scanner: scanner,
	}, nil
}

func (sm *scanManager) StartScan(targets []string) (*nmap.Run, error) {
	sm.scanner.SetTargets(targets...)

	result, err := sm.scanner.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run scan: %w", err)
	}

	return result, nil
}
