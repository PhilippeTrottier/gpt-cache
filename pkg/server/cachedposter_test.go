package server_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gpt-cache/pkg/server"
	"github.com/stretchr/testify/assert"
)

type stubBodyPoster struct {
	timesCalled int
	respBody    string
}

func (m *stubBodyPoster) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	m.timesCalled++
	stubResp := &http.Response{
		Body: io.NopCloser(strings.NewReader(m.respBody)),
	}
	return stubResp, nil
}

type stubErrServerClosedPoster struct{}

func (m *stubErrServerClosedPoster) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	stubResp := &http.Response{
		Body: io.NopCloser(strings.NewReader("stub response")),
	}
	return stubResp, http.ErrServerClosed
}

type stubErrShortBufferOnBodyReadPoster struct{}

func (m *stubErrShortBufferOnBodyReadPoster) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	stubResp := &http.Response{
		Body: io.NopCloser(stubFailReader{}),
	}
	return stubResp, nil
}

type stubFailReader struct{}

func (m stubFailReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrShortBuffer
}

func TestPost_ForwardsPosterError(t *testing.T) {
	s := server.NewCachedPoster(new(stubErrServerClosedPoster))

	_, err := s.Post("http://example.com", "some json")

	assert.ErrorIs(t, err, http.ErrServerClosed)
}

func TestPost_ForwardsBodyReadError(t *testing.T) {
	s := server.NewCachedPoster(new(stubErrShortBufferOnBodyReadPoster))

	_, err := s.Post("http://example.com", "some json")

	assert.ErrorIs(t, err, io.ErrShortBuffer)
}

func TestPost_ReturnsResponseBodyAndNil_OnCacheMiss(t *testing.T) {
	poster := &stubBodyPoster{
		respBody: "stub response",
	}
	s := server.NewCachedPoster(poster)

	resp, err := s.Post("http://example.com", "some json")

	assert.NoError(t, err)
	assert.Equal(t, "stub response", string(resp))
}

func TestPost_ReturnsResponseBodyAndNil_OnCacheHit(t *testing.T) {
	poster := &stubBodyPoster{
		respBody: "stub response",
	}
	s := server.NewCachedPoster(poster)

	s.Post("http://example.com", "some json")
	resp, err := s.Post("http://example.com", "some json")

	assert.NoError(t, err)
	assert.Equal(t, "stub response", string(resp))
}

func TestPost_InvokesPosterOnlyOnceForSameJsonAndURL(t *testing.T) {
	poster := new(stubBodyPoster)
	s := server.NewCachedPoster(poster)

	s.Post("http://example.com", "some json")
	s.Post("http://example.com", "some json")
	s.Post("http://example.com", "some json")

	assert.Equal(t, 1, poster.timesCalled)
}

func TestPost_InvokesPosterWhenRequestBodyDiffers(t *testing.T) {
	poster := new(stubBodyPoster)
	s := server.NewCachedPoster(poster)

	s.Post("http://example.com", "some json")
	s.Post("http://example.com", "some other json")

	assert.Equal(t, 2, poster.timesCalled)
}

func TestPost_InvokesPosterWhenURLDiffers(t *testing.T) {
	poster := new(stubBodyPoster)
	s := server.NewCachedPoster(poster)

	s.Post("http://one.com", "some json")
	s.Post("http://two.com", "some json")

	assert.Equal(t, 2, poster.timesCalled)
}
