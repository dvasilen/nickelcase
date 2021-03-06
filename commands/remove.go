package commands

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/passwd"
	"github.com/giacomocariello/nickelcase/uri"
)

// RemoveCommand : remove keys from nickelcases.
func RemoveCommand(c *cli.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return fmt.Errorf("Invalid number of arguments")
	}
	pwd, err := passwd.GetPassword(c, c.String("password"))
	if err != nil {
		return err
	}
	var outputUri string
	parsedData := make(map[string]interface{})
	if len(c.String("file")) > 0 {
		outputUri = c.String("file")
		err = uri.ReadMapFromURI(c, outputUri, uri.ReadDataFromEncryptedStream(pwd), parsedData)
		if err != nil {
			return err
		}
	} else {
		outputUri = c.String("output")
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
	}
	for _, key := range args {
		delete(parsedData, key)
	}
	return uri.WriteMapToURI(c, outputUri, parsedData, uri.WriteDataToEncryptedStream(pwd))
}
