package handler

import (
	"fmt"
	"io"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user == "" || pass == "" {
		w.Header().Set("Content-type", "text/plan")
		w.Header().Set("WWW-Authenticate", `Basic realm="my private area"`)
		w.WriteHeader(http.StatusUnauthorized)
		http.Error(w, "Not authorized", 401)
		return
	}
	url := r.URL.Query().Get("url")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "illegal url:", url)
		return
	}
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "access error:", url)
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "inline; filename=calendar.ics")
	io.Copy(w, resp.Body)
}
