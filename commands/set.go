package commands

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/passwd"
	"github.com/giacomocariello/nickelcase/uri"
)

// SetCommand : set a key in a nickelcase
func SetCommand(c *cli.Context) error {
	args := c.Args()
	if len(args) != 2 {
		return fmt.Errorf("Invalid number of arguments")
	}
	key := args[0]
	value := args[1]
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
	parsedData[key] = value
	return uri.WriteMapToURI(c, outputUri, parsedData, uri.WriteDataToEncryptedStream(pwd))
}
