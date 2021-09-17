package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.getUrlsInfo)

	return mux
}

func (app *application) getUrlsInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	var urls []string

	err := json.NewDecoder(r.Body).Decode(&urls)

	if err != nil {
		fmt.Println("Not valid json")
		http.Error(w, "Not valid json", 405)
		return
	}

	if len(urls) > 19 {
		fmt.Println("Too many values in json")
		app.clientError(w, http.StatusRequestEntityTooLarge)
		http.Error(w, "Too many values in json", 405)
		return
	}
	app.infoLog.Println(len(urls))

	time.Sleep(10 * time.Second)

	ans := make(map[string]interface{})

	for _, url := range urls {
		req, err := http.Get(url)
		if err != nil {
			app.serverError(w, err)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			app.serverError(w, err)
			return
		}

		ans[url] = string(body)
	}

	bytesRepresentation, err := json.Marshal(ans)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytesRepresentation)
}
