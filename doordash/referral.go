package doordash

import (
	"bytes"
	"encoding/json"
	"fmt"
	http "github.com/saucesteals/fhttp"
)

func (a *Account) CreateReferral() error {

	payload, err := json.Marshal(GQLPayload{
		OperationName: "createReferral",
		Variables: CreateReferralVars{
			ReferralLinkId: a.ReferralID,
			MfaAction:      "consumer_referral",
		},
		Query: CreateReferralQuery,
	})

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

	data := CreateReferralResponse{}

	if err != json.NewDecoder(res.Body).Decode(&data) {
		return err
	}

	if data.Data.CreateReferral.Referral.ReferralStatus != "REFERRAL_STATUS_PENDING" {
		return fmt.Errorf("doordash: invalid referral status: %q", data.Data.CreateReferral.Referral.ReferralStatus)
	}

	return nil
}

func (a *Account) GetSelfReferral() (string, error) {

	gqlPayload := GQLPayload{
		OperationName: "getReferralLinkForSender",
		Variables:     map[string]string{},
		Query:         GetReferralLinkQuery,
	}

	payload, err := json.Marshal(gqlPayload)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "https://www.doordash.com/graphql", bytes.NewReader(payload))

	if err != nil {
		return "", err
	}

	req.Header.Add("content-type", "application/json")

	res, err := a.client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if err != nil {
		return "", err
	}

	data := GetSelfReferralResponse{}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Data.GetReferralLinkForSender.ReferralLink.Id == "" {
		return "", fmt.Errorf("login: invalid self referral response: %v", data)
	}

	return data.Data.GetReferralLinkForSender.ReferralLink.Id, nil
}

type GetSelfReferralResponse = GraphQLResponse[struct {
	GetReferralLinkForSender struct {
		ReferralLink struct {
			Id string `json:"id"`
		} `json:"referralLink"`
	} `json:"getReferralLinkForSender"`
}]

type CreateReferralResponse = GraphQLResponse[struct {
	CreateReferral struct {
		Referral struct {
			ReferralStatus string `json:"referralStatus"`
		} `json:"referral"`
	} `json:"createReferral"`
}]

type CreateReferralVars struct {
	ReferralLinkId string `json:"referralLinkId"`
	MfaAction      string `json:"mfaAction"`
}
