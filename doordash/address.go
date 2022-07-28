package doordash

import (
	"bytes"
	"encoding/json"
	"fmt"
	http "github.com/saucesteals/fhttp"
)

var (
	address = SetAddressVars{
		GooglePlaceID:    "EiM5MTEgWm9vIERyaXZlLCBMb3MgQW5nZWxlcywgQ0EsIFVTQSIuKiwKFAoSCduxBZJ2wMKAEUrNyU4GTUkEEhQKEgkT2ifcXcfCgBH0CEYlb98v4g",
		City:             "Los Angeles",
		State:            "CA",
		ZipCode:          "90027",
		Street:           "Zoo Drive",
		Shortname:        "Zoo Dr",
		PrintableAddress: "Los Angeles Zoo, Zoo Drive, Los Angeles, CA, USA",
		Lat:              34.1483481,
		Lng:              -118.2840899,
		CountryCode:      "US",
		AddressTrackingData: AddressTrackingData{
			FormattedAddress: "Zoo Dr, Los Angeles, CA 90027, USA",
			GooglePlaceID:    "EiM5MTEgWm9vIERyaXZlLCBMb3MgQW5nZWxlcywgQ0EsIFVTQSIuKiwKFAoSCduxBZJ2wMKAEUrNyU4GTUkEEhQKEgkT2ifcXcfCgBH0CEYlb98v4g",
			Name:             "",
			Route:            "Zoo Drive",
			Subpremise:       "",
		},
	}
)

func (a *Account) SetAddress() error {

	payload, err := json.Marshal(GQLPayload{
		OperationName: "addConsumerAddress",
		Variables:     address,
		Query:         SetAddressQuery,
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

	if res.StatusCode != 200 {
		return fmt.Errorf("doordash: invalid setaddress status code: %q", res.Status)
	}

	return nil
}

type SetAddressVars struct {
	GooglePlaceID       string              `json:"googlePlaceId"`
	City                string              `json:"city"`
	State               string              `json:"state"`
	ZipCode             string              `json:"zipCode"`
	Street              string              `json:"street"`
	Shortname           string              `json:"shortname"`
	PrintableAddress    string              `json:"printableAddress"`
	Lat                 float64             `json:"lat"`
	Lng                 float64             `json:"lng"`
	CountryCode         string              `json:"countryCode"`
	AddressTrackingData AddressTrackingData `json:"addressTrackingData"`
}

type AddressTrackingData struct {
	FormattedAddress string `json:"formattedAddress"`
	GooglePlaceID    string `json:"googlePlaceId"`
	Name             string `json:"name"`
	Route            string `json:"route"`
	Subpremise       string `json:"subpremise"`
}
