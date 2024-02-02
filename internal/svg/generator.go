package svg

import (
	"bytes"
	"html/template"
	"regexp"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/scans"
	"github.com/google/uuid"
)

var svgTemplate = `
{{ $w := mul .Amount 200 }}
{{ $h := mul .Amount 100 }}
<svg viewBox="0 0 {{ $w }}px {{ $h }}px" width="{{ $w }}px" heigth="{{ $h }}px" xmlns="http://www.w3.org/2000/svg">
{{ $routers := 0 }}
	{{ range $i, $Host := .Run.Hosts }}

		{{ $router := couldBeRouter (index $Host.Addresses 0).Addr }}
		{{ $color := "currentColor" }}
		{{ $idxNoRouter := sub $i $routers}}
		{{ $x := add 20 (mul $idxNoRouter 210) }}
		{{ $y := 70 }}

		{{ if $router }}
			{{ $color = "red" }}
			{{ $x = add 20 (mul $routers 210) }}
			{{ $y = 20 }}
			{{ $routers = inc $routers }}
		{{ end }}

		<rect x="{{ $x }}" y="{{ $y }}" width="40" height="40" rx="5" fill="{{ $color }}" />
		<text x="{{ add $x 45 }}" y="{{ add $y 16 }}" fill="currentColor" font-size="16px" font-family="monospace">{{ (index $Host.Addresses 0).Addr }}</text>
		{{ if $Host.OS.Matches }}
		<text x="{{ add $x 45 }}" y="{{ add $y 32 }}" fill="currentColor" font-size="16px" font-family="monospace">{{ (index $Host.OS.Matches 0).Name }}</text>
		{{ else }}
		<text x="{{ add $x 45 }}" y="{{ add $y 32 }}" fill="currentColor" font-size="16px" font-family="monospace">Unknown OS</text>
		{{ end }}
	{{ end }}
</svg>
`

func RunToSvg(sm *scans.ScanManager, id uuid.UUID) string {

	if sm.Scans[id].Svg != "" {
		return sm.Scans[id].Svg
	}

	tpl, err := template.New("svg").Funcs(template.FuncMap{
		"add":           add,
		"inc":           increment,
		"sub":           subtract,
		"mul":           multiply,
		"div":           divide,
		"couldBeRouter": scans.IPCouldBeRouter,
	}).Parse(svgTemplate)
	if err != nil {
		panic(err)
	}

	out := bytes.NewBuffer([]byte{})
	err = tpl.Execute(out, map[string]any{
		"Run":    sm.Scans[id].Result,
		"Amount": len(sm.Scans[id].Result.Hosts),
	})
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`\s+`)
	sm.Scans[id].Svg = re.ReplaceAllLiteralString(out.String(), " ")
	return sm.Scans[id].Svg
}

func add(in ...int) int {
	solution := 0
	for _, i := range in {
		solution += i
	}
	return solution
}

func increment(in int) int {
	return in + 1
}

func subtract(in ...int) int {
	solution := in[0]
	for _, i := range in[1:] {
		solution -= i
	}
	return solution
}

func multiply(in ...int) int {
	solution := 1
	for _, i := range in {
		solution *= i
	}
	return solution
}

func divide(in ...int) int {
	solution := in[0]
	for _, i := range in[1:] {
		solution /= i
	}
	return solution
}
