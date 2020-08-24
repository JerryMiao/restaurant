package serve

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"restaurant/pkg"
	"strings"
)

var serachURL = "https://api.yelp.com/v3/businesses/search?latitude=@&longitude=@&radius=@"

type resultModel struct {
	Total      int
	Businesses []busines
}
type busines struct {
	Rating       int
	Price        string
	Phone        int
	ID           string
	Alias        string
	IsClosed     bool
	ReviewCount  int
	Name         string
	URL          string
	ImageURL     string
	Distance     float64
	Transactions []string
}

// GetBusinesses GetBusinesses
// Feed back corporate information within a certain distance through your own location
func GetBusinesses(w http.ResponseWriter, r *http.Request) {
	if pkg.Rds.Get(r.Header.Get("token")).Val() == "" {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "token verification failed"})
		return
	}
	r.ParseForm()
	latitude, longitude, radius := r.Form["latitude"][0], r.Form["longitude"][0], r.Form["radius"][0]
	res, err := getBusinesses(latitude, longitude, radius)
	if err != nil {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": err.Error()})
		return

	}
	SendJSON(w, http.StatusOK, SendMap{"success": true, "distance": res})
}

func getBusinesses(latitude, longitude, radius string) (*resultModel, error) {
	res := &resultModel{}
	client := &http.Client{}
	serachURL = strings.Replace(serachURL, "@", latitude, 1)
	serachURL = strings.Replace(serachURL, "@", longitude, 1)
	serachURL = strings.Replace(serachURL, "@", radius, 1)
	req, err := http.NewRequest("GET", serachURL, nil)
	if err != nil {
		return res, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	json.Unmarshal(b, &res)
	return res, nil
}
