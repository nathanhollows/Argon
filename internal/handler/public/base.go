package public

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/nathanhollows/Argon/internal/helpers"
	"github.com/nathanhollows/Argon/internal/models"
	"gitlab.com/golang-commonmark/markdown"
	"gorm.io/gorm"
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
	"sub": func(a, b int) int {
		return a - b
	},
	"url": func(s ...string) string {
		return helpers.URL(s...)
	},
	"currentYear": func() int {
		return time.Now().UTC().Year()
	},
	"stylesheetversion": func() string {
		file, err := os.Stat("web/static/css/style.css")
		if err != nil {
			fmt.Println(err)
		}
		modifiedtime := file.ModTime().Nanosecond()
		return fmt.Sprint(modifiedtime)
	},
	"unescape": func(s string) template.HTML {
		return template.HTML(s)
	},
}

func parseMD(page string, db *gorm.DB) template.HTML {
	md := markdown.New(
		markdown.XHTMLOutput(true),
		markdown.HTML(true),
		markdown.Breaks(true))

	page = regexp.MustCompile("==(.*)==").ReplaceAllString(page, "<mark>$1</mark>")
	page = regexp.MustCompile(":::([^:::]*):::").ReplaceAllString(page, `<article>
$1</article>`)
	regMedia := regexp.MustCompile(`\[\[\w+:(\d+)\]\]`)
	mediaCodes := regMedia.FindAllStringSubmatch(page, -1)

	for _, shortcode := range mediaCodes {
		var media = models.Media{}
		db.Where("id = ?", shortcode[1]).Find(&media)
		page = strings.Replace(page, shortcode[0], string(media.ToHTML()), 1)
	}

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
	if data["siteTitle"] == nil {
		data["siteTitle"] = "QR Trail"
	}
	err := parse(patterns...).ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return err
}
