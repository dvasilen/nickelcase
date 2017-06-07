package commands

import (
	"bytes"
	"fmt"

	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/passwd"
	"github.com/giacomocariello/nickelcase/uri"
)

func ListCommand(c *cli.Context) error {
	pwd, err := passwd.GetPassword(c.String("password"))
	if err != nil {
		return err
	}
	parsedData := make(map[string]interface{})
	sources := c.StringSlice("encrypted-input")
	if len(sources) > 0 {
		for _, src := range sources {
			err = uri.ReadMapFromURI(c, src, uri.ReadDataFromEncryptedStream(pwd), parsedData)
			if err != nil {
				return err
			}
		}
	} else {
		err = uri.ReadMapFromURI(c, "", uri.ReadDataFromEncryptedStream(pwd), parsedData)
		if err != nil {
			return err
		}
	}
	buf := new(bytes.Buffer)
	for key := range parsedData {
		if _, err = fmt.Fprintf(buf, "%s\n", key); err != nil {
			return err
		}
	}
	return uri.WriteDataToURI(c, c.String("output"), buf.Bytes(), uri.WriteDataToPlaintextStream)
}
