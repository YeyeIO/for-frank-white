# - DoorDash Chain Generator -


# Usage
```shell
doorchain accounts [chains=1]
  -doordashPassword string
        Password for all accounts (default env:DOORDASH_PASSWORD)
  -imapAddr string
        Imap addr for email verification (default env:IMAP_ADDR)
  -imapPass string
        Imap password (default env:IMAP_PASSWORD)
  -imapUser string
        Imap username (default env:IMAP_USER)
  -out string
        Output directory (default env:OUT_DIR)
  -referral string
        Referral ID for first account (default env:REFERRAL_ID)
  -region string
        Region for all accounts (default env:REGION)
```

- All of these command line options can be set through environment variables (example [here](https://github.com/saucesteals/doordash/blob/main/env.example))
- You can optionally create a `.env` file in the same directory as the binary (or wherever you are planning to run it from) and follow the example from [here](https://github.com/saucesteals/doordash/blob/main/env.example)
- You will still need to provide an argument for the amount of accounts (and chains if needed) when executing the binary


Example: *Creating 3 chains with 30 accounts in each one*
```sh
doorchain 30 3 
```


# Installation

## Github

- Grab a release directly from [Github Releases](https://github.com/saucesteals/doordash/releases)
- [Un-tar](https://www.google.com/search?q=how+to+open+a+tar.gz+file) & move the `doorchain` binary to a proper location (ex. `/usr/local/bin`)
- If your shell can't find `doorchain`, add the location to $PATH

## Manual

---

_Prerequisites_

- https://go.dev/doc/install

---

- Pull doordash from the repo

```sh
go get -d github.com/saucesteals/doordash
```

- Build the doorchain binary

```sh
go install github.com/saucesteals/doordash/cmd/doorchain
```
