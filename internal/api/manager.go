package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/db"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/db/models"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/logging"
	"github.com/Ullaakut/nmap/v3"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ApiManager struct {
	Database *db.Database
	Logger   *logging.Logger
}

func NewApiManager(database *db.Database, logger *logging.Logger) *ApiManager {
	return &ApiManager{
		Database: database,
		Logger:   logger,
	}
}

func (a *ApiManager) HandleApiRequest(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if strings.HasPrefix(r.URL.Path, "/api/scan") {
		id := uuid.New().String()
		w.Write([]byte(id))

		go func() {
			scanner, err := nmap.NewScanner(
				context.Background(),
				nmap.WithTargets("scanme.nmap.org"),
				nmap.WithPorts("1-1000"),
				nmap.WithServiceInfo(),
				nmap.WithVerbosity(3),
				nmap.WithOSDetection(),
				nmap.WithFilterHost(func(h nmap.Host) bool {
					return h.Status.State != "down"
				}),
			)
			if err != nil {
				log.Fatalf("unable to create nmap scanner: %v", err)
			}

			progress := make(chan float32, 1)

			go func() {
				for p := range progress {
					fmt.Printf("Progress: %v %%\n", p)
				}
			}()

			result, warnings, err := scanner.Progress(progress).Run()
			if len(*warnings) > 0 {
				log.Printf("run finished with warnings: %s\n", *warnings)
			}
			if err != nil {
				log.Fatalf("unable to run nmap scan: %v", err)
			}

			fmt.Printf("Nmap done: %d hosts up scanned in %.2f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
			a.Database.Data[id] = models.Scan{
				Subnets: map[string]models.SubnetScan{
					"test": {
						Subnet: "test",
						Result: result,
					},
				},
				OSScan:    true,
				PortScan:  true,
				StartTime: 0,
				EndTime:   0,
			}
			err = a.Database.SaveToFile("test.json")
			if err != nil {
				panic(err)
			}

			a.Logger.Logf(id, "Scan finished")
		}()
	}
}
