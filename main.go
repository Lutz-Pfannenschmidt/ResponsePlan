package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/htmx"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/scans"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/svg"
	"github.com/Lutz-Pfannenschmidt/yagll"
	"github.com/google/uuid"

	"github.com/akamensky/argparse"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

//go:embed templates/*
var templates embed.FS

//go:embed cdn/*
var cdnFs embed.FS
var cdn, _ = fs.Sub(cdnFs, "cdn")

var devMode = false
var upgrader = websocket.Upgrader{}
var scanManager = scans.NewScanManager()

func attachTemplateFunctions(t *template.Template) *template.Template {
	return t.Funcs(template.FuncMap{
		"component": func(name string, data ...string) template.HTML {

			input := ""
			if len(data) >= 1 {
				input = strings.TrimSpace(data[0])
			}

			var dat interface{}

			if input != "" {
				err := json.Unmarshal([]byte(input), &dat)
				if err != nil {
					yagll.Errorf("Error parsing data: %s", err.Error())
					yagll.Errorf("Data: %s", input)
					return template.HTML("Error parsing data")
				}
			}

			nameSplit := strings.Split(name, "/")
			t, err := attachTemplateFunctions(template.New(nameSplit[len(nameSplit)-1])).ParseFS(templates, "templates/"+name)
			if err != nil {
				yagll.Errorf("Error parsing template: %s", err.Error())
				return template.HTML("Error parsing template")
			}

			var tpl bytes.Buffer

			if err := t.Execute(&tpl, dat); err != nil {
				yagll.Errorf("Error executing template: %s", err.Error())
				return template.HTML("Error executing template")
			}

			response := template.HTML(tpl.String())
			return response
		},
		"string": func(data any) string {
			text, err := json.Marshal(data)
			if err != nil {
				yagll.Errorf("Error parsing data")
				return "Error passing data"
			}
			return string(text)
		},
		"svg": func() template.HTML {
			for id := range scanManager.Scans {
				return template.HTML(svg.OverwriteRunToSvg(scanManager, id))
			}
			return template.HTML("No scans")
		},
	})
}

func getErrorHandler(error int, message string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err := attachTemplateFunctions(template.New("index.html")).ParseFS(templates, "templates/index.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}

		errorData := map[string]interface{}{
			"Error":        error,
			"ErrorMessage": message,
			"Title":        "Error " + strconv.Itoa(error),
		}

		tpl.Execute(w, errorData)
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	tpl, err := attachTemplateFunctions(template.New("index.html")).ParseFS(templates, "templates/index.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}

	testData := map[string]interface{}{
		"Name":  "Lutz Pfannenschmidt",
		"Title": "Home",
		"Test":  []int{1, 2, 3, 4, 5},
		"Dev":   devMode,
		"Next":  true,
	}

	tpl.Execute(w, testData)
}

func IndexWithComponent(title string, component string) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		tpl, err := attachTemplateFunctions(template.New("index.html")).ParseFS(templates, "templates/index.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}

		testData := map[string]interface{}{
			"Title":     title,
			"Component": component,
		}

		tpl.Execute(w, testData)
	}
}

func Component(title string, component string) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		split := strings.Split(component, "/")
		fname := split[len(split)-1]
		tpl, err := attachTemplateFunctions(template.New(fname)).ParseFS(templates, "templates/"+component)

		if err != nil {
			panic(err)
		}

		testData := map[string]interface{}{
			"Title":     title,
			"Component": component,
		}

		tpl.Execute(w, testData)
	}
}

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
		yagll.Debugf("%s: %s", key, data[key])
	}

	id := scanManager.StartScan(&scans.ScanConfig{
		Targets:  data["ipRange"],
		Ports:    scans.TransformPortRange(data["ports"]),
		OSScan:   data["osDetection"] == "true",
		TopPorts: data["portMode"] == "top",
	}, func(id uuid.UUID) {
		yagll.Debugf("Scan finished: %s", id.String())
	})

	w.Write([]byte(id.String()))
	yagll.Debugf("Scan started: %s", id.String())
}

func ScansHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := upgrader.Upgrade(w, r, nil)
	c.WriteMessage(websocket.TextMessage, []byte(`<div hx-swap-oob="afterbegin:#runningScans"><div class="alert alert-info"><span class="loading loading-ring"></span><span>New message arrived.</span></div></div>`))
	yagll.Debugf("Connection established with %s", r.RemoteAddr)
	if err != nil {
		yagll.Debugf("upgrade: %s", err)
		return
	}
	defer c.Close()
	defer yagll.Debugln("Connection closed")
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			yagll.Debugf("read: %s", err)
			break
		}
		yagll.Debugf("recv: %s", message)
		err = c.WriteMessage(mt, []byte(`<div id='content'>resp from server</div>`))
		if err != nil {
			yagll.Debugf("write: %s", err)
			break
		}
	}
}

func MakeDeviceInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	tpl, err := attachTemplateFunctions(template.New("deviceInfo.html")).ParseFS(templates, "templates/components/deviceInfo.html")
	if err != nil {
		yagll.Errorf("Error parsing template: %s", err.Error())
		return
	}

	jsonHost, err := json.Marshal(host)
	if err != nil {
		yagll.Errorf("Error parsing json: %s", err.Error())
		return
	}

	testData := map[string]interface{}{
		"Title": "Device Info",
		"Host":  host,
		"Json":  string(jsonHost),
	}

	tpl.Execute(w, testData)
}

func main() {
	parser := argparse.NewParser("ResponsePlan", "A simple web application for incidence response.")

	memory := parser.Flag("m", "memory", &argparse.Options{Help: "Will disable saving data to file"})
	port := parser.Int("p", "port", &argparse.Options{Help: "The port to run Responseplan on", Default: 1337})
	devFlag := parser.Flag("d", "dev", &argparse.Options{Help: "Enable development mode (additional logging, expose to lan)"})
	expose := parser.Flag("e", "expose", &argparse.Options{Help: "Expose ResponsePlan to lan"})
	outfile := parser.String("o", "out", &argparse.Options{Help: "The file to save data to", Default: "data.responseplan"})
	infile := parser.String("i", "in", &argparse.Options{Help: "The file to load data from", Default: "data.responseplan"})

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
	}

	router := httprouter.New()

	// Serve the CDN
	router.ServeFiles("/cdn/*filepath", http.FS(cdn))

	// Routes
	router.GET("/", IndexWithComponent("Graph", "components/graph.html"))
	router.GET("/graph", IndexWithComponent("Graph", "components/graph.html"))
	router.GET("/analytics", IndexWithComponent("Analytics", "components/analytics.html"))
	router.GET("/history", IndexWithComponent("History", "components/history.html"))

	// HTMX component routes
	router.GET("/x/", Component("Graph", "components/graph.html"))
	router.GET("/x/graph", Component("Graph", "components/graph.html"))
	router.GET("/x/analytics", Component("Analytics", "components/analytics.html"))
	router.GET("/x/history", Component("History", "components/history.html"))
	router.GET("/x/deviceInfo/:uuid/:idx", MakeDeviceInfo)

	// Websocket routes
	router.GET("/ws/scans", ScansHandler)

	// API routes
	router.POST("/api/startScan", StartScan)

	// Error handling
	router.NotFound = http.HandlerFunc(getErrorHandler(404, "Page not found"))
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		yagll.Errorf("Panic: %s", err)
		getErrorHandler(500, "Internal server error")(w, r)
	}
	router.MethodNotAllowed = http.HandlerFunc(getErrorHandler(405, "Method not allowed"))

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
			os.Exit(0)
		}
	}()

	url := "127.0.0.1:" + strconv.Itoa(*port)
	if *expose || devMode {
		url = "0.0.0.0:" + strconv.Itoa(*port)
	}
	yagll.Infof("Starting server on port %d", *port)
	yagll.Infoln(yagll.Red + "Server running on http://" + url + yagll.Reset)
	err = http.ListenAndServe(url, router)
	if err != nil {
		yagll.Errorf("Error starting server: %s", err.Error())
	}
	yagll.Infoln("Shutting down")
}
