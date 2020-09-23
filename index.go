package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/emersion/go-ical"
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
	url := "https://apps.fujisan.co.jp/desknets/cgi-bin/dneoical/dneoical.php"
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
	err = fixical(resp.Body, w)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "fix ical error:", err)
	}
}

func fixical(r io.Reader, w io.Writer) error {
	dec := ical.NewDecoder(r)
	cal, err := dec.Decode()
	if err != nil {
		return err
	}

	newcal := ical.NewCalendar()
	newcal.Props.SetText(ical.PropVersion, cal.Props.Get(ical.PropVersion).Value)
	newcal.Props.SetText(ical.PropProductID, cal.Props.Get(ical.PropProductID).Value)

	for _, event := range cal.Events() {
		event.Props.SetText(ical.PropUID, fmt.Sprint(time.Now().UnixNano()))
		event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())
		newcal.Children = append(newcal.Children, event.Component)
	}
	return ical.NewEncoder(w).Encode(newcal)
}
