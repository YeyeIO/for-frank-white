package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/emersion/go-imap/client"
	_ "github.com/joho/godotenv/autoload"
	"github.com/saucesteals/doordash/doordash"
	"github.com/saucesteals/doordash/emailverify"
	"github.com/saucesteals/doordash/internal/faker"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	imapAddr     = flag.String("imapAddr", os.Getenv("IMAP_ADDR"), "Imap addr for email verification")
	imapUsername = flag.String("imapUser", os.Getenv("IMAP_USERNAME"), "Imap username")
	imapPassword = flag.String("imapPass", os.Getenv("IMAP_PASSWORD"), "Imap password")
	region       = flag.String("region", os.Getenv("REGION"), "Region for all accounts")
	referral     = flag.String("referral", os.Getenv("DOORDASH_REFERRAL_ID"), "Referral ID for first account")
	password     = flag.String("doordashPassword", os.Getenv("DOORDASH_PASSWORD"), "Password for all accounts")
	outDir       = flag.String("out", os.Getenv("OUT_PATH"), "output directory")
	accountCount = 0
	chainCount   = 1
)

var (
	generatePhoneNumber func() string
	emailClient         *emailverify.Client
	xsrfToken           string
	gmailUser           string
)

func makeAccounts(count int) {
	out, err := makeOutputFile()

	if err != nil {
		log.Printf("Error while generting output file %q", err)
		return
	}

	lastRef := *referral

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < count; {
		log.Printf("Starting %d/%d", i+1, count)

		account := doordash.NewAccount(xsrfToken, emailClient, doordash.AccountDetails{
			Email:       fmt.Sprintf("%s+%s@gmail.com", gmailUser, faker.RandHex(7)),
			Password:    *password,
			PhoneNumber: generatePhoneNumber(),
			FirstName:   faker.FirstName(),
			LastName:    faker.LastName(),
			ReferralID:  lastRef,
			CountryCode: *region,
		})

		log.Printf("Email: %q; CountryCode: %q; Phone: %q",
			account.Email, account.CountryCode, account.PhoneNumber)

		if err := account.CreateAndRefer(context.Background()); err != nil {
			log.Printf("error while creating and referring: %q", err)
			continue
		}

		ref, err := account.GetSelfReferral()

		if err != nil {
			log.Printf("error while geting self referral code: %q", err)
			continue
		}

		lastRef = ref

		if err := appendAccount(out, count-i, account.AccountDetails); err != nil {
			log.Println(err)
			continue
		}
		i++
	}

}

func usage() {
	fmt.Fprint(os.Stderr, "usage: doorchain accounts chains\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	if err := ensureDir(*outDir); err != nil {
		log.Fatalln(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(chainCount)

	for i := 0; i < chainCount; i++ {
		go func() {
			makeAccounts(accountCount)
			wg.Done()
		}()
	}

	wg.Wait()
}

func init() {
	var err error

	flag.Usage = usage
	log.SetFlags(0)
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
	}

	accountCount, err = strconv.Atoi(flag.Arg(0))

	if err != nil {
		log.Fatalf("invalid accounts count: %s", err)
	}

	if accountCount < 1 {
		flag.Usage()
	}

	if err := testImap(); err != nil {
		log.Fatalln(err)
	}

	if flag.NArg() > 1 {
		chainCount, err = strconv.Atoi(flag.Arg(1))

		if err != nil {
			log.Fatalf("invalid chains count: %q", err)
		}
	}

	if *referral == "" {
		log.Fatalf("invalid referral id: %q", *referral)
	}

	switch *region {
	case "HR":
		generatePhoneNumber = faker.HrPhoneNumber
	case "US":
		generatePhoneNumber = faker.UsPhoneNumber
	case "AU":
		generatePhoneNumber = faker.AuPhoneNumber
	default:
		log.Printf("invalid region: %q", *region)
		return
	}

	emailClient = &emailverify.Client{
		Addr:     *imapAddr,
		Username: *imapUsername,
		Password: *imapPassword,
	}

	gmailUser = strings.SplitN(*imapUsername, "@", 2)[0]

	xsrfToken, err = doordash.GenerateXSRFToken()

	if err != nil {
		log.Fatalln(err)
	}

	log.SetFlags(log.Ltime | log.Lmsgprefix)

	log.SetPrefix("| ")

}

func testImap() error {
	c, err := client.DialTLS(*imapAddr, nil)

	if err != nil {
		return err
	}

	if err := c.Login(*imapUsername, *imapPassword); err != nil {
		return err
	}

	return c.Logout()
}

func makeOutputFile() (string, error) {

	if err := ensureDir(*outDir); err != nil {
		return "", err
	}

	return path.Join(*outDir, fmt.Sprintf("%s_%s.txt", time.Now().Format("Jan02_15_04_05"), faker.RandHex(3))), nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil

}

func appendAccount(out string, num int, account doordash.AccountDetails) error {

	f, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%d. %s:%s\n", num, account.Email, account.Password)); err != nil {
		return err
	}

	return nil
}
