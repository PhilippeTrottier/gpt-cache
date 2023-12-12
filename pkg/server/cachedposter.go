package server

import (
	"io"
	"net/http"
	"strings"
)

type Poster interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type CachedPoster struct {
	HTTPPoster Poster
	cache      map[cacheKey][]byte
}

type cacheKey struct {
	url      string
	jsonBody string
}

func (s *CachedPoster) Initialize() {
	s.cache = make(map[cacheKey][]byte)
}

// Post should not be called concurrently.
func (s *CachedPoster) Post(url string, jsonBody string) ([]byte, error) {
	ck := cacheKey{ url, jsonBody }
	if answer, ok := s.cache[ck]; ok {
		return answer, nil
	}

	resp, err := s.HTTPPoster.Post(url, "application/json", strings.NewReader(jsonBody))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	answer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.cache[ck] = answer
	return answer, nil
}
