package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
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

	ans := make(map[string]interface{})

	var wg sync.WaitGroup
	tokens := make(chan struct{}, 4)
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			tokens <- struct{}{}
			req, err := http.Get(url)
			<-tokens
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
		}(url)
	}
	wg.Wait()

	bytesRepresentation, err := json.Marshal(ans)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytesRepresentation)
}
