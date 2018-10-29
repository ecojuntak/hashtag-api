package cmds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/ecojuntak/hashtag-api/data"
	"github.com/gorilla/mux"
)

func AllHashtagHandler(w http.ResponseWriter, r *http.Request) {
	hashtags := data.GetAll()
	payload, _ := json.Marshal(hashtags)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
}

func SingleHashtagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	feed_ids := data.GetFeedIds(name)
	json_str, _ := json.Marshal(feed_ids)

	queryParam := url.Values{"ids": {string(json_str[:])}}

	http.Get("http://localhost:8000/feeds/hashtag?" + queryParam.Encode())

	response, _ := http.Get("http://localhost:8000/feeds/hashtag?" + queryParam.Encode())

	payload, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err.Error)
	}

	fmt.Println(response.Body.Read)

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal("Hashtag Service")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
}

func StartREST() {
	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/hashtags", AllHashtagHandler)
	r.HandleFunc("/hashtags/{name}", SingleHashtagHandler)

	fmt.Println("REST server run on localhost:8088")
	log.Fatal(http.ListenAndServe(":8088", r))
}
