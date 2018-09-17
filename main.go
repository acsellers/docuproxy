package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/acsellers/docusign"
)

var (
	svc  *docusign.Service
	demo *docusign.Template
	user = flag.String("username", "", "UserName for docusign")
	pass = flag.String("password", "", "Password for docusign")
	acct = flag.String("accountid", "", "AccountID for docusign")
	key  = flag.String("key", "", "Integrator Key for docusign")
	host = flag.String("host", "demo.docusign.net", "Host (demo or prod)")
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
}

func main() {
	flag.Parse()
	Connect()

	env := &docusign.Envelope{}
	err := json.NewDecoder(os.Stdin).Decode(env)
	if err != nil {
		log.Fatal("Cannot decode input")
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
