package admin

import (
	"fmt"
	"net/http"

	"github.com/nathanhollows/Argon/internal/handler"
	"github.com/nathanhollows/Argon/internal/models"
	"gorm.io/gorm/clause"
)

// DataDump generates a CSV of all scan data
func DataDump(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/csv")

	scans := []models.ScanEvent{}
	env.DB.Model(scans).Preload(clause.Associations).Find(&scans)

	fmt.Fprint(w, "time,user,code,title,userAgent\r\n")
	for _, scan := range scans {
		fmt.Fprintf(w, "\"%v\",\"%v\",\"%v\",\"%v\",\"%v\"\r\n", scan.CreatedAt.String(), scan.UserID, scan.PageCode, scan.Page.Title, scan.UserAgent)
	}

	return nil
}
