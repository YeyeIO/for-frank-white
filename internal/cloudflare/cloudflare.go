package cloudflare

import (
	"crypto/tls"
	ua "github.com/EDDYCJY/fake-useragent"
	http "github.com/saucesteals/fhttp"
)

type RoundTripper struct {
	T  http.RoundTripper
	UA string
}

func NewTlsConfig() *tls.Config {
	return &tls.Config{
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		ClientSessionCache: tls.NewLRUClientSessionCache(0),
	}
}

func NewRoundTripper(t http.RoundTripper) http.RoundTripper {
	if t == nil {
		t = &http.Transport{
			TLSClientConfig:   NewTlsConfig(),
			ForceAttemptHTTP2: true,
		}
	} else {
		t.(*http.Transport).TLSClientConfig = NewTlsConfig()
	}

	return &RoundTripper{
		T:  t,
		UA: ua.Firefox(),
	}
}

func (t *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	}

	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", t.UA)
	}

	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-fetch-dest", "document")

	if req.Header.Get("accept-encoding") == "" {
		req.Header.Set("accept-encoding", "gzip, deflate, br")
	}

	if req.Header.Get("accept-language") == "" {
		req.Header.Set("accept-language", "en-US,en;q=0.9")
	}

	return t.T.RoundTrip(req)

}
