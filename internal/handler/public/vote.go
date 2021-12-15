package public

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nathanhollows/Argon/internal/handler"
	"github.com/nathanhollows/Argon/internal/models"
)

// Vote stores the poll results and returns summary data
func Vote(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return handler.StatusError{Code: http.StatusMethodNotAllowed, Err: errors.New("method must be POST")}
	}

	type results struct {
		Question string
		Result   string
	}

	type counts struct {
		Leave  int64
		Repair int64
	}
	var count counts

	var res results
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return handler.StatusError{Code: http.StatusBadRequest, Err: errors.New("could not read data")}
	}

	poll := models.Poll{}
	session, err := env.Session.Get(r, "uid")
	if err != nil {
		return handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("something went wrong")}
	}

	var voted int64 = 0
	env.DB.Model(&poll).Where("user_id = ? AND poll = ?", session.Values["id"], res.Question).Count(&voted)
	if voted == 0 {
		question, _ := strconv.Atoi(res.Question)
		poll.Poll = uint(question)
		poll.Result = res.Result
		poll.UserID = fmt.Sprint(session.Values["id"])
		env.DB.Save(&poll)
	}

	env.DB.Model(&poll).Where("poll = ? AND result = 'leave'", res.Question).Count(&count.Leave)
	env.DB.Model(&poll).Where("poll = ? AND result = 'repair'", res.Question).Count(&count.Repair)

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(count)
	fmt.Fprint(w, string(response))
	return nil
}
