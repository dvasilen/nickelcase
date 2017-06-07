package uri

import (
	"net/http"
	"time"

	"github.com/urfave/cli"
)

func getHTTPClient(c *cli.Context) *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}
