package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/gpt-cache/pkg/api"
	"github.com/gpt-cache/pkg/caching"
)

type impl struct {
	url string
	cp  *caching.CachedPoster
}

func (i *impl) PostForward(w http.ResponseWriter, r *http.Request) {
	jsonBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rawJson, err := i.cp.Post(i.url, string(jsonBody))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // todo forward error code
		return
	}

	w.WriteHeader(http.StatusOK) // todo return 201 if not in cache
	w.Write(rawJson)
}

func main() {
	url := flag.String("url", "https://postman-echo.com/post", "url to forward to")
	port := flag.Uint("port", 8080, "port to listen on")
	flag.Parse()

	cp := caching.NewCachedPoster(new(http.Client))

	h := api.Handler(&impl{*url, cp})
	http.Handle("/", h)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
