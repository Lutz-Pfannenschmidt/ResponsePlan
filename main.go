package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/api"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/db"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/httpstring"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/logging"
	ws "github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/websocket"
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

	parser := argparse.NewParser("ResponsePlan", "A simple web application for incidence response.")

	keepData := parser.Flag("k", "keep", &argparse.Options{Help: "Save the data in a database"})
	port := parser.Int("p", "port", &argparse.Options{Help: "The port to run Responseplan on", Default: 1337})
	debugFlag := parser.Flag("d", "debug", &argparse.Options{Help: "For additional logging"})

	database := db.NewDatabase()
	if *keepData {
		database.LoadFromFile("test.json")
	}

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	logger := logging.NewLogger(*debugFlag)

	logger.Debug("Debugging enabled.")
	if *debugFlag && *port != 1337 {
		logger.Debug("Debugging overrides provided port to default (1337).")
		*port = 1337
	}

	assetsFs, err := fs.Sub(static, "web/dist/assets")
	if err != nil {
		panic(err)
	}

	wsManager := ws.NewConnectionManager()

	wsManager.On("ping", func(conn *websocket.Conn, message []byte) {
		fmt.Println("received:", string(message), "from", conn.RemoteAddr().String())
		conn.WriteMessage(websocket.TextMessage, []byte("pong"))
	})

	router := httprouter.New()
	apiManager := api.NewApiManager(database, logger)

	if *debugFlag {
		router.HandleOPTIONS = true
		router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Access-Control-Request-Method") != "" {
				// Set CORS headers
				header := w.Header()
				header.Set("Access-Control-Allow-Origin", "*")
				header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			}

			// Adjust status code to 204
			w.WriteHeader(http.StatusNoContent)
		})
	}

	// WebSocket endpoint handler
	router.GET("/ws", wsManager.HandleWebSocket)

	// API endpoint handler
	router.GET("/api/*path", apiManager.HandleApiRequest)

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
