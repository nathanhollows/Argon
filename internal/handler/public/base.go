package public

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
)

var funcs = template.FuncMap{
	"uppercase": func(v string) string {
		return strings.ToUpper(v)
	},
	"divide": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b)
	},
	"progress": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b) * 100
	},
	"add": func(a, b int) int {
		return a + b
	},
}

func parseMD(page string) template.HTML {
	md := markdown.New(
		markdown.XHTMLOutput(true),
		markdown.HTML(true),
		markdown.Breaks(true))

	return template.HTML(md.RenderToString([]byte(page)))
}

func parse(patterns ...string) *template.Template {
	patterns = append(patterns, "layout.html", "flash.html")
	for i := 0; i < len(patterns); i++ {
		patterns[i] = "web/public/" + patterns[i]
	}
	return template.Must(template.New("base").Funcs(funcs).ParseFiles(patterns...))
}

func render(w http.ResponseWriter, data map[string]interface{}, patterns ...string) error {
	w.Header().Set("Content-Type", "text/html")
	if data["title"] != nil {
		data["title"] = fmt.Sprint(data["title"]) + " | QR Trail"
	} else {
		data["title"] = "QR Trail"
	}
	err := parse(patterns...).ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return err
}
