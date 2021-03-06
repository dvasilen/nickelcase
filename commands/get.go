package commands

import (
	"bytes"
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/passwd"
	"github.com/giacomocariello/nickelcase/uri"
)

// GetCommand : get value of a nickelcase map key
func GetCommand(c *cli.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return fmt.Errorf("Invalid number of arguments")
	}
	pwd, err := passwd.GetPassword(c, c.String("password"))
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
	for _, key := range args {
		value, ok := parsedData[key]
		if !ok {
			fmt.Fprintf(os.Stderr, "Key %s is not present in the archive\n", key)
		} else if _, err = fmt.Fprintf(buf, "%s\n", value); err != nil {
			return err
		}
	}
	return uri.WriteDataToURI(c, c.String("output"), buf.Bytes(), uri.WriteDataToPlaintextStream)
}
