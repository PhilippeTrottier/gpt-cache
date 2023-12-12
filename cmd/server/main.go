package main

import (
	"io"
	"net/http"

	"github.com/gpt-cache/pkg/api"
	"github.com/gpt-cache/pkg/server"
)

type impl struct {
	url string
	cp  *server.CachedPoster
}

func (i *impl) PostForward(w http.ResponseWriter, r *http.Request) {
	jsonBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rawJson, err := i.cp.Post(i.url, string(jsonBody))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(rawJson)
}

func main() {
	cp := server.NewCachedPoster(new(http.Client))

	h := api.Handler(&impl{"https://postman-echo.com/post", cp})
	http.Handle("/", h)
	http.ListenAndServe(":8080", nil)
}
