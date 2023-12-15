package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/httpstring"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/logging"
	ws "github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/websocket"
	"github.com/Ullaakut/nmap/v3"
	"github.com/akamensky/argparse"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

//go:embed web/dist/assets*
var static embed.FS

//go:embed web/dist/index.html
var index []byte

//go:embed web/dist/favicon.ico
var favicon []byte

func main() {

	parser := argparse.NewParser("print", "Prints provided string to stdout")

	// keepData := parser.Flag("k", "keep", &argparse.Options{Help: "Save the data in a database."})
	port := parser.Int("p", "port", &argparse.Options{Help: "The port to run Responseplan on.", Default: 1337})
	debugFlag := parser.Flag("d", "debug", &argparse.Options{Help: "For additional logging."})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	fmt.Println(*port)

	logger := logging.NewLogger(*debugFlag)

	logger.Debug("Debugging enabled.")

	assetsFs, err := fs.Sub(static, "web/dist/assets")
	if err != nil {
		panic(err)
	}

	manager := ws.NewConnectionManager()

	manager.On("ping", func(conn *websocket.Conn, message []byte) {
		fmt.Println("recieved:", string(message), "from", conn.RemoteAddr().String())
		conn.WriteMessage(websocket.TextMessage, []byte("pong"))
	})

	manager.On("scan", func(conn *websocket.Conn, message []byte) {
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
		str, err := json.Marshal(result)
		if err != nil {
			log.Fatalf("unable to marshal result: %v", err)
		}

		conn.WriteMessage(websocket.TextMessage, []byte(str))

		fmt.Printf("Nmap done: %d hosts up scanned in %.2f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
	})

	router := httprouter.New()

	// WebSocket endpoint handler
	router.GET("/ws", manager.HandleWebSocket)

	// Static file server
	router.ServeFiles("/assets/*filepath", http.FS(assetsFs))
	router.GET("/", httpstring.StringHandlerRouter(index))
	router.GET("/favicon.ico", httpstring.StringHandlerRouter(favicon))

	// 404 handler
	router.NotFound = httpstring.StringHandlerFunc(index)

	// Start the HTTP server
	logger.Log("Serving at http://localhost:" + strconv.Itoa(*port) + " ...")
	http.ListenAndServe(":"+strconv.Itoa(*port), router)

}
