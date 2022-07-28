package doordash

import (
	"context"
	"github.com/emersion/go-imap"
	"github.com/saucesteals/doordash/emailverify"
	"github.com/saucesteals/doordash/internal/simplejar"
	http "github.com/saucesteals/fhttp"
)

type AccountDetails struct {
	Email, Password, ReferralID string
	PhoneNumber, CountryCode    string
	FirstName, LastName         string
}

type Account struct {
	AccountDetails
	Imap      *emailverify.Client
	xsrfToken string
	client    *http.Client
	jar       *simplejar.Jar
}

const clientId = "1666519390426295040"
const authState = "5965aa35-1138-47ad-a1f0-69ee06b48b54"

func NewAccount(xsrfToken string, imap *emailverify.Client, details AccountDetails) *Account {
	a := &Account{AccountDetails: details, Imap: imap, xsrfToken: xsrfToken}

	jar := simplejar.NewJar()

	jar.SetCookies(nil, []*http.Cookie{
		{Name: "XSRF-TOKEN", Value: xsrfToken},
		{Name: "authState", Value: authState},
	})

	a.client = newHttpClient()

	a.client.Jar = jar

	return a
}

func (a *Account) CreateAndRefer(ctx context.Context) error {

	if err := a.SignUp(); err != nil {
		return err
	}

	if err := a.GeneratePasscode(); err != nil {
		return err
	}

	criteria := imap.NewSearchCriteria()
	criteria.Header.Add("SUBJECT", "Your verification code")
	criteria.Header.Add("TO", a.Email)

	code, err := a.Imap.GetLatestCode(ctx, criteria)

	if err != nil {
		return err
	}

	if err := a.VerifyPasscode(code); err != nil {
		return err
	}

	if err := a.CreateReferral(); err != nil {
		return err
	}

	if err := a.SetAddress(); err != nil {
		return err
	}

	return nil
}
