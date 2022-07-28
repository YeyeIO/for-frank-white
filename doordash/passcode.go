package doordash

import (
	"bytes"
	"encoding/json"
	"fmt"
	http "github.com/saucesteals/fhttp"
)

func (a *Account) GeneratePasscode() error {

	payload, err := json.Marshal(GQLPayload{
		OperationName: "generatePasscodeBFF",
		Variables: GeneratePasscodeVars{
			Channel:    "email",
			Action:     "consumer_referral",
			Experience: "doordash",
			Language:   "en-US",
			MfaDetail:  map[string]string{},
		},
		Query: GeneratePasscodeQuery,
	})

	if err != nil {
		return err
	}

	res, err := a.client.Post("https://www.doordash.com/graphql", "application/json", bytes.NewReader(payload))

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("doordash: bad generatepasscode status: %q", res.Status)
	}

	return nil
}

func (a *Account) VerifyPasscode(passcode string) error {

	gqlPayload := GQLPayload{
		OperationName: "verifyPasscodeBFF",
		Variables: VerifyPasscodeVars{
			Code:      passcode,
			Action:    "consumer_referral",
			MfaDetail: map[string]string{},
		},
		Query: VerifyPasscodeQuery,
	}

	payload, err := json.Marshal(gqlPayload)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://www.doordash.com/graphql", bytes.NewReader(payload))

	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")

	res, err := a.client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

type GeneratePasscodeVars struct {
	Channel    string      `json:"channel"`
	Action     string      `json:"action"`
	Experience string      `json:"experience"`
	Language   string      `json:"language"`
	MfaDetail  interface{} `json:"mfaDetail"`
}

type VerifyPasscodeVars struct {
	Code      string      `json:"code"`
	Action    string      `json:"action"`
	MfaDetail interface{} `json:"mfaDetail"`
}
