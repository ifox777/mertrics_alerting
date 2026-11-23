package agent

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendGauge(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{"sucess", 200, false},
		{"request error", 400, true},
		{"response error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					assert.Equal(t, r.URL.Path, "/update/gauge/Alloc/21.5")
					assert.Equal(t, r.Header.Get("Content-Type"), "text/plain")

					w.WriteHeader(tt.status)
				}))

			defer server.Close()
			address := strings.TrimPrefix(server.URL, "http://")
			err := SendGauge(address, "Alloc", 21.5)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "send gauge Alloc to "+address)
				assert.Contains(t, err.Error(), "HTTP status "+strconv.Itoa(tt.status))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendCounter(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{"sucess", 200, false},
		{"request error", 400, true},
		{"response error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					assert.Equal(t, r.URL.Path, "/update/counter/Alloc/2")
					assert.Equal(t, r.Header.Get("Content-Type"), "text/plain")

					w.WriteHeader(tt.status)
				}))
			defer server.Close()
			address := strings.TrimPrefix(server.URL, "http://")
			err := SendCounter(address, "Alloc", 2)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "send counter Alloc to "+address)
				assert.Contains(t, err.Error(), "HTTP status "+strconv.Itoa(tt.status))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendGaugeNetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	))

	address := strings.TrimPrefix(server.URL, "http://")
	server.Close()

	err := SendGauge(address, "Alloc", 21.5)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "send gauge Alloc to "+address+": ")
	var requestErr *url.Error
	require.ErrorAs(t, err, &requestErr)
}

func TestSendCounterNetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	))

	address := strings.TrimPrefix(server.URL, "http://")
	server.Close()

	err := SendCounter(address, "Alloc", 2)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "send counter Alloc to "+address+": ")
	var requestErr *url.Error
	require.ErrorAs(t, err, &requestErr)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type trackedResponseBody struct {
	io.Reader
	closeCalls int
}

func (b *trackedResponseBody) Close() error {
	b.closeCalls++
	return nil
}

func TestSendClosesResponseBody(t *testing.T) {
	senders := []struct {
		name string
		send func(string) error
	}{
		{"gauge", func(address string) error { return SendGauge(address, "Alloc", 21.5) }},
		{"counter", func(address string) error { return SendCounter(address, "PollCount", 2) }},
	}

	for _, sender := range senders {
		t.Run(sender.name, func(t *testing.T) {
			for _, status := range []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError} {
				t.Run(strconv.Itoa(status), func(t *testing.T) {
					body := &trackedResponseBody{Reader: strings.NewReader("response")}
					previousClient := http.DefaultClient
					t.Cleanup(func() { http.DefaultClient = previousClient })
					// The senders use http.Post, so replace its client for this serial test.
					http.DefaultClient = &http.Client{
						Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
							return &http.Response{
								StatusCode:    status,
								Header:        make(http.Header),
								Body:          body,
								ContentLength: int64(len("response")),
								Request:       req,
							}, nil
						}),
					}

					err := sender.send("metrics.test:8080")
					// Check cleanup before a fatal assertion can stop the test.
					assert.Equal(t, 1, body.closeCalls, "response body must be closed once before returning")
					if status == http.StatusOK {
						require.NoError(t, err)
					} else {
						require.Error(t, err)
					}
				})
			}
		})
	}
}
