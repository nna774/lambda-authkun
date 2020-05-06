package adapter

import (
	"fmt"
	"net/http"
)

// NewCallbackHandler is
func NewCallbackHandler(proxyPass string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(http.MethodGet, proxyPass+"?"+r.URL.RawQuery, nil)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // Don't follow
			},
		}
		req.Header.Add("Cookie", w.Header().Get("Cookie"))
		resp, err := client.Do(req)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			w.Write([]byte(fmt.Sprintf("got status: %d", resp.StatusCode)))
			return
		}

		for _, v := range resp.Header.Values("Set-Cookie") {
			w.Header().Add("Set-Cookie", v)
		}
		w.Header().Add("Location", resp.Header.Get("Location"))
		w.WriteHeader(resp.StatusCode)
	}
}
