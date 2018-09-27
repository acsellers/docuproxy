package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/acsellers/docusign"
)

var (
	svc  *docusign.Service
	demo *docusign.Template
	user = flag.String("username", "", "UserName for docusign")
	pass = flag.String("password", "", "Password for docusign")
	acct = flag.String("accountid", "", "AccountID for docusign")
	key  = flag.String("key", "", "Integrator Key for docusign")
	host = flag.String("host", "demo.docusign.net", "Host (demo.docusign.net or www.docusign.net)")
)

func Connect() {
	cfg := &docusign.Config{
		UserName:      *user,
		Password:      *pass,
		IntegratorKey: *key,
		AccountId:     *acct,
		Host:          *host,
	}
	svc = docusign.New(cfg, "")
	if !strings.Contains(*host, "demo") {
		li, err := svc.LoginInformation(
			context.Background(),
			docusign.LoginInfoParam{Name: "api_password", Value: "true"},
		)
		if err != nil {
			log.Fatal("Login Info: ", err)
		}
		if len(li.LoginAccounts) == 0 {
			log.Fatal("No Login Accounts Available?")
		}
		u, _ := url.Parse(li.LoginAccounts[0].BaseUrl)
		svc = docusign.New(&docusign.Config{
			UserName:      li.LoginAccounts[0].UserId,
			Password:      li.ApiPassword,
			IntegratorKey: *key,
			AccountId:     li.LoginAccounts[0].AccountId,
			Host:          u.Host,
		}, "")
	}
}

func main() {
	flag.Parse()
	Connect()

	env := &docusign.Envelope{}
	err := json.NewDecoder(os.Stdin).Decode(env)
	if err != nil {
		log.Fatal("Cannot decode input: ", err)
	}
	saved, err := svc.EnvelopeCreate(
		context.Background(),
		env,
	)
	if err != nil {
		fmt.Println(saved)
		log.Fatal("Env Create: ", err)
	}
	json.NewEncoder(os.Stdout).Encode(saved)
}
