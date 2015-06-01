package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/search", searchHandler)
	log.Println("Start handling http://localhost/search...")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func searchHandler(w http.ResponseWriter, req *http.Request) {
	vals := req.URL.Query()

	tlat, err := strconv.ParseFloat(vals.Get("tlat"), 64)
	if err != nil {
		log.Printf("Unable to parse tlat: %v", err)
		return
	}
	tlon, err := strconv.ParseFloat(vals.Get("tlon"), 64)
	if err != nil {
		log.Printf("Unable to parse tlon: %v", err)
		return
	}
	vlat, err := strconv.ParseFloat(vals.Get("vlat"), 64)
	if err != nil {
		log.Printf("Unable to parse vlat: %v", err)
		return
	}
	vlon, err := strconv.ParseFloat(vals.Get("vlon"), 64)
	if err != nil {
		log.Printf("Unable to parse vlon: %v", err)
		return
	}
	// fmt.Println(tlat, tlon, vlat, vlon)

	if tlat > vlat {
		tlat, vlat = vlat, tlat
	}
	if tlon > vlon {
		tlon, vlon = vlon, tlon
	}
	// fmt.Println(tlat, tlon, vlat, vlon)

	body := requestToOverpass(tlat, tlon, vlat, vlon)
	if body != nil {
		res := parseJSON(body)
		if res != nil {
			// fmt.Printf("%v", res)

			for _, v := range res {
				fmt.Fprintf(w, "%f,%f\n", v.Lat, v.Lon)
			}
		}
	}
}

func requestToOverpass(tlat, tlon, vlat, vlon float64) []byte {
	url := "http://overpass-api.de/api/interpreter"
	query := fmt.Sprintf(`[out:json];
        node["amenity"="bicycle_parking"]
            (%f,%f,%f,%f);
        out;`, tlat, tlon, vlat, vlon)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(query)))
	if err != nil {
		log.Printf("Got error from Overpass API server: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Response body reading error: %v", err)
		return nil
	}

	return body
}

type results struct {
	Arr []coord `json:"elements"`
}

type coord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func parseJSON(body []byte) []coord {
	res := results{}

	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("Parsing response from Overpass API error: %v", err)
		return nil
	}

	return res.Arr
}
