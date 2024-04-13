package htmx

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/scans"
	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/svg"
	"github.com/Lutz-Pfannenschmidt/yagll"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type Renderer struct {
	devMode *bool
	fs      *embed.FS
	sm      *scans.ScanManager
}

func NewRenderer(devMode *bool, fs *embed.FS, sm *scans.ScanManager) *Renderer {
	return &Renderer{
		devMode: devMode,
		fs:      fs,
		sm:      sm,
	}
}

// ServeComponent serves a component with the given title and data.
// data is optional and can be used to pass data to the component.
// The component is expected to be in the "components" directory.
func (r *Renderer) ServeComponent(title, component string, data ...map[string]any) func(w http.ResponseWriter, rq *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, rq *http.Request, ps httprouter.Params) {
		tpl, err := r.recursiveParseComponent("index.html").ParseFS(r.fs, "templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			yagll.Errorf("Error parsing component template: %s", err)
			return
		}

		buf := &bytes.Buffer{}

		data := map[string]interface{}{
			"Title":     title,
			"Component": "components/" + component,
		}

		err = tpl.Execute(buf, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			yagll.Errorf("Error executing component template: %s", err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(buf.Bytes())
	}
}

// ServeComponentX serves a component with the given data.
// data is optional and can be used to pass data to the component.
// The component is expected to be in the "components" directory.
func (r *Renderer) ServeComponentX(component string, data ...map[string]any) func(w http.ResponseWriter, rq *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, rq *http.Request, ps httprouter.Params) {
		tpl, err := r.recursiveParseComponent(component).ParseFS(r.fs, "templates/components/"+component)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			yagll.Errorf("Error parsing component template: %s", err)
			return
		}

		buf := &bytes.Buffer{}

		dat := map[string]interface{}{}
		if len(data) > 0 {
			dat = data[0]
		}

		err = tpl.Execute(buf, dat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			yagll.Errorf("Error executing component template: %s", err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(buf.Bytes())
	}
}

func (r *Renderer) ServeJSON(data interface{}) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		str, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			yagll.Errorf("Error marshalling JSON: %s", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(str)
	}
}

func (r *Renderer) ServeErrorPage(errNum int, msg string) func(w http.ResponseWriter, rq *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		tpl, err := r.recursiveParseComponent("index.html").ParseFS(r.fs, "templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			yagll.Errorf("Error parsing component template: %s", err)
			return
		}

		buf := &bytes.Buffer{}

		data := map[string]interface{}{
			"Error":        err,
			"ErrorMessage": msg,
			"Title":        "Error " + strconv.Itoa(errNum),
		}

		err = tpl.Execute(buf, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			yagll.Errorf("Error executing component template: %s", err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(buf.Bytes())
	}
}

func (r *Renderer) RenderComponent(component string, icon bool, data ...string) template.HTML {
	var dat interface{}

	component = strings.ReplaceAll(component, "components/", "")
	component = strings.ReplaceAll(component, "icons/", "")

	folder := "components"
	if icon {
		folder = "icons"
	}

	if len(data) > 0 {
		err := json.Unmarshal([]byte(data[0]), &dat)
		if err != nil {
			yagll.Errorf("Error unmarshalling component data: %s", err)
			return template.HTML("Error unmarshalling component data")
		}
	}

	tpl, err := r.recursiveParseComponent(component).ParseFS(r.fs, "templates/"+folder+"/"+component)
	if err != nil {
		yagll.Errorf("Error parsing component template: %s", err)
		return template.HTML("Error parsing component template")
	}

	buf := &bytes.Buffer{}

	err = tpl.Execute(buf, dat)
	if err != nil {
		yagll.Errorf("Error executing component template: %s", err)
		return template.HTML("Error executing component template")
	}

	return template.HTML(buf.String())
}

// RenderGraph renders a scan based on the given ID.
// If no ID is given, the latest scan is rendered.
// If multiple IDs are given, the first one is used.
func (r *Renderer) RenderGraph(id ...string) template.HTML {
	if len(id) > 0 {
		uid, err := uuid.Parse(id[0])
		if err != nil {
			yagll.Errorf("Error parsing UUID: %s", err)
			return template.HTML("Error parsing UUID")
		}

		return template.HTML(svg.OverwriteRunToSvg(r.sm, uid))
	} else {
		if len(r.sm.Scans) == 0 {
			return template.HTML("No scans")
		}

		var latestScanId uuid.UUID
		var latestScanTime int64
		for id, scan := range r.sm.Scans {
			if scan != nil && scan.EndTime > latestScanTime {
				latestScanTime = scan.EndTime
				latestScanId = id
			}
		}

		return template.HTML(svg.OverwriteRunToSvg(r.sm, latestScanId))
	}

}

func (r *Renderer) recursiveParseComponent(component string) *template.Template {

	split := strings.Split(component, "/")
	fname := "index.html"
	if len(split) >= 1 {
		fname = split[len(split)-1]
	}

	tpl := template.New(fname)
	tpl.Funcs(template.FuncMap{
		"component": func(component string, data ...string) template.HTML {
			return r.RenderComponent(component, false, data...)
		},
		"icon": func(component string, data ...string) template.HTML {
			return r.RenderComponent(component, true, data...)
		},
		"string":   mustMarshal,
		"svg":      r.RenderGraph,
		"allScans": func() map[uuid.UUID]*scans.Scan { return r.sm.Scans },
	})

	return tpl
}

func mustMarshal(data interface{}) string {
	str, err := json.Marshal(data)
	if err != nil {
		yagll.Errorf("Error marshalling data: %s", err)
		return "Error marshalling data"
	}

	return string(str)
}

func (r *Renderer) ServeRedirect(path string) func(w http.ResponseWriter, rq *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, path, http.StatusSeeOther)
	}
}

func (r *Renderer) Redirect(path string, w http.ResponseWriter, rq *http.Request) {
	http.Redirect(w, rq, path, http.StatusSeeOther)
}
