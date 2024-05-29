package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/htmx"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/scans"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/tty2web"
	"github.com/Lutz-Pfannenschmidt/yagll"
	"github.com/google/uuid"

	"github.com/akamensky/argparse"
	"github.com/julienschmidt/httprouter"
)

//go:embed templates/*
var templates embed.FS

//go:embed cdn/*
var cdnFs embed.FS
var cdn, _ = fs.Sub(cdnFs, "cdn")

var devMode = false
var scanManager = scans.NewScanManager()
var renderer *htmx.Renderer

var tty2webPath = "tty2web"
var useTty2web = true

var maxSshRetries *int

func StartScan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	yagll.Debugln("Starting scan")

	data, err := htmx.ParseBody(r.Body)
	if err != nil {
		yagll.Errorf("Error parsing body: %s", err.Error())
		panic(err)
	}

	keys := []string{"ipRange", "osDetection", "portMode", "ports"}

	yagll.Debugln("Data:")
	for _, key := range keys {
		if _, ok := data[key]; !ok {
			data[key] = ""
			yagll.Debugf("Missing key: %s", key)
		}
	}

	id := scanManager.StartScan(&scans.ScanConfig{
		Targets:  data["ipRange"],
		Ports:    scans.TransformPortRange(data["ports"]),
		OSScan:   data["osDetection"] == "true",
		TopPorts: data["portMode"] == "top",
	}, func(id uuid.UUID) {
		yagll.Debugf("Scan finished: %s", id.String())
	})

	renderer.Redirect("/graph", w, r)
	yagll.Debugf("Scan started: %s", id.String())
}

func RunningScans(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	msg := ""
	for _, scan := range scanManager.Scans {
		if scan.EndTime == 0 {
			msg += `<div class="alert alert-info !animate-none"><span class="loading loading-ring"></span><span>` + scan.Config.Targets + `:` + scan.Config.Ports + `</span></div>`
		}
	}
	w.Write([]byte(msg))
}

func MakeDeviceInfoHandler(jsonOnly bool) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uuid, err := uuid.Parse(ps.ByName("uuid"))
		if err != nil {
			yagll.Errorf("Error parsing uuid: %s", err.Error())
			return
		}

		idx, err := strconv.Atoi(ps.ByName("idx"))
		if err != nil {
			yagll.Errorf("Error parsing idx: %s", err.Error())
			return
		}

		scan, ok := scanManager.Scans[uuid]
		if !ok {
			yagll.Errorf("Scan not found: %s", uuid.String())
			return
		}

		if idx < 0 || idx >= len(scan.Result.Hosts) {
			yagll.Errorf("Index out of range: %d", idx)
			return
		}

		host := scan.Result.Hosts[idx]

		if jsonOnly {
			jsonHost, err := json.MarshalIndent(host, "", "  ")
			if err != nil {
				yagll.Errorf("Error parsing json: %s", err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonHost)
			return
		}

		jsonHost, err := json.MarshalIndent(host, "", "  ")
		if err != nil {
			yagll.Errorf("Error parsing json: %s", err.Error())
			return
		}

		data := map[string]interface{}{
			"Title": "Device Info",
			"Host":  host,
			"Json":  template.HTML(string(jsonHost)),
			"UUID":  uuid.String(),
			"IDX":   idx,
		}

		stringData, err := json.Marshal(data)
		if err != nil {
			yagll.Errorf("Error marshalling data: %s", err.Error())
			return
		}

		html := renderer.RenderComponent("deviceInfo.html", false, string(stringData))
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}
}

func addServerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the Server header
		w.Header().Set("Server", "ResponsePlan")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func StartSSH(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ip := ps.ByName("ip")
	user := ps.ByName("user")
	w.Header().Set("Content-Type", "text/html")

	yagll.Infof("Starting SSH to %s@%s", user, ip)

	port := strconv.Itoa(tty2web.StartSSH(ip, user))

	w.Write([]byte(`<meta http-equiv="refresh" content="0; url=http://127.0.0.1:` + port + `">redirecting to http://127.0.0.1:` + port))
	yagll.Infof("Started SSH on %s", port)
}

func main() {
	parser := argparse.NewParser("ResponsePlan", "A simple web application for incidence response.")

	memory := parser.Flag("m", "memory", &argparse.Options{Help: "Will disable saving data to file"})
	port := parser.Int("p", "port", &argparse.Options{Help: "The port to run Responseplan on", Default: 1337})
	devFlag := parser.Flag("d", "dev", &argparse.Options{Help: "Enable development mode (additional logging, expose to lan)"})
	expose := parser.Flag("e", "expose", &argparse.Options{Help: "Expose ResponsePlan to lan"})
	outfile := parser.String("o", "out", &argparse.Options{Help: "The file to save data to", Default: "data.responseplan"})
	infile := parser.String("i", "in", &argparse.Options{Help: "The file to load data from", Default: "data.responseplan"})
	customTty2webPath := parser.String("", "tty2web", &argparse.Options{Help: "The path to the tty2web binary", Default: "tty2web"})
	maxSshRetries = parser.Int("", "tryssh", &argparse.Options{Help: "The maximum number of retries to find a free port for SSH", Default: 10})
	// custom := parser.String("s", "scan", &argparse.Options{Help: "Provide a custon nmap command (e.g. 'ResponsePlan -s 'nmap -sS -p 80,443')"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *devFlag {
		devMode = true
		yagll.Debugln("Debug mode enabled")
	}
	yagll.Toggle(yagll.DEBUG, devMode)

	if !*memory {
		scanManager.LoadFromFile(*infile)
		yagll.Debugf("Loaded %d scans from file", len(scanManager.Scans))
		scanManager.AutoSave(*outfile)
	}

	if *customTty2webPath != tty2webPath {
		tty2webPath = *customTty2webPath
		yagll.Infof("Using custom tty2web path: %s", tty2webPath)
		if !tty2web.CheckCustomInstall(tty2webPath) {
			yagll.Errorf(yagll.Red+"Custom tty2web path not found: %s"+yagll.Reset, tty2webPath)
			useTty2web = false
		}
	} else if !tty2web.CheckInstall() {
		yagll.Errorln(yagll.Red + "NOTE: please install https://github.com/kost/tty2web to use th webSSH feature!" + yagll.Reset)
		useTty2web = false
	}

	router := httprouter.New()
	renderer = htmx.NewRenderer(&devMode, &templates, scanManager)

	// Serve the CDN
	router.ServeFiles("/cdn/*filepath", http.FS(cdn))

	// Routes
	router.GET("/", renderer.ServeRedirect("/graph"))
	router.GET("/graph", renderer.ServeComponent("Graph", "graph.html"))
	router.GET("/analytics", renderer.ServeComponent("Analytics", "analytics.html"))
	router.GET("/history", renderer.ServeComponent("History", "history.html"))

	// HTMX component routes
	router.GET("/x/", renderer.ServeComponentX("graph.html"))
	router.GET("/x/graph", renderer.ServeComponentX("graph.html"))
	router.GET("/x/analytics", renderer.ServeComponentX("analytics.html"))
	router.GET("/x/history", renderer.ServeComponentX("history.html"))
	router.GET("/x/deviceInfo/:uuid/:idx", MakeDeviceInfoHandler(false))
	router.GET("/x/runningScans", RunningScans)

	// API routes
	router.POST("/api/startScan", StartScan)
	router.GET("/api/deviceJson/:uuid/:idx", MakeDeviceInfoHandler(true))
	router.GET("/api/startSSH/:user/:ip", StartSSH)

	// Error handling
	router.NotFound = http.HandlerFunc(renderer.ServeErrorPage(404, "Page not found"))
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		yagll.Errorf("Panic: %s", err)
		renderer.ServeErrorPage(500, "Internal server error")(w, r)
	}
	router.MethodNotAllowed = http.HandlerFunc(renderer.ServeErrorPage(405, "Method not allowed"))

	yagll.Debugln("Setting up SIGINT handler")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Print("\n\033[F\033[K")
			yagll.Debugln("Received signal: " + sig.String())
			if !*memory {
				scanManager.SaveToFile(*outfile)
				yagll.Debugf("Saved %d scans to file", len(scanManager.Scans))
				yagll.Infoln("Done saving data to data.responseplan")
			}
			if useTty2web {
				tty2web.KillAll()
				yagll.Infoln("Killed all tty2web instances")
			}
			os.Exit(0)
		}
	}()

	url := "127.0.0.1:" + strconv.Itoa(*port)
	if *expose || devMode {
		url = "0.0.0.0:" + strconv.Itoa(*port)
	}
	yagll.Infof("Starting server on port %d", *port)
	yagll.Infoln(yagll.Green + "Server running on http://" + url + yagll.Reset)
	err = http.ListenAndServe(url, addServerHeader(router))
	if err != nil {
		yagll.Errorf("Error starting server: %s", err.Error())
	}
	yagll.Infoln("Shutting down")
}
