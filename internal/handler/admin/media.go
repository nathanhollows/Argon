package admin

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nathanhollows/Argon/internal/handler"
	"github.com/nathanhollows/Argon/internal/models"
)

// Media manages assests, go figure.
// E.g. images, videos, audio
func Media(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	data := make(map[string]interface{})

	var images = []models.Media{}
	env.DB.Where("type = 'image'").Order("created_at DESC").Find(&images)
	data["images"] = images

	var audio = []models.Media{}
	env.DB.Where("type = 'audio'").Order("created_at DESC").Find(&audio)
	data["audio"] = audio

	data["title"] = "Media"

	return render(w, data, "media/index.html")
}

// Upload saves audio and images to the server
// Only accepts POST requests
func Upload(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	if r.Method != http.MethodPost {
		return handler.StatusError{Code: http.StatusMethodNotAllowed, Err: errors.New("Request must be POST")}
	}

	err := r.ParseMultipartForm(20000000)
	if err != nil {
		log.Println(err.Error())
		return handler.StatusError{Code: http.StatusBadRequest, Err: errors.New("The file could not be uploaded")}
	}

	formdata := r.MultipartForm
	files := formdata.File["file"]

	if len(files) == 0 {
		return handler.StatusError{Code: http.StatusBadRequest, Err: errors.New("No files available to read. Try again")}
	}

	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			return handler.StatusError{Code: http.StatusBadRequest, Err: errors.New("no files available to read")}
		}

		filetype := strings.Split(files[i].Header.Get("Content-Type"), "/")[0]
		format := strings.Split(files[i].Header.Get("Content-Type"), "/")[1]
		title := strings.TrimSuffix(files[i].Filename, filepath.Ext(files[i].Filename))
		filename := fmt.Sprint(time.Now().Nanosecond(), "-", files[i].Filename)
		filename = strings.Replace(filename, " ", "-", -1)
		filepath := fmt.Sprint("web/static/uploads/", filetype, "/", filename)
		data["type"] = filetype

		hash := sha256.New()
		tr := io.TeeReader(file, hash)

		if filetype != "audio" && filetype != "image" {
			return handler.StatusError{Code: http.StatusNotAcceptable, Err: errors.New("no files available to read")}
		}

		out, err := os.Create(filepath)

		defer out.Close()
		if err != nil {
			return handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("unable to create the file")}
		}

		_, err = io.Copy(out, tr) // file not files[i] !

		if err != nil {
			return handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("unable to write the file")}
		}

		filehash := fmt.Sprintf("%x", hash.Sum(nil))

		media := models.Media{}
		env.DB.Model(media).Where("hash = ?", filehash).Limit(1).Find(&media)
		if media.File == "" {
			media = models.Media{
				Title:  title,
				File:   filename,
				Type:   filetype,
				Format: format,
				Hash:   filehash,
			}
			env.DB.Create(&media)
		} else {
			path := filepath
			os.Remove(path)
		}
		data["url"] = media.URL()
		data["shortcode"] = media.Shortcode()

		if len(formdata.Value["page"]) != 0 {
			page := models.Page{}
			env.DB.Where("code = ?", formdata.Value["page"][0]).Find(&page)
			if page.Code != "" {
				if filetype == "image" {
					page.Cover = media
				}
				env.DB.Updates(&page)
			}
		}

		if filetype == "image" {
			cmd := exec.Command("convert", filename, "-resize", "576x576>", "-quality", "95", "-define", "png:compression-filter=5", "small/"+filename)
			cmd.Dir = "./web/static/uploads/image/"
			cmd.Run()
			cmd = exec.Command("convert", filename, "-resize", "1000x1000>", "-quality", "95", "-define", "png:compression-filter=5", "medium/"+filename)
			cmd.Dir = "./web/static/uploads/image/"
			cmd.Run()
			cmd = exec.Command("convert", filename, "-resize", "2000x2000>", "-quality", "95", "-define", "png:compression-filter=5", "large/"+filename)
			cmd.Dir = "./web/static/uploads/image/"
			cmd.Run()
			data["small"] = media.ImgURL("small")
			data["medium"] = media.ImgURL("mediun")
			data["large"] = media.ImgURL("large")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(data)
			fmt.Fprint(w, string(response))
			return nil
		}

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(data)
		fmt.Fprint(w, string(response))
	}
	return nil
}
