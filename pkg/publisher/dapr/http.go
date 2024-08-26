package publisher

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"
)

type DaprHttp struct {
	client              *http.Client
	host                string
	httpPort            string
	pubsubComponentName string
	pubsubTopic         string
	url                 string
}

func (d DaprHttp) Publish(ctx context.Context, msg []byte) error {
	if d.client == nil {
		d.client = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	if d.url == "" {
		d.url = d.host + ":" + d.httpPort +
			"/v1.0/publish/" + d.pubsubComponentName + "/" + d.pubsubTopic
	}

	res, err := d.client.Post(d.url, "application/json", bytes.NewReader(msg))
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	defer res.Body.Close()

	return nil
}
