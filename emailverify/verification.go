package emailverify

import (
	"context"
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"io"
	"regexp"
)

var (
	errCodeNotFound = errors.New("email-verify: code not found in email")
	codeRe          = regexp.MustCompile(`>[0-9]{6}<`)
)

type Client struct {
	Username string
	Password string
	Addr     string
}

func (v *Client) GetLatestCode(ctx context.Context, criteria *imap.SearchCriteria) (string, error) {

	c, err := client.DialTLS(v.Addr, nil)

	if err != nil {
		return "", err
	}

	defer c.Logout()

	if err := c.Login(v.Username, v.Password); err != nil {
		return "", err
	}

	var sequences []uint32

out:
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:

			if _, err := c.Select("INBOX", false); err != nil {
				return "", err
			}

			sequences, err = c.Search(criteria)

			if err != nil {
				return "", err
			}

			if len(sequences) != 0 {
				break out
			}
		}
	}

	seqset := &imap.SeqSet{}
	seqset.AddNum(sequences[len(sequences)-1])

	messages := make(chan *imap.Message, 1)

	if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}, messages); err != nil {
		return "", err
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case msg := <-messages:
		for _, buff := range msg.Body {
			body, err := io.ReadAll(buff)
			if err != nil {
				return "", err
			}
			code := codeRe.Find(body)
			if len(code) == 0 {
				continue
			}

			return string(code[1:7]), nil
		}
	}

	return "", errCodeNotFound
}
