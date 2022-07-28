package doordash

import (
	"fmt"
	"github.com/saucesteals/doordash/internal/cloudflare"
	simplejar "github.com/saucesteals/doordash/internal/simplejar"
	http "github.com/saucesteals/fhttp"
	"net/url"
)

func newHttpClient() *http.Client {
	client := &http.Client{}

	client.Transport = cloudflare.NewRoundTripper(nil)

	return client
}

func GenerateXSRFToken() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://identity.doordash.com/auth/user/signup", nil)

	if err != nil {
		return "", err
	}

	q := url.Values{
		"client_id":     {clientId},
		"intl":          {"en-US"},
		"layout":        {"referrals_cx"},
		"meta":          {fmt.Sprintf("https://www.doordash.com/consumer/referred/")},
		"prompt":        {"none"},
		"redirect_uri":  {"https://www.doordash.com/post-login/"},
		"response_type": {"code"},
		"scope":         {"*"},
		"state":         {"/consumer/referred//?action=Login"},
	}
	req.URL.RawQuery = q.Encode()

	client := newHttpClient()

	jar := simplejar.NewJar()
	client.Jar = jar

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	xsrfToken, err := jar.Get("XSRF-TOKEN")

	if err != nil {
		return "", err
	}

	return xsrfToken.Value, nil
}
