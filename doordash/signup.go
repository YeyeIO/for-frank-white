package doordash

import (
	"bytes"
	"encoding/json"
	"fmt"
	http "github.com/saucesteals/fhttp"
	"io/ioutil"
	"strings"
)

func (a *Account) postLogin(code string) error {
	postLogin := GQLPayload{
		OperationName: "postLoginQuery",
		Query:         PostLoginQuery,
		Variables: PostLoginVars{
			Action: "Signup",
			Code:   code,
		},
	}

	postLoginPayload, err := json.Marshal(postLogin)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://www.doordash.com/graphql", bytes.NewReader(postLoginPayload))

	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")

	res, err := a.client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("doordash: invalid postLogin status: %q", res.Status)
	}

	return nil
}

func (a *Account) SignUp() error {

	payload := &SignUpPayload{
		ClientID:     clientId,
		CountryCode:  a.CountryCode,
		Email:        a.Email,
		FirstName:    a.FirstName,
		LastName:     a.LastName,
		Password:     a.Password,
		PhoneNumber:  a.PhoneNumber,
		RedirectURI:  "https://www.doordash.com/post-login/",
		ResponseType: "code",
		Scope:        "*",
		State:        fmt.Sprintf("/consumer/referred/%s/?action=Login", a.ReferralID),
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://identity.doordash.com/signup", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	req.Header.Add("x-xsrf-token", a.xsrfToken)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept-language", "en-US")

	res, err := a.client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	data := &SignUpResponse{}

	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, data); err != nil {
		return err
	}

	if data.RedirectUri == "" {
		return fmt.Errorf("doordash: invalid signup response: %q", string(body))
	}

	res, err = a.client.Get(data.RedirectUri)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	code := strings.Split(strings.Split(data.RedirectUri, "code=")[1], "&")[0]

	return a.postLogin(code)
}

type SignUpPayload struct {
	ClientID     string      `json:"clientId"`
	CountryCode  string      `json:"countryCode"`
	DeviceID     interface{} `json:"deviceId"`
	Email        string      `json:"email"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Password     string      `json:"password"`
	PhoneNumber  string      `json:"phoneNumber"`
	RedirectURI  string      `json:"redirectUri"`
	ResponseType string      `json:"responseType"`
	Scope        string      `json:"scope"`
	State        string      `json:"state"`
}

type SignUpResponse struct {
	DeviceId    string `json:"deviceId"`
	RedirectUri string `json:"redirectUri"`
}

type PostLoginVars struct {
	Action string `json:"action"`
	Code   string `json:"code"`
}
