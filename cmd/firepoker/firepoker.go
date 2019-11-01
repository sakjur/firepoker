package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sakjur/firepoker/sms"

	"github.com/sakjur/firepoker/config"
)

func main() {
	c := flag.String("c", "/etc/firepoker/config.toml", "Destination of the configuration file")
	recipient := flag.String("t", "", "Recipient of the SMS")
	flag.Parse()
	args := flag.Args()

	if *recipient == "" {
		log.Fatalf("recipient is empty, use -t to set a recipient")
	}

	message := strings.Join(args, " ")

	if message == "" {
		msg, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read message from stdin")
		}

		message = string(msg)
	}

	f, err := os.Open(*c)
	if err != nil {
		log.Fatalf("got error when reading file %s: %v", *c, err)
	}

	cfg, err := config.Read(f)
	if err != nil {
		log.Fatalf("got error when reading configuration: %v", err)
	}

	err = cfg.Providers.Elks.Send(
		sms.Message{
			Content: message,
			Target:  sms.Phonenumber(*recipient),
		},
	)
	if err != nil {
		log.Fatalf("got error when sending sms: %v", err)
	}
}
