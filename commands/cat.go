package commands

import (
	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/passwd"
	"github.com/giacomocariello/nickelcase/uri"
)

func CatCommand(c *cli.Context) error {
	var pwd string
	var err error
	plaintextSources := c.StringSlice("plaintext-input")
	encryptedSources := c.StringSlice("encrypted-input")
	isLoad := c.Command.FullName() == "load"
	if len(encryptedSources) > 0 {
		pwd, err = passwd.GetPassword(c.String("password"))
		if err != nil {
			return err
		}
	} else if isLoad {
		pwd, err = passwd.GetNewPassword(c.String("password"))
		if err != nil {
			return err
		}
	}
	parsedData := make(map[string]interface{})
	for _, src := range encryptedSources {
		err = uri.ReadMapFromURI(src, uri.ReadDataFromEncryptedStream(pwd), parsedData)
		if err != nil {
			return err
		}
	}
	for _, src := range plaintextSources {
		err = uri.ReadMapFromURI(src, uri.ReadDataFromPlaintextStream(), parsedData)
		if err != nil {
			return err
		}
	}
	if len(plaintextSources) == 0 && len(encryptedSources) == 0 {
		if isLoad {
			err = uri.ReadMapFromURI("", uri.ReadDataFromPlaintextStream(), parsedData)
			if err != nil {
				return err
			}
		} else {
			err = uri.ReadMapFromURI("", uri.ReadDataFromEncryptedStream(pwd), parsedData)
			if err != nil {
				return err
			}
		}
	}
	if isLoad {
		return uri.WriteMapToURI(c.Args().Get(0), parsedData, uri.WriteDataToEncryptedStream(pwd))
	} else {
		return uri.WriteMapToURI(c.Args().Get(0), parsedData, uri.WriteDataToPlaintextStream)
	}
}
