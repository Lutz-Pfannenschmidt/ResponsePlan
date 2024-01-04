package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/db"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/db/models"
	"github.com/Ullaakut/nmap/v3"
	log "github.com/amoghe/distillog"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ApiManager struct {
	Database *db.Database
	debug    bool
}

func NewApiManager(database *db.Database, debug bool) *ApiManager {
	return &ApiManager{
		Database: database,
		debug:    debug,
	}
}

func (a *ApiManager) HandleApiRequest(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if a.debug {
		header := w.Header()
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
	}

	if strings.HasPrefix(r.URL.Path, "/api/scan") {
		id := uuid.New().String()
		w.Write([]byte(id))

		a.Database.Data[id] = models.Scan{
			Subnets:   &map[string]models.SubnetScan{},
			OSScan:    true,
			PortScan:  true,
			StartTime: time.Now().Unix(),
			EndTime:   0,
		}

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
				log.Errorln("unable to create nmap scanner: %v", err)
				return
			}

			// progress := make(chan float32, 1)

			// go func() {
			// 	for p := range progress {
			// 		fmt.Printf("Progress: %v %%\n", p)
			// 	}
			// }()

			// result, warnings, err := scanner.Progress(progress).Run()
			result, warnings, err := scanner.Run()
			if len(*warnings) > 0 {
				log.Infoln("run finished with warnings: %s\n", *warnings)
			}
			if err != nil {
				log.Errorln("unable to run nmap scan: %v", err)
				return
			}

			fmt.Printf("Nmap done: %d hosts up scanned in %.2f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)

			if entry, ok := a.Database.Data[id]; ok {
				entry.EndTime = time.Now().Unix()
				entry.Subnets = &map[string]models.SubnetScan{
					"test": {
						Subnet: "test",
						Result: result,
					},
				}

				a.Database.Data[id] = entry
			}

			err = a.Database.SaveToFile("test.json")
			if err != nil {
				panic(err)
			}

			log.Infoln(id, "Scan finished")
		}()
	}
}
