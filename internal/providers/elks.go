package providers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/sakjur/firepoker/sms"
)

const elksApi = "https://api.46elks.com/a1/%s"

type Elks struct {
	Sender string `toml:"sender"`
	Key    string `toml:"key"`
	Secret string `toml:"secret"`
}

type elksResponse struct {
	Status  string `json:"status"`
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
	Cost    int    `json:"cost"`
	Parts   int    `json:"parts"`
}

func (e Elks) Send(message sms.Message) error {
	err := message.Valid()
	if err != nil {
		return err
	}

	return e.send(message)
}

func (e Elks) send(m sms.Message) error {
	d := url.Values{
		"from":    []string{e.Sender},
		"to":      []string{string(m.Target)},
		"message": []string{m.Content},
	}.Encode()
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(elksApi, "sms"), strings.NewReader(d))
	if err != nil {
		return err
	}

	basicAuth := fmt.Sprintf("%s:%s", e.Key, e.Secret)
	req.Header.Add(
		"authorization",
		fmt.Sprintf("basic %s", base64.RawStdEncoding.EncodeToString([]byte(basicAuth))),
	)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode-res.StatusCode%100 == 200 {
		payload := &elksResponse{}
		err := json.NewDecoder(res.Body).Decode(payload)
		if err != nil {
			log.Println("failed decoding response: %v", err)
		}

		log.Printf("%s -> %s :: %d parts, cost %d", payload.From, payload.To, payload.Parts, payload.Cost)
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		log.Printf("%s", body)
	}

	return nil
}
